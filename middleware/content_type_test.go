package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

func TestContentTypeMiddleWare(t *testing.T) {
	muxer := mux.NewRouter()
	muxer.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	n := negroni.New()
	n.Use(NewContentType())
	n.UseHandler(muxer)

	server := httptest.NewServer(n)
	defer server.Close()

	client := new(http.Client)
	req, err := http.NewRequest("GET", server.URL, nil)

	if err != nil {
		t.Error(err)
	}

	response, err := client.Do(req)
	if err != nil {
		t.Error(err)
	}

	contentType := response.Header.Get(HttpContentTypeKey)

	if contentType != HttpContentTypeValue {
		t.Errorf("Expected %v, got %v", HttpContentTypeValue, contentType)
	}
}
