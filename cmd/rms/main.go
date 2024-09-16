package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/xyproto/synth"
)

const version = "0.0.1"

func main() {
	// Define flags
	outputFile := flag.String("o", "combined.wav", "Specify the output file")
	showVersion := flag.Bool("version", false, "Show the version and exit")
	showHelp := flag.Bool("help", false, "Show help")

	// Parse flags
	flag.Parse()

	// Show version and exit if --version is passed
	if *showVersion {
		fmt.Printf("rms version %s\n", version)
		return
	}

	// Show help and exit if --help is passed
	if *showHelp {
		flag.Usage()
		return
	}

	// Expect at least two input files
	if flag.NArg() < 2 {
		fmt.Println("Usage: rms [options] <input1.wav> <input2.wav> [additional input files...]")
		flag.Usage()
		return
	}

	// Load the first input file to initialize the combined samples and sample rate
	inputFiles := flag.Args()
	firstFile := inputFiles[0]
	combined, sampleRate, err := synth.LoadWav(firstFile, true)
	if err != nil {
		log.Fatalf("Failed to load %s: %v", firstFile, err)
	}

	// Find the highest frequency across all files (initialize to 0) and track the loudest sample
	highestFrequency := 0.0
	loudestPeak := synth.FindPeakAmplitude(combined)

	// Process additional files and mix them using RMSMixing
	for _, inputFile := range inputFiles {
		// Load the next file
		wave, sr, err := synth.LoadWav(inputFile, true)
		if err != nil {
			log.Fatalf("Failed to load %s: %v", inputFile, err)
		}

		// Ensure the sample rate matches
		if sr != sampleRate {
			log.Fatalf("Sample rate mismatch between %s and %s", firstFile, inputFile)
		}

		// Determine the highest frequency in the current file
		currentHighestFrequency := synth.AnalyzeHighestFrequency(wave, sr)
		if currentHighestFrequency > highestFrequency {
			highestFrequency = currentHighestFrequency
		}

		// Find the peak amplitude in the current file and track the loudest peak
		peak := synth.FindPeakAmplitude(wave)
		if peak > loudestPeak {
			loudestPeak = peak
		}

		// Pad the shorter sample with zeros
		combined, wave = synth.PadSamples(combined, wave)

		// Mix the current combined samples with the newly loaded samples using RMSMixing
		combined, err = synth.RMSMixing(combined, wave)
		if err != nil {
			log.Fatalf("Error during RMS mixing of %s: %v", inputFile, err)
		}
	}

	// Apply low-pass filter using the highest detected frequency
	fmt.Printf("Applying low-pass filter with cutoff frequency: %.2f Hz\n", highestFrequency)
	combined = synth.LowPassFilter(combined, highestFrequency, sampleRate)

	// Normalize the final combined samples to the loudest input sample's peak
	fmt.Printf("Normalizing loudness to the loudest peak: %f\n", loudestPeak)
	combined = synth.NormalizeSamples(combined, loudestPeak)

	// Save the final combined result to the output file
	if err := synth.SaveToWav(*outputFile, combined, sampleRate); err != nil {
		log.Fatalf("Failed to save %s: %v", *outputFile, err)
	}

	fmt.Printf("Successfully mixed %d files into %s\n", len(inputFiles), *outputFile)
}
