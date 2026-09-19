package api

import (
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

func apiSongAdditionalAssetGetHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	query := r.URL.Query()
	songID := query.Get("songId")
	storageName := query.Get("storageName")
	storagePath := query.Get("storagePath")
	hashName := query.Get("hashName")
	fileName := query.Get("fileName")
	if songID == "" ||
		storageName == "" ||
		storagePath == "" ||
		hashName == "" ||
		fileName == "" {
		http.Error(
			w,
			"Missing parameters",
			http.StatusBadRequest,
		)
		return
	}
	id, err := strconv.Atoi(songID)
	if err != nil {
		http.Error(
			w,
			"Invalid song id",
			http.StatusBadRequest,
		)
		return
	}
	cleanName := filepath.Clean(fileName)
	if strings.Contains(
		cleanName,
		"..",
	) {
		http.Error(
			w,
			"Invalid file name",
			http.StatusBadRequest,
		)
		return
	}
	filePath := filepath.Join(
		storagePath,
		storageName,
		hashName,
		strconv.Itoa(id),
		"additional",
		cleanName,
	)
	http.ServeFile(
		w,
		r,
		filePath,
	)
}
