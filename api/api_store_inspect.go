package api

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
)

type StorageInspectRequest struct {
	StoragePath string                  `json:"storagePath"`
	Storages    []StorageInspectItemReq `json:"storages"`
}
type StorageInspectItemReq struct {
	StorageName string `json:"storageName"`
	HashName    string `json:"hashName"`
}
type StorageInspectResponse struct {
	Storages []StorageInspectItem `json:"storages"`
}
type StorageInspectItem struct {
	StorageName   string `json:"storageName"`
	HashName      string `json:"hashName"`
	Status        string `json:"status"`
	Reason        string `json:"reason,omitempty"`
	DeployVersion *int   `json:"deployVersion,omitempty"`
}

func apiStorageInspectHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Method not allowed",
		})
		return
	}
	var request StorageInspectRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid request",
		})
		return
	}
	response := StorageInspectResponse{
		Storages: []StorageInspectItem{},
	}
	packagePath := filepath.Join(
		request.StoragePath,
		"package.json",
	)
	nodeModulesPath := filepath.Join(
		request.StoragePath,
		"node_modules",
	)
	if _, err := os.Stat(packagePath); err != nil {
		for _, storage := range request.Storages {
			response.Storages = append(
				response.Storages,
				StorageInspectItem{
					StorageName: storage.StorageName,
					HashName:    storage.HashName,
					Status:      "inactive",
					Reason:      "Missing package.json",
				},
			)
		}
		json.NewEncoder(w).Encode(response)
		return
	}
	if _, err := os.Stat(nodeModulesPath); err != nil {
		for _, storage := range request.Storages {
			response.Storages = append(
				response.Storages,
				StorageInspectItem{
					StorageName: storage.StorageName,
					HashName:    storage.HashName,
					Status:      "inactive",
					Reason:      "Missing node_modules",
				},
			)
		}
		json.NewEncoder(w).Encode(response)
		return
	}
	for _, storage := range request.Storages {
		item := StorageInspectItem{
			StorageName: storage.StorageName,
			HashName:    storage.HashName,
			Status:      "inactive",
		}
		storageDir := filepath.Join(
			request.StoragePath,
			storage.StorageName,
		)
		if _, err := os.Stat(storageDir); err != nil {
			item.Reason = "Missing storage directory"
			response.Storages = append(
				response.Storages,
				item,
			)
			continue
		}
		hashDir := filepath.Join(
			storageDir,
			storage.HashName,
		)
		if _, err := os.Stat(hashDir); err != nil {
			item.Reason = "Missing hash directory"
			response.Storages = append(
				response.Storages,
				item,
			)
			continue
		}
		manifestPath := filepath.Join(
			hashDir,
			"manifest.json",
		)
		data, err := os.ReadFile(manifestPath)
		if err != nil {
			item.Reason = "Missing manifest"
			response.Storages = append(
				response.Storages,
				item,
			)
			continue
		}
		var manifest SongManifest
		if err := json.Unmarshal(
			data,
			&manifest,
		); err != nil {
			item.Reason = "Invalid manifest"
			response.Storages = append(
				response.Storages,
				item,
			)
			continue
		}
		item.Status = "active"
		version := manifest.DeployVersion
		item.DeployVersion = &version
		response.Storages = append(
			response.Storages,
			item,
		)
	}
	json.NewEncoder(w).Encode(response)
}
