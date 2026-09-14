package api

import (
	"net/http"
	"path/filepath"
	"strconv"
)

func apiSongAssetsGetHandler(w http.ResponseWriter, r *http.Request) {

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

	switch fileName {

	case "audio.mp3", "cover.jpg":

	default:

		http.Error(
			w,
			"Unsupported file",
			http.StatusBadRequest,
		)

		return
	}

	filePath := filepath.Join(
		storagePath,
		storageName,
		hashName,
		strconv.Itoa(id),
		fileName,
	)

	http.ServeFile(
		w,
		r,
		filePath,
	)
}
