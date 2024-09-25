// audioeffects.go
// Package audioeffects provides a collection of audio processing effects for manipulating audio samples.
package audioeffects

import (
	"math"
)

// Biquad represents a biquad filter with its coefficients and state variables.
type Biquad struct {
	a0, a1, a2 float64 // Feedforward coefficients
	b1, b2     float64 // Feedback coefficients
	x1, x2     float64 // Previous input samples
	y1, y2     float64 // Previous output samples
}

// NewBiquad creates a new Biquad filter based on the specified type and parameters.
// filterType can be "low-pass", "high-pass", "band-pass", "notch", "all-pass", etc.
// freq is the center frequency, Q is the quality factor, and sampleRate is the sampling rate in Hz.
func NewBiquad(filterType string, freq, Q, sampleRate float64) *Biquad {
	omega := 2 * math.Pi * freq / sampleRate
	sinOmega := math.Sin(omega)
	cosOmega := math.Cos(omega)
	alpha := sinOmega / (2 * Q)

	var a0, a1, a2, b1, b2 float64

	switch filterType {
	case "low-pass":
		a0 = 1 + alpha
		a1 = -2 * cosOmega
		a2 = 1 - alpha
		b1 = 2 * cosOmega
		b2 = 1 - alpha
	case "high-pass":
		a0 = 1 + alpha
		a1 = -2 * cosOmega
		a2 = 1 - alpha
		b1 = 2 * cosOmega
		b2 = 1 - alpha
	case "band-pass":
		a0 = 1 + alpha
		a1 = -2 * cosOmega
		a2 = 1 - alpha
		b1 = 2 * cosOmega
		b2 = 1 - alpha
	case "notch":
		a0 = 1 + alpha
		a1 = -2 * cosOmega
		a2 = 1 - alpha
		b1 = 2 * cosOmega
		b2 = 1 - alpha
	case "all-pass":
		// All-pass filter coefficients
		a0 = 1 + alpha
		a1 = -2 * cosOmega
		a2 = 1 - alpha
		b1 = 2 * cosOmega
		b2 = -(1 - alpha)
	default:
		// Default to pass-through if filter type is unknown
		a0 = 1
		a1 = 0
		a2 = 0
		b1 = 0
		b2 = 0
	}

	// Normalize coefficients
	return &Biquad{
		a0: a0 / a0, // This will always be 1.0; kept for consistency
		a1: a1 / a0,
		a2: a2 / a0,
		b1: b1 / a0,
		b2: b2 / a0,
		x1: 0,
		x2: 0,
		y1: 0,
		y2: 0,
	}
}

// UpdateParameters updates the filter coefficients based on new parameters.
// This allows dynamic modification of the filter in real-time.
func (b *Biquad) UpdateParameters(filterType string, freq, Q, sampleRate float64) {
	omega := 2 * math.Pi * freq / sampleRate
	sinOmega := math.Sin(omega)
	cosOmega := math.Cos(omega)
	alpha := sinOmega / (2 * Q)

	var a0, a1, a2, b1, b2 float64

	switch filterType {
	case "low-pass":
		a0 = 1 + alpha
		a1 = -2 * cosOmega
		a2 = 1 - alpha
		b1 = 2 * cosOmega
		b2 = 1 - alpha
	case "high-pass":
		a0 = 1 + alpha
		a1 = -2 * cosOmega
		a2 = 1 - alpha
		b1 = 2 * cosOmega
		b2 = 1 - alpha
	case "band-pass":
		a0 = 1 + alpha
		a1 = -2 * cosOmega
		a2 = 1 - alpha
		b1 = 2 * cosOmega
		b2 = 1 - alpha
	case "notch":
		a0 = 1 + alpha
		a1 = -2 * cosOmega
		a2 = 1 - alpha
		b1 = 2 * cosOmega
		b2 = 1 - alpha
	case "all-pass":
		// All-pass filter coefficients
		a0 = 1 + alpha
		a1 = -2 * cosOmega
		a2 = 1 - alpha
		b1 = 2 * cosOmega
		b2 = -(1 - alpha)
	default:
		// Default to pass-through if filter type is unknown
		a0 = 1
		a1 = 0
		a2 = 0
		b1 = 0
		b2 = 0
	}

	// Normalize coefficients
	b.a0 = a0 / a0 // This will always be 1.0; kept for consistency
	b.a1 = a1 / a0
	b.a2 = a2 / a0
	b.b1 = b1 / a0
	b.b2 = b2 / a0
}

