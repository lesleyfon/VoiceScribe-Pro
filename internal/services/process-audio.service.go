package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

type TranscriptionResponse struct {
	Transcription   string  `json:"transcription"`
	RequestDuration float64 `json:"request_duration"`
}
type ProcessAudioRequestBody struct {
	Content bytes.Buffer `json:"content"`
	Message string       `json:"message"`
}

func ProcessAudio(audioData []byte) (*TranscriptionResponse, error) {

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("audio_file", "audio.wav")

	if err != nil {
		return nil, err
	}

	_, err = io.Copy(part, bytes.NewReader(audioData))
	if err != nil {
		return nil, err
	}

	writer.Close()

	var DOCKER_URL_BASE_PATH = os.Getenv("INTERNAL_LOCAL_DOCKER_URL_BASE_PATH")

	if DOCKER_URL_BASE_PATH == "" {
		log.Println("INTERNAL_LOCAL_DOCKER_URL_BASE_PATH is not set, using default value")
		return nil, fmt.Errorf("INTERNAL_LOCAL_DOCKER_URL_BASE_PATH is not set")
	}

	log.Println("DOCKER_URL_BASE_PATH", DOCKER_URL_BASE_PATH)

	// Make a request to python ml service
	request, err := http.NewRequest("POST", DOCKER_URL_BASE_PATH+"/process-audio", body)

	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", writer.FormDataContentType())

	response, err := http.DefaultClient.Do(request)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	var bodyBytes []byte
	if response.StatusCode != http.StatusOK {
		// Read the body for debugging
		bodyBytes, _ = io.ReadAll(response.Body)
		return nil, fmt.Errorf("bad status: %s, body: %s", response.Status, string(bodyBytes))
	}

	var result TranscriptionResponse
	err = json.NewDecoder(response.Body).Decode(&result)
	log.Print(result.RequestDuration)
	if err != nil {
		log.Println("Error decoding response:", err)
		return nil, err
	}

	return &result, nil
}
