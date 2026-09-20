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
	err = json.Unmarshal(data, &manifest)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid manifest",
		})
		return
	}
	oldVersion := manifest.DeployVersion
	newVersion := oldVersion + 1
	manifest.DeployVersion = newVersion
	newData, _ := json.MarshalIndent(
		manifest,
		"",
		"  ",
	)
	err = os.WriteFile(
		manifestPath,
		newData,
		0644,
	)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Unable to update manifest",
		})
		return
	}
	deployDir := filepath.Join(
		storagePath,
		storageName,
	)
	cmd := exec.Command(
		"npm",
		"run",
		"deploy",
	)
	cmd.Dir = deployDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		manifest.DeployVersion = oldVersion
		rollback, _ := json.MarshalIndent(
			manifest,
			"",
			"  ",
		)
		os.WriteFile(
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
		"output":        string(output),
	})
}
