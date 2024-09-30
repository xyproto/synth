package synth

import (
	"errors"
	"io"
	"math/rand"
)

type SoundType int

const ( // SoundType
	Kick = iota
	Clap
	Snare
	ClosedHH
	OpenHH
	Rimshot
	Tom
	Percussion
	Ride
	Crash
	Bass
	Xylophone
	Lead
)

// NewSettings creates a new Settings instance with default values for a percussive sound
func NewSettings(output io.WriteSeeker, startFreq, endFreq, duration float64, sampleRate, bitDepth, channels int) (*Settings, error) {
	if sampleRate <= 0 || duration <= 0 || channels <= 0 {
		return nil, errors.New("invalid sample rate, duration or channels")
	}
	return &Settings{
		SoundType:                  Kick, // by default
		SampleRate:                 sampleRate,
		BitDepth:                   bitDepth, // Audio bit depth (8, 16, 24, or 32 bits)
		Channels:                   channels,
		Output:                     output,                           // Output writer
		StartFreq:                  startFreq,                        // Starting frequency (Hz)
		EndFreq:                    endFreq,                          // Ending frequency (Hz)
		Duration:                   duration,                         // Duration in seconds
		WaveformType:               WaveSine,                         // Default waveform type
		NoiseAmount:                0.0,                              // Default noise amount
		Attack:                     0.005,                            // Attack time in seconds
		Decay:                      0.2,                              // Decay time in seconds
		Sustain:                    0.1,                              // Sustain level
		Release:                    0.1,                              // Release time in seconds
		Drive:                      0.2,                              // Drive (distortion) amount
		FilterCutoff:               8000,                             // Filter cutoff frequency (Hz)
		Sweep:                      0.7,                              // Pitch sweep amount
		PitchDecay:                 0.4,                              // Pitch envelope decay time
		NumOscillators:             1,                                // Default to 1 oscillator
		OscillatorLevels:           []float64{1.0},                   // Default oscillator level
		SaturatorAmount:            0.3,                              // Saturation amount
		FilterResonance:            1.0,                              // Filter resonance
		FilterBands:                []float64{200.0, 1000.0, 3000.0}, // Multi-band filter frequencies
		FadeDuration:               0.01,                             // Fade in/out duration in seconds
		SmoothFrequencyTransitions: true,                             // Enable smooth frequency transitions by default
	}, nil
}

