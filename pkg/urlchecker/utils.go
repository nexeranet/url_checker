package urlchecker

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func DecodeRequestBody(request *http.Request, target interface{}) error {
	body, err := io.ReadAll(request.Body)
	if err != nil {
		return err
	}
	defer func() {
		_ = request.Body.Close()
	}()

	return json.Unmarshal(body, target)
}

func Response(response http.ResponseWriter, statusCode int, data interface{}) (n int, err error) {
	if response == nil {
		err = errors.New("response writer is nil")
		return
	}
	response.Header().Set("Content-Type", "application/json")
	body, err := encodeResponseBody(data)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	response.WriteHeader(statusCode)
	return response.Write(body)
}

func BadRequestResponse(response http.ResponseWriter, body interface{}) (err error) {
	_, err = Response(response, http.StatusBadRequest, body)
	return
}

func SuccessResponse(response http.ResponseWriter, body interface{}) (err error) {
	_, err = Response(response, http.StatusOK, body)
	return
}

func encodeResponseBody(data interface{}) ([]byte, error) {
	switch value := data.(type) {
	case nil:
		return nil, nil
	case []byte:
		return value, nil
	default:
		return json.Marshal(value)
	}
}

type ErrorResponse struct {
	Message string `json:"error"`
}
