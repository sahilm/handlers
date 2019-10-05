package handler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("RequestsHandler", func() {
	var (
		handler        RequestsHandler
		nextHandler    http.Handler
		startFunc      RequestStartFunc
		endFunc        RequestEndFunc
		times          []time.Time
		responseString string
		recorder       *httptest.ResponseRecorder
	)
	BeforeEach(func() {
		responseString = "foo"
		nextHandler = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			_, err := fmt.Fprint(w, responseString)
			Expect(err).ToNot(HaveOccurred())
		})
		startTime, err := time.Parse(time.RFC3339Nano, "2019-10-05T21:04:05.123+00:00")
		Expect(err).ToNot(HaveOccurred())
		endTime := startTime.Add(100 * time.Millisecond)
		times = []time.Time{
			startTime,
			endTime,
		}
		recorder = httptest.NewRecorder()
	})

	It("should invoke callbacks with correct metadata", func() {
		var (
			startFuncCalled int
			endFuncCalled   int
		)

		startFunc = func(r *http.Request, metadata RequestMetadata) {
			startFuncCalled++
			Expect(metadata).To(MatchFields(IgnoreExtras, Fields{
				"StartTimestamp": Equal(times[0]),
				"RemoteAddr":     Equal("127.0.0.1"),
				"ExecutionTime":  BeNumerically("==", 0),
			}))
		}

		endFunc = func(w http.ResponseWriter, r *http.Request, metadata RequestMetadata) {
			endFuncCalled++
			Expect(metadata).To(MatchAllFields(Fields{
				"StartTimestamp": Equal(times[0]),
				"EndTimestamp":   Equal(times[1]),
				"RemoteAddr":     Equal("127.0.0.1"),
				"ExecutionTime":  Equal(100 * time.Millisecond),
			}))
		}

		handler = RequestsHandler{
			OnRequestStartFunc: startFunc,
			OnRequestEndFunc:   endFunc,
			Next:               nextHandler,
			clock:              fakeClock(times),
		}

		request := httptest.NewRequest("GET", "/test", nil)
		request.RemoteAddr = "127.0.0.1:443"
		handler.ServeHTTP(recorder, request)
		Expect(startFuncCalled).To(Equal(1))
		Expect(endFuncCalled).To(Equal(1))
	})

	It("should delegate to next handler", func() {
		handler = RequestsHandler{
			OnRequestStartFunc: func(r *http.Request, metadata RequestMetadata) {},
			OnRequestEndFunc:   func(w http.ResponseWriter, r *http.Request, metadata RequestMetadata) {},
			Next:               nextHandler,
			clock:              fakeClock(times),
		}
		request := httptest.NewRequest("GET", "/test", nil)
		handler.ServeHTTP(recorder, request)
		Expect(recorder.Code).To(Equal(http.StatusOK))
		bytes, err := ioutil.ReadAll(recorder.Body)
		Expect(err).ToNot(HaveOccurred())
		Expect(bytes).To(Equal([]byte(responseString)))
	})

	It("should parse X-Forwarded-For header", func() {
		handler = RequestsHandler{
			OnRequestStartFunc: func(r *http.Request, metadata RequestMetadata) {
				Expect(metadata).To(MatchFields(IgnoreExtras, Fields{
					"RemoteAddr": Equal("192.168.0.1:443"),
				}))
			},
			OnRequestEndFunc: func(w http.ResponseWriter, r *http.Request, metadata RequestMetadata) {
				Expect(metadata).To(MatchFields(IgnoreExtras, Fields{
					"RemoteAddr": Equal("192.168.0.1:443"),
				}))
			},
			Next:  nextHandler,
			clock: fakeClock(times),
		}
		request := httptest.NewRequest("GET", "/test", nil)
		request.RemoteAddr = "127.0.0.1:443"
		request.Header.Add("X-Forwarded-For", "192.168.0.1:443")
		handler.ServeHTTP(recorder, request)
	})
})

func fakeClock(times []time.Time) clock {
	return func() time.Time {
		t := times[0]
		times = times[1:]
		return t
	}
}
