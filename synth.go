package synth

import (
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"os/exec"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
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
	alpha := 2 * math.Pi * cutoff / float64(sampleRate)
	prev := 0.0
	for i, sample := range samples {
		filtered[i] = prev + alpha*(sample-prev)
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

// SaveToWav saves the waveform to a wav file
func SaveToWav(filename string, samples []float64, sampleRate int) error {
	outFile, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating wav file: %v", err)
	}
	defer outFile.Close()

	enc := wav.NewEncoder(outFile, sampleRate, 16, 1, 1)

	buf := &audio.IntBuffer{
		Format: &audio.Format{SampleRate: sampleRate, NumChannels: 1},
		Data:   make([]int, len(samples)),
	}

	for i, sample := range samples {
		buf.Data[i] = int(sample * math.MaxInt16) // Convert to 16-bit PCM
	}

	if err := enc.Write(buf); err != nil {
		return fmt.Errorf("error writing wav file: %v", err)
	}

	return enc.Close()
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

// LinearSummation mixes multiple audio samples by adding them together.
// It automatically clamps the sum to avoid overflow and distortion.
func LinearSummation(samples ...[]float64) ([]float64, error) {
	if len(samples) == 0 {
		return nil, errors.New("no samples provided")
	}

	numSamples := len(samples[0])
	combined := make([]float64, numSamples)

	for i := 0; i < numSamples; i++ {
		sum := float64(0)
		for _, sample := range samples {
			if len(sample) != numSamples {
				return nil, errors.New("mismatched sample lengths")
			}
			sum += sample[i]
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
	duration := float64(l) / float64(sampleRate)
	frequency := float64(zeroCrossings) / (2 * duration)
	return frequency
}

// NormalizeSamples scales the samples so the peak amplitude matches the given max amplitude
func NormalizeSamples(samples []float64, targetPeak float64) []float64 {
	currentPeak := FindPeakAmplitude(samples)
	if currentPeak == 0 {
		return samples
	}

	scale := targetPeak / currentPeak
	l := len(samples)

	normalizedSamples := make([]float64, l)
	for i := 0; i < l; i++ {
		normalized := samples[i] * scale
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

// PlayWav plays a WAV file using mpv or ffmpeg
func PlayWav(filePath string) error {
	cmd := exec.Command("mpv", filePath)
	err := cmd.Start()
	if err != nil {
		// Fallback to ffmpeg if mpv is not available
		cmd = exec.Command("ffmpeg", "-i", filePath, "-f", "null", "-")
		err = cmd.Start()
		if err != nil {
			return fmt.Errorf("error playing sound with both mpv and ffmpeg: %v", err)
		}
	}
	cmd.Wait()
	return nil
}
