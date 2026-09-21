package api

import (
	"encoding/json"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

func apiStorageDeployHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	query := r.URL.Query()
	storagePath := query.Get("storagePath")
	storageName := query.Get("storageName")
	hashName := query.Get("hashName")
	if storagePath == "" || storageName == "" || hashName == "" {
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Missing parameters",
		})
		return
	}
	/*
		syncResult, err := SynchronizeStorage(
			storagePath,
			storageName,
			hashName,
		)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Storage synchronization failed: " + err.Error(),
			})
			return
		}*/

	manifestPath := filepath.Join(
		storagePath,
		storageName,
		hashName,
		"manifest.json",
	)
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Manifest not found",
		})
		return
	}
	var manifest SongManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid manifest",
		})
		return
	}
	oldVersion := manifest.DeployVersion
	newVersion := oldVersion + 1
	manifest.DeployVersion = newVersion
	newData, err := json.MarshalIndent(
		manifest,
		"",
		"  ",
	)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Unable to encode manifest",
		})
		return
	}
	if err := os.WriteFile(
		manifestPath,
		newData,
		0644,
	); err != nil {
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Unable to update manifest",
		})
		return
	}
	cmd := exec.Command(
		"npx",
		"wrangler",
		"pages",
		"deploy",
		storageName,
		"--project-name",
		storageName,
	)
	cmd.Dir = storagePath
	output, err := cmd.CombinedOutput()
	if err != nil {
		manifest.DeployVersion = oldVersion
		rollback, _ := json.MarshalIndent(
			manifest,
			"",
			"  ",
		)
		_ = os.WriteFile(
			manifestPath,
			rollback,
			0644,
		)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":         "Deploy failed",
			"output":        string(output),
			"deployVersion": oldVersion,
		})
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":        "success",
		"deployVersion": newVersion,
		//"sync":          syncResult,
		"output": string(output),
	})
}
