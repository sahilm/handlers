package handler

import (
	"encoding/hex"
	"math/rand"
	"net/http"
	"runtime"
	"time"
)

type Stack struct {
	File       string
	LineNumber int
	FuncName   string
}

type RecoveryFunc func(w http.ResponseWriter, req *http.Request, panicMessage interface{}, stackTrace []Stack, id string)

func RecoveryHandler(recoveryFunc RecoveryFunc, next http.Handler) http.Handler {
	h := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				recoveryFunc(w, r, err, stackTrace(), randomHexEncodedString(7))
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

func randomHexEncodedString(length int) string {
	rand.Seed(time.Now().UnixNano())
	buff := make([]byte, length)
	rand.Read(buff)
	return hex.EncodeToString(buff)
}
