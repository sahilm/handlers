package logrushandler_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	"github.com/sahilm/handlers/logrushandler"
	"github.com/sirupsen/logrus"
	logrustest "github.com/sirupsen/logrus/hooks/test"
)

var _ = Describe("RequestsHandler", func() {
	var (
		responseString string
		nextHandler    http.Handler
		request        *http.Request
		recorder       *httptest.ResponseRecorder
		logger         *logrus.Logger
		hook           *logrustest.Hook
	)

	BeforeEach(func() {
		responseString = "Not found!"
		nextHandler = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			_, err := fmt.Fprint(w, responseString)
			Expect(err).ToNot(HaveOccurred())
		})

		logger = logrus.New()
		logger.SetOutput(GinkgoWriter)
		hook = logrustest.NewLocal(logger)
		request = httptest.NewRequest("GET", "/something/1/else", nil)
		request.Header.Add("User-Agent", "007")
		request.Header.Add("Referer", "test-referer")
		recorder = httptest.NewRecorder()
	})

	It("should log requests and delegate to the next handler", func() {
		loggerEntry := logger.WithFields(logrus.Fields{})
		handler := logrushandler.NewRequestsHandler(loggerEntry, nextHandler)
		handler.ServeHTTP(recorder, request)

		Expect(hook.Entries).To(HaveLen(1))
		logEntry := hook.LastEntry()
		Expect(logEntry.Level).To(Equal(logrus.InfoLevel))
		Expect(logEntry.Data).To(MatchAllKeys(Keys{
			"startTimestamp": Not(BeEmpty()),
			"endTimestamp":   Not(BeEmpty()),
			"runtime":        BeNumerically(">", 0),
			"remoteAddr":     Not(BeEmpty()),
			"status":         Equal(http.StatusNotFound),
			"proto":          Equal("HTTP/1.1"),
			"referer":        Equal("test-referer"),
			"userAgent":      Equal("007"),
			"method":         Equal(request.Method),
		}))
		Expect(logEntry.Message).To(Equal("GET /something/1/else"))

		Expect(recorder.Code).To(Equal(http.StatusNotFound))
		bytes, err := ioutil.ReadAll(recorder.Body)
		Expect(err).ToNot(HaveOccurred())
		Expect(bytes).To(Equal([]byte(responseString)))
	})

	It("should log request IDs if present", func() {
		loggerEntry := logger.WithFields(logrus.Fields{})
		handler := logrushandler.NewRequestsHandler(loggerEntry, nextHandler)
		request.Header.Add("X-Request-Id", "abcd")
		handler.ServeHTTP(recorder, request)

		Expect(hook.Entries).To(HaveLen(1))
		Expect(hook.LastEntry().Data["request-id"]).To(Equal("abcd"))
	})
})