// Process applies the biquad filter to a single sample.
func (b *Biquad) Process(sample float64) float64 {
	y := b.a0*sample + b.a1*b.x1 + b.a2*b.x2 - b.b1*b.y1 - b.b2*b.y2
	b.x2 = b.x1
	b.x1 = sample
	b.y2 = b.y1
	b.y1 = y
	return y
}

// CompressorSettings defines the settings for a compressor.
type CompressorSettings struct {
	Threshold float64 // Threshold level
	Ratio     float64 // Compression ratio
	Attack    float64 // Attack time in seconds
	Release   float64 // Release time in seconds
}

// FadeIn applies a fade-in effect to the samples.
// duration: Duration of the fade-in in seconds.
// sampleRate: Sampling rate in Hz.
func FadeIn(samples []float64, duration float64, sampleRate int) []float64 {
	faded := make([]float64, len(samples))
	fadeSamples := int(duration * float64(sampleRate))
	if fadeSamples > len(samples) {
		fadeSamples = len(samples)
	}
	for i := 0; i < len(samples); i++ {
		if i < fadeSamples {
			faded[i] = samples[i] * float64(i) / float64(fadeSamples)
		} else {
			faded[i] = samples[i]
		}
	}
	return faded
}

// FadeOut applies a fade-out effect to the samples.
// duration: Duration of the fade-out in seconds.
// sampleRate: Sampling rate in Hz.
func FadeOut(samples []float64, duration float64, sampleRate int) []float64 {
	faded := make([]float64, len(samples))
	fadeSamples := int(duration * float64(sampleRate))
	startFade := len(samples) - fadeSamples
	if startFade < 0 {
		startFade = 0
	}
	for i := 0; i < len(samples); i++ {
		if i >= startFade {
			faded[i] = samples[i] * float64(len(samples)-i) / float64(fadeSamples)
		} else {
			faded[i] = samples[i]
		}
	}
	return faded
}

// LowPassFilter applies a low-pass filter to the samples.
// cutoffFreq: Cutoff frequency in Hz.
// sampleRate: Sampling rate in Hz.
func LowPassFilter(samples []float64, cutoffFreq float64, sampleRate int) []float64 {
	filter := NewBiquad("low-pass", cutoffFreq, 0.707, float64(sampleRate))
	filtered := make([]float64, len(samples))
	for i, sample := range samples {
		filtered[i] = filter.Process(sample)
	}
	return filtered
}

// HighPassFilter applies a high-pass filter to the samples.
// cutoffFreq: Cutoff frequency in Hz.
// sampleRate: Sampling rate in Hz.
func HighPassFilter(samples []float64, cutoffFreq float64, sampleRate int) []float64 {
	filter := NewBiquad("high-pass", cutoffFreq, 0.707, float64(sampleRate))
	filtered := make([]float64, len(samples))
	for i, sample := range samples {
		filtered[i] = filter.Process(sample)
	}
	return filtered
}

// BandPassFilter applies a band-pass filter to the samples.
// lowFreq: Lower cutoff frequency in Hz.
// highFreq: Upper cutoff frequency in Hz.
// sampleRate: Sampling rate in Hz.
func BandPassFilter(samples []float64, lowFreq, highFreq float64, sampleRate int) []float64 {
	// Calculate center frequency and Q factor
	centerFreq := (lowFreq + highFreq) / 2
	BandWidth := highFreq - lowFreq
	Q := centerFreq / BandWidth
	filter := NewBiquad("band-pass", centerFreq, Q, float64(sampleRate))
	filtered := make([]float64, len(samples))
	for i, sample := range samples {
		filtered[i] = filter.Process(sample)
	}
	return filtered
}

// NoiseGate reduces the volume of audio signals that fall below a certain threshold.
// threshold: Threshold level below which audio is attenuated.
// attack: Attack time in seconds.
// release: Release time in seconds.
// sampleRate: Sampling rate in Hz.
func NoiseGate(samples []float64, threshold, attack, release float64, sampleRate int) []float64 {
	gated := make([]float64, len(samples))
	gain := 1.0

	attackCoeff := math.Exp(-1.0 / (attack * float64(sampleRate)))
	releaseCoeff := math.Exp(-1.0 / (release * float64(sampleRate)))

	for i, sample := range samples {
		absSample := math.Abs(sample)

		if absSample > threshold {
			gain = attackCoeff*gain + (1.0-attackCoeff)*1.0
		} else {
			gain = releaseCoeff*gain + (1.0-releaseCoeff)*0.0
		}

		gated[i] = sample * gain
	}

	return gated
}

