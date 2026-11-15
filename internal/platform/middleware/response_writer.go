package middleware

import "net/http"

type responseWriter struct {
	http.ResponseWriter

	status    int
	bytes     int
	committed bool
}

func getResponseWriter(w http.ResponseWriter) *responseWriter {
	if writer, ok := w.(*responseWriter); ok {
		return writer
	}
	return &responseWriter{ResponseWriter: w}
}

func (w *responseWriter) WriteHeader(status int) {
	if w.committed {
		return
	}
	w.status = status
	w.committed = true
	w.ResponseWriter.WriteHeader(status)
}

func (w *responseWriter) Write(b []byte) (int, error) {
	if !w.committed {
		w.WriteHeader(http.StatusOK)
	}
	n, err := w.ResponseWriter.Write(b)
	w.bytes += n
	return n, err
}

func (w *responseWriter) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}

func (w *responseWriter) Committed() bool {
	return w.committed
}

func (w *responseWriter) Status() int {
	if w.status == 0 {
		return http.StatusOK
	}
	return w.status
}

func (w *responseWriter) BytesWritten() int {
	return w.bytes
}
