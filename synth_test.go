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
	expected := createTestWaveform(0.75, 10) // Averaging 0.5 and 1.0

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

	// Open the file for writing
	file, err := os.Create(filename)
	if err != nil {
		t.Fatalf("Failed to create WAV file: %v", err)
	}
	defer file.Close()

	err = SaveToWav(file, samples, 44100, 16, 1)
	if err != nil {
		t.Fatalf("Failed to save WAV file: %v", err)
	}

	_, sampleRate, err := LoadWav(filename, false)
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

func TestWavSaveAndLoad(t *testing.T) {
	// Create a test waveform with values in the [-1.0, 1.0] range
	samples := createTestWaveform(0.5, 100)
	filename := "test_wav_output.wav"
	defer os.Remove(filename)

	// Open the file for writing
	file, err := os.Create(filename)
	if err != nil {
		t.Fatalf("Failed to create WAV file: %v", err)
	}
	defer file.Close()

	// Save the waveform to a WAV file
	err = SaveToWav(file, samples, 44100, 16, 1)
	if err != nil {
		t.Fatalf("Failed to save WAV file: %v", err)
	}

	// Load the waveform back from the WAV file
	loadedSamples, sampleRate, err := LoadWav(filename, false)
	if err != nil {
		t.Fatalf("Failed to load WAV file: %v", err)
	}

	// Ensure the sample rate matches
	if sampleRate != 44100 {
		t.Errorf("Expected sample rate of 44100, got %d", sampleRate)
	}

	// Ensure the samples match within a small tolerance (due to potential floating-point inaccuracies)
	tolerance := 0.0001
	for i, sample := range samples {
		if math.Abs(sample-loadedSamples[i]) > tolerance {
			t.Errorf("Loaded sample at index %d differs from original: expected %f, got %f", i, sample, loadedSamples[i])
		}
	}
}

// TestLimiter tests if the Limiter correctly clamps values to the [-1, 1] range
func TestLimiter(t *testing.T) {
	// Create a test waveform with values that exceed [-1, 1]
	samples := []float64{1.5, -1.2, 0.5, -0.8, 2.0}

	limited := Limiter(samples)

	for i, v := range limited {
		if v > 1.0 {
			t.Errorf("Limiter failed at index %d: expected value <= 1.0, got %f", i, v)
		}
		if v < -1.0 {
			t.Errorf("Limiter failed at index %d: expected value >= -1.0, got %f", i, v)
		}
	}
}

// TestSaveTo tests the SaveTo function
func TestSaveTo(t *testing.T) {
	cfg := &Settings{
		StartFreq:        100.0,
		EndFreq:          50.0,
		SampleRate:       44100,
		BitDepth:         16,
		Duration:         1.0,
		OscillatorLevels: []float64{1.0},
	}

	// Test saving the generated waveform to a .wav file
	filename, err := cfg.GenerateAndSaveTo("kick", ".")
	defer os.Remove(filename) // Cleanup after test

	if err != nil {
		t.Fatalf("SaveTo failed: %v", err)
	}

	// Check if the file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Fatalf("Expected file to exist: %s", filename)
	}
}

// TestGenerateKick checks if GenerateKickWaveform correctly generates non-zero length samples
func TestGenerateKick(t *testing.T) {
	cfg := &Settings{
		StartFreq:        100.0,
		EndFreq:          50.0,
		Duration:         1.0,
		SampleRate:       44100,
		BitDepth:         16,
		OscillatorLevels: []float64{1.0},
	}

	samples, err := cfg.GenerateKick()
	if err != nil {
		t.Fatalf("GenerateKick failed: %v", err)
	}

	if len(samples) == 0 {
		t.Fatalf("Expected non-zero length waveform, got %d samples", len(samples))
	}

	// Check that the samples fall within the expected [-1.0, 1.0] range
	for i, sample := range samples {
		if sample > 1.0 || sample < -1.0 {
			t.Errorf("Sample at index %d is out of range [-1, 1]: %f", i, sample)
		}
	}
}

// TestSaveToWavEmptySamples ensures SaveToWav does not save zero-length files
func TestSaveToWavEmptySamples(t *testing.T) {
	samples := []float64{} // Zero-length samples
	filename := "test_empty_output.wav"
	defer os.Remove(filename)

	// Open the file for writing
	file, err := os.Create(filename)
	if err != nil {
		t.Fatalf("Failed to create WAV file: %v", err)
	}
	defer file.Close()

	err = SaveToWav(file, samples, 44100, 16, 1)
	if err == nil {
		t.Fatalf("Expected error when saving zero-length waveform, but got nil")
	}
}

