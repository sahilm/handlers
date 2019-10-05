package logrushandler_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"time"

	nested "github.com/antonfisher/nested-logrus-formatter"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sahilm/handlers/logrushandler"
	"github.com/sirupsen/logrus"
	logrustest "github.com/sirupsen/logrus/hooks/test"
)

var _ = Describe("LogrusRecoveryHandler", func() {
	var (
		responseString       string
		panicMessage         string
		nextHandler          http.Handler
		panickingNextHandler http.Handler
		request              *http.Request
		recorder             *httptest.ResponseRecorder
		logger               *logrus.Logger
		hook                 *logrustest.Hook
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

		logger = logrus.New()
		logger.SetOutput(GinkgoWriter)
		hook = logrustest.NewLocal(logger)
		logger.SetFormatter(&nested.Formatter{
			TimestampFormat: time.RFC3339,
		})

		request = httptest.NewRequest("GET", "/", nil)
		recorder = httptest.NewRecorder()
	})

	It("should log panics and respond with 500", func() {
		handler := logrushandler.RecoveryHandler{
			Logger: logger.WithFields(logrus.Fields{}),
		}
		handler.Handler(panickingNextHandler).ServeHTTP(recorder, request)
		Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
		Expect(hook.Entries).To(HaveLen(1))
		Expect(hook.LastEntry().Level).To(Equal(logrus.ErrorLevel))
		Expect(hook.LastEntry().Message).To(ContainSubstring("panic"))
	})

	It("should log nothing if there are no panics", func() {
		handler := logrushandler.RecoveryHandler{
			Logger: logger.WithFields(logrus.Fields{}),
		}
		handler.Handler(nextHandler).ServeHTTP(recorder, request)
		Expect(hook.Entries).To(HaveLen(0))
	})

	It("should delegate to next handler", func() {
		handler := logrushandler.RecoveryHandler{
			Logger: logger.WithFields(logrus.Fields{}),
		}
		handler.Handler(nextHandler).ServeHTTP(recorder, request)
		Expect(recorder.Code).To(Equal(http.StatusOK))
		bytes, err := ioutil.ReadAll(recorder.Body)
		Expect(err).ToNot(HaveOccurred())
		Expect(bytes).To(Equal([]byte(responseString)))
	})
})
