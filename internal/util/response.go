package util

import (
	"encoding/json"
	"net/http"
)

func JsonResponse(w http.ResponseWriter, statusCode int, data interface{}) error {
	responseData := map[string]interface{}{
		"statusCode": statusCode,
		"message":    http.StatusText(statusCode),
		"data":       data,
	}

	// Marshal data into JSON
	jsonData, err := json.Marshal(responseData)
	if err != nil {
		return err
	}

	// Set content type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Set status code
	w.WriteHeader(statusCode)

	// Write JSON response
	_, err = w.Write(jsonData)
	if err != nil {
		return err
	}

	return nil
}
