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

func (rh RequestsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := rh.clock()

	metadata := RequestMetadata{
		StartTimestamp: start,
		RemoteAddr:     remoteAddr(r),
	}

	rh.OnRequestStartFunc(r, metadata)

	rh.Next.ServeHTTP(w, r)

	end := rh.clock()
	metadata.EndTimestamp = end
	metadata.ExecutionTime = end.Sub(start)
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
