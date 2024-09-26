package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/xyproto/playsample"
	"github.com/xyproto/synth"
)

const version = "0.0.1"

const bitDepth = 16

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
	combined, sampleRate, err := playsample.LoadWav(firstFile, true)
	if err != nil {
		log.Fatalf("Failed to load %s: %v", firstFile, err)
	}

	// Initialize loudest peak
	loudestPeak := synth.FindPeakAmplitude(combined)

	// Process additional files and mix them using LinearSummation
	for _, inputFile := range inputFiles[1:] { // Start from the second file
		// Load the next file
		wave, sr, err := playsample.LoadWav(inputFile, true)
		if err != nil {
			log.Fatalf("Failed to load %s: %v", inputFile, err)
		}

		// Ensure the sample rate matches
		if sr != sampleRate {
			log.Fatalf("Sample rate mismatch between %s and %s", firstFile, inputFile)
		}

		// Find the peak amplitude in the current file and track the loudest peak
		peak := synth.FindPeakAmplitude(wave)
		if peak > loudestPeak {
			loudestPeak = peak
		}

		// Pad the shorter sample with zeros
		combined, wave = synth.PadSamples(combined, wave)

		// Mix the current combined samples with the newly loaded samples using LinearSummation
		combined, err = synth.LinearSummation(combined, wave)
		if err != nil {
			log.Fatalf("Error during linear summation mixing of %s: %v", inputFile, err)
		}
	}

	// Normalize the final combined samples to the loudest input sample's peak
	fmt.Printf("Normalizing loudness to the loudest peak: %f\n", loudestPeak)
	combined = synth.NormalizeSamples(combined, loudestPeak)

	// Apply a quick fade-out to the end of the combined samples
	fadeDuration := 0.01 // Fade-out duration in seconds (10 milliseconds)
	combined = synth.ApplyQuadraticFadeOut(combined, fadeDuration, sampleRate)

	// Open the output file for writing
	outFile, err := os.Create(*outputFile)
	if err != nil {
		log.Fatalf("Failed to create output file %s: %v", *outputFile, err)
	}
	defer outFile.Close()

	// Save the final combined result to the output file using an io.WriteSeeker
	const channels = 1
	if err := playsample.SaveToWav(outFile, combined, sampleRate, bitDepth, channels); err != nil {
		log.Fatalf("Failed to save %s: %v", *outputFile, err)
	}

	fmt.Printf("Successfully mixed %d files into %s\n", len(inputFiles), *outputFile)
}
