package handler_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sahim/handlers/handler"
)

var _ = Describe("RecoveryHandler", func() {
	var (
		responseString       string
		panicMessage         string
		nextHandler          http.Handler
		panickingNextHandler http.Handler
		recoveryFunc         handler.RecoveryFunc
		request              *http.Request
		recorder             *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		responseString = "All good!"
		panicMessage = "I died"

		nextHandler = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			_, err := fmt.Fprint(w, responseString)
			Expect(err).ToNot(HaveOccurred())
		})

		panickingNextHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			panic(panicMessage)
		})

		recoveryFunc = func(w http.ResponseWriter, req *http.Request, panicMessage interface{}, stackTrace []handler.Stack, id string) {
			w.WriteHeader(http.StatusInternalServerError)
			for _, s := range stackTrace {
				_, err := fmt.Fprintf(w, "[%s] %s:%d %s()\n", id, s.File, s.LineNumber, s.FuncName)
				Expect(err).ToNot(HaveOccurred())
			}
		}

		request = httptest.NewRequest("GET", "/", nil)
		recorder = httptest.NewRecorder()
	})

	It("should delegate to nextHandler", func() {
		recoveryHandler := handler.RecoveryHandler(recoveryFunc, nextHandler)
		recoveryHandler.ServeHTTP(recorder, request)
		Expect(recorder.Code).To(Equal(http.StatusOK))
		bytes, err := ioutil.ReadAll(recorder.Body)
		Expect(err).ToNot(HaveOccurred())
		Expect(bytes).To(Equal([]byte(responseString)))
	})

	It("should trap panics", func() {
		recoveryHandler := handler.RecoveryHandler(recoveryFunc, panickingNextHandler)
		recoveryHandler.ServeHTTP(recorder, request)
		Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
		bytes, err := ioutil.ReadAll(recorder.Body)
		Expect(err).ToNot(HaveOccurred())
		Expect(string(bytes)).To(ContainSubstring("runtime.gopanic()"))
	})
})