// NewRandom generates random settings for the given sound type.
// It handles each SoundType by initializing a Settings instance with randomized parameters suitable for that sound.
func NewRandom(soundType SoundType, output io.WriteSeeker, sampleRate, bitDepth, channels int) *Settings {
	switch soundType {
	case Kick:
		cfg, _ := NewSettings(output, 50.0, 30.0, 1.0, sampleRate, bitDepth, channels)
		cfg.SoundType = soundType
		cfg.Attack = rand.Float64()*0.02 + 0.001       // 0.001 to 0.021 seconds
		cfg.Decay = rand.Float64()*0.5 + 0.1           // 0.1 to 0.6 seconds
		cfg.Sustain = rand.Float64() * 0.5             // 0.0 to 0.5
		cfg.Release = rand.Float64()*0.5 + 0.1         // 0.1 to 0.6 seconds
		cfg.Drive = rand.Float64() * 0.7               // 0.0 to 0.7
		cfg.FilterCutoff = 1500 + rand.Float64()*5000  // 1500 to 6500 Hz
		cfg.FilterResonance = 0.5 + rand.Float64()*1.5 // 0.5 to 2.0
		cfg.Sweep = rand.Float64() * 1.5               // 0.0 to 1.5
		cfg.PitchDecay = rand.Float64() * 1.5          // 0.0 to 1.5
		cfg.FadeDuration = rand.Float64() * 0.1        // 0.0 to 0.1 seconds
		cfg.SmoothFrequencyTransitions = rand.Float64() >= 0.1
		if rand.Float64() < 0.2 { // 20% chance to use varied waveforms
			cfg.WaveformType = rand.Intn(7)
		} else {
			cfg.WaveformType = rand.Intn(2) // Sine or Triangle
		}
		return cfg
	case Clap:
		cfg, _ := NewSettings(output, 300.0, 200.0, 0.3, sampleRate, bitDepth, channels)
		cfg.SoundType = soundType
		cfg.NoiseAmount = rand.Float64()*0.5 + 0.5     // 0.5 to 1.0
		cfg.Attack = rand.Float64()*0.01 + 0.001       // 0.001 to 0.011 seconds
		cfg.Decay = rand.Float64()*0.2 + 0.05          // 0.05 to 0.25 seconds
		cfg.Sustain = rand.Float64() * 0.3             // 0.0 to 0.3
		cfg.Release = rand.Float64()*0.2 + 0.05        // 0.05 to 0.25 seconds
		cfg.Drive = rand.Float64() * 0.5               // 0.0 to 0.5
		cfg.FilterCutoff = 4000 + rand.Float64()*3000  // 4000 to 7000 Hz
		cfg.FilterResonance = 0.8 + rand.Float64()*1.2 // 0.8 to 2.0
		cfg.Sweep = rand.Float64() * 1.2               // 0.0 to 1.2
		cfg.PitchDecay = rand.Float64() * 1.2          // 0.0 to 1.2
		cfg.FadeDuration = rand.Float64() * 0.05       // 0.0 to 0.05 seconds
		cfg.SmoothFrequencyTransitions = rand.Float64() >= 0.2
		cfg.WaveformType = rand.Intn(4) // Sine, Triangle, Sawtooth, Square
		return cfg
	case Snare:
		cfg, _ := NewSettings(output, 300.0, 150.0, 0.5, sampleRate, bitDepth, channels)
		cfg.SoundType = soundType
		cfg.NoiseAmount = 0.5 + rand.Float64()*0.5     // 0.5 to 1.0
		cfg.Attack = rand.Float64() * 0.01             // 0.0 to 0.01 seconds
		cfg.Decay = 0.05 + rand.Float64()*0.2          // 0.05 to 0.25 seconds
		cfg.Sustain = 0.0                              // No sustain for snare
		cfg.Release = 0.05 + rand.Float64()*0.2        // 0.05 to 0.25 seconds
		cfg.Drive = rand.Float64() * 0.7               // 0.0 to 0.7
		cfg.FilterCutoff = 5000 + rand.Float64()*3000  // 5000 to 8000 Hz
		cfg.FilterResonance = 1.0 + rand.Float64()*1.0 // 1.0 to 2.0
		cfg.FadeDuration = rand.Float64() * 0.05       // 0.0 to 0.05 seconds
		cfg.SmoothFrequencyTransitions = rand.Float64() >= 0.1
		cfg.WaveformType = rand.Intn(4) // Sine, Triangle, Sawtooth, Square
		return cfg
	case ClosedHH:
		cfg, _ := NewSettings(output, 8000.0, 5000.0, 0.1, sampleRate, bitDepth, channels)
		cfg.SoundType = soundType
		cfg.NoiseAmount = 0.3 + rand.Float64()*0.4     // 0.3 to 0.7
		cfg.Attack = rand.Float64() * 0.005            // 0.0 to 0.005 seconds
		cfg.Decay = 0.05 + rand.Float64()*0.15         // 0.05 to 0.2 seconds
		cfg.Sustain = 0.0                              // No sustain for percussive sounds
		cfg.Release = 0.05 + rand.Float64()*0.15       // 0.05 to 0.2 seconds
		cfg.Drive = rand.Float64() * 0.5               // 0.0 to 0.5
		cfg.FilterCutoff = 6000 + rand.Float64()*4000  // 6000 to 10000 Hz
		cfg.FilterResonance = 0.5 + rand.Float64()*1.5 // 0.5 to 2.0
		cfg.FadeDuration = rand.Float64() * 0.02       // 0.0 to 0.02 seconds
		cfg.SmoothFrequencyTransitions = rand.Float64() >= 0.15
		cfg.WaveformType = rand.Intn(4) // Sine, Triangle, Sawtooth, Square
		return cfg
	case OpenHH:
		cfg, _ := NewSettings(output, 10000.0, 7000.0, 0.3, sampleRate, bitDepth, channels)
		cfg.SoundType = soundType
		cfg.NoiseAmount = 0.4 + rand.Float64()*0.4     // 0.4 to 0.8
		cfg.Attack = rand.Float64() * 0.01             // 0.0 to 0.01 seconds
		cfg.Decay = 0.1 + rand.Float64()*0.2           // 0.1 to 0.3 seconds
		cfg.Sustain = 0.0                              // No sustain
		cfg.Release = 0.1 + rand.Float64()*0.2         // 0.1 to 0.3 seconds
		cfg.Drive = rand.Float64() * 0.6               // 0.0 to 0.6
		cfg.FilterCutoff = 7000 + rand.Float64()*3000  // 7000 to 10000 Hz
		cfg.FilterResonance = 1.0 + rand.Float64()*1.0 // 1.0 to 2.0
		cfg.FadeDuration = rand.Float64() * 0.03       // 0.0 to 0.03 seconds
		cfg.SmoothFrequencyTransitions = rand.Float64() >= 0.1
		cfg.WaveformType = rand.Intn(4) // Sine, Triangle, Sawtooth, Square
		return cfg
	case Rimshot:
		cfg, _ := NewSettings(output, 4000.0, 2000.0, 0.2, sampleRate, bitDepth, channels)
		cfg.SoundType = soundType
		cfg.NoiseAmount = 0.6 + rand.Float64()*0.4     // 0.6 to 1.0
		cfg.Attack = rand.Float64() * 0.005            // 0.0 to 0.005 seconds
		cfg.Decay = 0.05 + rand.Float64()*0.15         // 0.05 to 0.2 seconds
		cfg.Sustain = 0.0                              // No sustain
		cfg.Release = 0.05 + rand.Float64()*0.15       // 0.05 to 0.2 seconds
		cfg.Drive = rand.Float64() * 0.6               // 0.0 to 0.6
		cfg.FilterCutoff = 8000 + rand.Float64()*2000  // 8000 to 10000 Hz
		cfg.FilterResonance = 1.0 + rand.Float64()*1.0 // 1.0 to 2.0
		cfg.FadeDuration = rand.Float64() * 0.02       // 0.0 to 0.02 seconds
		cfg.SmoothFrequencyTransitions = rand.Float64() >= 0.1
		cfg.WaveformType = rand.Intn(4) // Sine, Triangle, Sawtooth, Square
		return cfg
	case Tom:
		cfg, _ := NewSettings(output, 200.0, 100.0, 0.4, sampleRate, bitDepth, channels)
		cfg.SoundType = soundType
		cfg.Attack = rand.Float64()*0.01 + 0.001       // 0.001 to 0.011 seconds
		cfg.Decay = 0.1 + rand.Float64()*0.3           // 0.1 to 0.4 seconds
		cfg.Sustain = 0.0                              // No sustain
		cfg.Release = 0.1 + rand.Float64()*0.3         // 0.1 to 0.4 seconds
		cfg.Drive = rand.Float64() * 0.5               // 0.0 to 0.5
		cfg.FilterCutoff = 3000 + rand.Float64()*2000  // 3000 to 5000 Hz
		cfg.FilterResonance = 0.8 + rand.Float64()*1.2 // 0.8 to 2.0
		cfg.FadeDuration = rand.Float64() * 0.025      // 0.0 to 0.025 seconds
		cfg.SmoothFrequencyTransitions = rand.Float64() >= 0.15
		cfg.WaveformType = rand.Intn(4) // Sine, Triangle, Sawtooth, Square
		return cfg
	case Percussion:
		cfg, _ := NewSettings(output, 1000.0, 500.0, 0.2, sampleRate, bitDepth, channels)
		cfg.SoundType = soundType
		cfg.NoiseAmount = 0.4 + rand.Float64()*0.4     // 0.4 to 0.8
		cfg.Attack = rand.Float64() * 0.005            // 0.0 to 0.005 seconds
		cfg.Decay = 0.05 + rand.Float64()*0.15         // 0.05 to 0.2 seconds
		cfg.Sustain = 0.0                              // No sustain
		cfg.Release = 0.05 + rand.Float64()*0.15       // 0.05 to 0.2 seconds
		cfg.Drive = rand.Float64() * 0.5               // 0.0 to 0.5
		cfg.FilterCutoff = 6000 + rand.Float64()*3000  // 6000 to 9000 Hz
		cfg.FilterResonance = 1.0 + rand.Float64()*1.0 // 1.0 to 2.0
		cfg.FadeDuration = rand.Float64() * 0.03       // 0.0 to 0.03 seconds
		cfg.SmoothFrequencyTransitions = rand.Float64() >= 0.1
		cfg.WaveformType = rand.Intn(4) // Sine, Triangle, Sawtooth, Square
		return cfg
	case Ride:
		cfg, _ := NewSettings(output, 12000.0, 7000.0, 0.4, sampleRate, bitDepth, channels)
		cfg.SoundType = soundType
		cfg.NoiseAmount = 0.5 + rand.Float64()*0.5     // 0.5 to 1.0
		cfg.Attack = rand.Float64() * 0.01             // 0.0 to 0.01 seconds
		cfg.Decay = 0.1 + rand.Float64()*0.3           // 0.1 to 0.4 seconds
		cfg.Sustain = 0.2 + rand.Float64()*0.3         // 0.2 to 0.5
		cfg.Release = 0.1 + rand.Float64()*0.3         // 0.1 to 0.4 seconds
		cfg.Drive = rand.Float64() * 0.6               // 0.0 to 0.6
		cfg.FilterCutoff = 8000 + rand.Float64()*4000  // 8000 to 12000 Hz
		cfg.FilterResonance = 1.0 + rand.Float64()*1.0 // 1.0 to 2.0
		cfg.FadeDuration = rand.Float64() * 0.03       // 0.0 to 0.03 seconds
		cfg.SmoothFrequencyTransitions = rand.Float64() >= 0.15
		cfg.WaveformType = rand.Intn(4) // Sine, Triangle, Sawtooth, Square
		return cfg
	case Crash:
		cfg, _ := NewSettings(output, 15000.0, 10000.0, 0.3, sampleRate, bitDepth, channels)
		cfg.SoundType = soundType
		cfg.NoiseAmount = 0.6 + rand.Float64()*0.4     // 0.6 to 1.0
		cfg.Attack = rand.Float64() * 0.005            // 0.0 to 0.005 seconds
		cfg.Decay = 0.1 + rand.Float64()*0.25          // 0.1 to 0.35 seconds
		cfg.Sustain = 0.3 + rand.Float64()*0.4         // 0.3 to 0.7
		cfg.Release = 0.1 + rand.Float64()*0.3         // 0.1 to 0.4 seconds
		cfg.Drive = rand.Float64() * 0.7               // 0.0 to 0.7
		cfg.FilterCutoff = 10000 + rand.Float64()*5000 // 10000 to 15000 Hz
		cfg.FilterResonance = 1.2 + rand.Float64()*1.0 // 1.2 to 2.2
		cfg.FadeDuration = rand.Float64() * 0.04       // 0.0 to 0.04 seconds
		cfg.SmoothFrequencyTransitions = rand.Float64() >= 0.2
		cfg.WaveformType = rand.Intn(4) // Sine, Triangle, Sawtooth, Square
		return cfg
	case Bass:
		cfg, _ := NewSettings(output, 60.0, 30.0, 1.0, sampleRate, bitDepth, channels)
		cfg.SoundType = soundType
		cfg.Attack = rand.Float64()*0.01 + 0.001       // 0.001 to 0.011 seconds
		cfg.Decay = 0.1 + rand.Float64()*0.4           // 0.1 to 0.5 seconds
		cfg.Sustain = rand.Float64() * 0.5             // 0.0 to 0.5
		cfg.Release = 0.1 + rand.Float64()*0.4         // 0.1 to 0.5 seconds
		cfg.Drive = rand.Float64() * 0.7               // 0.0 to 0.7
		cfg.FilterCutoff = 150.0 + rand.Float64()*500  // 150.0 to 650.0 Hz
		cfg.FilterResonance = 0.5 + rand.Float64()*1.5 // 0.5 to 2.0
		cfg.Sweep = rand.Float64() * 1.0               // 0.0 to 1.0
		cfg.PitchDecay = rand.Float64() * 1.0          // 0.0 to 1.0
		cfg.FadeDuration = rand.Float64() * 0.05       // 0.0 to 0.05 seconds
		cfg.SmoothFrequencyTransitions = rand.Float64() >= 0.2
		if rand.Float64() < 0.3 { // 30% chance to use varied waveforms
			cfg.WaveformType = rand.Intn(7)
		} else {
			cfg.WaveformType = rand.Intn(2) // Sine or Triangle
		}
		return cfg
	case Xylophone:
		cfg, _ := NewSettings(output, 1000.0, 500.0, 0.2, sampleRate, bitDepth, channels)
		cfg.SoundType = soundType
		cfg.Attack = rand.Float64()*0.005 + 0.001      // 0.001 to 0.006 seconds
		cfg.Decay = 0.05 + rand.Float64()*0.15         // 0.05 to 0.2 seconds
		cfg.Sustain = 0.0                              // No sustain
		cfg.Release = 0.05 + rand.Float64()*0.15       // 0.05 to 0.2 seconds
		cfg.Drive = rand.Float64() * 0.5               // 0.0 to 0.5
		cfg.FilterCutoff = 8000 + rand.Float64()*2000  // 8000 to 10000 Hz
		cfg.FilterResonance = 1.0 + rand.Float64()*1.0 // 1.0 to 2.0
		cfg.FadeDuration = rand.Float64() * 0.02       // 0.0 to 0.02 seconds
		cfg.SmoothFrequencyTransitions = rand.Float64() >= 0.1
		cfg.WaveformType = rand.Intn(4) // Sine, Triangle, Sawtooth, Square
		return cfg
	case Lead:
		cfg, _ := NewSettings(output, 880.0, 440.0, 0.5, sampleRate, bitDepth, channels)
		cfg.SoundType = soundType
		cfg.Attack = rand.Float64()*0.02 + 0.005       // 0.005 to 0.025 seconds
		cfg.Decay = 0.2 + rand.Float64()*0.5           // 0.2 to 0.7 seconds
		cfg.Sustain = 0.3 + rand.Float64()*0.7         // 0.3 to 1.0
		cfg.Release = 0.2 + rand.Float64()*0.5         // 0.2 to 0.7 seconds
		cfg.Drive = rand.Float64() * 0.8               // 0.0 to 0.8
		cfg.FilterCutoff = 5000 + rand.Float64()*5000  // 5000 to 10000 Hz
		cfg.FilterResonance = 0.7 + rand.Float64()*1.3 // 0.7 to 2.0
		cfg.Sweep = rand.Float64() * 2.0               // 0.0 to 2.0
		cfg.PitchDecay = rand.Float64() * 2.0          // 0.0 to 2.0
		cfg.FadeDuration = rand.Float64() * 0.05       // 0.0 to 0.05 seconds
		cfg.SmoothFrequencyTransitions = rand.Float64() >= 0.2
		cfg.WaveformType = rand.Intn(4) // Sine, Triangle, Sawtooth, Square
		return cfg
	default:
		// Handle unknown SoundType by returning a default Settings with randomized parameters
		cfg, _ := NewSettings(output, 500.0, 250.0, 0.5, sampleRate, bitDepth, channels)
		cfg.SoundType = soundType
		cfg.Attack = rand.Float64()*0.02 + 0.001       // 0.001 to 0.021 seconds
		cfg.Decay = rand.Float64()*0.3 + 0.1           // 0.1 to 0.4 seconds
		cfg.Sustain = rand.Float64() * 0.5             // 0.0 to 0.5
		cfg.Release = rand.Float64()*0.3 + 0.1         // 0.1 to 0.4 seconds
		cfg.Drive = rand.Float64() * 0.6               // 0.0 to 0.6
		cfg.FilterCutoff = 3000 + rand.Float64()*7000  // 3000 to 10000 Hz
		cfg.FilterResonance = 1.0 + rand.Float64()*1.0 // 1.0 to 2.0
		cfg.Sweep = rand.Float64() * 1.5               // 0.0 to 1.5
		cfg.PitchDecay = rand.Float64() * 1.5          // 0.0 to 1.5
		cfg.FadeDuration = rand.Float64() * 0.05       // 0.0 to 0.05 seconds
		cfg.SmoothFrequencyTransitions = rand.Float64() >= 0.1
		cfg.WaveformType = rand.Intn(4) // Sine, Triangle, Sawtooth, Square
		return cfg
	}
}

