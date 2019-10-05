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

func RecoveryHandler(recoveryFunc RecoveryFunc, next http.Handler) http.Handler {
	h := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				recoveryFunc(w, r, err, stackTrace())
			}
		}()
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(h)
}

func stackTrace() []Stack {
	var traces []Stack
	for skip := 2; ; skip++ {
		pc, file, lineNumber, ok := runtime.Caller(skip)
		if !ok {
			break
		}
		if file[len(file)-1] == 'c' {
			continue
		}
		f := runtime.FuncForPC(pc)
		traces = append(traces, Stack{
			file,
			lineNumber,
			f.Name(),
		})
	}
	return traces
}
