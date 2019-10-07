package logrushandler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/sahilm/handlers/handler"
	"github.com/sirupsen/logrus"
)

type RecoveryHandler struct {
	Logger *logrus.Entry
	hrh    handler.RecoveryHandler
}

func NewRecoveryHandler(logger *logrus.Entry, next http.Handler) RecoveryHandler {
	rh := RecoveryHandler{Logger: logger}
	rh.hrh = handler.RecoveryHandler{
		OnRecoveryFunc: rh.recoveryFunc,
		Next:           next,
	}
	return rh
}

func (rh RecoveryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rh.hrh.ServeHTTP(w, r)
}

func (rh RecoveryHandler) recoveryFunc(w http.ResponseWriter, req *http.Request, panicMessage interface{},
	stackTrace []handler.Stack) {

	w.WriteHeader(http.StatusInternalServerError)

	var sb strings.Builder
	for _, s := range stackTrace {
		_, err := fmt.Fprintf(&sb, "%s:%d %s()\n", s.File, s.LineNumber, s.FuncName)
		if err != nil {
			rh.Logger.WithFields(logrus.Fields{
				"err": err,
			}).Error("failed to print stack trace")
			return
		}
	}

	rh.Logger.WithFields(logrus.Fields{
		"panic": panicMessage,
	}).Error(sb.String())
}
