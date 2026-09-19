package api

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type SongManifest struct {
	DeployVersion int            `json:"deployversion"`
	Songs         []ManifestSong `json:"songs"`
}
type ManifestSong struct {
	ID      int            `json:"id"`
	Version int            `json:"version"`
	Assets  ManifestAssets `json:"assets"`
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
		_ = json.Unmarshal(
			data,
			&manifest,
		)
	} else {
		manifest.DeployVersion = 0
	}
	songDir := filepath.Join(
		storagePath,
		storageName,
		hashName,
		strconv.Itoa(songID),
	)
	audioExists := false
	coverExists := false
	if _, err := os.Stat(
		filepath.Join(
			songDir,
			"audio.mp3",
		),
	); err == nil {
		audioExists = true
	}
	if _, err := os.Stat(
		filepath.Join(
			songDir,
			"cover.jpg",
		),
	); err == nil {
		coverExists = true
	}
	additionalDir := filepath.Join(
		songDir,
		"additional",
	)
	additional := 0
	if entries, err := os.ReadDir(
		additionalDir,
	); err == nil {
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			if strings.HasSuffix(
				strings.ToLower(entry.Name()),
				".jpg",
			) {
				additional++
			}
		}
	}
	for i := range manifest.Songs {
		if manifest.Songs[i].ID != songID {
			continue
		}
		manifest.Songs[i].Version++
		manifest.Songs[i].Assets.Audio = audioExists
		manifest.Songs[i].Assets.Cover = coverExists
		manifest.Songs[i].Assets.Additional = additional
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
	manifest.Songs = append(
		manifest.Songs,
		ManifestSong{
			ID:      songID,
			Version: 0,
			Assets: ManifestAssets{
				Audio:      audioExists,
				Cover:      coverExists,
				Additional: additional,
			},
		},
	)
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
