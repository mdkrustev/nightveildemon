package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func apiStorageSynchronizeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	query := r.URL.Query()
	storagePath := query.Get("storagePath")
	storageName := query.Get("storageName")
	hashName := query.Get("hashName")
	if storagePath == "" || storageName == "" || hashName == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{
			"status":  "error",
			"message": "Missing required storage parameters.",
		})
		return
	}
	result, err := SynchronizeStorage(
		storagePath,
		storageName,
		hashName,
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]any{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	message := ""
	if result.AddedSongs == 0 && result.UpdatedSongs == 0 {
		message = fmt.Sprintf(
			"Storage is already synchronized. Deploy version is v%d.",
			result.DeployVersion,
		)
	} else {
		message = "Synchronization completed."
		if result.AddedSongs > 0 {
			message += fmt.Sprintf(
				" Added %d new songs.",
				result.AddedSongs,
			)
		}
		if result.UpdatedSongs > 0 {
			message += fmt.Sprintf(
				" Updated %d songs.",
				result.UpdatedSongs,
			)
		}
		message += fmt.Sprintf(
			" Deploy version is v%d.",
			result.DeployVersion,
		)
	}
	json.NewEncoder(w).Encode(map[string]any{
		"status":  "success",
		"message": message,
	})
}
