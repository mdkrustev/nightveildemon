package main

import (
	"nightveil-demon/media"
	"os/exec"
)

func CheckFFmpeg() bool {

	cmd := exec.Command(
		media.FFmpegPath(),
		"-version",
	)

	err := cmd.Run()

	return err == nil
}
