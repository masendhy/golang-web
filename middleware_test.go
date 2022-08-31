package golangweb

import (
	"fmt"
	"net/http"
	"testing"
)

type LogMiddleware struct {
	Handler http.Handler
}

func (middleware *LogMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Before Execute Handler")
	middleware.Handler.ServeHTTP(w, r)
	fmt.Println("After Execute Handler")
}

func TestMiddleware(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Handler executed")
		fmt.Fprint(w, "Hello Middleware")
	})

	LogMiddleware := &LogMiddleware{
		Handler: mux,
	}
	server := http.Server{
		Addr:    "localhost:8888",
		Handler: LogMiddleware,
	}

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
