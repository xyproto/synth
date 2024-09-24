//go:build sdl2

package synth

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
)

type AudioDeviceKey int

// Player represents a struct with methods for playing audio
type Player struct {
	Initialized         bool
	PlayingAudioDevices map[AudioDeviceKey]sdl.AudioDeviceID
}

var mut sync.Mutex

// NewPlayer tries to initialize SDL2 Audio and returns a Player struct
func NewPlayer() *Player {
	var player Player
	player.PlayingAudioDevices = make(map[AudioDeviceKey]sdl.AudioDeviceID)
	if err := sdl.Init(sdl.INIT_AUDIO); err != nil {
		log.Printf("could not initialize SDL: %v", err)
		return &player
	}
	player.Initialized = true
	return &player
}

// Close calls sdl.Quit and sets the Player struct as uninitialized
func (player *Player) Close() {
	mut.Lock()
	defer mut.Unlock()
	for _, audioDeviceID := range player.PlayingAudioDevices {
		sdl.CloseAudioDevice(audioDeviceID)
	}
	player.PlayingAudioDevices = make(map[AudioDeviceKey]sdl.AudioDeviceID)
	player.Initialized = false
	sdl.Quit()
}

// PlayKick generates and plays the current kick drum sound
func (player *Player) GeneratePlay(t string, cfg *Settings) (AudioDeviceKey, time.Duration, error) {
	if !player.Initialized {
		return -1, 0, errors.New("SDL2 Audio needs to be initialized first")
	}
	samples, err := cfg.Generate(t)
	if err != nil {
		return -1, 0, err
	}
	return player.PlayWaveform(samples, cfg.SampleRate, cfg.BitDepth, cfg.Channels)
}

// PlayWav plays a WAV file using SDL2 and SDL_mixer
func (player *Player) PlayWav(filePath string) error {
	if !player.Initialized {
		return errors.New("SDL2 Audio needs to be initialized first")
	}
	// Open the audio device
	if err := mix.OpenAudio(44100, sdl.AUDIO_S16SYS, 2, 1024); err != nil {
		return fmt.Errorf("could not open audio: %v", err)
	}
	defer mix.CloseAudio()
	// Load the WAV file
	music, err := mix.LoadMUS(filePath)
	if err != nil {
		return fmt.Errorf("could not load music file: %v", err)
	}
	defer music.Free()
	if err := music.Play(1); err != nil {
		return fmt.Errorf("could not play music: %v", err)
	}
	// Wait until the music is finished playing
	for mix.PlayingMusic() {
		sdl.Delay(100)
	}
	return nil
}

