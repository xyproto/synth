package synth

import (
	"math"
)

// ApplyEnvelope applies an ADSR envelope to the waveform
func ApplyEnvelope(samples []float64, attack, decay, sustain, release float64, sampleRate int) []float64 {
	adsr := make([]float64, len(samples))
	totalDuration := float64(len(samples)) / float64(sampleRate)
	for i := range samples {
		t := float64(i) / float64(sampleRate)
		var envelope float64
		if t < attack {
			envelope = t / attack
		} else if t < attack+decay {
			envelope = 1 - (t-attack)/decay*(1-sustain)
		} else if t < totalDuration-release {
			envelope = sustain
		} else if t < totalDuration {
			envelope = sustain * (1 - (t-(totalDuration-release))/release)
		} else {
			envelope = 0
		}
		adsr[i] = samples[i] * envelope
	}
	return adsr
}

// ApplyEnvelopeAtTime generates the ADSR envelope at a specific time
func (cfg *Settings) ApplyEnvelopeAtTime(t float64) float64 {
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

func ApplyPitchModulation(samples []float64, modFreq, modDepth float64, sampleRate int) []float64 {
	modulated := make([]float64, len(samples))
	for i := 0; i < len(samples); i++ {
		t := float64(i) / float64(sampleRate)
		mod := math.Sin(2*math.Pi*modFreq*t) * modDepth
		value := samples[i] * math.Pow(2, mod)
		// Clamp the value
		if value > 1.0 {
			value = 1.0
		} else if value < -1.0 {
			value = -1.0
		}
		modulated[i] = value
	}
	return modulated
}

// ApplyPanning applies stereo panning to the samples. pan should be in the range [-1, 1], where -1 is full left and 1 is full right.
func ApplyPanning(samples []float64, pan float64) ([]float64, []float64) {
	leftChannel := make([]float64, len(samples))
	rightChannel := make([]float64, len(samples))
	leftGain := math.Cos((pan + 1) * math.Pi / 4) // Equal-power panning
	rightGain := math.Sin((pan + 1) * math.Pi / 4)

	for i := range samples {
		leftChannel[i] = samples[i] * leftGain
		rightChannel[i] = samples[i] * rightGain
	}

	return leftChannel, rightChannel
}

// ApplyFrequencyModulation applies frequency modulation to a waveform using a modulator frequency and depth
func ApplyFrequencyModulation(samples []float64, modFreq, modDepth float64, sampleRate int) []float64 {
	modulated := make([]float64, len(samples))
	carrierPhase := 0.0
	modulatorPhase := 0.0
	for i := range samples {
		carrierFreq := modDepth * math.Sin(2*math.Pi*modFreq*modulatorPhase)
		carrierPhase += carrierFreq / float64(sampleRate)
		modulated[i] = math.Sin(2*math.Pi*carrierPhase) * samples[i]
		modulatorPhase += 1.0 / float64(sampleRate)
	}
	return modulated
}

// ApplyFadeIn applies a fade-in to the start of the samples using the specified fade curve
func ApplyFadeIn(samples []float64, fadeDuration float64, sampleRate int, curve FadeCurve) []float64 {
	fadeSamples := int(fadeDuration * float64(sampleRate))
	if fadeSamples > len(samples) {
		fadeSamples = len(samples)
	}
	for i := 0; i < fadeSamples; i++ {
		t := float64(i) / float64(fadeSamples)
		multiplier := curve(t)
		samples[i] *= multiplier
	}
	return samples
}

// ApplyFadeOut applies a fade-out to the end of the samples using the specified fade curve
func ApplyFadeOut(samples []float64, fadeDuration float64, sampleRate int, curve FadeCurve) []float64 {
	totalSamples := len(samples)
	fadeSamples := int(fadeDuration * float64(sampleRate))
	if fadeSamples > totalSamples {
		fadeSamples = totalSamples
	}
	for i := 0; i < fadeSamples; i++ {
		t := float64(i) / float64(fadeSamples)
		multiplier := curve(1.0 - t)
		index := totalSamples - fadeSamples + i
		if index >= 0 && index < totalSamples {
			samples[index] *= multiplier
		}
	}
	return samples
}

// ApplyChorus applies a chorus effect to the samples
func ApplyChorus(samples []float64, sampleRate int, delaySec float64, depth float64, rate float64) []float64 {
	delaySamples := int(delaySec * float64(sampleRate))
	modulated := make([]float64, len(samples))
	for i := 0; i < len(samples); i++ {
		t := float64(i) / float64(sampleRate)
		modulation := depth * math.Sin(2*math.Pi*rate*t)
		delay := int(float64(delaySamples) * (1 + modulation))
		index := i - delay
		if index >= 0 && index < len(samples) {
			modulated[i] = (samples[i] + samples[index]) / 2
		} else {
			modulated[i] = samples[i]
		}
	}
	return modulated
}
