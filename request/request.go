package request

import (
	"encoding/json"
	"io"
	"net/http"
	"sdk/response"
	"sdk/validation"
	"sdk/model"
)

func Decode[T any](body io.ReadCloser) (T, error){
	var payload T
	err := json.NewDecoder(body).Decode(&payload)
	if err != nil {
		return payload, err
	}

	return  payload, nil
}

func HandleBody[T any](w *http.ResponseWriter, r *http.Request) (*T, error) {
	body, err := Decode[T](r.Body)

	if err != nil {
		response.WriteJSON(*w, model.ErrorResponse{
			Message: err.Error(),
		}, http.StatusBadRequest)
		return nil, err
	}

	err = validation.IsValid(body)

	if err != nil {
		response.WriteJSON(*w, model.ErrorResponse{
			Message: err.Error(),
		}, http.StatusBadRequest)
		return nil, err
	}

	return &body, nil
}