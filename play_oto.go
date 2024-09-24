//go:build !sdl2

package synth

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/ebitengine/oto/v3"
)

type AudioDeviceKey int

// Player represents a struct with methods for playing audio
type Player struct {
	Initialized         bool
	PlayingAudioDevices map[AudioDeviceKey]*oto.Player
	ctx                 *oto.Context
}

var mut sync.Mutex

var theContext *oto.Context

// NewPlayer initializes Oto Audio with system defaults and returns a Player struct
func NewPlayer() *Player {
	if theContext == nil {
		// Create Oto context options using common default values
		op := &oto.NewContextOptions{
			SampleRate:   48000,                   // Adjust sample rate to 48000 for better sound quality
			ChannelCount: 2,                       // Stereo (2 channels)
			Format:       oto.FormatSignedInt16LE, // Use 16-bit signed integer format for better consistency with SDL2
		}

		// Initialize Oto context
		ctx, readyChan, err := oto.NewContext(op)
		if err != nil {
			fmt.Printf("could not initialize Oto: %v", err)
			return nil
		}

		// Wait for the audio device to be ready
		<-readyChan

		theContext = ctx
	}

	player := &Player{
		Initialized:         true,
		PlayingAudioDevices: make(map[AudioDeviceKey]*oto.Player),
		ctx:                 theContext,
	}

	return player
}

// Close does not close the Oto context, but stops all playing audio and clears resources.
func (player *Player) Close() {
	mut.Lock()
	defer mut.Unlock()

	// Stop all playing audio
	for _, audioPlayer := range player.PlayingAudioDevices {
		audioPlayer.Close()
	}

	player.PlayingAudioDevices = make(map[AudioDeviceKey]*oto.Player)

	player.Initialized = false
}

// PlayKick generates and plays the current kick drum sound
func (player *Player) GeneratePlay(t string, cfg *Settings) (AudioDeviceKey, time.Duration, error) {
	if !player.Initialized {
		return -1, 0, errors.New("Oto Audio needs to be initialized first")
	}
	samples, err := cfg.Generate(t)
	if err != nil {
		return -1, 0, err
	}
	return player.PlayWaveform(samples, cfg.SampleRate, cfg.BitDepth, cfg.Channels)
}

// PlayWaveform plays raw waveform samples using Oto
func (player *Player) PlayWaveform(samples []float64, sampleRate, bitDepth, channels int) (AudioDeviceKey, time.Duration, error) {
	if !player.Initialized {
		return -1, 0, errors.New("Oto Audio needs to be initialized first")
	}

	// Convert the float64 samples into the correct format for playback
	buf := new(bytes.Buffer)
	sampleSize := 0 // To track how many bytes per sample for duration calculation

	switch bitDepth {
	case 8:
		sampleSize = 1 // 8-bit
		for _, sample := range samples {
			sampleInt8 := int8(sample * 127)
			if err := binary.Write(buf, binary.LittleEndian, sampleInt8); err != nil {
				return -1, 0, fmt.Errorf("error converting float64 to int8: %v", err)
			}
		}
	case 16:
		sampleSize = 2 // 16-bit
		for _, sample := range samples {
			sampleInt16 := int16(sample * 32767)
			if err := binary.Write(buf, binary.LittleEndian, sampleInt16); err != nil {
				return -1, 0, fmt.Errorf("error converting float64 to int16: %v", err)
			}
		}
	case 32:
		sampleSize = 4 // 32-bit
		for _, sample := range samples {
			sampleFloat32 := float32(sample)
			if err := binary.Write(buf, binary.LittleEndian, sampleFloat32); err != nil {
				return -1, 0, fmt.Errorf("error converting float64 to float32: %v", err)
			}
		}
	default:
		return -1, 0, fmt.Errorf("unsupported bit depth: %d", bitDepth)
	}

	// Create a new Oto Player and write the buffer to it
	audioPlayer := player.ctx.NewPlayer(buf)

	// Generate a new AudioDeviceKey and store the player
	mut.Lock()
	defer mut.Unlock()

	var audioDeviceKey AudioDeviceKey = -1
	for i := AudioDeviceKey(0); i < 1024; i++ {
		if _, found := player.PlayingAudioDevices[i]; !found {
			player.PlayingAudioDevices[i] = audioPlayer
			audioDeviceKey = i
			break
		}
	}

	// Start playback
	audioPlayer.Play()

	// Adjust silence threshold to be more forgiving for Oto playback
	const silenceThreshold = 0.01

	lastNonSilentIndex := len(samples) - 1
	for i := len(samples) - 1; i >= 0; i-- {
		if samples[i] > silenceThreshold || samples[i] < -silenceThreshold {
			lastNonSilentIndex = i
			break
		}
	}

	adjustedSamples := lastNonSilentIndex + 1
	adjustedDuration := time.Duration(float64(adjustedSamples*sampleSize*channels) / float64(sampleRate) * float64(time.Second))

	return audioDeviceKey, adjustedDuration, nil
}

// WaitClose will wait for the given duration and then stop the audio player
func (player *Player) WaitClose(audioDeviceKey AudioDeviceKey, t time.Duration) {
	time.Sleep(t)

	mut.Lock()
	defer mut.Unlock()

	if audioPlayer, found := player.PlayingAudioDevices[audioDeviceKey]; found {
		audioPlayer.Close()
		delete(player.PlayingAudioDevices, audioDeviceKey)
	}
}

// WaitClosePlus100 will wait for the given duration + 100 milliseconds and then stop the audio player
func (player *Player) WaitClosePlus100(audioDeviceKey AudioDeviceKey, t time.Duration) {
	time.Sleep(t + 100*time.Millisecond)

	mut.Lock()
	defer mut.Unlock()

	if audioPlayer, found := player.PlayingAudioDevices[audioDeviceKey]; found {
		audioPlayer.Close()
		delete(player.PlayingAudioDevices, audioDeviceKey)
	}
}