// StereoDelay implements separate delay effects on the left and right channels.
// delayTimeLeft: Delay time for the left channel in seconds.
// delayTimeRight: Delay time for the right channel in seconds.
// feedback: Feedback factor (0.0 to less than 1.0).
// mix: Mixing proportion of the delayed signal (0.0 to 1.0).
func StereoDelay(left, right []float64, sampleRate int, delayTimeLeft, delayTimeRight float64, feedback, mix float64) ([]float64, []float64) {
	delaySamplesLeft := int(delayTimeLeft * float64(sampleRate))
	delaySamplesRight := int(delayTimeRight * float64(sampleRate))
	if delaySamplesLeft <= 0 {
		delaySamplesLeft = 1
	}
	if delaySamplesRight <= 0 {
		delaySamplesRight = 1
	}
	delayedLeft := make([]float64, len(left))
	delayedRight := make([]float64, len(right))
	bufferLeft := make([]float64, delaySamplesLeft)
	bufferRight := make([]float64, delaySamplesRight)
	bufferIndexLeft := 0
	bufferIndexRight := 0

	for i := 0; i < len(left); i++ {
		// Left channel
		delayedSampleLeft := bufferLeft[bufferIndexLeft]
		delayedLeft[i] = left[i]*(1-mix) + delayedSampleLeft*mix
		bufferLeft[bufferIndexLeft] = left[i] + delayedSampleLeft*feedback
		bufferIndexLeft = (bufferIndexLeft + 1) % delaySamplesLeft

		// Right channel
		delayedSampleRight := bufferRight[bufferIndexRight]
		delayedRight[i] = right[i]*(1-mix) + delayedSampleRight*mix
		bufferRight[bufferIndexRight] = right[i] + delayedSampleRight*feedback
		bufferIndexRight = (bufferIndexRight + 1) % delaySamplesRight
	}

	return delayedLeft, delayedRight
}

// Expander increases the dynamic range by amplifying quieter signals and attenuating louder ones.
// threshold: Threshold level below which expansion is applied.
// ratio: Expansion ratio.
// attack: Attack time in seconds.
// release: Release time in seconds.
// sampleRate: Sampling rate in Hz.
func Expander(samples []float64, threshold, ratio, attack, release float64, sampleRate int) []float64 {
	expanded := make([]float64, len(samples))
	envelope := 0.0

	attackCoeff := math.Exp(-1.0 / (attack * float64(sampleRate)))
	releaseCoeff := math.Exp(-1.0 / (release * float64(sampleRate)))

	for i, sample := range samples {
		absSample := math.Abs(sample)

		if absSample > envelope {
			envelope = attackCoeff*envelope + (1.0-attackCoeff)*absSample
		} else {
			envelope = releaseCoeff*envelope + (1.0-releaseCoeff)*absSample
		}

		var gain float64
		if envelope < threshold {
			gain = envelope * ratio / threshold
			if gain > 1.0 {
				gain = 1.0
			}
		} else {
			gain = 1.0
		}

		expanded[i] = sample * gain
	}

	return expanded
}

// SoftClip applies a soft clipping distortion to a single sample.
// drive controls the amount of distortion.
func SoftClip(sample, drive float64) float64 {
	return (3 + drive) * sample / (1 + drive*math.Abs(sample))
}

// SoftClippingDistortion applies soft clipping distortion to the samples.
// drive controls the amount of distortion.
func SoftClippingDistortion(samples []float64, drive float64) []float64 {
	distorted := make([]float64, len(samples))
	for i, sample := range samples {
		distorted[i] = SoftClip(sample, drive)
		// Clamp to [-1, 1]
		if distorted[i] > 1.0 {
			distorted[i] = 1.0
		} else if distorted[i] < -1.0 {
			distorted[i] = -1.0
		}
	}
	return distorted
}

// SidechainCompressor applies compression to the target signal based on the trigger signal.
// target: The primary audio signal to be compressed.
// trigger: The secondary audio signal used to control the compression.
// threshold: Threshold level above which compression is applied.
// ratio: Compression ratio.
// attack: Attack time in seconds.
// release: Release time in seconds.
// sampleRate: Sampling rate in Hz.
func SidechainCompressor(target, trigger []float64, threshold, ratio, attack, release float64, sampleRate int) []float64 {
	if len(target) != len(trigger) {
		// Handle error: target and trigger must be the same length
		return target
	}

	compressed := make([]float64, len(target))
	gain := 1.0

	attackCoeff := math.Exp(-1.0 / (attack * float64(sampleRate)))
	releaseCoeff := math.Exp(-1.0 / (release * float64(sampleRate)))

	for i := 0; i < len(target); i++ {
		absTrigger := math.Abs(trigger[i])

		// Envelope follower
		if absTrigger > gain {
			gain = attackCoeff*gain + (1.0-attackCoeff)*absTrigger
		} else {
			gain = releaseCoeff*gain + (1.0-releaseCoeff)*absTrigger
		}

		// Gain computation based on threshold and ratio
		if gain > threshold {
			gain = threshold + (gain-threshold)/ratio
			gain /= gain // Normalize gain to prevent amplification
			if gain < 0.0 {
				gain = 0.0
			}
		} else {
			gain = 1.0
		}

		compressed[i] = target[i] * gain
	}

	return compressed
}

