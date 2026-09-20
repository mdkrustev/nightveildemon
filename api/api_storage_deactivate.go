package api

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
)

type StorageDeactivateRequest struct {
	StoragePath string `json:"storagePath"`
	StorageName string `json:"storageName"`
	HashName    string `json:"hashName"`
}

func apiStorageDeactivateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Method not allowed",
		})
		return
	}

	var body StorageDeactivateRequest

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Invalid request body",
		})
		return
	}

	if body.StoragePath == "" || body.StorageName == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Missing required fields",
		})
		return
	}

	storageDir := filepath.Join(
		body.StoragePath,
		body.StorageName,
	)

	if _, err := os.Stat(storageDir); os.IsNotExist(err) {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "Storage already removed",
		})
		return
	}

	if err := os.RemoveAll(storageDir); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"storage": body.StorageName,
	})
}