func (player *Player) PlayWaveform(samples []float64, sampleRate, bitDepth, channels int) (AudioDeviceKey, time.Duration, error) {
	if !player.Initialized {
		return -1, 0, errors.New("SDL2 Audio needs to be initialized first")
	}
	var desired, obtained sdl.AudioSpec
	desired.Freq = int32(sampleRate)
	desired.Channels = uint8(channels)
	desired.Samples = 4096
	switch bitDepth {
	case 8:
		desired.Format = sdl.AUDIO_S8
	case 16:
		desired.Format = sdl.AUDIO_S16LSB
	case 24:
		desired.Format = sdl.AUDIO_S32LSB
	case 32:
		desired.Format = sdl.AUDIO_F32SYS
	default:
		return -1, 0, fmt.Errorf("unsupported bit depth: %d", bitDepth)
	}
	audioDeviceID, err := sdl.OpenAudioDevice("", false, &desired, &obtained, 0)
	if err != nil {
		return -1, 0, fmt.Errorf("could not open audio device: %v", err)
	}
	mut.Lock()
	defer mut.Unlock()
	var audioDeviceKey AudioDeviceKey = -1
	var i AudioDeviceKey
	for i = 0; i < 1024; i++ {
		if _, found := player.PlayingAudioDevices[i]; !found {
			player.PlayingAudioDevices[i] = audioDeviceID
			audioDeviceKey = i
			break
		}
	}
	if obtained.Format != desired.Format || obtained.Channels != desired.Channels || obtained.Freq != desired.Freq {
		return audioDeviceKey, 0, fmt.Errorf("obtained audio spec does not match desired spec")
	}
	buf := new(bytes.Buffer)
	switch bitDepth {
	case 8:
		for _, sample := range samples {
			sampleInt8 := int8(sample * 127)
			if err := binary.Write(buf, binary.LittleEndian, sampleInt8); err != nil {
				return audioDeviceKey, 0, fmt.Errorf("error converting float64 to int8: %v", err)
			}
		}
	case 16:
		for _, sample := range samples {
			sampleInt16 := int16(sample * 32767)
			if err := binary.Write(buf, binary.LittleEndian, sampleInt16); err != nil {
				return audioDeviceKey, 0, fmt.Errorf("error converting float64 to int16: %v", err)
			}
		}
	case 24:
		for _, sample := range samples {
			sampleInt32 := int32(sample * 8388607)
			bytes24 := make([]byte, 3)
			bytes24[0] = byte(sampleInt32 & 0xFF)
			bytes24[1] = byte((sampleInt32 >> 8) & 0xFF)
			bytes24[2] = byte((sampleInt32 >> 16) & 0xFF)
			if _, err := buf.Write(bytes24); err != nil {
				return audioDeviceKey, 0, fmt.Errorf("error converting float64 to 24-bit: %v", err)
			}
		}
	case 32:
		for _, sample := range samples {
			sampleFloat32 := float32(sample)
			if err := binary.Write(buf, binary.LittleEndian, sampleFloat32); err != nil {
				return audioDeviceKey, 0, fmt.Errorf("error converting float64 to float32: %v", err)
			}
		}
	default:
		return audioDeviceKey, 0, fmt.Errorf("unsupported bit depth: %d", bitDepth)
	}
	if err := sdl.QueueAudio(audioDeviceID, buf.Bytes()); err != nil {
		return audioDeviceKey, 0, fmt.Errorf("could not queue audio data: %v", err)
	}
	sdl.PauseAudioDevice(audioDeviceID, false)

	const silenceThreshold = 0.001

	lastNonSilentIndex := len(samples) - 1
	for i := len(samples) - 1; i >= 0; i-- {
		if samples[i] > silenceThreshold || samples[i] < -silenceThreshold {
			lastNonSilentIndex = i
			break
		}
	}

	adjustedSamples := lastNonSilentIndex + 1
	adjustedDuration := time.Duration(float64(adjustedSamples) / float64(sampleRate) * float64(time.Second))

	return audioDeviceKey, adjustedDuration, nil
}

// WaitClose will wait for the given duration using sdl.Delay and then close the audio device
func (player *Player) WaitClose(audioDeviceKey AudioDeviceKey, t time.Duration) {
	// sdl.GetQueuedAudioSize(audioDeviceID) > 0
	sdl.Delay(uint32(t.Milliseconds()))
	mut.Lock()
	defer mut.Unlock()
	for k, audioDeviceID := range player.PlayingAudioDevices {
		if k == audioDeviceKey {
			sdl.CloseAudioDevice(audioDeviceID)
			delete(player.PlayingAudioDevices, audioDeviceKey)
			break
		}
	}
}

// WaitClosePlus100 will wait for the given duration + 100 milliseconds, using sdl.Delay, and then close the audio device
func (player *Player) WaitClosePlus100(audioDeviceKey AudioDeviceKey, t time.Duration) {
	sdl.Delay(uint32(t.Milliseconds()) + 100)
	mut.Lock()
	defer mut.Unlock()
	for k, audioDeviceID := range player.PlayingAudioDevices {
		if k == audioDeviceKey {
			sdl.CloseAudioDevice(audioDeviceID)
			delete(player.PlayingAudioDevices, audioDeviceKey)
			break
		}
	}
}