// Compressor applies dynamic range compression to the samples.
// threshold: Threshold level above which compression is applied.
// ratio: Compression ratio.
// attack: Attack time in seconds.
// release: Release time in seconds.
// sampleRate: Sampling rate in Hz.
func Compressor(samples []float64, threshold, ratio, attack, release float64, sampleRate int) []float64 {
	compressed := make([]float64, len(samples))
	gain := 1.0

	attackCoeff := math.Exp(-1.0 / (attack * float64(sampleRate)))
	releaseCoeff := math.Exp(-1.0 / (release * float64(sampleRate)))

	for i, sample := range samples {
		absSample := math.Abs(sample)

		if absSample > threshold {
			gain = attackCoeff*gain + (1.0-attackCoeff)*(threshold/absSample)
		} else {
			gain = releaseCoeff*gain + (1.0-releaseCoeff)*1.0
		}

		compressed[i] = sample * gain
	}

	return compressed
}

// Envelope applies an ADSR envelope to the samples.
// attack: Attack time in seconds.
// decay: Decay time in seconds.
// sustainLevel: Sustain level (0.0 to 1.0).
// release: Release time in seconds.
// sampleRate: Sampling rate in Hz.
func Envelope(samples []float64, attack, decay, sustainLevel, release float64, sampleRate int) []float64 {
	enveloped := make([]float64, len(samples))
	state := "attack"
	currentLevel := 0.0
	attackSamples := int(attack * float64(sampleRate))
	decaySamples := int(decay * float64(sampleRate))
	releaseSamples := int(release * float64(sampleRate))
	sustainSamples := len(samples) - attackSamples - decaySamples - releaseSamples
	if sustainSamples < 0 {
		sustainSamples = 0
		releaseSamples = len(samples) - attackSamples - decaySamples
		if releaseSamples < 0 {
			releaseSamples = 0
		}
	}

	for i := 0; i < len(samples); i++ {
		switch state {
		case "attack":
			if attackSamples > 0 {
				currentLevel = float64(i) / float64(attackSamples)
				if i >= attackSamples {
					currentLevel = 1.0
					state = "decay"
				}
			} else {
				currentLevel = 1.0
				state = "decay"
			}
		case "decay":
			if decaySamples > 0 {
				currentLevel = 1.0 - ((1.0 - sustainLevel) * float64(i-attackSamples) / float64(decaySamples))
				if i >= attackSamples+decaySamples {
					currentLevel = sustainLevel
					state = "sustain"
				}
			} else {
				currentLevel = sustainLevel
				state = "sustain"
			}
		case "sustain":
			currentLevel = sustainLevel
			if i >= attackSamples+decaySamples+sustainSamples {
				state = "release"
			}
		case "release":
			if releaseSamples > 0 {
				currentLevel = sustainLevel * (1.0 - float64(i-attackSamples-decaySamples-sustainSamples)/float64(releaseSamples))
				if i >= attackSamples+decaySamples+sustainSamples+releaseSamples {
					currentLevel = 0.0
				}
			} else {
				currentLevel = 0.0
			}
		}

		enveloped[i] = samples[i] * currentLevel
	}

	return enveloped
}

// Panning adjusts the stereo balance of the samples.
// pan: Panning value where -1.0 is full left, 0.0 is center, and 1.0 is full right.
// Returns left and right channel samples.
func Panning(samples []float64, pan float64) ([]float64, []float64) {
	left := make([]float64, len(samples))
	right := make([]float64, len(samples))
	// Clamp pan to [-1, 1]
	if pan < -1.0 {
		pan = -1.0
	} else if pan > 1.0 {
		pan = 1.0
	}
	// Calculate pan angles
	theta := (pan + 1.0) * (math.Pi / 4) // Map pan from [-1,1] to [0, pi/2]
	leftCoeff := math.Cos(theta)
	rightCoeff := math.Sin(theta)
	for i, sample := range samples {
		left[i] = sample * leftCoeff
		right[i] = sample * rightCoeff
	}
	return left, right
}

