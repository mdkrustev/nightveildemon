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

}
