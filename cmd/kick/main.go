package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/xyproto/synth"
)

const version = "0.0.3"

// parseCommaSeparatedFloats parses a comma-separated string of floats
func parseCommaSeparatedFloats(input string) []float64 {
	var result []float64
	for _, s := range strings.Split(input, ",") {
		value, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
		if err != nil {
			fmt.Println("Error parsing value:", s)
			continue
		}
		result = append(result, value)
	}
	return result
}

func main() {
	// Drum style flags
	kick606 := flag.Bool("606", false, "Generate a kick.wav like a 606 kick drum")
	kick707 := flag.Bool("707", false, "Generate a kick.wav like a 707 kick drum")
	kick808 := flag.Bool("808", false, "Generate a kick.wav like an 808 kick drum")
	kick909 := flag.Bool("909", false, "Generate a kick.wav like a 909 kick drum")
	kickLinnDrum := flag.Bool("linn", false, "Generate a kick.wav like a LinnDrum kick drum")
	kickDeepHouse := flag.Bool("deephouse", false, "Generate a deep house kick drum")
	kickExperimental := flag.Bool("experimental", false, "Generate a kick.wav with experimental-style characteristics")

	// Sound customization flags
	noiseType := flag.String("noise", "none", "Type of noise to mix in (none, white, pink, brown)")
	noiseAmount := flag.Float64("noiseamount", 0.0, "Amount of noise to mix in (0.0 to 1.0)")
	length := flag.Float64("length", 1000, "Length of the kick drum sample in milliseconds")
	quality := flag.Int("quality", 96, "Sample rate in kHz (44, 48, 96, or 192)")
	bitDepth := flag.Int("bitdepth", 16, "Bit depth of the audio (8, 16, 24 or 32)")
	channels := flag.Int("channels", 1, "Channels (1 or 2)")
	waveform := flag.Int("waveform", synth.WaveSine, "Waveform type (0: Sine, 1: Triangle, 2: Sawtooth, 3: Square)")
	attack := flag.Float64("attack", 0.003, "Attack time in seconds")
	decay := flag.Float64("decay", 0.3, "Decay time in seconds")
	sustain := flag.Float64("sustain", 0.1, "Sustain level (0.0 to 1.0)")
	release := flag.Float64("release", 0.25, "Release time in seconds") // Increased release time
	sweep := flag.Float64("sweep", 0.8, "Pitch sweep rate")
	filterCutoff := flag.Float64("filter", 5000.0, "Low-pass filter cutoff frequency (Hz)")
	pitchDecay := flag.Float64("pitchdecay", 0.2, "Pitch envelope decay time")
	drive := flag.Float64("drive", 0.1, "Amount of distortion/drive")
	numOscillators := flag.Int("numoscillators", 1, "Number of oscillators for layering")
	oscillatorLevels := flag.String("oscillatorlevels", "1.0", "Comma-separated levels for each oscillator")
	saturatorAmount := flag.Float64("saturator", 0.3, "Amount of saturation to apply")
	filterBands := flag.String("filterbands", "200,1000,3000", "Comma-separated multi-band filter cutoff frequencies")
	outputFile := flag.String("o", "kick.wav", "Output file path")
	playKick := flag.Bool("p", false, "Play the generated kick") // Added -p flag
	showVersion := flag.Bool("version", false, "Show the current version")
	showHelp := flag.Bool("help", false, "Display this help")

	flag.Parse()

	if *showHelp {
		flag.Usage()
		return
	}

	if *showVersion {
		fmt.Println("kick version", version)
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
	case *kick606:
		cfg, err = synth.New606Kick(nil, *length/1000.0, sampleRate, *bitDepth, *channels)
		fmt.Println("Generating 606 kick with a punchy, shorter sound.")
	case *kick707:
		cfg, err = synth.New707Kick(nil, *length/1000.0, sampleRate, *bitDepth, *channels)
		fmt.Println("Generating 707 kick with a classic, shorter punchy sound.")
	case *kick808:
		cfg, err = synth.New808Kick(nil, *length/1000.0, sampleRate, *bitDepth, *channels)
		fmt.Println("Generating 808 kick with deep sub-bass and smooth characteristics.")
	case *kick909:
		cfg, err = synth.New909Kick(nil, *length/1000.0, sampleRate, *bitDepth, *channels)
		fmt.Println("Generating 909 kick with punchy, mid-range presence and quick decay.")
	case *kickLinnDrum:
		cfg, err = synth.NewLinnDrumKick(nil, *length/1000.0, sampleRate, *bitDepth, *channels)
		fmt.Println("Generating LinnDrum kick with an iconic, punchy sound.")
	case *kickDeepHouse:
		cfg, err = synth.NewDeepHouseKick(nil, *length/1000.0, sampleRate, *bitDepth, *channels)
		fmt.Println("Generating Deep House kick with smooth, warm bass.")
	case *kickExperimental:
		cfg, err = synth.NewExperimentalKick(nil, *length/1000.0, sampleRate, *bitDepth, *channels)
		fmt.Println("Generating experimental-style kick with unique texture.")
	default:
		cfg, err = synth.NewSettings(nil, 150.0, 40.0, *length/1000.0, sampleRate, *bitDepth, *channels)
		fmt.Println("Generating default kick with user-defined characteristics.")
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
	cfg.Sweep = *sweep
	cfg.FilterCutoff = *filterCutoff
	cfg.PitchDecay = *pitchDecay
	cfg.Drive = *drive
	cfg.NumOscillators = *numOscillators
	cfg.OscillatorLevels = parseCommaSeparatedFloats(*oscillatorLevels)
	cfg.SaturatorAmount = *saturatorAmount
	cfg.FilterBands = parseCommaSeparatedFloats(*filterBands)
	cfg.FadeDuration = *release // Match fade duration to release time
	cfg.SmoothFrequencyTransitions = true

	// Generate the kick drum waveform
	samples, err := cfg.GenerateKick()
	if err != nil {
		fmt.Println("Failed to generate kick:", err)
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
			fmt.Println("Invalid noise type. Choose from: none, white, pink, brown.")
			os.Exit(1)
		}

		// Mix noise samples into the kick samples
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
		fmt.Println("Failed to save kick to file:", err)
		return
	}

	fmt.Println("Kick drum sound generated and written to", *outputFile)

	// Play the kick if -p flag is provided
	if *playKick {
		fmt.Println("Playing the generated kick drum sound...")
		// Use PlayWaveform to play the samples directly
		player := synth.NewPlayer()
		defer player.Close()
		if audioDeviceKey, playbackDuration, err := player.PlayWaveform(samples, sampleRate, *bitDepth, *channels); err != nil {
			fmt.Println("Failed to play kick:", err)
		} else {
			player.WaitClosePlus100(audioDeviceKey, playbackDuration)
		}
	}

}
