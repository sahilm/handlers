package handler

import (
	"net/http"
	"runtime"
)

type Stack struct {
	File       string
	LineNumber int
	FuncName   string
}

type RecoveryFunc func(w http.ResponseWriter, req *http.Request, panicMessage interface{}, stackTrace []Stack)

type RecoveryHandler struct {
	OnRecoveryFunc RecoveryFunc
	Next           http.Handler
}

func (rh RecoveryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			if rh.OnRecoveryFunc != nil {
				rh.OnRecoveryFunc(w, r, err, stackTrace())
			}
		}
	}()
	rh.Next.ServeHTTP(w, r)
}

func stackTrace() []Stack {
	var traces []Stack
	pc := make([]uintptr, 100)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	for {
		frame, more := frames.Next()
		traces = append(traces, Stack{
			frame.File,
			frame.Line,
			frame.Func.Name(),
		})
		if !more {
			break
		}
	}
	return traces
}
