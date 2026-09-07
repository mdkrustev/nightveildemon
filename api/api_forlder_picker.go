package api

import (
	"encoding/json"
	"net/http"
	"os/exec"
	"strings"
)

func apiFolderPickerHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	cmd := exec.Command(
		"osascript",
		"-e",
		`POSIX path of (choose folder with prompt "Избери директория")`,
	)

	out, err := cmd.Output()

	if err != nil {
		json.NewEncoder(w).Encode(
			map[string]interface{}{
				"cancelled": true,
			},
		)
		return
	}

	path := strings.TrimSpace(string(out))

	json.NewEncoder(w).Encode(
		map[string]interface{}{
			"path": path,
		},
	)
}
