# cmd/server
<!-- В данной директории будет содержаться код Сервера, который скомпилируется в бинарное приложение.
package main -->

import (
	"net/http"
)

type UrlMap map[string] string

func main() {
    if err := run(); err != nil {
        panic(err)
    }
}

func get_short_url(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
        w.WriteHeader(http.StatusBadRequest)
        return
    }
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "tetx/plan")
	_, _ = w.Write([]byte(`
      {
        "response": {
          "text": "http://localhost:8080/EwHXdJfB"
        }
      }
    `))	
}

func get_url(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
        w.WriteHeader(http.StatusBadRequest)
        return
    }
	w.WriteHeader(http.StatusTemporaryRedirect)
	w.Header().Set("Content-Type", "tetx/plan")
	w.Header().Set("Location", "https://practicum.yandex.ru/")
}

func run() error {
	mux := http.NewServeMux()
	mux.HandleFunc(`/update/{type}/{value}`, get_short_url)
	mux.HandleFunc(`/{id}`, get_url)
    return http.ListenAndServe(`:8080`, http.HandlerFunc(webhook))
}