// NewSnareSettings creates a new Settings instance with predefined values optimized for snare drum sounds.
func NewSnareSettings(output io.WriteSeeker, sampleRate, bitDepth, channels int) (*Settings, error) {
	if sampleRate <= 0 || bitDepth <= 0 || channels <= 0 {
		return nil, errors.New("invalid sample rate, bit depth, or channels")
	}
	return &Settings{
		SoundType:                  Snare,
		SampleRate:                 sampleRate,
		BitDepth:                   bitDepth,   // e.g., 16
		Channels:                   channels,   // e.g., 1 for mono, 2 for stereo
		Output:                     output,     // Output writer
		StartFreq:                  300.0,      // Starting frequency in Hz
		EndFreq:                    100.0,      // Ending frequency in Hz (rapid pitch drop)
		Duration:                   0.5,        // Short duration for snare snap
		WaveformType:               WaveSquare, // Square waveform for a brighter sound
		NoiseAmount:                0.7,        // Significant noise for snare rattle
		Attack:                     0.005,      // Very quick attack
		Decay:                      0.1,        // Short decay
		Sustain:                    0.0,        // No sustain
		Release:                    0.05,       // Short release
		Drive:                      0.3,        // Moderate drive for slight distortion
		FilterCutoff:               5000.0,     // High cutoff to retain brightness
		Sweep:                      0.5,        // Moderate pitch sweep
		PitchDecay:                 0.3,        // Pitch envelope decay
		NumOscillators:             1,          // Single oscillator
		OscillatorLevels:           []float64{1.0},
		SaturatorAmount:            0.2,                              // Mild saturation
		FilterResonance:            1.0,                              // Standard resonance
		FilterBands:                []float64{500.0, 2000.0, 6000.0}, // Example multi-band frequencies
		FadeDuration:               0.01,                             // Quick fade to prevent clicks
		SmoothFrequencyTransitions: true,                             // Enable smooth transitions
	}, nil
}