// Tremolo applies an amplitude modulation effect to the samples.
// rate: Modulation rate in Hz.
// depth: Modulation depth (0.0 to 1.0).
func Tremolo(samples []float64, sampleRate int, rate, depth float64) []float64 {
	modulated := make([]float64, len(samples))
	lfoPhase := 0.0
	lfoIncrement := rate / float64(sampleRate)
	for i, sample := range samples {
		mod := 1.0 - depth + depth*math.Sin(2*math.Pi*lfoPhase)
		modulated[i] = sample * mod
		lfoPhase += lfoIncrement
		if lfoPhase >= 1.0 {
			lfoPhase -= 1.0
		}
	}
	return modulated
}

// Flanger applies a flanger effect to the samples.
// baseDelay: Base delay time in seconds.
// modDepth: Modulation depth in seconds.
// modRate: Modulation rate in Hz.
// feedback: Feedback factor (0.0 to less than 1.0).
// mix: Mixing proportion of the delayed signal (0.0 to 1.0).
func Flanger(samples []float64, sampleRate int, baseDelay, modDepth, modRate, feedback, mix float64) []float64 {
	flanged := make([]float64, len(samples))
	bufferSize := int((baseDelay+modDepth)*float64(sampleRate)) + 2
	buffer := make([]float64, bufferSize)
	bufferIndex := 0
	lfoPhase := 0.0
	lfoIncrement := modRate / float64(sampleRate)

	for i, sample := range samples {
		// Calculate current delay in samples
		lfoValue := math.Sin(2 * math.Pi * lfoPhase)
		currentDelay := baseDelay + modDepth*lfoValue
		delaySamples := int(currentDelay * float64(sampleRate))
		readIndex := (bufferIndex - delaySamples + bufferSize) % bufferSize

		// Get delayed sample
		delayed := buffer[readIndex]

		// Apply feedback
		buffer[bufferIndex] = sample + delayed*feedback

		// Mix dry and wet signals
		flanged[i] = sample*(1.0-mix) + delayed*mix

		// Increment indices and phase
		bufferIndex = (bufferIndex + 1) % bufferSize
		lfoPhase += lfoIncrement
		if lfoPhase >= 1.0 {
			lfoPhase -= 1.0
		}
	}

	return flanged
}

// Phaser applies a phaser effect to the samples.
// rate: Modulation rate in Hz.
// depth: Modulation depth.
// feedback: Feedback factor (0.0 to less than 1.0).
// sampleRate: Sampling rate in Hz.
func Phaser(samples []float64, sampleRate int, rate, depth, feedback float64) []float64 {
	phased := make([]float64, len(samples))
	phaseIncrement := rate / float64(sampleRate)
	phase := 0.0

	// Initialize two all-pass filters for a basic phaser
	// Note: A proper phaser would require multiple all-pass filters with phase shifts
	allPass1 := NewBiquad("all-pass", 1000.0, 0.7, float64(sampleRate))
	allPass2 := NewBiquad("all-pass", 1500.0, 0.7, float64(sampleRate))

	for i, sample := range samples {
		// Sweep the center frequency with LFO
		sweep := math.Sin(2 * math.Pi * phase)
		centerFreq := 1000.0 + depth*1000.0*sweep // Example frequency sweep from 0 to 2000 Hz
		if centerFreq < 20.0 {                    // Prevent frequencies below human hearing
			centerFreq = 20.0
		} else if centerFreq > float64(sampleRate)/2 {
			centerFreq = float64(sampleRate) / 2
		}

		// Update filter parameters dynamically
		allPass1.UpdateParameters("all-pass", centerFreq, 0.7, float64(sampleRate))
		allPass2.UpdateParameters("all-pass", centerFreq, 0.7, float64(sampleRate))

		// Handle feedback sample
		var feedbackSample float64
		if i > 0 {
			feedbackSample = feedback * phased[i-1]
		} else {
			feedbackSample = 0.0
		}

		// Apply all-pass filters with feedback
		out1 := allPass1.Process(sample + feedbackSample)
		out2 := allPass2.Process(out1 + feedbackSample)
		phased[i] = out2

		// Increment phase
		phase += phaseIncrement
		if phase >= 1.0 {
			phase -= 1.0
		}
	}

	return phased
}

