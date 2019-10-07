package logrushandler_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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
		responseString = "All good!"
		nextHandler = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			_, err := fmt.Fprint(w, responseString)
			Expect(err).ToNot(HaveOccurred())
		})

		logger = logrus.New()
		logger.SetOutput(GinkgoWriter)
		hook = logrustest.NewLocal(logger)
		request = httptest.NewRequest("GET", "/something/1/else", nil)
		recorder = httptest.NewRecorder()
	})

	It("should log requests and delegate to the next handler", func() {
		loggerEntry := logger.WithFields(logrus.Fields{})
		handler := logrushandler.NewRequestsHandler(loggerEntry, nextHandler)
		handler.ServeHTTP(recorder, request)

		Expect(hook.Entries).To(HaveLen(1))

		Expect(recorder.Code).To(Equal(http.StatusOK))
		bytes, err := ioutil.ReadAll(recorder.Body)
		Expect(err).ToNot(HaveOccurred())
		Expect(bytes).To(Equal([]byte(responseString)))
	})
})
