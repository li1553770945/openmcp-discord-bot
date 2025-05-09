package handler

import (
	"net/http"
)

func GithubRelease(w http.ResponseWriter, req *http.Request) {
	_, err := w.Write([]byte("Hello World!"))
	if err != nil {
		return
	}
	return
}