// RingModulation applies ring modulation to the samples.
// carrierFreq: Frequency of the carrier oscillator in Hz.
// sampleRate: Sampling rate in Hz.
func RingModulation(samples []float64, carrierFreq float64, sampleRate int) []float64 {
	modulated := make([]float64, len(samples))
	carrierPhase := 0.0
	carrierIncrement := carrierFreq / float64(sampleRate)
	for i, sample := range samples {
		carrier := math.Sin(2 * math.Pi * carrierPhase)
		modulated[i] = sample * carrier
		carrierPhase += carrierIncrement
		if carrierPhase >= 1.0 {
			carrierPhase -= 1.0
		}
	}
	return modulated
}

// WahWah applies a wah-wah effect to the samples.
// baseFreq: Base center frequency of the band-pass filter in Hz.
// sweepFreq: Frequency of the sweep (LFO) in Hz.
// Q: Quality factor of the band-pass filter.
func WahWah(samples []float64, sampleRate int, baseFreq, sweepFreq, Q float64) []float64 {
	phased := make([]float64, len(samples))
	filter := NewBiquad("band-pass", baseFreq, Q, float64(sampleRate))
	lfoPhase := 0.0
	lfoIncrement := sweepFreq / float64(sampleRate)

	for i, sample := range samples {
		// Sweep the center frequency with LFO
		sweep := math.Sin(2 * math.Pi * lfoPhase)
		centerFreq := baseFreq + 500.0*sweep // Example sweep range: ±500Hz
		if centerFreq < 20.0 {
			centerFreq = 20.0
		} else if centerFreq > float64(sampleRate)/2 {
			centerFreq = float64(sampleRate) / 2
		}
		// Update filter parameters
		filter.UpdateParameters("band-pass", centerFreq, Q, float64(sampleRate))
		// Apply filter
		phased[i] = filter.Process(sample)
		// Increment LFO phase
		lfoPhase += lfoIncrement
		if lfoPhase >= 1.0 {
			lfoPhase -= 1.0
		}
	}

	return phased
}

// StereoWidening enhances the stereo image by adjusting the amplitudes of left and right channels.
// left: Left channel samples.
// right: Right channel samples.
// width: Width factor (0.0 to 1.0), where 0.0 is original stereo and 1.0 is maximum widening.
func StereoWidening(left, right []float64, width float64) ([]float64, []float64) {
	widenedLeft := make([]float64, len(left))
	widenedRight := make([]float64, len(right))
	for i := 0; i < len(left); i++ {
		widenedLeft[i] = left[i] * (1.0 + width)
		widenedRight[i] = right[i] * (1.0 + width)
		// Clamp to [-1, 1]
		if widenedLeft[i] > 1.0 {
			widenedLeft[i] = 1.0
		} else if widenedLeft[i] < -1.0 {
			widenedLeft[i] = -1.0
		}
		if widenedRight[i] > 1.0 {
			widenedRight[i] = 1.0
		} else if widenedRight[i] < -1.0 {
			widenedRight[i] = -1.0
		}
	}
	return widenedLeft, widenedRight
}

// MultibandCompression applies compression independently across different frequency bands.
// bands: Slice of frequency bands, each defined by Low and High cutoff frequencies.
// compressors: Slice of CompressorSettings structs corresponding to each band.
// sampleRate: Sampling rate in Hz.
func MultibandCompression(samples []float64, bands []struct {
	Low  float64
	High float64
}, compressors []CompressorSettings, sampleRate int) []float64 {
	if len(bands) != len(compressors) {
		// Handle error: bands and compressors must have the same length
		return samples
	}

	// Split into bands
	splitBands := make([][]float64, len(bands))
	for i, band := range bands {
		splitBands[i] = BandPassFilter(samples, band.Low, band.High, sampleRate)
	}

	// Compress each band
	for i := range splitBands {
		settings := compressors[i]
		splitBands[i] = Compressor(splitBands[i], settings.Threshold, settings.Ratio, settings.Attack, settings.Release, sampleRate)
	}

	// Recombine bands
	recombined := make([]float64, len(samples))
	for _, band := range splitBands {
		for i := range recombined {
			recombined[i] += band[i]
		}
	}

	// Normalize to prevent clipping
	max := 0.0
	for _, sample := range recombined {
		if math.Abs(sample) > max {
			max = math.Abs(sample)
		}
	}
	if max > 1.0 {
		for i := range recombined {
			recombined[i] /= max
		}
	}

	return recombined
}

