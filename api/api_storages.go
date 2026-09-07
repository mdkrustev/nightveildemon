package api

import (
	"encoding/json"
	"net/http"
)

// GET /api/files
func apiStoragesHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	files := []string{
		"example.mp4",
		"image.png",
		"audio.mp3",
	}

	json.NewEncoder(w).Encode(
		map[string]interface{}{
			"files": files,
		},
	)
}
