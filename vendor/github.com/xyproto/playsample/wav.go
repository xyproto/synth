package playsample

import (
	"fmt"
	"io"
	"math"
	"os"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
)

// SaveToWav saves the samples to a WAV file using the provided parameters.
func SaveToWav(w io.WriteSeeker, samples []float64, sampleRate, bitDepth, channels int) error {
	if len(samples) == 0 {
		return fmt.Errorf("cannot save empty waveform: no samples provided")
	}
	if bitDepth != 8 && bitDepth != 16 && bitDepth != 24 && bitDepth != 32 {
		return fmt.Errorf("bitdepth should be 8, 16, 24, or 32, not %d", bitDepth)
	}
	if channels <= 0 {
		return fmt.Errorf("channels should be greater than 0, got %d", channels)
	}
	if len(samples)%channels != 0 {
		return fmt.Errorf("number of samples (%d) is not a multiple of channels (%d)", len(samples), channels)
	}

	const audioFormat = 1
	enc := wav.NewEncoder(w, sampleRate, bitDepth, channels, audioFormat)

	intData := make([]int, len(samples))
	maxIntValue := float64(int(1)<<(bitDepth-1)) - 1
	minIntValue := -maxIntValue - 1

	for i := 0; i < len(samples); i++ {
		sample := samples[i]
		if sample > 1.0 {
			sample = 1.0
		} else if sample < -1.0 {
			sample = -1.0
		}
		scaled := sample * maxIntValue
		scaled = math.Max(math.Min(scaled, maxIntValue), minIntValue)
		intData[i] = int(math.Round(scaled))
	}

	buf := &audio.IntBuffer{
		Format:         &audio.Format{SampleRate: sampleRate, NumChannels: channels},
		Data:           intData,
		SourceBitDepth: bitDepth,
	}

	if err := enc.Write(buf); err != nil {
		return fmt.Errorf("error writing WAV file: %v", err)
	}
	return enc.Close()
}

// LoadWav loads a WAV file and returns the samples as []float64 and the sample rate.
// If monoToStereo is true and the audio is mono, the samples will be duplicated to create stereo output.
func LoadWav(filename string, monoToStereo bool) ([]float64, int, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, 0, fmt.Errorf("error opening WAV file: %v", err)
	}
	defer f.Close()

	decoder := wav.NewDecoder(f)
	if !decoder.IsValidFile() {
		return nil, 0, fmt.Errorf("invalid WAV file")
	}

	intBuf, err := decoder.FullPCMBuffer()
	if err != nil {
		return nil, 0, fmt.Errorf("error decoding WAV file: %v", err)
	}

	sampleRate := intBuf.Format.SampleRate
	channels := intBuf.Format.NumChannels
	bitDepth := intBuf.SourceBitDepth
	intData := intBuf.Data

	maxIntValue := float64(int(1) << (bitDepth - 1))
	samples := make([]float64, len(intData))

	for i, v := range intData {
		samples[i] = float64(v) / maxIntValue
	}

	if channels == 1 && monoToStereo {
		numSamples := len(samples)
		stereoSamples := make([]float64, numSamples*2)
		for i := 0; i < numSamples; i++ {
			sample := samples[i]
			stereoSamples[2*i] = sample   // Left channel
			stereoSamples[2*i+1] = sample // Right channel
		}
		samples = stereoSamples
		channels = 2
	}

	return samples, sampleRate, nil
}