// PitchShift shifts the pitch of the samples by the specified number of semitones.
// Note: This implementation uses simple resampling and changes the duration of the audio.
// For high-quality pitch shifting without altering duration, consider implementing advanced algorithms like Phase Vocoder or WSOLA.
func PitchShift(samples []float64, semitones float64) []float64 {
	rate := math.Pow(2, semitones/12)
	newLength := int(float64(len(samples)) / rate)
	shifted := make([]float64, newLength)
	for i := 0; i < newLength; i++ {
		srcIndex := float64(i) * rate
		lower := int(math.Floor(srcIndex))
		upper := lower + 1
		if upper >= len(samples) {
			upper = len(samples) - 1
		}
		frac := srcIndex - float64(lower)
		shifted[i] = samples[lower]*(1-frac) + samples[upper]*frac
	}
	return shifted
}

// FrequencyModulation applies frequency modulation to the samples.
// It generates an FM signal based on the input samples acting as the modulator.
// carrierFreq is the base frequency of the carrier oscillator.
// modDepth controls the intensity of the frequency modulation.
func FrequencyModulation(samples []float64, carrierFreq, modDepth, sampleRate float64) []float64 {
	modulated := make([]float64, len(samples))
	carrierPhase := 0.0

	for i, sample := range samples {
		// Calculate instantaneous frequency
		instantFreq := carrierFreq + modDepth*sample

		// Update phase
		carrierPhase += (instantFreq / sampleRate) * 2 * math.Pi

		// Wrap phase to [0, 2π] to prevent numerical issues
		if carrierPhase > 2*math.Pi {
			carrierPhase -= 2 * math.Pi
		} else if carrierPhase < 0 {
			carrierPhase += 2 * math.Pi
		}

		// Generate FM signal
		modulated[i] = math.Sin(carrierPhase)
	}

	return modulated
}

// PitchModulation applies pitch modulation (vibrato) to the samples.
// modFreq is the frequency of the LFO in Hz.
// modDepth is the maximum pitch deviation in seconds.
// sampleRate is the sampling rate in Hz.
func PitchModulation(samples []float64, modFreq, modDepth float64, sampleRate int) []float64 {
	modulated := make([]float64, len(samples))
	lfoPhase := 0.0
	lfoIncrement := modFreq / float64(sampleRate)
	maxDelay := int(modDepth * float64(sampleRate)) // Convert depth to samples

	// Create a circular buffer for delay
	buffer := make([]float64, maxDelay*2)
	bufferIndex := 0

	for i := 0; i < len(samples); i++ {
		// Calculate LFO value
		lfoValue := math.Sin(2 * math.Pi * lfoPhase)
		// Calculate current delay in samples
		currentDelay := int(math.Abs(lfoValue) * float64(maxDelay))
		readIndex := (bufferIndex - currentDelay + len(buffer)) % len(buffer)

		// Read delayed sample
		modulated[i] = buffer[readIndex]

		// Write current sample to buffer
		buffer[bufferIndex] = samples[i]
		bufferIndex = (bufferIndex + 1) % len(buffer)

		// Increment LFO phase
		lfoPhase += lfoIncrement
		if lfoPhase >= 1.0 {
			lfoPhase -= 1.0
		}
	}

	return modulated
}

// Reverb applies a more realistic reverb effect by using multiple parallel delay lines.
// delayTimes: Slice of delay times in seconds for each parallel delay line.
// decays: Slice of decay factors corresponding to each delay line.
// Ensure that len(delayTimes) == len(decays).
func Reverb(samples []float64, sampleRate int, delayTimes []float64, decays []float64) []float64 {
	if len(delayTimes) != len(decays) {
		// Mismatch in the number of delay lines and decays
		return samples
	}

	reverb := make([]float64, len(samples))
	numDelays := len(delayTimes)
	buffers := make([][]float64, numDelays)
	bufferIndices := make([]int, numDelays)

	for i := 0; i < numDelays; i++ {
		delaySamples := int(delayTimes[i] * float64(sampleRate))
		if delaySamples <= 0 {
			delaySamples = 1
		}
		buffers[i] = make([]float64, delaySamples)
		bufferIndices[i] = 0
	}

	for i, sample := range samples {
		var wet float64 = 0.0
		for j := 0; j < numDelays; j++ {
			delayed := buffers[j][bufferIndices[j]]
			wet += delayed * decays[j]
			// Update buffer with current sample plus feedback
			buffers[j][bufferIndices[j]] = sample + delayed*decays[j]
			// Increment buffer index
			bufferIndices[j] = (bufferIndices[j] + 1) % len(buffers[j])
		}
		// Mix dry and wet signals
		reverb[i] = sample + wet
		// Optional: Clamp to prevent clipping
		if reverb[i] > 1.0 {
			reverb[i] = 1.0
		} else if reverb[i] < -1.0 {
			reverb[i] = -1.0
		}
	}

	return reverb
}

