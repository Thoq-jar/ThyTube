package features

import (
	"fmt"
	"os"
	"os/exec"
)

func Download(url string) error {
	outputTemplate := "./download/%(title)s.%(ext)s"

	args := []string{
		"-f", "best[ext=mp4]",
		"--recode-video", "mp4",
		"-o", outputTemplate,
		url,
	}

	command := exec.Command("yt-dlp", args...)
	command.Stderr = os.Stderr
	command.Stdout = os.Stdout

	err := command.Run()
	if err != nil {
		return fmt.Errorf("yt-dlp command failed: %w", err)
	}

	return nil
}
