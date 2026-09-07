package api

import (
	"encoding/json"
	"net/http"
	"os"
)

func apiCheckFolderHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	path := r.URL.Query().Get("path")

	if path == "" {

		json.NewEncoder(w).Encode(
			map[string]interface{}{
				"exists": false,
			},
		)

		return
	}

	_, err := os.Stat(path)

	exists := err == nil

	json.NewEncoder(w).Encode(
		map[string]interface{}{
			"exists": exists,
			"path":   path,
		},
	)
}
