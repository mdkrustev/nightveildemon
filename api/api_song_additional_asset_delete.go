package api

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func apiSongAdditionalAssetDeleteHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	query := r.URL.Query()
	songID := query.Get("songId")
	storageName := query.Get("storageName")
	storagePath := query.Get("storagePath")
	hashName := query.Get("hashName")
	fileName := query.Get("fileName")
	if songID == "" || storageName == "" || storagePath == "" || hashName == "" || fileName == "" {
		json.NewEncoder(w).Encode(map[string]string{"error": "Missing parameters"})
		return
	}
	id, err := strconv.Atoi(songID)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid song id"})
		return
	}
	additionalDir := filepath.Join(storagePath, storageName, hashName, strconv.Itoa(id), "additional")
	target := filepath.Join(additionalDir, fileName)
	if _, err := os.Stat(target); err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": "File not found"})
		return
	}
	if err := os.Remove(target); err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": "Unable to delete file"})
		return
	}
	entries, err := os.ReadDir(additionalDir)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": "Unable to read directory"})
		return
	}
	type item struct {
		Name  string
		Index int
	}
	files := []item{}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasSuffix(strings.ToLower(name), ".jpg") {
			continue
		}
		base := strings.TrimSuffix(strings.ToLower(name), ".jpg")
		if !strings.HasPrefix(base, "img") {
			continue
		}
		num, err := strconv.Atoi(strings.TrimPrefix(base, "img"))
		if err != nil {
			continue
		}
		files = append(files, item{
			Name:  name,
			Index: num,
		})
	}
	sort.Slice(files, func(i, j int) bool {
		return files[i].Index < files[j].Index
	})
	for i, file := range files {
		oldPath := filepath.Join(additionalDir, file.Name)
		tmpPath := filepath.Join(additionalDir, "rename_"+strconv.Itoa(i+1)+".tmp")
		if err := os.Rename(oldPath, tmpPath); err != nil {
			json.NewEncoder(w).Encode(map[string]string{"error": "Unable to prepare rename"})
			return
		}
	}
	for i := range files {
		tmpPath := filepath.Join(additionalDir, "rename_"+strconv.Itoa(i+1)+".tmp")
		newPath := filepath.Join(additionalDir, "img"+strconv.Itoa(i+1)+".jpg")
		if err := os.Rename(tmpPath, newPath); err != nil {
			json.NewEncoder(w).Encode(map[string]string{"error": "Unable to finalize rename"})
			return
		}
	}
	err = UpdateSongManifest(storagePath, storageName, hashName, id)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": "Unable to update manifest: " + err.Error()})
		return
	}
	json.NewEncoder(w).Encode(map[string]any{
		"status": "success",
		"songId": id,
	})
}
