package handler

import (
	"net/http"

	"cise.ufl.edu/simpleboard/pkg/response"
)


func Health(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}
