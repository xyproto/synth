package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/xyproto/synth"
)

var (
	version     = "0.0.1"
	quality     int
	bitDepth    int
	channels    int
	duration    time.Duration
	baseFreq    float64
	showVersion bool
	showHelp    bool
	playSound   bool // Added -p flag variable
)

func main() {
	// Initialize flags in the main function
	flag.BoolVar(&showVersion, "version", false, "Show version information")
	flag.BoolVar(&showHelp, "help", false, "Show help information")
	flag.BoolVar(&playSound, "p", false, "Play the generated sound") // Added -p flag
	flag.IntVar(&quality, "quality", 96, "Sample rate in kHz (44, 48, 96, or 192)")
	flag.IntVar(&bitDepth, "bitdepth", 16, "Bit depth of the audio (16 or 24)")
	flag.IntVar(&channels, "channels", 1, "Channels (1 or 2)")

	flag.DurationVar(&duration, "duration", 10*time.Second, "Duration of the audio (e.g., 10s, 5m)")
	flag.Float64Var(&baseFreq, "freq", 55.0, "Base frequency for the bass sound (in Hz)")

	flag.Parse()

	if showVersion {
		fmt.Printf("Bass Synth Generator, version %s\n", version)
		os.Exit(0)
	}

	if showHelp {
		flag.Usage()
		os.Exit(0)
	}

	var sampleRate int
	switch quality {
	case 44:
		sampleRate = 44100
	case 48:
		sampleRate = 48000
	case 96:
		sampleRate = 96000
	case 192:
		sampleRate = 192000
	default:
		fmt.Println("Invalid sample rate. Choose 44, 48, 96, or 192 kHz.")
		os.Exit(1)
	}

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

	// Save the generated sound to a .wav file
	outFile, err := os.Create("sweep.wav")
	if err != nil {
		log.Fatalf("Failed to create output file: %v", err)
	}
	defer outFile.Close()

	if err := synth.SaveToWav(outFile, limited, sampleRate, bitDepth, channels); err != nil {
		fmt.Printf("Error saving WAV file: %v\n", err)
		return
	}

	fmt.Println("Successfully generated 'sweep.wav'")

	// Play the sound if -p flag is provided
	if playSound {
		fmt.Println("Playing the generated sound...")
		player := synth.NewPlayer()
		defer player.Close()
		if audioDeviceKey, playbackDuration, err := player.PlayWaveform(limited, sampleRate, bitDepth, channels); err != nil {
			fmt.Printf("Error playing sound: %v\n", err)
		} else {
			player.WaitClosePlus100(audioDeviceKey, playbackDuration)
		}
	}
}
