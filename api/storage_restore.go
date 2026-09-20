package api

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func RestoreStorageAssets(
	storagePath string,
	storageName string,
	hashName string,
) error {

	baseURL := "https://" + storageName + ".pages.dev/" + hashName

	manifestURL := baseURL + "/manifest.json"

	resp, err := http.Get(manifestURL)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil
	}

	manifestData, err := io.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	var manifest SongManifest

	err = json.Unmarshal(
		manifestData,
		&manifest,
	)

	if err != nil {
		return err
	}

	localManifestPath := filepath.Join(
		storagePath,
		storageName,
		hashName,
		"manifest.json",
	)

	err = os.WriteFile(
		localManifestPath,
		manifestData,
		0644,
	)

	if err != nil {
		return err
	}

	for _, song := range manifest.Songs {

		songDir := filepath.Join(
			storagePath,
			storageName,
			hashName,
			strconv.Itoa(song.ID),
		)

		err = os.MkdirAll(
			songDir,
			0755,
		)

		if err != nil {
			return err
		}

		if song.Assets.Audio {

			err = downloadStorageFile(
				baseURL+"/"+
					strconv.Itoa(song.ID)+
					"/audio.mp3",
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

			err = downloadStorageFile(
				baseURL+"/"+
					strconv.Itoa(song.ID)+
					"/cover.jpg",
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

			err = os.MkdirAll(
				additionalDir,
				0755,
			)

			if err != nil {
				return err
			}

			for i := 1; i <= song.Assets.Additional; i++ {

				err = downloadStorageFile(
					baseURL+"/"+
						strconv.Itoa(song.ID)+
						"/additional/img"+
						strconv.Itoa(i)+
						".jpg",
					filepath.Join(
						additionalDir,
						"img"+strconv.Itoa(i)+".jpg",
					),
				)

				if err != nil {
					return err
				}

			}

		}

	}

	return nil
}

func downloadStorageFile(
	url string,
	target string,
) error {

	resp, err := http.Get(url)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil
	}

	file, err := os.Create(
		target,
	)

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