// TestSaveToWavNonEmptySamples ensures SaveToWav correctly saves non-zero length waveforms
func TestSaveToWavNonEmptySamples(t *testing.T) {
	samples := createTestWaveform(0.5, 100) // Non-zero length waveform
	filename := "test_output_non_empty.wav"
	defer os.Remove(filename)

	// Open the file for writing
	file, err := os.Create(filename)
	if err != nil {
		t.Fatalf("Failed to create WAV file: %v", err)
	}
	defer file.Close()

	err = SaveToWav(file, samples, 44100, 16, 1)
	if err != nil {
		t.Fatalf("Failed to save non-zero length waveform: %v", err)
	}

	// Check if the file exists and is non-zero in size
	fileInfo, err := os.Stat(filename)
	if err != nil {
		t.Fatalf("Failed to stat file: %v", err)
	}

	if fileInfo.Size() == 0 {
		t.Fatalf("Expected non-zero file size, got %d bytes", fileInfo.Size())
	}
}

func TestBandPassFilter(t *testing.T) {
	samples := createTestWaveform(1.0, 100)
	filtered := BandPassFilter(samples, 500.0, 5000.0, 44100)

	if len(filtered) != len(samples) {
		t.Errorf("Expected filtered waveform length of %d, got %d", len(samples), len(filtered))
	}
}

func TestApplyPitchModulation(t *testing.T) {
	samples := createSineWave(440.0, 100, 44100)
	modulated := ApplyPitchModulation(samples, 5.0, 0.2, 44100)

	if len(modulated) != len(samples) {
		t.Errorf("Expected modulated waveform length of %d, got %d", len(samples), len(modulated))
	}
}

func TestApplyPanning(t *testing.T) {
	samples := createTestWaveform(0.5, 100)
	left, right := ApplyPanning(samples, 0.5)

	if len(left) != len(right) || len(left) != len(samples) {
		t.Errorf("Expected panned waveform lengths to match the input, got %d and %d", len(left), len(right))
	}
}

func TestGenerateNoise(t *testing.T) {
	length := 100
	whiteNoise := GenerateNoise(NoiseWhite, length, 0.5)
	pinkNoise := GenerateNoise(NoisePink, length, 0.5)
	brownNoise := GenerateNoise(NoiseBrown, length, 0.5)

	if len(whiteNoise) != length || len(pinkNoise) != length || len(brownNoise) != length {
		t.Errorf("Expected noise of length %d, got different lengths", length)
	}
}

func TestApplyFrequencyModulation(t *testing.T) {
	samples := createSineWave(440.0, 100, 44100)
	modulated := ApplyFrequencyModulation(samples, 5.0, 0.2, 44100)

	if len(modulated) != len(samples) {
		t.Errorf("Expected modulated waveform length of %d, got %d", len(samples), len(modulated))
	}
}

func TestHighPassFilter(t *testing.T) {
	samples := createTestWaveform(1.0, 100)
	filtered := HighPassFilter(samples, 1000.0, 44100)

	if len(filtered) != len(samples) {
		t.Errorf("Expected filtered waveform length of %d, got %d", len(samples), len(filtered))
	}

	// Check that the high frequencies are retained
	if filtered[0] != samples[0] {
		t.Errorf("Expected first sample to be unchanged, got %f", filtered[0])
	}
}

func TestSchroederReverb(t *testing.T) {
	samples := createTestWaveform(0.5, 1000)
	combDelays := []int{1557, 1617, 1491, 1422}
	allPassDelays := []int{225, 556}
	reverb, err := SchroederReverb(samples, 0.85, combDelays, allPassDelays)
	if err != nil {
		t.Fatalf("SchroederReverb failed: %v", err)
	}

	if len(reverb) != len(samples) {
		t.Errorf("Expected reverb waveform length of %d, got %d", len(samples), len(reverb))
	}
}

func TestApplyChorus(t *testing.T) {
	samples := createTestWaveform(0.5, 100)
	sampleRate := 44100
	delaySec := 0.005
	depth := 0.5
	rate := 1.5
	modulated := ApplyChorus(samples, sampleRate, delaySec, depth, rate)

	if len(modulated) != len(samples) {
		t.Errorf("Expected chorus-modulated waveform length of %d, got %d", len(samples), len(modulated))
	}
}

