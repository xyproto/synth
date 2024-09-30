package synth

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/xyproto/playsample"
)

func (soundType SoundType) String() string {
	switch soundType {
	case Kick:
		return "kick"
	case Clap:
		return "clap"
	case Snare:
		return "snare"
	case ClosedHH:
		return "closedHH"
	case OpenHH:
		return "openHH"
	case Rimshot:
		return "rimshot"
	case Tom:
		return "tom"
	case Percussion:
		return "percussion"
	case Ride:
		return "ride"
	case Crash:
		return "crash"
	case Bass:
		return "bass"
	case Xylophone:
		return "xylophone"
	case Lead:
		return "lead"
	default:
		return "unknown"
	}
}

// GenerateAndSaveTo generates samples for a given type (e.g., "kick", "snare") and saves it to a specified directory, avoiding filename collisions.
func (cfg *Settings) GenerateAndSaveTo(soundType SoundType, directory string) (string, error) {
	n := 1
	var fileName string
	for {
		// Construct the file path with an incrementing number based on the type
		fileName = filepath.Join(directory, fmt.Sprintf("%s%d.wav", soundType, n))
		if _, err := os.Stat(fileName); os.IsNotExist(err) {
			break
		}
		n++
	}
	// Create the new file
	file, err := os.Create(fileName)
	if err != nil {
		return "", err
	}
	defer file.Close()
	// Set the file as the output for the sound generation
	cfg.Output = file
	// Generate the samples for the requested type
	samples, err := cfg.Generate(soundType)
	if err != nil {
		return "", err
	}
	// Save the generated samples to the WAV file
	err = playsample.SaveToWav(file, samples, cfg.SampleRate, cfg.BitDepth, cfg.Channels)
	if err != nil {
		return "", fmt.Errorf("error saving to wav file: %v", err)
	}
	return fileName, nil
}
