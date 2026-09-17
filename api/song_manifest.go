package api

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
)

type SongManifest struct {
	Songs []ManifestSong `json:"songs"`
}

type ManifestSong struct {
	ID     int            `json:"id"`
	Assets ManifestAssets `json:"assets"`
}

type ManifestAssets struct {
	Audio      bool `json:"audio"`
	Cover      bool `json:"cover"`
	Additional int  `json:"additional"`
}

func UpdateSongManifest(
	storagePath string,
	storageName string,
	hashName string,
	songID int,
) error {

	manifestPath := filepath.Join(
		storagePath,
		storageName,
		hashName,
		"manifest.json",
	)

	var manifest SongManifest

	if data, err := os.ReadFile(manifestPath); err == nil {
		_ = json.Unmarshal(data, &manifest)
	}

	songDir := filepath.Join(
		storagePath,
		storageName,
		hashName,
		strconv.Itoa(songID),
	)

	audioExists := false
	coverExists := false

	if _, err := os.Stat(filepath.Join(songDir, "audio.mp3")); err == nil {
		audioExists = true
	}

	if _, err := os.Stat(filepath.Join(songDir, "cover.jpg")); err == nil {
		coverExists = true
	}

	found := false

	for i := range manifest.Songs {

		if manifest.Songs[i].ID != songID {
			continue
		}

		manifest.Songs[i].Assets.Audio = audioExists
		manifest.Songs[i].Assets.Cover = coverExists

		found = true
		break
	}

	if !found {

		manifest.Songs = append(
			manifest.Songs,
			ManifestSong{
				ID: songID,
				Assets: ManifestAssets{
					Audio:      audioExists,
					Cover:      coverExists,
					Additional: 0,
				},
			},
		)
	}

	data, err := json.MarshalIndent(
		manifest,
		"",
		"  ",
	)

	if err != nil {
		return err
	}

	return os.WriteFile(
		manifestPath,
		data,
		0644,
	)
}