func TestGenerateSweepWaveform(t *testing.T) {
	cfg := &Settings{
		StartFreq:    100.0,
		EndFreq:      1000.0,
		SampleRate:   44100,
		BitDepth:     16,
		Duration:     1.0,
		WaveformType: WaveSine,
	}

	samples, err := cfg.GenerateSweepWaveform()
	if err != nil {
		t.Fatalf("GenerateSweepWaveform failed: %v", err)
	}

	if len(samples) != int(cfg.Duration*float64(cfg.SampleRate)) {
		t.Errorf("Expected sweep waveform length of %d, got %d", int(cfg.Duration*float64(cfg.SampleRate)), len(samples))
	}
}

func TestColor(t *testing.T) {
	cfg := &Settings{
		WaveformType: WaveSine,
		Attack:       0.1,
		Decay:        0.2,
		Sustain:      0.7,
		Release:      0.3,
		Drive:        0.5,
		FilterCutoff: 5000,
		Sweep:        0.2,
		PitchDecay:   0.1,
	}

	color := cfg.Color()
	// Since the color is based on a hash, we can't predict it, but we can check if it's properly formed
	if color.A != 255 {
		t.Errorf("Expected alpha value of 255, got %d", color.A)
	}
}

func TestCopySettings(t *testing.T) {
	original := &Settings{
		StartFreq:        440.0,
		EndFreq:          880.0,
		Duration:         1.0,
		SampleRate:       44100,
		BitDepth:         16,
		WaveformType:     WaveSine,
		OscillatorLevels: []float64{1.0, 0.8},
	}

	copy2 := CopySettings(original)
	if original == copy2 {
		t.Error("CopySettings did not create a new instance")
	}

	// Modify the original and check if the copy remains unchanged
	original.OscillatorLevels[0] = 0.5
	if copy2.OscillatorLevels[0] != 1.0 {
		t.Error("CopySettings did not perform a deep copy of OscillatorLevels")
	}
}

func TestGenerateNoiseTypes(t *testing.T) {
	length := 1000
	amount := 0.5
	tolerance := 0.05 // Allow a 10% tolerance
	noiseTypes := []int{NoiseWhite, NoisePink, NoiseBrown}

	for _, noiseType := range noiseTypes {
		noise := GenerateNoise(noiseType, length, amount)
		if len(noise) != length {
			t.Errorf("Expected noise length of %d, got %d for noise type %d", length, len(noise), noiseType)
		}

		for i, v := range noise {
			if v < -(amount+tolerance) || v > (amount+tolerance) {
				t.Errorf("Noise value at index %d out of range [%.6f, %.6f]: %f", i, -(amount + tolerance), amount+tolerance, v)
			}
		}
	}
}

func TestApplyPanningExtremes(t *testing.T) {
	samples := createTestWaveform(0.5, 100)
	tolerance := 1e-6

	// Full left
	_, right := ApplyPanning(samples, -1.0)
	for i := range samples {
		if math.Abs(right[i]) > tolerance {
			t.Errorf("Expected right channel to be silent at full left pan, got %f at index %d", right[i], i)
		}
	}

	// Full right
	left, _ := ApplyPanning(samples, 1.0)
	for i := range samples {
		if math.Abs(left[i]) > tolerance {
			t.Errorf("Expected left channel to be silent at full right pan, got %f at index %d", left[i], i)
		}
	}
}

func TestApplyFrequencyModulationBounds(t *testing.T) {
	samples := createSineWave(440.0, 1000, 44100)
	modFreq := 5.0
	modDepth := 0.1
	modulated := ApplyFrequencyModulation(samples, modFreq, modDepth, 44100)

	if len(modulated) != len(samples) {
		t.Errorf("Expected modulated waveform length of %d, got %d", len(samples), len(modulated))
	}

	// Check that the modulated samples are within the [-1, 1] range
	for i, v := range modulated {
		if v < -1.0 || v > 1.0 {
			t.Errorf("Modulated sample at index %d out of range [-1, 1]: %f", i, v)
		}
	}
}

func TestApplyPitchModulationBounds(t *testing.T) {
	samples := createSineWave(440.0, 1000, 44100)
	modFreq := 5.0
	modDepth := 0.1
	modulated := ApplyPitchModulation(samples, modFreq, modDepth, 44100)

	if len(modulated) != len(samples) {
		t.Errorf("Expected modulated waveform length of %d, got %d", len(samples), len(modulated))
	}

	tolerance := 0.01 // Allow a small tolerance
	for i, v := range modulated {
		if v < -1.0-tolerance || v > 1.0+tolerance {
			t.Errorf("Pitch-modulated sample at index %d out of range [-1, 1] with tolerance %f: %f", i, tolerance, v)
		}
	}
}

