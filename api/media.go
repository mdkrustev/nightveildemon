package api

import (
	"os"
	"path/filepath"
	"runtime"
)

func ffmpegPath() string {

	if runtime.GOOS == "windows" {

		if _, err := os.Stat("./bin/ffmpeg.exe"); err == nil {

			return "./bin/ffmpeg.exe"

		}

	} else {

		if _, err := os.Stat("./bin/ffmpeg"); err == nil {

			return "./bin/ffmpeg"

		}

	}

	exe, err := os.Executable()

	if err != nil {

		return "ffmpeg"

	}

	appDir := filepath.Dir(exe)

	switch runtime.GOOS {

	case "windows":

		return filepath.Join(
			appDir,
			"ffmpeg.exe",
		)

	case "darwin":

		return filepath.Join(
			appDir,
			"../Resources/ffmpeg",
		)

	case "linux":

		return filepath.Join(
			appDir,
			"ffmpeg",
		)

	}

	return "ffmpeg"

}
