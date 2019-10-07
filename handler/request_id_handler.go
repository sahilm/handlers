package handler

import (
	"net/http"

	"github.com/google/uuid"
)

type IDGenerator func() string

type RequestIDHandler struct {
	IDGenerator IDGenerator
	Next        http.Handler
}

func (ri RequestIDHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get(RequestIDHeader) == "" {
		r.Header.Set(RequestIDHeader, ri.IDGenerator())
	}
	ri.Next.ServeHTTP(w, r)
}

func NewUUIDRequestIDHandler(next http.Handler) RequestIDHandler {
	idGenfn := func() string {
		return uuid.New().String()
	}
	return RequestIDHandler{
		IDGenerator: idGenfn,
		Next:        next,
	}
}