// New606Kick creates a 606-style kick drum
func New606Kick(output io.WriteSeeker, duration float64, sampleRate, bitDepth, channels int) (*Settings, error) {
	cfg, err := NewSettings(output, 80.0, 50.0, duration, sampleRate, bitDepth, channels)
	if err != nil {
		return nil, err
	}
	cfg.SoundType = Kick
	cfg.WaveformType = WaveSine
	cfg.Attack = 0.002
	cfg.Decay = 0.25
	cfg.Sustain = 0.0
	cfg.Release = 0.15
	cfg.Drive = 0.3
	cfg.FilterCutoff = 4000
	cfg.FilterResonance = 1.2
	cfg.Sweep = 0.8
	cfg.PitchDecay = 0.6
	cfg.FadeDuration = 0.008
	cfg.SmoothFrequencyTransitions = true
	return cfg, nil
}

// New707Kick creates a 707-style kick drum
func New707Kick(output io.WriteSeeker, duration float64, sampleRate, bitDepth, channels int) (*Settings, error) {
	cfg, err := NewSettings(output, 85.0, 55.0, duration, sampleRate, bitDepth, channels)
	if err != nil {
		return nil, err
	}
	cfg.SoundType = Kick
	cfg.WaveformType = WaveTriangle
	cfg.Attack = 0.003
	cfg.Decay = 0.3
	cfg.Sustain = 0.0
	cfg.Release = 0.2
	cfg.Drive = 0.4
	cfg.FilterCutoff = 4500
	cfg.FilterResonance = 1.0
	cfg.Sweep = 0.7
	cfg.PitchDecay = 0.5
	cfg.FadeDuration = 0.01
	cfg.SmoothFrequencyTransitions = true
	return cfg, nil
}

