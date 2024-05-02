package psi

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func Use(middlewares ...func(http.Handler) http.Handler) func(Router) error {
	return func(r Router) error {
		r.Use(middlewares...)
		return nil
	}
}

type loggedResponse struct {
	Inner   http.ResponseWriter
	Status  int
	Written int
}

func (l *loggedResponse) Header() http.Header {
	return l.Inner.Header()
}

func (l *loggedResponse) Write(b []byte) (n int, err error) {
	n, err = l.Inner.Write(b)
	l.Written += n
	return
}

func (l *loggedResponse) WriteHeader(statusCode int) {
	l.Inner.WriteHeader(statusCode)
	l.Status = statusCode
}

func (l *loggedResponse) Flush() {
	if flusher, ok := l.Inner.(http.Flusher); ok {
		flusher.Flush()
	}
}

func LogRecoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(innerW http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		w := &loggedResponse{Inner: innerW}

		recStr := ""
		func() {
			defer func() {
				if reced := recover(); reced != nil {
					w.WriteHeader(http.StatusInternalServerError)
					switch reced := reced.(type) {
					case error:
						recStr = fmt.Sprintf("\t| %T: %s\n", reced, reced.Error())
					default:
						recStr = fmt.Sprintf("\t| %T: %v\n", reced, reced)
					}
				}
			}()
			next.ServeHTTP(w, r)
		}()

		elapsed := time.Since(startTime)
		log.Printf(
			"%s -> %s %q -> %d %s in %s\n%s",
			r.RemoteAddr,
			r.Method,
			r.URL.Path,
			w.Status,
			bytesizeString(w.Written),
			elapsed,
			recStr,
		)
	})
}

func bytesizeString(size int) string {
	postfixes := [...]string{"B", "KB", "MB", "GB", "TB"}
	postfixIdx := 0

	for size >= 1000 {
		size /= 1000
		postfixIdx++
	}

	return fmt.Sprintf("%d%s", size, postfixes[postfixIdx])
}
