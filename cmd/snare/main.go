// cmd/snare/main.go

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
	// Snare style flags
	snare606 := flag.Bool("606", false, "Generate a snare.wav like a 606 snare drum")
	snare707 := flag.Bool("707", false, "Generate a snare.wav like a 707 snare drum")
	snare808 := flag.Bool("808", false, "Generate a snare.wav like an 808 snare drum")
	snare909 := flag.Bool("909", false, "Generate a snare.wav like a 909 snare drum")
	snareLinn := flag.Bool("linn", false, "Generate a snare.wav like a LinnDrum snare")
	snareExperimental := flag.Bool("experimental", false, "Generate a snare.wav with experimental-style characteristics")

	// Sound customization flags
	noiseAmount := flag.Float64("noiseamount", 0.5, "Amount of noise to mix in (0.0 to 1.0)")
	length := flag.Float64("length", 2000, "Length of the snare drum sample in milliseconds")
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
		cfg, err = synth.New606(synth.Snare, nil, *length/1000.0, sampleRate, *bitDepth, *channels)
		fmt.Println("Generating 606 snare with a sharp, crisp sound.")
	case *snare707:
		cfg, err = synth.New707(synth.Snare, nil, *length/1000.0, sampleRate, *bitDepth, *channels)
		fmt.Println("Generating 707 snare with a classic electronic sound.")
	case *snare808:
		cfg, err = synth.New808(synth.Snare, nil, *length/1000.0, sampleRate, *bitDepth, *channels)
		fmt.Println("Generating 808 snare with a deep, resonant sound.")
	case *snare909:
		cfg, err = synth.New909(synth.Snare, nil, *length/1000.0, sampleRate, *bitDepth, *channels)
		fmt.Println("Generating 909 snare with a punchy, bright sound.")
	case *snareLinn:
		cfg, err = synth.NewLinn(synth.Snare, nil, *length/1000.0, sampleRate, *bitDepth, *channels)
		fmt.Println("Generating LinnDrum snare with a classic digital sound.")
	case *snareExperimental:
		cfg, err = synth.NewExperimental(synth.Snare, nil, *length/1000.0, sampleRate, *bitDepth, *channels)
		fmt.Println("Generating experimental-style snare with unique characteristics.")
	default:
		cfg, err = synth.NewSnareSettings(nil, sampleRate, *bitDepth, *channels)
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
	cfg.PitchDecay = *pitchDecay

	// Generate the snare drum waveform
	samples, err := cfg.GenerateSnare()
	if err != nil {
		fmt.Println("Failed to generate snare:", err)
		return
	}

	// Generate noise samples if needed
	if *noiseAmount > 0.0 {
		numSamples := len(samples)

		noiseSamples := synth.GenerateWhiteNoise(numSamples, *noiseAmount)

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
		fmt.Println("Failed to save snare to file:", err)
		return
	}

	fmt.Println("Snare drum sound generated and written to", *outputFile)

	// Play the snare if -p flag is provided
	if *playSnare {
		fmt.Println("Playing the generated snare drum sound...")
		player := playsample.NewPlayer()
		defer player.Close()
		if err := player.PlayWaveform(samples, sampleRate, *bitDepth, *channels); err != nil {
			fmt.Println("Failed to play snare:", err)
			return
		}
		// TODO: Figure out why this is needed
		time.Sleep(time.Duration(int64(*length * 1e6)))
	}

}
