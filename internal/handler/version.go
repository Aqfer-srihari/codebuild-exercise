package handler

import (
	"net/http"
	"question5updation/internal/response"
)

func VersionHandler(r *http.Request, _ map[string]string) response.APIResponse {
	return response.APIResponse{
		Status: http.StatusOK,
		Data:   map[string]string{"version": "1"},
	}
}
