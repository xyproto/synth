package synth

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"image/color"
	"io"
	"math"
)

// Constants for waveform types
const (
	WaveSine = iota
	WaveTriangle
	WaveSawtooth
	WaveSquare
)

// Constants for noise types
const (
	NoiseNone = iota
	NoiseWhite
	NoisePink
	NoiseBrown
)

// Settings holds the configuration for generating a sound
type Settings struct {
	StartFreq                  float64
	EndFreq                    float64
	SampleRate                 int
	Duration                   float64
	WaveformType               int
	Attack                     float64
	Decay                      float64
	Sustain                    float64
	Release                    float64
	Drive                      float64
	FilterCutoff               float64
	FilterResonance            float64
	Sweep                      float64
	PitchDecay                 float64
	NoiseType                  int
	NoiseAmount                float64
	Output                     io.WriteSeeker
	NumOscillators             int
	OscillatorLevels           []float64
	SaturatorAmount            float64
	FilterBands                []float64
	BitDepth                   int
	FadeDuration               float64
	SmoothFrequencyTransitions bool
}

// SawtoothOscillator generates a sawtooth waveform at a specific frequency
func SawtoothOscillator(freq float64, length int, sampleRate int) []float64 {
	osc := make([]float64, length)
	for i := range osc {
		osc[i] = 2 * (float64(i)/float64(sampleRate)*freq - math.Floor(0.5+float64(i)/float64(sampleRate)*freq))
	}
	return osc
}

// DetunedOscillators generates multiple detuned sawtooth oscillators and combines them
func DetunedOscillators(freq float64, detune []float64, length int, sampleRate int) []float64 {
	numOsc := len(detune)
	combined := make([]float64, length)
	for _, d := range detune {
		osc := SawtoothOscillator(freq*(1+d), length, sampleRate)
		for i := range combined {
			combined[i] += osc[i] / float64(numOsc) // Average to avoid high amplitudes
		}
	}
	return combined
}

// ApplyEnvelope applies an ADSR envelope to the waveform
func ApplyEnvelope(samples []float64, attack, decay, sustain, release float64, sampleRate int) []float64 {
	adsr := make([]float64, len(samples))
	for i := range samples {
		t := float64(i) / float64(sampleRate)
		if t < attack {
			adsr[i] = samples[i] * (t / attack)
		} else if t < attack+decay {
			adsr[i] = samples[i] * (1 - (t-attack)/decay*(1-sustain))
		} else if t < float64(len(samples))/float64(sampleRate)-release {
			adsr[i] = samples[i] * sustain
		} else {
			adsr[i] = samples[i] * (1 - (t-(float64(len(samples))/float64(sampleRate)-release))/release*sustain)
		}
	}
	return adsr
}

// LowPassFilter applies a basic low-pass filter to the samples
func LowPassFilter(samples []float64, cutoff float64, sampleRate int) []float64 {
	filtered := make([]float64, len(samples))
	rc := 1.0 / (2.0 * math.Pi * cutoff)
	dt := 1.0 / float64(sampleRate)
	alpha := dt / (rc + dt)

	prev := samples[0]       // Initialize with the first sample
	filtered[0] = samples[0] // The first sample remains the same

	for i := 1; i < len(samples); i++ {
		filtered[i] = alpha*samples[i] + (1-alpha)*prev
		prev = filtered[i]
	}
	return filtered
}

// Drive applies a simple drive effect by scaling and clipping
func Drive(samples []float64, gain float64) []float64 {
	driven := make([]float64, len(samples))
	for i, sample := range samples {
		driven[i] = sample * gain
		if driven[i] > 1 {
			driven[i] = 1
		} else if driven[i] < -1 {
			driven[i] = -1
		}
	}
	return driven
}

// Limiter ensures the signal doesn't exceed [-1, 1] range
func Limiter(samples []float64) []float64 {
	limited := make([]float64, len(samples))
	for i, sample := range samples {
		if sample > 1 {
			limited[i] = 1
		} else if sample < -1 {
			limited[i] = -1
		} else {
			limited[i] = sample
		}
	}
	return limited
}