// New808Kick creates an 808-style kick drum with improved parameters
func New808Kick(output io.WriteSeeker, duration float64, sampleRate, bitDepth, channels int) (*Settings, error) {
	cfg, err := NewSettings(output, 55.0, 30.0, duration, sampleRate, bitDepth, channels)
	if err != nil {
		return nil, err
	}
	cfg.SoundType = Kick
	cfg.WaveformType = WaveSine
	cfg.Attack = 0.0 // Instant attack for a punchier sound
	cfg.Decay = 0.6
	cfg.Sustain = 0.0
	cfg.Release = 0.4
	cfg.Drive = 0.5
	cfg.FilterCutoff = 2000 // Lower cutoff for deeper bass
	cfg.FilterResonance = 1.5
	cfg.Sweep = 1.0
	cfg.PitchDecay = 0.8
	cfg.FadeDuration = 0.005
	cfg.SmoothFrequencyTransitions = true
	return cfg, nil
}

// New909Kick creates a 909-style kick drum with enhanced punch
func New909Kick(output io.WriteSeeker, duration float64, sampleRate, bitDepth, channels int) (*Settings, error) {
	cfg, err := NewSettings(output, 60.0, 50.0, duration, sampleRate, bitDepth, channels)
	if err != nil {
		return nil, err
	}
	cfg.SoundType = Kick
	cfg.WaveformType = WaveTriangle
	cfg.Attack = 0.0
	cfg.Decay = 0.25
	cfg.Sustain = 0.0
	cfg.Release = 0.2
	cfg.Drive = 0.7
	cfg.FilterCutoff = 2500
	cfg.FilterResonance = 1.4
	cfg.Sweep = 0.8
	cfg.PitchDecay = 0.6
	cfg.FadeDuration = 0.005
	cfg.SmoothFrequencyTransitions = true
	return cfg, nil
}

