package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type StorageSyncResult struct {
	DeployVersion int  `json:"deployVersion"`
	AddedSongs    int  `json:"addedSongs"`
	UpdatedSongs  int  `json:"updatedSongs"`
	IsUpToDate    bool `json:"isUpToDate"`
}

func SynchronizeStorage(storagePath string, storageName string, hashName string) (*StorageSyncResult, error) {
	localManifestPath := filepath.Join(storagePath, storageName, hashName, "manifest.json")
	remoteManifestURL := fmt.Sprintf("https://%s.pages.dev/%s/manifest.json", storageName, hashName)
	remoteManifest, err := fetchRemoteManifest(remoteManifestURL)
	if err != nil {
		return nil, err
	}
	localManifest, err := loadLocalManifest(localManifestPath)
	if err != nil {
		return nil, err
	}
	localSongs := map[int]ManifestSong{}
	for _, song := range localManifest.Songs {
		localSongs[song.ID] = song
	}
	addedSongs := 0
	updatedSongs := 0
	resultSongs := []ManifestSong{}
	for _, remoteSong := range remoteManifest.Songs {
		localSong, exists := localSongs[remoteSong.ID]
		if !exists {
			if err := syncSongAssets(storagePath, storageName, hashName, remoteSong); err != nil {
				return nil, err
			}
			addedSongs++
			resultSongs = append(resultSongs, remoteSong)
			continue
		}
		if remoteSong.Version > localSong.Version {
			if err := replaceSongAssets(storagePath, storageName, hashName, remoteSong); err != nil {
				return nil, err
			}
			updatedSongs++
			resultSongs = append(resultSongs, remoteSong)
			continue
		}
		if localSong.Version > remoteSong.Version {
			if err := repairMissingAssets(storagePath, storageName, hashName, localSong); err != nil {
				return nil, err
			}
			resultSongs = append(resultSongs, localSong)
			continue
		}
		if err := repairMissingAssets(storagePath, storageName, hashName, localSong); err != nil {
			return nil, err
		}
		resultSongs = append(resultSongs, localSong)
	}
	localManifest.DeployVersion = remoteManifest.DeployVersion
	localManifest.Songs = resultSongs
	data, err := json.MarshalIndent(localManifest, "", "  ")
	if err != nil {
		return nil, err
	}
	if err := os.WriteFile(localManifestPath, data, 0644); err != nil {
		return nil, err
	}
	return &StorageSyncResult{
		DeployVersion: remoteManifest.DeployVersion,
		AddedSongs:    addedSongs,
		UpdatedSongs:  updatedSongs,
		IsUpToDate:    addedSongs == 0 && updatedSongs == 0,
	}, nil
}

func fetchRemoteManifest(url string) (SongManifest, error) {
	var manifest SongManifest
	resp, err := http.Get(url)
	if err != nil {
		return manifest, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return manifest, fmt.Errorf("remote manifest unavailable: %s", resp.Status)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return manifest, err
	}
	err = json.Unmarshal(data, &manifest)
	return manifest, err
}

func loadLocalManifest(path string) (SongManifest, error) {
	var manifest SongManifest
	data, err := os.ReadFile(path)
	if err != nil {
		return manifest, err
	}
	err = json.Unmarshal(data, &manifest)
	return manifest, err
}

func syncSongAssets(storagePath string, storageName string, hashName string, song ManifestSong) error {
	songDir := filepath.Join(storagePath, storageName, hashName, fmt.Sprintf("%d", song.ID))
	if err := os.MkdirAll(songDir, 0755); err != nil {
		return err
	}
	return downloadSongAssets(storageName, hashName, songDir, song)
}

func replaceSongAssets(storagePath string, storageName string, hashName string, song ManifestSong) error {
	songDir := filepath.Join(storagePath, storageName, hashName, fmt.Sprintf("%d", song.ID))
	if err := os.RemoveAll(songDir); err != nil {
		return err
	}
	return syncSongAssets(storagePath, storageName, hashName, song)
}

func repairMissingAssets(storagePath string, storageName string, hashName string, song ManifestSong) error {
	songDir := filepath.Join(storagePath, storageName, hashName, fmt.Sprintf("%d", song.ID))
	return downloadSongAssetsIfMissing(storageName, hashName, songDir, song)
}