func TestApplyDriveNoGain(t *testing.T) {
	samples := createTestWaveform(0.5, 100)
	driven := Drive(samples, 1.0) // Gain of 1.0 should not alter the samples

	for i, v := range driven {
		if v != samples[i] {
			t.Errorf("Expected driven sample to equal original at index %d, got %f", i, v)
		}
	}
}

func TestApplyDriveHighGain(t *testing.T) {
	samples := createTestWaveform(0.5, 100)
	gain := 10.0
	driven := Drive(samples, gain)

	// All samples should be clipped to 1.0
	for i, v := range driven {
		if v != 1.0 {
			t.Errorf("Expected driven sample to be clipped at 1.0 at index %d, got %f", i, v)
		}
	}
}

func TestLimiterNoClipping(t *testing.T) {
	samples := createTestWaveform(0.5, 100)
	limited := Limiter(samples)

	for i, v := range limited {
		if v != samples[i] {
			t.Errorf("Expected limiter to not alter sample at index %d, got %f", i, v)
		}
	}
}

func TestLimiterClipping(t *testing.T) {
	samples := createTestWaveform(1.5, 100) // Samples exceeding 1.0
	limited := Limiter(samples)

	for i, v := range limited {
		if v != 1.0 {
			t.Errorf("Expected limiter to clip sample at index %d to 1.0, got %f", i, v)
		}
	}
}

func TestFindPeakAmplitudeZeroSamples(t *testing.T) {
	samples := []float64{}
	peak := FindPeakAmplitude(samples)
	if peak != 0.0 {
		t.Errorf("Expected peak amplitude of 0.0 for empty samples, got %f", peak)
	}
}

func TestNormalizeSamplesZeroPeak(t *testing.T) {
	samples := createTestWaveform(0.0, 100)
	normalized := NormalizeSamples(samples, 1.0)

	for i, v := range normalized {
		if v != 0.0 {
			t.Errorf("Expected normalized sample to be 0.0 at index %d, got %f", i, v)
		}
	}
}

func TestPadSamplesEqualLength(t *testing.T) {
	wave1 := createTestWaveform(0.5, 100)
	wave2 := createTestWaveform(0.5, 100)

	padded1, padded2 := PadSamples(wave1, wave2)
	if len(padded1) != 100 || len(padded2) != 100 {
		t.Errorf("Expected padded waves to have length 100, got %d and %d", len(padded1), len(padded2))
	}
}

func TestPadSamplesFirstShorter(t *testing.T) {
	wave1 := createTestWaveform(0.5, 50)
	wave2 := createTestWaveform(0.5, 100)

	padded1, padded2 := PadSamples(wave1, wave2)
	if len(padded1) != 100 || len(padded2) != 100 {
		t.Errorf("Expected padded waves to have length 100, got %d and %d", len(padded1), len(padded2))
	}
}

func TestPadSamplesSecondShorter(t *testing.T) {
	wave1 := createTestWaveform(0.5, 100)
	wave2 := createTestWaveform(0.5, 50)

	padded1, padded2 := PadSamples(wave1, wave2)
	if len(padded1) != 100 || len(padded2) != 100 {
		t.Errorf("Expected padded waves to have length 100, got %d and %d", len(padded1), len(padded2))
	}
}

func TestAnalyzeHighestFrequencyZeroSamples(t *testing.T) {
	samples := []float64{}
	frequency := AnalyzeHighestFrequency(samples, 44100)
	if frequency != 0.0 {
		t.Errorf("Expected frequency of 0.0 for empty samples, got %f", frequency)
	}
}

func TestAnalyzeHighestFrequencyConstantSignal(t *testing.T) {
	samples := createTestWaveform(1.0, 1000)
	frequency := AnalyzeHighestFrequency(samples, 44100)
	if frequency != 0.0 {
		t.Errorf("Expected frequency of 0.0 for constant signal, got %f", frequency)
	}
}

func TestWeightedSummationZeroWeights(t *testing.T) {
	wave1 := createTestWaveform(0.5, 100)
	wave2 := createTestWaveform(0.5, 100)
	weights := []float64{0.0, 0.0}

	result, err := WeightedSummation(weights, wave1, wave2)
	if err != nil {
		t.Fatalf("WeightedSummation failed: %v", err)
	}

	for i, v := range result {
		if v != 0.0 {
			t.Errorf("Expected summed sample to be 0.0 at index %d, got %f", i, v)
		}
	}
}

