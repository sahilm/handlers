package logrushandler

import (
	"net/http"

	"github.com/sahilm/handlers/handler"
	"github.com/sirupsen/logrus"
)

const ISO8601Format = "2006-01-02T15:04:05Z0700"

type RequestsHandler struct {
	Logger *logrus.Entry
	hrh    handler.RequestsHandler
}

func NewRequestsHandler(logger *logrus.Entry, next http.Handler) RequestsHandler {
	rh := RequestsHandler{Logger: logger}
	rh.hrh = handler.NewRequestsHandler(rh.onRequestStart, rh.onRequestEnd, next)
	return rh
}

func (rh RequestsHandler) onRequestStart(r *http.Request, metadata handler.RequestMetadata) {
}

func (rh RequestsHandler) onRequestEnd(w http.ResponseWriter, r *http.Request, metadata handler.RequestMetadata) {
	rh.Logger.WithFields(logrus.Fields{
		"requestTime": metadata.StartTimestamp.Format(ISO8601Format),
		"remoteAddr":  metadata.RemoteAddr,
		"proto":       r.Proto,
		"referer":     r.Referer(),
		"userAgent":   r.UserAgent(),
		"method":      r.Method,
	}).Info(r.RequestURI)
}

func (rh RequestsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rh.hrh.ServeHTTP(w, r)
}
