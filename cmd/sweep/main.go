package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/xyproto/synth"
)

var (
	version     = "0.0.1"
	sampleRate  int
	duration    time.Duration
	baseFreq    float64
	showVersion bool
	showHelp    bool
)

func init() {
	flag.BoolVar(&showVersion, "version", false, "Show version information")
	flag.BoolVar(&showHelp, "help", false, "Show help information")
	flag.IntVar(&sampleRate, "samplerate", 44100, "Sample rate (in Hz)")
	flag.DurationVar(&duration, "duration", 10*time.Second, "Duration of the audio (e.g., 10s, 5m)")
	flag.Float64Var(&baseFreq, "freq", 55.0, "Base frequency for the bass sound (in Hz)")
}

func main() {
	flag.Parse()

	if showVersion {
		fmt.Printf("Bass Synth Generator, version %s\n", version)
		os.Exit(0)
	}

	if showHelp {
		flag.Usage()
		os.Exit(0)
	}

	// Open the output file for writing as an io.WriteSeeker
	outFile, err := os.Create("sweep.wav")
	if err != nil {
		log.Fatalf("Failed to create output file: %v", err)
	}
	defer outFile.Close()
	var out io.WriteSeeker = outFile

	// Calculate the length of the waveform
	length := sampleRate * int(duration.Seconds())

	// Detune settings for the oscillators
	detune := []float64{-0.01, -0.005, 0.0, 0.005, 0.01}

	// Generate detuned sawtooth oscillators
	bassOscillators := synth.DetunedOscillators(baseFreq, detune, length, sampleRate)

	// Apply an ADSR envelope
	env := synth.ApplyEnvelope(bassOscillators, 0.1, 0.4, 0.6, 0.7, sampleRate)

	// Apply a low-pass filter to smooth the high frequencies
	filtered := synth.LowPassFilter(env, 200, sampleRate)

	// Apply drive (distortion)
	driven := synth.Drive(filtered, 1.2)

	// Apply a limiter to prevent clipping
	limited := synth.Limiter(driven)

	// Save the generated sound to a .wav file using the io.WriteSeeker
	if err := synth.SaveToWav(out, limited, sampleRate); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println("Successfully generated 'sweep.wav'")
}
