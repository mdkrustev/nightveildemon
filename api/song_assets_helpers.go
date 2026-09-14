package api

import (
	"fmt"
	"os"
	"os/exec"

	"nightveil-demon/media"
)

func SaveCleanAudio(input string, output string) error {

	cmd := exec.Command(

		media.FFmpegPath(),

		"-y",

		"-i",
		input,

		// взима само аудио потока
		"-map",
		"0:a:0",

		// маха всички metadata
		"-map_metadata",
		"-1",

		// маха attached pictures / video streams
		"-vn",

		// без прекодиране
		"-c:a",
		"copy",

		output,
	)

	out, err := cmd.CombinedOutput()

	if err != nil {

		return fmt.Errorf(
			"ffmpeg audio cleanup failed: %s",
			string(out),
		)

	}

	return nil
}

func ConvertCoverToJPG(input string, output string) error {

	cmd := exec.Command(

		media.FFmpegPath(),

		"-y",

		"-i",
		input,

		"-q:v",
		"2",

		output,
	)

	out, err := cmd.CombinedOutput()

	if err != nil {

		return fmt.Errorf(
			"cover conversion failed: %s",
			string(out),
		)

	}

	return nil

}

func EnsureDir(path string) error {

	info, err := os.Stat(path)

	if err == nil {

		if !info.IsDir() {

			return fmt.Errorf(
				"path exists but is not directory: %s",
				path,
			)

		}

		return nil

	}

	if !os.IsNotExist(err) {

		return err

	}

	return os.MkdirAll(
		path,
		0755,
	)

}
