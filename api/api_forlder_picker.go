package api

import (
	"encoding/json"
	"net/http"
	"os/exec"
	"runtime"
	"strings"
)

func apiFolderPickerHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	var cmd *exec.Cmd

	switch runtime.GOOS {

	case "darwin":

		cmd = exec.Command(
			"osascript",
			"-e",
			`POSIX path of (choose folder with prompt "Избери директория")`,
		)

	case "windows":

		cmd = exec.Command(
			"powershell",
			"-NoProfile",
			"-Command",
			`Add-Type -AssemblyName System.Windows.Forms;$dialog=New-Object System.Windows.Forms.FolderBrowserDialog;if($dialog.ShowDialog() -eq [System.Windows.Forms.DialogResult]::OK){Write-Output $dialog.SelectedPath}`,
		)

	default:

		json.NewEncoder(w).Encode(
			map[string]interface{}{
				"error": "Unsupported operating system",
			},
		)

		return
	}

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

	if path == "" {

		json.NewEncoder(w).Encode(
			map[string]interface{}{
				"cancelled": true,
			},
		)

		return
	}

	json.NewEncoder(w).Encode(
		map[string]interface{}{
			"path": path,
		},
	)
}
