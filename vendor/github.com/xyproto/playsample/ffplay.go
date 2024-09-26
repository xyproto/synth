package playsample

import (
	"fmt"
	"os/exec"
)

// FFPlayWav plays a WAV file using ffplay
func FFPlayWav(filePath string) error {
	cmd := exec.Command("ffplay", "-nodisp", "-autoexit", filePath)
	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("error playing sound with ffplay: %v", err)
	}
	return cmd.Wait()
}

// FFPlayWav plays a WAV file using ffplay and by specifying a sample rate
func FFPlayWavWithSampleRate(filePath string, sampleRate int) error {
	cmd := exec.Command("ffplay", "-nodisp", "-autoexit", "-ar", fmt.Sprintf("%d", sampleRate), filePath)
	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("error playing sound with ffplay: %v", err)
	}
	return cmd.Wait()
}