// NewLinnDrumKick creates a LinnDrum-style kick drum sound
func NewLinnDrumKick(output io.WriteSeeker, duration float64, sampleRate, bitDepth, channels int) (*Settings, error) {
	cfg, err := NewSettings(output, 75.0, 45.0, duration, sampleRate, bitDepth, channels)
	if err != nil {
		return nil, err
	}
	cfg.SoundType = Kick
	cfg.WaveformType = WaveSquare
	cfg.Attack = 0.004
	cfg.Decay = 0.35
	cfg.Sustain = 0.0
	cfg.Release = 0.25
	cfg.Drive = 0.5
	cfg.FilterCutoff = 3000
	cfg.FilterResonance = 1.0
	cfg.Sweep = 0.6
	cfg.PitchDecay = 0.5
	cfg.FadeDuration = 0.015
	cfg.SmoothFrequencyTransitions = true
	return cfg, nil
}

// NewDeepHouseKick creates a Deep House kick drum with richer bass
func NewDeepHouseKick(output io.WriteSeeker, duration float64, sampleRate, bitDepth, channels int) (*Settings, error) {
	cfg, err := NewSettings(output, 50.0, 30.0, duration, sampleRate, bitDepth, channels)
	if err != nil {
		return nil, err
	}
	cfg.SoundType = Kick
	cfg.WaveformType = WaveSine
	cfg.Attack = 0.005
	cfg.Decay = 1.0
	cfg.Sustain = 0.0
	cfg.Release = 0.8
	cfg.Drive = 0.6
	cfg.FilterCutoff = 1500
	cfg.FilterResonance = 1.2
	cfg.Sweep = 0.9
	cfg.PitchDecay = 0.7
	cfg.FadeDuration = 0.02
	cfg.SmoothFrequencyTransitions = true
	return cfg, nil
}