// Chorus applies a chorus effect to the samples.
// sampleRate: Sampling rate in Hz.
// delay: Delay time in seconds.
// depth: Modulation depth in seconds.
// rate: Modulation rate in Hz.
// mix: Mixing proportion of the delayed signal (0.0 to 1.0).
func Chorus(samples []float64, sampleRate int, delay, depth, rate, mix float64) []float64 {
	chorused := make([]float64, len(samples))
	bufferSize := int((delay+depth)*float64(sampleRate)) + 2
	buffer := make([]float64, bufferSize)
	bufferIndex := 0
	lfoPhase := 0.0
	lfoIncrement := rate / float64(sampleRate)

	for i, sample := range samples {
		// Calculate current delay
		lfoValue := math.Sin(2 * math.Pi * lfoPhase)
		currentDelay := delay + depth*lfoValue
		delaySamples := int(currentDelay * float64(sampleRate))
		readIndex := (bufferIndex - delaySamples + bufferSize) % bufferSize

		// Get delayed sample
		delayed := buffer[readIndex]

		// Write current sample to buffer
		buffer[bufferIndex] = sample + delayed*0.5 // Feedback factor (example)

		// Mix dry and wet signals
		chorused[i] = sample*(1.0-mix) + delayed*mix

		// Increment indices and phase
		bufferIndex = (bufferIndex + 1) % bufferSize
		lfoPhase += lfoIncrement
		if lfoPhase >= 1.0 {
			lfoPhase -= 1.0
		}
	}

	return chorused
}

// Bitcrusher reduces the bit depth and/or sample rate of the samples to create a lo-fi effect.
// bitDepth: Number of bits per sample (e.g., 8 for 8-bit).
// sampleRateReduction: Factor by which to reduce the sample rate (e.g., 2 to halve the sample rate).
func Bitcrusher(samples []float64, bitDepth int, sampleRateReduction int) []float64 {
	if bitDepth < 1 {
		bitDepth = 1
	}
	if bitDepth > 16 {
		bitDepth = 16
	}
	if sampleRateReduction < 1 {
		sampleRateReduction = 1
	}

	step := 1.0 / math.Pow(2, float64(bitDepth))
	bitcrushed := make([]float64, len(samples))

	for i, sample := range samples {
		// Quantize the sample
		bitcrushed[i] = math.Round(sample/step) * step

		// Reduce sample rate by keeping every nth sample
		if sampleRateReduction > 1 && i%sampleRateReduction != 0 {
			bitcrushed[i] = 0.0
		}

		// Clamp to [-1, 1]
		if bitcrushed[i] > 1.0 {
			bitcrushed[i] = 1.0
		} else if bitcrushed[i] < -1.0 {
			bitcrushed[i] = -1.0
		}
	}

	return bitcrushed
}

// Drive applies a distortion effect to a single sample.
// drive controls the intensity of the distortion.
func Drive(sample, drive float64) float64 {
	if drive > 0 {
		return sample * (1 + drive) / (1 + drive*math.Abs(sample))
	}
	return sample
}

// Limiter ensures the signal doesn't exceed the range [-1.0, 1.0].
func Limiter(samples []float64) []float64 {
	limited := make([]float64, len(samples))
	for i, sample := range samples {
		if sample > 1.0 {
			limited[i] = 1.0
		} else if sample < -1.0 {
			limited[i] = -1.0
		} else {
			limited[i] = sample
		}
	}
	return limited
}

// NormalizeSamples scales the samples so that the peak amplitude matches the given target peak.
func NormalizeSamples(samples []float64, targetPeak float64) []float64 {
	if targetPeak <= 0 {
		// Invalid target peak, return samples unmodified
		return samples
	}

	currentPeak := 0.0
	for _, sample := range samples {
		abs := math.Abs(sample)
		if abs > currentPeak {
			currentPeak = abs
		}
	}

	if currentPeak == 0 {
		// All samples are zero, return unmodified
		return samples
	}

	scale := targetPeak / currentPeak
	normalizedSamples := make([]float64, len(samples))
	for i, sample := range samples {
		normalizedSamples[i] = sample * scale
		// Clamp the values to the range [-1.0, 1.0] after scaling
		if normalizedSamples[i] > 1.0 {
			normalizedSamples[i] = 1.0
		} else if normalizedSamples[i] < -1.0 {
			normalizedSamples[i] = -1.0
		}
	}
	return normalizedSamples
}
