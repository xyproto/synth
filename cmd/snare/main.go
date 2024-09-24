// cmd/snare/main.go

package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/xyproto/synth"
)

const version = "0.0.1"

func main() {
	// Snare style flags
	snare606 := flag.Bool("606", false, "Generate a snare.wav like a 606 snare drum")
	snare707 := flag.Bool("707", false, "Generate a snare.wav like a 707 snare drum")
	snare808 := flag.Bool("808", false, "Generate a snare.wav like an 808 snare drum")
	snare909 := flag.Bool("909", false, "Generate a snare.wav like a 909 snare drum")
	snareLinnDrum := flag.Bool("linn", false, "Generate a snare.wav like a LinnDrum snare drum")
	snareExperimental := flag.Bool("experimental", false, "Generate a snare.wav with experimental-style characteristics")

	// Sound customization flags
	noiseType := flag.String("noise", "white", "Type of noise to mix in (white, pink, brown)")
	noiseAmount := flag.Float64("noiseamount", 0.5, "Amount of noise to mix in (0.0 to 1.0)")
	length := flag.Float64("length", 500, "Length of the snare drum sample in milliseconds")
	quality := flag.Int("quality", 96, "Sample rate in kHz (44, 48, 96, or 192)")
	bitDepth := flag.Int("bitdepth", 16, "Bit depth of the audio (8, 16, 24 or 32)")
	channels := flag.Int("channels", 1, "Channels (1 or 2)")
	waveform := flag.Int("waveform", synth.WaveSine, "Waveform type (0: Sine, 1: Triangle, 2: Sawtooth, 3: Square)")
	attack := flag.Float64("attack", 0.005, "Attack time in seconds")
	decay := flag.Float64("decay", 0.2, "Decay time in seconds")
	sustain := flag.Float64("sustain", 0.1, "Sustain level (0.0 to 1.0)")
	release := flag.Float64("release", 0.1, "Release time in seconds")
	filterCutoff := flag.Float64("filter", 8000.0, "Filter cutoff frequency (Hz)")
	drive := flag.Float64("drive", 0.2, "Amount of distortion/drive")
	outputFile := flag.String("o", "snare.wav", "Output file path")
	playSnare := flag.Bool("p", false, "Play the generated snare")
	showVersion := flag.Bool("version", false, "Show the current version")
	showHelp := flag.Bool("help", false, "Display this help")

	flag.Parse()

	if *showHelp {
		flag.Usage()
		return
	}

	if *showVersion {
		fmt.Println("snare version", version)
		return
	}

	var sampleRate int
	switch *quality {
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

	var cfg *synth.Settings
	var err error
	switch {
	case *snare606:
		cfg, err = synth.New606Snare(nil, *length/1000.0, sampleRate, *bitDepth, *channels)
		fmt.Println("Generating 606 snare with a sharp, crisp sound.")
	case *snare707:
		cfg, err = synth.New707Snare(nil, *length/1000.0, sampleRate, *bitDepth, *channels)
		fmt.Println("Generating 707 snare with a classic electronic sound.")
	case *snare808:
		cfg, err = synth.New808Snare(nil, *length/1000.0, sampleRate, *bitDepth, *channels)
		fmt.Println("Generating 808 snare with a deep, resonant sound.")
	case *snare909:
		cfg, err = synth.New909Snare(nil, *length/1000.0, sampleRate, *bitDepth, *channels)
		fmt.Println("Generating 909 snare with a punchy, bright sound.")
	case *snareLinnDrum:
		cfg, err = synth.NewLinnDrumSnare(nil, *length/1000.0, sampleRate, *bitDepth, *channels)
		fmt.Println("Generating LinnDrum snare with a classic digital sound.")
	case *snareExperimental:
		cfg, err = synth.NewExperimentalSnare(nil, *length/1000.0, sampleRate, *bitDepth, *channels)
		fmt.Println("Generating experimental-style snare with unique characteristics.")
	default:
		cfg, err = synth.NewSettings(nil, 200.0, 100.0, *length/1000.0, sampleRate, *bitDepth, *channels)
		fmt.Println("Generating default snare with user-defined characteristics.")
	}

	if err != nil {
		fmt.Println("Error creating config:", err)
		os.Exit(1)
	}

	// Apply user-defined settings
	cfg.WaveformType = *waveform
	cfg.Attack = *attack
	cfg.Decay = *decay
	cfg.Sustain = *sustain
	cfg.Release = *release
	cfg.FilterCutoff = *filterCutoff
	cfg.Drive = *drive
	cfg.SmoothFrequencyTransitions = true

	// Generate the snare drum waveform
	samples, err := cfg.GenerateSnare()
	if err != nil {
		fmt.Println("Failed to generate snare:", err)
		return
	}

	// Generate noise samples if needed
	if strings.ToLower(*noiseType) != "none" && *noiseAmount > 0.0 {
		numSamples := len(samples)
		var noiseSamples []float64

		switch strings.ToLower(*noiseType) {
		case "white":
			noiseSamples = synth.GenerateWhiteNoise(numSamples, *noiseAmount)
		case "pink":
			noiseSamples = synth.GeneratePinkNoise(numSamples, *noiseAmount)
		case "brown":
			noiseSamples = synth.GenerateBrownNoise(numSamples, *noiseAmount)
		default:
			fmt.Println("Invalid noise type. Choose from: white, pink, brown.")
			os.Exit(1)
		}

		// Apply band-pass filter to noise to shape it for snare characteristics
		noiseSamples = synth.BandPassFilter(noiseSamples, 150.0, 8000.0, sampleRate)

		// Mix noise samples into the snare samples
		for i := 0; i < numSamples; i++ {
			samples[i] += noiseSamples[i]
		}

		// Apply limiter to prevent clipping
		samples = synth.Limiter(samples)
	}

	// Apply a fade-out at the end to prevent crackling noise
	samples = synth.ApplyFadeOut(samples, cfg.Release, sampleRate, synth.QuadraticFade)

	// Open the output file for writing
	outFile, err := os.Create(*outputFile)
	if err != nil {
		fmt.Println("Failed to create output file:", err)
		return
	}
	defer outFile.Close()

	// Save the waveform to the output file
	if err := synth.SaveToWav(outFile, samples, sampleRate, *bitDepth, *channels); err != nil {
		fmt.Println("Failed to save snare to file:", err)
		return
	}

	fmt.Println("Snare drum sound generated and written to", *outputFile)

	// Play the snare if -p flag is provided
	if *playSnare {
		fmt.Println("Playing the generated snare drum sound...")
		// Use PlayWaveform to play the samples directly
		player := synth.NewPlayer()
		defer player.Close()
		if audioDeviceKey, playbackDuration, err := player.PlayWaveform(samples, sampleRate, *bitDepth, *channels); err != nil {
			fmt.Println("Failed to play snare:", err)
		} else {
			player.WaitClosePlus100(audioDeviceKey, playbackDuration)
		}
	}

}
