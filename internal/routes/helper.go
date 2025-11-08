package routes

import (
	"encoding/json"
	"net/http"
)

// func returnJsonError(err error, responseCode int) (int, []byte) {
// 	resp := ErrorResponse{
// 		Error: err.Error(),
// 	}

// 	data, marshalErr := json.Marshal(resp)
// 	if marshalErr != nil {
// 		return http.StatusInternalServerError, nil
// 	}

// 	return responseCode, data
// }

func returnJsonResponse[T any](resp T, responseCode int) (int, []byte) {
	data, marshalErr := json.Marshal(resp)
	if marshalErr != nil {
		return http.StatusInternalServerError, nil
	}

	return responseCode, []byte(data)
}
