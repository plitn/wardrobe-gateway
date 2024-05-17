package handler

import "net/http"

type Service interface {
	LoginUser(w http.ResponseWriter, r *http.Request)
}
