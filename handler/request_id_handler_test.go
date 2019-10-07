package handler_test

import (
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sahilm/handlers/handler"
)

var _ = Describe("UUIDRequestIdHandler", func() {
	var (
		nextHandler http.Handler
	)

	BeforeEach(func() {
		nextHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add(handler.RequestIDHeader, r.Header.Get(handler.RequestIDHeader))
			w.WriteHeader(http.StatusOK)
		})
	})

	When("there is no request id header in request", func() {
		It("should set request ID header to a UUID in the request", func() {
			idHandler := handler.NewUUIDRequestIDHandler(nextHandler)
			idMap := make(map[string]struct{})
			sentinel := struct{}{}
			for i := 0; i < 10; i++ {
				recorder := httptest.NewRecorder()
				request := httptest.NewRequest("GET", "/test", nil)
				idHandler.ServeHTTP(recorder, request)
				Expect(recorder.Code).To(Equal(http.StatusOK))
				idMap[recorder.Header().Get(handler.RequestIDHeader)] = sentinel
			}
			Expect(len(idMap)).To(Equal(10))
		})
	})

	When("there is an existing request ID header", func() {
		It("should not change it", func() {
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "/test", nil)
			request.Header.Set(handler.RequestIDHeader, "abcd")
			idHandler := handler.NewUUIDRequestIDHandler(nextHandler)
			idHandler.ServeHTTP(recorder, request)
			requestID := recorder.Header().Get(handler.RequestIDHeader)
			Expect(requestID).To(Equal("abcd"))
		})
	})
})
