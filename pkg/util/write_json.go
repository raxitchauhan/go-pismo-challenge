package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gofrs/uuid"
)

var (
	ErrEmptyHTTPStatus   = errors.New("http status must be set to a value other than 0")
	ErrEmptyErrorMessage = errors.New("error message cannot be empty")
)

type (
	ErrorResponse struct {
		Errors []ErrorDescription `json:"errors"`
	}

	ErrorDescription struct {
		ID      string      `json:"id"`
		Code    string      `json:"code"`
		Status  int         `json:"status"`
		Title   string      `json:"title"`
		Details string      `json:"detail"`
		Source  *FieldError `json:"source,omitempty"`
		// Trace   string      `json:"trace,omitempty"`
	}

	// FieldError describes an error for a specific field, usually provided upon the request
	FieldError struct {
		Field   string `json:"field"`
		Message string `json:"message"`
	}
)

// WriteJSON encodes the provided data into JSON format and writes it into the given reader with
// Content-Type set to "application/json; charset=utf-8"
func WriteJSON(w http.ResponseWriter, status int, data interface{}) (err error) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if status != http.StatusNoContent {
		err = json.NewEncoder(w).Encode(data)
	}

	return err
}

// WriteJSONError sets the fields and writes them into the given reader as JSON with
// Content-Type set to "application/json; charset=utf-8"
func WriteJSONError(
	w http.ResponseWriter,
	status int,
	errDesc ErrorDescription,
	sources ...FieldError,
) error {
	if errDesc.Title == "" {
		return ErrEmptyErrorMessage
	}

	var errResps []ErrorDescription
	if len(sources) > 0 {
		errResps = make([]ErrorDescription, 0, len(sources))

		for i := range sources {
			resp, err := configureErrorResponse(errDesc, &sources[i])
			if err != nil {
				return err
			}
			errResps = append(errResps, resp)
		}
	} else {
		resp, err := configureErrorResponse(errDesc, nil)
		if err != nil {
			return err
		}
		errResps = []ErrorDescription{resp}
	}

	return WriteJSON(w, status, ErrorResponse{Errors: errResps})
}

func configureErrorResponse(resp ErrorDescription, source *FieldError) (ErrorDescription, error) {
	if resp.ID == "" {
		id, err := uuid.NewV4()
		if err != nil {
			return ErrorDescription{}, fmt.Errorf("failed to generate error response UUID: %w", err)
		}
		resp.ID = id.String()
	}

	resp.Source = source

	return resp, nil
}
