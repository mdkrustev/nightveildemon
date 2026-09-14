package api

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func apiSongAssetsUploadHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	err := r.ParseMultipartForm(
		200 << 20,
	)

	if err != nil {

		json.NewEncoder(w).Encode(
			map[string]string{
				"error": "Invalid multipart form",
			},
		)

		return
	}

	songID := r.FormValue("songId")
	storageName := r.FormValue("storageName")
	storagePath := r.FormValue("storagePath")
	hashName := r.FormValue("hashName")

	if songID == "" ||
		storageName == "" ||
		storagePath == "" ||
		hashName == "" {

		json.NewEncoder(w).Encode(
			map[string]string{
				"error": "Missing required data",
			},
		)

		return
	}

	id, err := strconv.Atoi(songID)

	if err != nil {

		json.NewEncoder(w).Encode(
			map[string]string{
				"error": "Invalid song id",
			},
		)

		return
	}

	// storagePath/storageName/hashName

	hashDir := filepath.Join(
		storagePath,
		storageName,
		hashName,
	)

	info, err := os.Stat(hashDir)

	if err != nil ||
		!info.IsDir() {

		json.NewEncoder(w).Encode(
			map[string]string{
				"error": "Hash storage directory does not exist",
			},
		)

		return
	}

	// hashName/songId

	songDir := filepath.Join(
		hashDir,
		strconv.Itoa(id),
	)

	err = EnsureDir(
		songDir,
	)

	if err != nil {

		json.NewEncoder(w).Encode(
			map[string]string{
				"error": "Unable to create song directory",
			},
		)

		return
	}

	/*
		AUDIO
	*/

	audioFile, _, err := r.FormFile(
		"audio",
	)

	if err != nil {

		json.NewEncoder(w).Encode(
			map[string]string{
				"error": "Audio file missing",
			},
		)

		return
	}

	defer audioFile.Close()

	tempAudio, err := os.CreateTemp(
		"",
		"nightveil-audio-*",
	)

	if err != nil {

		json.NewEncoder(w).Encode(
			map[string]string{
				"error": "Cannot create temp audio",
			},
		)

		return
	}

	tempAudioPath := tempAudio.Name()

	defer os.Remove(
		tempAudioPath,
	)

	_, err = io.Copy(
		tempAudio,
		audioFile,
	)

	tempAudio.Close()

	if err != nil {

		json.NewEncoder(w).Encode(
			map[string]string{
				"error": "Unable to save audio",
			},
		)

		return
	}

	audioPath := filepath.Join(
		songDir,
		"audio.mp3",
	)

	err = SaveCleanAudio(
		tempAudioPath,
		audioPath,
	)

	if err != nil {

		json.NewEncoder(w).Encode(
			map[string]string{
				"error": err.Error(),
			},
		)

		return
	}

	/*
		COVER
	*/

	coverFile, _, err := r.FormFile(
		"cover",
	)

	if err == nil {

		defer coverFile.Close()

		tempCover, err := os.CreateTemp(
			"",
			"nightveil-cover-*",
		)

		if err != nil {

			json.NewEncoder(w).Encode(
				map[string]string{
					"error": "Cannot create temp cover",
				},
			)

			return
		}

		tempCoverPath := tempCover.Name()

		defer os.Remove(
			tempCoverPath,
		)

		_, err = io.Copy(
			tempCover,
			coverFile,
		)

		tempCover.Close()

		if err != nil {

			json.NewEncoder(w).Encode(
				map[string]string{
					"error": "Unable to save cover",
				},
			)

			return
		}

		coverPath := filepath.Join(
			songDir,
			"cover.jpg",
		)

		err = ConvertCoverToJPG(
			tempCoverPath,
			coverPath,
		)

		if err != nil {

			json.NewEncoder(w).Encode(
				map[string]string{
					"error": err.Error(),
				},
			)

			return
		}

	}

	json.NewEncoder(w).Encode(
		map[string]interface{}{

			"status": "success",

			"songId": id,

			"path": songDir,
		},
	)

}
