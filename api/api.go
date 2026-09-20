package api

import (
	"net/http"
)

type APIResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func RegisterAPIRoutes(mux *http.ServeMux) {

	mux.HandleFunc("/api/storages", apiStoragesHandler)
	mux.HandleFunc("/api/folder-picker", apiFolderPickerHandler)
	mux.HandleFunc("/api/check-folder", apiCheckFolderHandler)
	mux.HandleFunc("/api/storage-activate", apiStorageActivateHandler)
	mux.HandleFunc("/api/storage-deactivate", apiStorageDeactivateHandler)
	mux.HandleFunc("/api/song-assets-upload", apiSongAssetsUploadHandler)
	mux.HandleFunc("/api/song-assets-get", apiSongAssetsGetHandler)
	mux.HandleFunc("/api/song-additional-assets-upload", apiSongAdditionalAssetsUploadHandler)
	mux.HandleFunc("/api/song-additional-asset-get", apiSongAdditionalAssetGetHandler)
	mux.HandleFunc("/api/song-additional-assets-get", apiSongAdditionalAssetsGetHandler)
	mux.HandleFunc("/api/song-additional-asset-delete", apiSongAdditionalAssetDeleteHandler)
	mux.HandleFunc("/api/storage-inspect", apiStorageInspectHandler)
	mux.HandleFunc("/api/storage-deploy", apiStorageDeployHandler)

}
