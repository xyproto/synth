# Audio Effects

Apply audio effects to samples.

Note that this package is a bit experimental and a work in progress.

Audio samples are passed in as `[]float64`. The sample rate is typically 44100 or 48000.

### Example use

Write a .wav file that has a very simple clap-like sound, with some reverb applied:

```go
package main

import (
    "log"
    "math"
    "math/rand"
    "os"

    "github.com/go-audio/audio"
    "github.com/go-audio/wav"
    "github.com/xyproto/audioeffects"
)

func writeWav(filename string, samples []float64, sampleRate int, numChannels int) error {
    var (
        ch, intSample int
        intSamples    = make([]int, len(samples)*numChannels)
    )
    for i, sample := range samples {
        if sample > 1.0 {
            sample = 1.0
        } else if sample < -1.0 {
            sample = -1.0
        }
        intSample = int(math.Round(sample * 32767))
        for ch = 0; ch < numChannels; ch++ {
            intSamples[i*numChannels+ch] = intSample
        }
    }
    buf := &audio.IntBuffer{
        Data:           intSamples,
        Format:         &audio.Format{SampleRate: sampleRate, NumChannels: numChannels},
        SourceBitDepth: 16,
    }
    file, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer file.Close()
    enc := wav.NewEncoder(file, sampleRate, 16, numChannels, 1)
    if err := enc.Write(buf); err != nil {
        return err
    }
    return enc.Close()
}

func main() {
    sampleRate := 48000
    duration := 3.0

    numSamples := int(duration * float64(sampleRate))
    samples := make([]float64, numSamples)

    // Create white noise burst with quick decay
    for i := 0; i < numSamples; i++ {
        if i < int(float64(sampleRate)*0.05) { // 50 ms of noise
            samples[i] = (rand.Float64()*2 - 1) * (1.0 - float64(i)/(float64(sampleRate)*0.05))
        } else {
            samples[i] = 0.0
        }
    }

    // Apply Reverb with configurable mix
    delayTimes := []float64{0.2, 0.4, 0.6, 0.8} // in seconds
    decays := []float64{0.6, 0.4, 0.2, 0.1}     // decay factors
    mix := 0.3                                  // 30% wet signal, 70% dry signal

    reverbed := audioeffects.Reverb(samples, sampleRate, delayTimes, decays, mix)

    // Normalize to prevent clipping
    normalized := audioeffects.NormalizeSamples(reverbed, 0.8)

    const filename = "clap_reverb.wav"

    err := writeWav(filename, normalized, sampleRate, 1)
    if err != nil {
        log.Fatalf("Failed to write WAV to %s: %v", filename, err)
    }
    log.Printf("Successfully wrote %s\n", filename)
}
```

### General info

* Version: 0.9.1
* License: MIT
