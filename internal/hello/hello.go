package hello

import "net/http"

func Hello(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(`{"hello": "hej"}`))
	w.WriteHeader(http.StatusOK)
}