// GenerateKick generates the kick drum sound based on the settings
func (cfg *Settings) GenerateKick() error {
	numSamples := int(float64(cfg.SampleRate) * cfg.Duration)
	samples := make([]float64, numSamples)

	// Generate waveform based on the WaveformType
	for i := 0; i < numSamples; i++ {
		t := float64(i) / float64(cfg.SampleRate)
		frequency := cfg.StartFreq * math.Pow(cfg.EndFreq/cfg.StartFreq, t/cfg.Duration)
		var sample float64

		switch cfg.WaveformType {
		case WaveSine:
			sample = math.Sin(2 * math.Pi * frequency * t)
		case WaveTriangle:
			sample = 2*math.Abs(2*((t*frequency)-math.Floor((t*frequency)+0.5))) - 1
		case WaveSawtooth:
			sample = 2 * (t*frequency - math.Floor(0.5+t*frequency))
		case WaveSquare:
			sample = math.Copysign(1.0, math.Sin(2*math.Pi*frequency*t))
		}

		sample *= cfg.OscillatorLevels[0] // Apply the first oscillator level

		// Apply envelope (ADSR)
		sample *= cfg.ApplyEnvelope(t)

		// Apply drive (distortion)
		sample = cfg.ApplyDrive(sample)

		samples[i] = sample
	}

	// Apply limiter to the samples
	samples = Limiter(samples)

	// Save the result to a WAV file
	return SaveToWav("kick.wav", samples, cfg.SampleRate)
}

// ApplyEnvelope generates the ADSR envelope at a specific time
func (cfg *Settings) ApplyEnvelope(t float64) float64 {
	attack, decay, sustain, release := cfg.Attack, cfg.Decay, cfg.Sustain, cfg.Release
	duration := cfg.Duration

	if t < attack {
		return t / attack
	}
	if t < attack+decay {
		return 1.0 - (t-attack)/decay*(1.0-sustain)
	}
	if t < duration-release {
		return sustain
	}
	if t < duration {
		return sustain * (1.0 - (t-(duration-release))/release)
	}
	return 0.0
}

// ApplyDrive applies a drive (distortion) effect to the sample
func (cfg *Settings) ApplyDrive(sample float64) float64 {
	if cfg.Drive > 0 {
		return sample * (1 + cfg.Drive) / (1 + cfg.Drive*math.Abs(sample))
	}
	return sample
}

// LinearSummation mixes multiple audio samples by averaging them together.
// It automatically clamps the sum to avoid overflow and distortion.
func LinearSummation(samples ...[]float64) ([]float64, error) {
	if len(samples) == 0 {
		return nil, errors.New("no samples provided")
	}

	numSamples := len(samples[0])
	combined := make([]float64, numSamples)

	// Sum the samples from each input
	for i := 0; i < numSamples; i++ {
		sum := float64(0)
		for _, sample := range samples {
			if len(sample) != numSamples {
				return nil, errors.New("mismatched sample lengths")
			}
			sum += sample[i]
		}

		// Average the sum by dividing by the number of input samples
		average := sum / float64(len(samples))

		// Clamp the result to avoid overflow
		if average > 1 {
			combined[i] = 1
		} else if average < -1 {
			combined[i] = -1
		} else {
			combined[i] = average
		}
	}

	return combined, nil
}

// WeightedSummation mixes multiple audio samples by applying a weight to each sample.
// Each sample's amplitude is scaled by its corresponding weight before summing.
func WeightedSummation(weights []float64, samples ...[]float64) ([]float64, error) {
	if len(weights) != len(samples) {
		return nil, errors.New("number of weights must match number of samples")
	}

	if len(samples) == 0 {
		return nil, errors.New("no samples provided")
	}

	numSamples := len(samples[0])
	combined := make([]float64, numSamples)

	for i := 0; i < numSamples; i++ {
		sum := float64(0)
		for j, sample := range samples {
			if len(sample) != numSamples {
				return nil, errors.New("mismatched sample lengths")
			}
			sum += sample[i] * weights[j]
		}
		// Clamp the result to avoid overflow
		if sum > 1 {
			sum = 1
		} else if sum < -1 {
			sum = -1
		}
		combined[i] = sum
	}

	return combined, nil
}

// RMSMixing mixes audio samples using the Root Mean Square method.
func RMSMixing(samples ...[]float64) ([]float64, error) {
	if len(samples) == 0 {
		return nil, errors.New("no samples provided")
	}

	numSamples := len(samples[0])
	combined := make([]float64, numSamples)

	for i := 0; i < numSamples; i++ {
		sumSquares := float64(0)
		for _, sample := range samples {
			if len(sample) != numSamples {
				return nil, errors.New("mismatched sample lengths")
			}
			// Square the sample value and accumulate
			sumSquares += sample[i] * sample[i]
		}
		// Calculate RMS by taking the square root of the mean of squares
		rms := math.Sqrt(sumSquares / float64(len(samples)))

		// Clamp the result to [-1, 1] range
		if rms > 1 {
			rms = 1
		} else if rms < -1 {
			rms = -1
		}
		combined[i] = rms
	}

	return combined, nil
}

