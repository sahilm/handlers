package handler

import (
	"net/http"
	"strings"
	"time"
)

type RequestMetadata struct {
	StartTimestamp time.Time
	EndTimestamp   time.Time
	RemoteAddr     string
	ExecutionTime  time.Duration
	Status         int
}

type RequestStartFunc func(r *http.Request, metadata RequestMetadata)

type RequestEndFunc func(w http.ResponseWriter, r *http.Request, metadata RequestMetadata)

type clock func() time.Time

type RequestsHandler struct {
	OnRequestStartFunc RequestStartFunc
	OnRequestEndFunc   RequestEndFunc
	Next               http.Handler
	clock              clock
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewRequestsHandler(requestStartFunc RequestStartFunc, requestEndFunc RequestEndFunc, next http.Handler) RequestsHandler {
	return RequestsHandler{
		OnRequestStartFunc: requestStartFunc,
		OnRequestEndFunc:   requestEndFunc,
		Next:               next,
		clock:              time.Now,
	}
}

func (rh RequestsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := rh.clock()

	metadata := RequestMetadata{
		StartTimestamp: start,
		RemoteAddr:     remoteAddr(r),
	}

	rh.OnRequestStartFunc(r, metadata)

	lw := &loggingResponseWriter{w, http.StatusOK}
	rh.Next.ServeHTTP(lw, r)

	end := rh.clock()
	metadata.EndTimestamp = end
	metadata.ExecutionTime = end.Sub(start)
	metadata.Status = lw.statusCode
	rh.OnRequestEndFunc(w, r, metadata)
}

func remoteAddr(r *http.Request) string {
	remoteAddr := r.RemoteAddr
	if index := strings.LastIndex(remoteAddr, ":"); index != -1 {
		remoteAddr = remoteAddr[:index]
	}
	if s := r.Header.Get("X-Forwarded-For"); s != "" {
		remoteAddr = s
	}
	return remoteAddr
}

func (lw *loggingResponseWriter) WriteHeader(code int) {
	lw.statusCode = code
	lw.ResponseWriter.WriteHeader(code)
}
