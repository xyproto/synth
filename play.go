package synth

import (
	"errors"
	"fmt"
	"os"

	"github.com/xyproto/playsample"
)

var player = playsample.NewPlayer()

// GeneratePlay generates and plays the current kick drum sound
func GeneratePlay(t string, cfg *Settings) error {
	if !player.Initialized {
		return errors.New("A Player needs to be initialized first")
	}
	samples, err := cfg.Generate(t)
	if err != nil {
		return err
	}
	return player.PlayWaveform(samples, cfg.SampleRate, cfg.BitDepth, cfg.Channels)
}

// FFGeneratePlay generates a waveform of the given type, saves it to a temporary WAV file,
// plays it using ffplay, and then deletes the temporary file.
func FFGeneratePlay(t string, cfg *Settings) error {
	// Generate the kick drum waveform
	samples, err := cfg.Generate(t)
	if err != nil {
		return fmt.Errorf("error generating kick waveform: %v", err)
	}
	// Create a temporary WAV file
	tmpFile, err := os.CreateTemp("", "kickdrum_*.wav")
	if err != nil {
		return fmt.Errorf("error creating temporary wav file: %v", err)
	}
	defer os.Remove(tmpFile.Name()) // Ensure the file is removed after playing
	// Save the waveform to the temporary file using SaveToWav
	if err := playsample.SaveToWav(tmpFile, samples, cfg.SampleRate, cfg.BitDepth, cfg.Channels); err != nil {
		return fmt.Errorf("error saving wav file: %v", err)
	}
	// Play the temporary wav file using ffplay
	err = playsample.FFPlayWav(tmpFile.Name())
	if err != nil {
		return fmt.Errorf("error playing wav file: %v", err)
	}
	return nil
}