func TestWeightedSummationNegativeWeights(t *testing.T) {
	wave1 := createTestWaveform(0.5, 100)
	wave2 := createTestWaveform(0.5, 100)
	weights := []float64{-1.0, -1.0}

	result, err := WeightedSummation(weights, wave1, wave2)
	if err != nil {
		t.Fatalf("WeightedSummation failed: %v", err)
	}

	for i, v := range result {
		if v != -1.0 {
			t.Errorf("Expected summed sample to be -1.0 at index %d, got %f", i, v)
		}
	}
}

func TestApplyFadeInLinear(t *testing.T) {
	samples := make([]float64, 100)
	for i := range samples {
		samples[i] = 1.0
	}

	// Keep a copy of the original samples
	originalSamples := make([]float64, len(samples))
	copy(originalSamples, samples)

	fadeDuration := 0.5 // seconds
	sampleRate := 100   // for simplicity
	fadedSamples := ApplyFadeIn(samples, fadeDuration, sampleRate, LinearFade)
	expectedFadeSamples := int(fadeDuration * float64(sampleRate))
	for i := 0; i < expectedFadeSamples; i++ {
		expectedMultiplier := float64(i) / float64(expectedFadeSamples)
		expectedValue := originalSamples[i] * expectedMultiplier
		if math.Abs(fadedSamples[i]-expectedValue) > 1e-6 {
			t.Errorf("Linear fade-in not applied correctly at index %d: expected %f, got %f", i, expectedValue, fadedSamples[i])
		}
	}
}

func TestApplyFadeInQuadratic(t *testing.T) {
	samples := make([]float64, 100)
	for i := range samples {
		samples[i] = 1.0
	}

	// Keep a copy of the original samples
	originalSamples := make([]float64, len(samples))
	copy(originalSamples, samples)

	fadeDuration := 0.5 // seconds
	sampleRate := 100   // for simplicity
	fadedSamples := ApplyFadeIn(samples, fadeDuration, sampleRate, QuadraticFade)
	expectedFadeSamples := int(fadeDuration * float64(sampleRate))
	for i := 0; i < expectedFadeSamples; i++ {
		tVal := float64(i) / float64(expectedFadeSamples)
		expectedMultiplier := tVal * tVal
		expectedValue := originalSamples[i] * expectedMultiplier
		if math.Abs(fadedSamples[i]-expectedValue) > 1e-6 {
			t.Errorf("Quadratic fade-in not applied correctly at index %d: expected %f, got %f", i, expectedValue, fadedSamples[i])
		}
	}
}

func TestApplyFadeOutLinear(t *testing.T) {
	samples := make([]float64, 100)
	for i := range samples {
		samples[i] = 1.0
	}

	// Keep a copy of the original samples
	originalSamples := make([]float64, len(samples))
	copy(originalSamples, samples)

	fadeDuration := 0.5 // seconds
	sampleRate := 100   // for simplicity
	fadedSamples := ApplyFadeOut(samples, fadeDuration, sampleRate, LinearFade)
	expectedFadeSamples := int(fadeDuration * float64(sampleRate))
	totalSamples := len(samples)
	for i := 0; i < expectedFadeSamples; i++ {
		expectedMultiplier := 1.0 - float64(i)/float64(expectedFadeSamples)
		index := totalSamples - expectedFadeSamples + i
		expectedValue := originalSamples[index] * expectedMultiplier
		if math.Abs(fadedSamples[index]-expectedValue) > 1e-6 {
			t.Errorf("Linear fade-out not applied correctly at index %d: expected %f, got %f", index, expectedValue, fadedSamples[index])
		}
	}
}

func TestApplyFadeOutQuadratic(t *testing.T) {
	samples := make([]float64, 100)
	for i := range samples {
		samples[i] = 1.0
	}

	// Keep a copy of the original samples
	originalSamples := make([]float64, len(samples))
	copy(originalSamples, samples)

	fadeDuration := 0.5 // seconds
	sampleRate := 100   // for simplicity
	fadedSamples := ApplyFadeOut(samples, fadeDuration, sampleRate, QuadraticFade)
	expectedFadeSamples := int(fadeDuration * float64(sampleRate))
	totalSamples := len(samples)
	for i := 0; i < expectedFadeSamples; i++ {
		tVal := 1.0 - float64(i)/float64(expectedFadeSamples)
		expectedMultiplier := tVal * tVal
		index := totalSamples - expectedFadeSamples + i
		expectedValue := originalSamples[index] * expectedMultiplier
		if math.Abs(fadedSamples[index]-expectedValue) > 1e-6 {
			t.Errorf("Quadratic fade-out not applied correctly at index %d: expected %f, got %f", index, expectedValue, fadedSamples[index])
		}
	}
}
