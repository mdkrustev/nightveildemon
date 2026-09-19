package api

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func apiSongAdditionalAssetsGetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	query := r.URL.Query()
	songID := query.Get("songId")
	storageName := query.Get("storageName")
	storagePath := query.Get("storagePath")
	hashName := query.Get("hashName")
	if songID == "" || storageName == "" || storagePath == "" || hashName == "" {
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Missing parameters",
		})
		return
	}
	id, err := strconv.Atoi(songID)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid song id",
		})
		return
	}
	manifestPath := filepath.Join(
		storagePath,
		storageName,
		hashName,
		"manifest.json",
	)
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "success",
			"songId": id,
			"files":  []string{},
		})
		return
	}
	var manifest SongManifest
	err = json.Unmarshal(data, &manifest)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid manifest",
		})
		return
	}
	for _, song := range manifest.Songs {
		if song.ID != id {
			continue
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "success",
			"songId": id,
			"files":  song.Assets.Additional,
		})
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"songId": id,
		"files":  []string{},
	})
}
