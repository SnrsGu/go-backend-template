package router

import "net/http"

type StatusProbe struct {
	header     http.Header
	statusCode int
}

func newStatusProbe() *StatusProbe {
	return &StatusProbe{
		header: make(http.Header),
	}
}

func (s *StatusProbe) Header() http.Header {
	return s.header
}

func (s *StatusProbe) WriteHeader(statusCode int) {
	if s.statusCode == 0 {
		s.statusCode = statusCode
	}
}

func (s *StatusProbe) Write(data []byte) (int, error) {
	if s.statusCode == 0 {
		s.WriteHeader(http.StatusOK)
	}
	return len(data), nil
}