// NewExperimentalKick creates an experimental kick drum sound
func NewExperimentalKick(output io.WriteSeeker, duration float64, sampleRate, bitDepth, channels int) (*Settings, error) {
	cfg, err := NewSettings(output, 70.0, 20.0, duration, sampleRate, bitDepth, channels)
	if err != nil {
		return nil, err
	}
	cfg.SoundType = Kick
	cfg.WaveformType = WaveSawtooth
	cfg.Attack = 0.001
	cfg.Decay = 0.5
	cfg.Sustain = 0.0
	cfg.Release = 0.3
	cfg.Drive = 0.9
	cfg.FilterCutoff = 5000
	cfg.FilterResonance = 1.5
	cfg.Sweep = 1.2
	cfg.PitchDecay = 0.9
	cfg.FadeDuration = 0.01
	cfg.SmoothFrequencyTransitions = true
	return cfg, nil
}

// Snare Drum Functions

// New606Snare creates a snare drum sound similar to the Roland TR-606
func New606Snare(output io.WriteSeeker, duration float64, sampleRate, bitDepth, channels int) (*Settings, error) {
	cfg, err := NewSettings(output, 350.0, 200.0, duration, sampleRate, bitDepth, channels)
	if err != nil {
		return nil, err
	}
	cfg.SoundType = Snare
	cfg.WaveformType = WaveTriangle
	cfg.NoiseAmount = 0.6
	cfg.Attack = 0.005
	cfg.Decay = 0.15
	cfg.Sustain = 0.0
	cfg.Release = 0.1
	cfg.Drive = 0.3
	cfg.FilterCutoff = 7000
	cfg.FilterResonance = 1.2
	cfg.FadeDuration = 0.01
	cfg.SmoothFrequencyTransitions = true
	return cfg, nil
}

