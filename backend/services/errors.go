package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
)

var (
	Err_NOT_FOUND       = errors.New("NOT FOUND")
	Err_ALREADY_EXIST   = errors.New("ALREADY EXIST")
	Err_UNAUTHENTICATED = errors.New("NOT AUTHENTICATED")
	Err_UNAUTHORIZED    = errors.New("NOT AUTHORIZED")
)

func GraphQLError(w http.ResponseWriter, error string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	errBody := map[string]interface{}{
		"errors": []map[string]interface{}{
			{
				"message": error,
			},
		},
	}
	data, err := json.Marshal(errBody)
	if err != nil {
		log.Error().Err(err)
	}
	fmt.Fprintln(w, string(data))
}
