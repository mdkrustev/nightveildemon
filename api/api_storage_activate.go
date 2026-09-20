package api

import (
	"encoding/json"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

type StorageActivateRequest struct {
	Path          string `json:"path"`
	StorageName   string `json:"storageName"`
	HashDirectory string `json:"hashDirectory"`
}

func apiStorageActivateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var body StorageActivateRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{
			"success": false,
			"error":   "Invalid request body",
		})
		return
	}
	if body.Path == "" || body.StorageName == "" || body.HashDirectory == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{
			"success": false,
			"error":   "Missing required fields",
		})
		return
	}
	if _, err := os.Stat(body.Path); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{
			"success": false,
			"error":   "Storage path does not exist",
		})
		return
	}
	packagePath := filepath.Join(
		body.Path,
		"package.json",
	)
	if _, err := os.Stat(packagePath); os.IsNotExist(err) {
		packageJSON := `{
  "devDependencies": {
    "wrangler": "^4.128.0"
  }
}`
		if err := os.WriteFile(
			packagePath,
			[]byte(packageJSON),
			0644,
		); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]any{
				"success": false,
				"error":   err.Error(),
			})
			return
		}
	}
	storageDir := filepath.Join(
		body.Path,
		body.StorageName,
	)
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]any{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	hashDir := filepath.Join(
		storageDir,
		body.HashDirectory,
	)
	if err := os.MkdirAll(hashDir, 0755); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]any{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	manifestPath := filepath.Join(
		hashDir,
		"manifest.json",
	)
	if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
		manifest := SongManifest{
			DeployVersion: 0,
			Songs:         []ManifestSong{},
		}
		data, err := json.MarshalIndent(
			manifest,
			"",
			"  ",
		)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]any{
				"success": false,
				"error":   err.Error(),
			})
			return
		}
		if err := os.WriteFile(
			manifestPath,
			data,
			0644,
		); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]any{
				"success": false,
				"error":   err.Error(),
			})
			return
		}
	}
	nodeModulesPath := filepath.Join(
		body.Path,
		"node_modules",
	)
	if _, err := os.Stat(nodeModulesPath); os.IsNotExist(err) {
		cmd := exec.Command(
			"npm",
			"install",
		)
		cmd.Dir = body.Path
		output, err := cmd.CombinedOutput()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]any{
				"success": false,
				"error":   string(output),
			})
			return
		}
	}
	err := RestoreStorageAssets(
		body.Path,
		body.StorageName,
		body.HashDirectory,
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]any{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"path":    storageDir,
	})
}