// New707Snare creates a snare drum sound similar to the Roland TR-707
func New707Snare(output io.WriteSeeker, duration float64, sampleRate, bitDepth, channels int) (*Settings, error) {
	cfg, err := NewSettings(output, 300.0, 150.0, duration, sampleRate, bitDepth, channels)
	if err != nil {
		return nil, err
	}
	cfg.SoundType = Snare
	cfg.WaveformType = WaveSquare
	cfg.NoiseAmount = 0.5
	cfg.Attack = 0.001
	cfg.Decay = 0.2
	cfg.Sustain = 0.0
	cfg.Release = 0.15
	cfg.Drive = 0.4
	cfg.FilterCutoff = 6500
	cfg.FilterResonance = 1.0
	cfg.FadeDuration = 0.008
	cfg.SmoothFrequencyTransitions = true
	return cfg, nil
}

// New808Snare creates a snare drum sound similar to the Roland TR-808
func New808Snare(output io.WriteSeeker, duration float64, sampleRate, bitDepth, channels int) (*Settings, error) {
	cfg, err := NewSettings(output, 240.0, 120.0, duration, sampleRate, bitDepth, channels)
	if err != nil {
		return nil, err
	}
	cfg.SoundType = Snare
	cfg.WaveformType = WaveSine
	cfg.NoiseAmount = 0.7
	cfg.Attack = 0.002
	cfg.Decay = 0.25
	cfg.Sustain = 0.0
	cfg.Release = 0.2
	cfg.Drive = 0.5
	cfg.FilterCutoff = 6000
	cfg.FilterResonance = 1.5
	cfg.FadeDuration = 0.015
	cfg.SmoothFrequencyTransitions = true
	return cfg, nil
}

// New909Snare creates a snare drum sound similar to the Roland TR-909
func New909Snare(output io.WriteSeeker, duration float64, sampleRate, bitDepth, channels int) (*Settings, error) {
	cfg, err := NewSettings(output, 250.0, 130.0, duration, sampleRate, bitDepth, channels)
	if err != nil {
		return nil, err
	}
	cfg.SoundType = Snare
	cfg.WaveformType = WaveSawtooth
	cfg.NoiseAmount = 0.6
	cfg.Attack = 0.003
	cfg.Decay = 0.18
	cfg.Sustain = 0.0
	cfg.Release = 0.15
	cfg.Drive = 0.6
	cfg.FilterCutoff = 7000
	cfg.FilterResonance = 1.3
	cfg.FadeDuration = 0.012
	cfg.SmoothFrequencyTransitions = true
	return cfg, nil
}

// NewLinnDrumSnare creates a snare drum sound similar to the LinnDrum
func NewLinnDrumSnare(output io.WriteSeeker, duration float64, sampleRate, bitDepth, channels int) (*Settings, error) {
	cfg, err := NewSettings(output, 260.0, 140.0, duration, sampleRate, bitDepth, channels)
	if err != nil {
		return nil, err
	}
	cfg.SoundType = Snare
	cfg.WaveformType = WaveSquare
	cfg.NoiseAmount = 0.5
	cfg.Attack = 0.001
	cfg.Decay = 0.22
	cfg.Sustain = 0.0
	cfg.Release = 0.18
	cfg.Drive = 0.4
	cfg.FilterCutoff = 7500
	cfg.FilterResonance = 1.1
	cfg.FadeDuration = 0.01
	cfg.SmoothFrequencyTransitions = true
	return cfg, nil
}

// NewExperimentalSnare creates an experimental snare drum sound
func NewExperimentalSnare(output io.WriteSeeker, duration float64, sampleRate, bitDepth, channels int) (*Settings, error) {
	cfg, err := NewSettings(output, 400.0, 100.0, duration, sampleRate, bitDepth, channels)
	if err != nil {
		return nil, err
	}
	cfg.SoundType = Snare
	cfg.WaveformType = WaveSawtooth
	cfg.NoiseAmount = 0.8
	cfg.Attack = 0.002
	cfg.Decay = 0.3
	cfg.Sustain = 0.0
	cfg.Release = 0.25
	cfg.Drive = 0.7
	cfg.FilterCutoff = 5000
	cfg.FilterResonance = 1.7
	cfg.FadeDuration = 0.02
	cfg.SmoothFrequencyTransitions = true
	return cfg, nil
}
