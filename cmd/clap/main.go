// cmd/clap/main.go

package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/xyproto/playsample"
	"github.com/xyproto/synth"
)

const version = "0.0.1"

func main() {
	// Clap style flags
	clap606 := flag.Bool("606", false, "Generate a clap.wav like a 606 clap")
	clap707 := flag.Bool("707", false, "Generate a clap.wav like a 707 clap")
	clap808 := flag.Bool("808", false, "Generate a clap.wav like an 808 clap")
	clap909 := flag.Bool("909", false, "Generate a clap.wav like a 909 clap")
	clapLinnDrum := flag.Bool("linn", false, "Generate a clap.wav like a LinnDrum clap")
	clapExperimental := flag.Bool("experimental", false, "Generate a clap.wav with experimental-style characteristics")

	// Sound customization flags
	noiseAmount := flag.Float64("noiseamount", 0.5, "Amount of noise to mix in (0.0 to 1.0)")
	length := flag.Float64("length", 2000, "Length of the clap sample in milliseconds")
	quality := flag.Int("quality", 44, "Sample rate in kHz (44, 48, 96, or 192)")
	bitDepth := flag.Int("bitdepth", 16, "Bit depth of the audio (8, 16, 24 or 32)")
	channels := flag.Int("channels", 1, "Channels (1 or 2)")
	waveform := flag.Int("waveform", synth.WaveSquare, "Waveform type (0: Sine, 1: Triangle, 2: Sawtooth, 3: Square)")
	attack := flag.Float64("attack", 0.005, "Attack time in seconds")
	decay := flag.Float64("decay", 0.1, "Decay time in seconds")
	sustain := flag.Float64("sustain", 0.0, "Sustain level (0.0 to 1.0)")
	release := flag.Float64("release", 0.05, "Release time in seconds")
	filterCutoff := flag.Float64("filter", 5000.0, "Filter cutoff frequency (Hz)")
	drive := flag.Float64("drive", 0.3, "Amount of distortion/drive")
	pitchDecay := flag.Float64("pitchdecay", 0.3, "Pitch envelope decay")
	outputFile := flag.String("o", "clap.wav", "Output file path")
	playClap := flag.Bool("p", false, "Play the generated clap")
	showVersion := flag.Bool("version", false, "Show the current version")
	showHelp := flag.Bool("help", false, "Display this help")

	flag.Parse()

	if *showHelp {
		flag.Usage()
		return
	}

	if *showVersion {
		fmt.Println("clap version", version)
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
	case *clap606:
		cfg, err = synth.New606(synth.Clap, nil, *length/1000.0, sampleRate, *bitDepth, *channels)
		fmt.Println("Generating 606 clap with a sharp, crisp sound.")
	case *clap707:
		cfg, err = synth.New707(synth.Clap, nil, *length/1000.0, sampleRate, *bitDepth, *channels)
		fmt.Println("Generating 707 clap with a classic electronic sound.")
	case *clap808:
		cfg, err = synth.New808(synth.Clap, nil, *length/1000.0, sampleRate, *bitDepth, *channels)
		fmt.Println("Generating 808 clap with a deep, resonant sound.")
	case *clap909:
		cfg, err = synth.New909(synth.Clap, nil, *length/1000.0, sampleRate, *bitDepth, *channels)
		fmt.Println("Generating 909 clap with a punchy, bright sound.")
	case *clapLinnDrum:
		cfg, err = synth.NewLinn(synth.Clap, nil, *length/1000.0, sampleRate, *bitDepth, *channels)
		fmt.Println("Generating LinnDrum clap with a classic digital sound.")
	case *clapExperimental:
		cfg, err = synth.NewExperimental(synth.Clap, nil, *length/1000.0, sampleRate, *bitDepth, *channels)
		fmt.Println("Generating experimental-style clap with unique characteristics.")
	default:
		cfg, err = synth.NewClapSettings(nil, sampleRate, *bitDepth, *channels)
		fmt.Println("Generating default clap with user-defined characteristics.")
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
	cfg.PitchDecay = *pitchDecay

	// Generate the clap waveform
	samples, err := cfg.GenerateClap()
	if err != nil {
		fmt.Println("Failed to generate clap:", err)
		return
	}

	// Generate noise samples if needed
	if *noiseAmount > 0.0 {
		numSamples := len(samples)

		noiseSamples := synth.GenerateWhiteNoise(numSamples, *noiseAmount)

		// Apply band-pass filter to noise to shape it for clap characteristics
		noiseSamples = synth.BandPassFilter(noiseSamples, 150.0, 8000.0, sampleRate)

		// Mix noise samples into the clap samples
		for i := 0; i < numSamples; i++ {
			samples[i] += noiseSamples[i]
		}

		// Apply limiter to prevent clipping
		samples = synth.Limiter(samples)
	}

	// Apply a fade-out at the end to prevent crackling noise
	samples = synth.ApplyQuadraticFadeOut(samples, cfg.Release, sampleRate)

	// Open the output file for writing
	outFile, err := os.Create(*outputFile)
	if err != nil {
		fmt.Println("Failed to create output file:", err)
		return
	}
	defer outFile.Close()

	// Save the waveform to the output file
	if err := playsample.SaveToWav(outFile, samples, sampleRate, *bitDepth, *channels); err != nil {
		fmt.Println("Failed to save clap to file:", err)
		return
	}

	fmt.Println("Clap drum sound generated and written to", *outputFile)

	// Play the clap if -p flag is provided
	if *playClap {
		fmt.Println("Playing the generated clap sound...")
		player := playsample.NewPlayer()
		defer player.Close()
		if err := player.PlayWaveform(samples, sampleRate, *bitDepth, *channels); err != nil {
			fmt.Println("Failed to play clap:", err)
			return
		}
		// TODO: Figure out why this is needed
		time.Sleep(time.Duration(int64(*length * 1e6)))
	}

}