// AnalyzeHighestFrequency estimates the highest frequency in the audio signal
func AnalyzeHighestFrequency(samples []float64, sampleRate int) float64 {
	zeroCrossings := 0
	l := len(samples)
	for i := 1; i < l; i++ {
		if (samples[i-1] > 0 && samples[i] < 0) || (samples[i-1] < 0 && samples[i] > 0) {
			zeroCrossings++
		}
	}
	//fmt.Printf("Zero Crossings: %d\n", zeroCrossings)
	duration := float64(l) / float64(sampleRate)
	if duration == 0 {
		return 0
	}
	frequency := float64(zeroCrossings) / (2 * duration)
	return frequency
}

// NormalizeSamples scales the samples so the peak amplitude matches the given max amplitude
func NormalizeSamples(samples []float64, targetPeak float64) []float64 {
	currentPeak := FindPeakAmplitude(samples)
	if currentPeak == 0 {
		return samples // Avoid division by zero
	}
	scale := targetPeak / currentPeak
	normalizedSamples := make([]float64, len(samples))
	for i, sample := range samples {
		normalized := sample * scale
		// Clamp the values to the [-1, 1] range after scaling
		if normalized > 1 {
			normalizedSamples[i] = 1
		} else if normalized < -1 {
			normalizedSamples[i] = -1
		} else {
			normalizedSamples[i] = normalized
		}
	}
	return normalizedSamples
}

// FindPeakAmplitude returns the maximum absolute amplitude in the sample set
func FindPeakAmplitude(samples []float64) float64 {
	maxAmplitude := float64(0)
	for _, sample := range samples {
		if abs := math.Abs(sample); abs > maxAmplitude {
			maxAmplitude = abs
		}
	}
	return maxAmplitude
}

// PadSamples pads the shorter waveform with zeros to make both waveforms the same length.
func PadSamples(wave1, wave2 []float64) ([]float64, []float64) {
	length1 := len(wave1)
	length2 := len(wave2)

	if length1 == length2 {
		return wave1, wave2
	}

	// Pad the shorter waveform with zeros
	if length1 < length2 {
		paddedWave1 := make([]float64, length2)
		copy(paddedWave1, wave1)
		return paddedWave1, wave2
	}

	paddedWave2 := make([]float64, length1)
	copy(paddedWave2, wave2)
	return wave1, paddedWave2
}

// CopySettings creates a deep copy of a Settings struct
func CopySettings(cfg *Settings) *Settings {
	newCfg := *cfg
	newCfg.OscillatorLevels = append([]float64(nil), cfg.OscillatorLevels...) // Deep copy the slice
	return &newCfg
}

// GenerateKickWaveform generates the kick waveform and returns it as a slice of float64 samples (without writing to disk).
func (cfg *Settings) GenerateKickWaveform() ([]float64, error) {
	numSamples := int(float64(cfg.SampleRate) * cfg.Duration)
	samples := make([]float64, numSamples)

	for i := 0; i < numSamples; i++ {
		t := float64(i) / float64(cfg.SampleRate)
		frequency := cfg.StartFreq * math.Pow(cfg.EndFreq/cfg.StartFreq, t/cfg.Duration)
		var sample float64

		switch cfg.WaveformType {
		case WaveSine:
			sample = math.Sin(2 * math.Pi * frequency * t)
		case WaveTriangle:
			sample = 2*math.Abs(2*((t*frequency)-math.Floor((t*frequency)+0.5))) - 1
		case WaveSawtooth:
			sample = 2 * (t*frequency - math.Floor(0.5+t*frequency))
		case WaveSquare:
			sample = math.Copysign(1.0, math.Sin(2*math.Pi*frequency*t))
		}

		sample *= cfg.OscillatorLevels[0]
		sample *= cfg.ApplyEnvelope(t)
		sample = cfg.ApplyDrive(sample)
		samples[i] = sample
	}

	samples = Limiter(samples)
	return samples, nil
}

// Color returns a color that very approximately represents the current kick config
func (cfg *Settings) Color() color.RGBA {
	hasher := sha1.New()
	hasher.Write([]byte(fmt.Sprintf("%d%f%f%f%f%f%f%f%f", cfg.WaveformType, cfg.Attack, cfg.Decay, cfg.Sustain, cfg.Release, cfg.Drive, cfg.FilterCutoff, cfg.Sweep, cfg.PitchDecay)))
	hashBytes := hasher.Sum(nil)
	// Convert the first few bytes of the hash into an RGB color
	r := hashBytes[0]
	g := hashBytes[1]
	b := hashBytes[2]
	return color.RGBA{R: r, G: g, B: b, A: 255}
}
