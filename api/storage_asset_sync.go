package api

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func downloadSongAssets(storageName string, hashName string, songDir string, song ManifestSong) error {
	if err := os.MkdirAll(songDir, 0755); err != nil {
		return err
	}
	if song.Assets.Audio {
		err := downloadFile(
			fmt.Sprintf(
				"https://%s.pages.dev/%s/%d/audio.mp3",
				storageName,
				hashName,
				song.ID,
			),
			filepath.Join(
				songDir,
				"audio.mp3",
			),
		)
		if err != nil {
			return err
		}
	}
	if song.Assets.Cover {
		err := downloadFile(
			fmt.Sprintf(
				"https://%s.pages.dev/%s/%d/cover.jpg",
				storageName,
				hashName,
				song.ID,
			),
			filepath.Join(
				songDir,
				"cover.jpg",
			),
		)
		if err != nil {
			return err
		}
	}
	if song.Assets.Additional > 0 {
		additionalDir := filepath.Join(
			songDir,
			"additional",
		)
		if err := os.MkdirAll(additionalDir, 0755); err != nil {
			return err
		}
		for i := 1; i <= song.Assets.Additional; i++ {
			err := downloadFile(
				fmt.Sprintf(
					"https://%s.pages.dev/%s/%d/additional/img%d.jpg",
					storageName,
					hashName,
					song.ID,
					i,
				),
				filepath.Join(
					additionalDir,
					fmt.Sprintf(
						"img%d.jpg",
						i,
					),
				),
			)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func downloadSongAssetsIfMissing(storageName string, hashName string, songDir string, song ManifestSong) error {
	if err := os.MkdirAll(songDir, 0755); err != nil {
		return err
	}
	if song.Assets.Audio {
		target := filepath.Join(
			songDir,
			"audio.mp3",
		)
		if _, err := os.Stat(target); os.IsNotExist(err) {
			err := downloadFile(
				fmt.Sprintf(
					"https://%s.pages.dev/%s/%d/audio.mp3",
					storageName,
					hashName,
					song.ID,
				),
				target,
			)
			if err != nil {
				return err
			}
		}
	}
	if song.Assets.Cover {
		target := filepath.Join(
			songDir,
			"cover.jpg",
		)
		if _, err := os.Stat(target); os.IsNotExist(err) {
			err := downloadFile(
				fmt.Sprintf(
					"https://%s.pages.dev/%s/%d/cover.jpg",
					storageName,
					hashName,
					song.ID,
				),
				target,
			)
			if err != nil {
				return err
			}
		}
	}
	if song.Assets.Additional > 0 {
		additionalDir := filepath.Join(
			songDir,
			"additional",
		)
		if err := os.MkdirAll(additionalDir, 0755); err != nil {
			return err
		}
		for i := 1; i <= song.Assets.Additional; i++ {
			target := filepath.Join(
				additionalDir,
				fmt.Sprintf(
					"img%d.jpg",
					i,
				),
			)
			if _, err := os.Stat(target); os.IsNotExist(err) {
				err := downloadFile(
					fmt.Sprintf(
						"https://%s.pages.dev/%s/%d/additional/img%d.jpg",
						storageName,
						hashName,
						song.ID,
						i,
					),
					target,
				)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func downloadFile(url string, destination string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf(
			"asset download failed: %s %s",
			url,
			resp.Status,
		)
	}
	file, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(
		file,
		resp.Body,
	)
	return err
}
