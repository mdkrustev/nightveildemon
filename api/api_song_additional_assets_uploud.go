package api

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func apiSongAdditionalAssetsUploadHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	w.Header().Set(
		"Content-Type",
		"application/json",
	)
	err := r.ParseMultipartForm(
		500 << 20,
	)
	if err != nil {
		json.NewEncoder(w).Encode(
			map[string]string{
				"error": "Invalid multipart form",
			},
		)
		return
	}
	songID := r.FormValue("songId")
	storageName := r.FormValue("storageName")
	storagePath := r.FormValue("storagePath")
	hashName := r.FormValue("hashName")
	if songID == "" ||
		storageName == "" ||
		storagePath == "" ||
		hashName == "" {
		json.NewEncoder(w).Encode(
			map[string]string{
				"error": "Missing required data",
			},
		)
		return
	}
	id, err := strconv.Atoi(songID)
	if err != nil {
		json.NewEncoder(w).Encode(
			map[string]string{
				"error": "Invalid song id",
			},
		)
		return
	}
	additionalDir := filepath.Join(
		storagePath,
		storageName,
		hashName,
		strconv.Itoa(id),
		"additional",
	)
	err = EnsureDir(
		additionalDir,
	)
	if err != nil {
		json.NewEncoder(w).Encode(
			map[string]string{
				"error": "Unable to create additional directory",
			},
		)
		return
	}
	files := r.MultipartForm.File["images"]
	if len(files) == 0 {
		json.NewEncoder(w).Encode(
			map[string]string{
				"error": "No images uploaded",
			},
		)
		return
	}
	existingFiles := 0
	entries, err := os.ReadDir(
		additionalDir,
	)
	if err == nil {
		for _, entry := range entries {
			if !entry.IsDir() {
				existingFiles++
			}
		}
	}
	uploaded := []string{}
	for index, header := range files {
		file, err := header.Open()
		if err != nil {
			continue
		}
		temp, err := os.CreateTemp(
			"",
			"nightveil-additional-*",
		)
		if err != nil {
			file.Close()
			continue
		}
		_, err = io.Copy(
			temp,
			file,
		)
		file.Close()
		temp.Close()
		if err != nil {
			os.Remove(
				temp.Name(),
			)
			continue
		}
		fileName := "img" +
			strconv.Itoa(existingFiles+index+1) +
			".jpg"
		target := filepath.Join(
			additionalDir,
			fileName,
		)
		err = ConvertCoverToJPG(
			temp.Name(),
			target,
		)
		os.Remove(
			temp.Name(),
		)
		if err != nil {
			continue
		}
		uploaded = append(
			uploaded,
			fileName,
		)
	}
	err = UpdateSongManifest(
		storagePath,
		storageName,
		hashName,
		id,
	)
	if err != nil {
		json.NewEncoder(w).Encode(
			map[string]string{
				"error": "Unable to update manifest: " + err.Error(),
			},
		)
		return
	}
	json.NewEncoder(w).Encode(
		map[string]interface{}{
			"status": "success",
			"songId": id,
			"files":  uploaded,
		},
	)
}
