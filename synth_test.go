package synth

import (
	"math"
	"os"
	"testing"
)

func createTestWaveform(value float64, numSamples int) []float64 {
	waveform := make([]float64, numSamples)
	for i := 0; i < numSamples; i++ {
		waveform[i] = value
	}
	return waveform
}

func createSineWave(freq float64, numSamples int, sampleRate int) []float64 {
	sineWave := make([]float64, numSamples)
	for i := 0; i < numSamples; i++ {
		t := float64(i) / float64(sampleRate)
		sineWave[i] = math.Sin(2 * math.Pi * freq * t)
	}
	return sineWave
}

func TestAnalyzeHighestFrequency(t *testing.T) {
	// Create a sine wave of known frequency (e.g., 1000 Hz)
	frequency := 1000.0
	sampleRate := 44100
	numSamples := 44100 // 1 second worth of samples
	samples := createSineWave(frequency, numSamples, sampleRate)

	estimatedFrequency := AnalyzeHighestFrequency(samples, sampleRate)
	if math.Abs(estimatedFrequency-frequency) > 1.0 {
		t.Errorf("Expected frequency close to %.2f Hz, got %.2f Hz", frequency, estimatedFrequency)
	}
}
func TestLinearSummation(t *testing.T) {
	wave1 := createTestWaveform(0.5, 10)
	wave2 := createTestWaveform(1.0, 10)
	expected := createTestWaveform(1.0, 10) // Clamped to 1.0

	result, err := LinearSummation(wave1, wave2)
	if err != nil {
		t.Fatalf("Error in LinearSummation: %v", err)
	}

	for i, v := range result {
		if v != expected[i] {
			t.Errorf("LinearSummation failed at index %d: expected %f, got %f", i, expected[i], v)
		}
	}
}

// TestWeightedSummation checks if the weighted summation mixing works as expected
func TestWeightedSummation(t *testing.T) {
	wave1 := createTestWaveform(0.5, 10)
	wave2 := createTestWaveform(1.0, 10)
	weights := []float64{0.5, 0.5}
	expected := createTestWaveform(0.75, 10)

	result, err := WeightedSummation(weights, wave1, wave2)
	if err != nil {
		t.Fatalf("Error in WeightedSummation: %v", err)
	}

	for i, v := range result {
		if v != expected[i] {
			t.Errorf("WeightedSummation failed at index %d: expected %f, got %f", i, expected[i], v)
		}
	}
}

// TestRMSMixing checks if the RMS mixing works as expected
func TestRMSMixing(t *testing.T) {
	wave1 := createTestWaveform(0.5, 10)
	wave2 := createTestWaveform(1.0, 10)

	// Calculate the expected RMS value
	expectedRMS := math.Sqrt((0.5*0.5 + 1.0*1.0) / 2)

	result, err := RMSMixing(wave1, wave2)
	if err != nil {
		t.Fatalf("Error in RMSMixing: %v", err)
	}

	for i, v := range result {
		if v != expectedRMS {
			t.Errorf("RMSMixing failed at index %d: expected %f, got %f", i, expectedRMS, v)
		}
	}
}

// TestErrorCases tests that the functions handle error cases correctly
func TestErrorCases(t *testing.T) {
	// Mismatched sample lengths
	wave1 := createTestWaveform(0.5, 10)
	wave2 := createTestWaveform(0.5, 5)

	_, err := LinearSummation(wave1, wave2)
	if err == nil {
		t.Error("Expected error for mismatched sample lengths in LinearSummation")
	}

	_, err = WeightedSummation([]float64{0.5, 0.5}, wave1, wave2)
	if err == nil {
		t.Error("Expected error for mismatched sample lengths in WeightedSummation")
	}

	_, err = RMSMixing(wave1, wave2)
	if err == nil {
		t.Error("Expected error for mismatched sample lengths in RMSMixing")
	}

	// Mismatched weights
	wave3 := createTestWaveform(0.5, 10)
	_, err = WeightedSummation([]float64{0.5}, wave1, wave3)
	if err == nil {
		t.Error("Expected error for mismatched number of weights and samples in WeightedSummation")
	}
}

// TestWaveformGeneration tests SawtoothOscillator and DetunedOscillators
func TestWaveformGeneration(t *testing.T) {
	freq := 440.0
	length := 100
	sampleRate := 44100

	// Test SawtoothOscillator
	sawtooth := SawtoothOscillator(freq, length, sampleRate)
	if len(sawtooth) != length {
		t.Errorf("Expected sawtooth waveform length of %d, got %d", length, len(sawtooth))
	}

	// Test DetunedOscillators
	detune := []float64{-0.01, 0.0, 0.01}
	oscillators := DetunedOscillators(freq, detune, length, sampleRate)
	if len(oscillators) != length {
		t.Errorf("Expected detuned oscillators waveform length of %d, got %d", length, len(oscillators))
	}
}

// TestApplyEnvelope tests the ADSR envelope function
func TestApplyEnvelope(t *testing.T) {
	samples := createTestWaveform(1.0, 100)
	attack, decay, sustain, release := 0.1, 0.2, 0.5, 0.3
	sampleRate := 44100

	enveloped := ApplyEnvelope(samples, attack, decay, sustain, release, sampleRate)
	if len(enveloped) != len(samples) {
		t.Errorf("Expected enveloped waveform length of %d, got %d", len(samples), len(enveloped))
	}
}

// TestLowPassFilter tests the low-pass filter
func TestLowPassFilter(t *testing.T) {
	samples := createTestWaveform(1.0, 100)
	filtered := LowPassFilter(samples, 1000.0, 44100)

	if len(filtered) != len(samples) {
		t.Errorf("Expected filtered waveform length of %d, got %d", len(samples), len(filtered))
	}
}

// TestDrive tests the drive (distortion) function
func TestDrive(t *testing.T) {
	samples := createTestWaveform(1.0, 100)
	gain := 2.0
	driven := Drive(samples, gain)

	for i, v := range driven {
		if v > 1 || v < -1 {
			t.Errorf("Drive failed at index %d: expected values between -1 and 1, got %f", i, v)
		}
	}
}

// TestSaveAndLoadWav tests saving and loading WAV files
func TestSaveAndLoadWav(t *testing.T) {
	samples := createTestWaveform(1.0, 100)
	filename := "test_output.wav"
	defer os.Remove(filename)

	err := SaveToWav(filename, samples, 44100)
	if err != nil {
		t.Fatalf("Failed to save WAV file: %v", err)
	}

	_, sampleRate, err := LoadWav(filename)
	if err != nil {
		t.Fatalf("Failed to load WAV file: %v", err)
	}
	if sampleRate != 44100 {
		t.Errorf("Expected sample rate of 44100, got %d", sampleRate)
	}
}

// TestPadSamples tests padding of samples
func TestPadSamples(t *testing.T) {
	wave1 := createTestWaveform(0.5, 3)
	wave2 := createTestWaveform(0.5, 2)

	padded1, padded2 := PadSamples(wave1, wave2)

	if len(padded1) != len(padded2) {
		t.Errorf("Expected padded waves to have the same length, got %d and %d", len(padded1), len(padded2))
	}
}

// TestNormalizeSamples tests normalizing of samples
func TestNormalizeSamples(t *testing.T) {
	samples := createTestWaveform(0.1, 3)
	targetPeak := 1.0
	normalized := NormalizeSamples(samples, targetPeak)

	peak := FindPeakAmplitude(normalized)
	if peak != targetPeak {
		t.Errorf("Expected peak amplitude of %f, got %f", targetPeak, peak)
	}
}

// TestFindPeakAmplitude tests finding the peak amplitude
func TestFindPeakAmplitude(t *testing.T) {
	samples := createTestWaveform(0.5, 3)
	expectedPeak := 0.5

	peak := FindPeakAmplitude(samples)
	if peak != expectedPeak {
		t.Errorf("Expected peak amplitude of %f, got %f", expectedPeak, peak)
	}
}
