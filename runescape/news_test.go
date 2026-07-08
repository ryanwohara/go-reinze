package runescape

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_WriteNewsToDbWithEmptyNews(t *testing.T) {
	// A failed fetch makes getNews return an empty slice; writing it must
	// be a no-op, not a panic.
	writeNewsToDb(nil, []string{})
}

func Test_GetNewsUnreachableHost(t *testing.T) {
	news := getNews("http://127.0.0.1:1/", "h4 a")

	if len(news) != 0 {
		t.Errorf("Expecting empty slice, received %v", news)
	}
}

func Test_GetNewsNoArticle(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "<html><body><p>maintenance page, no article element</p></body></html>")
	}))
	defer srv.Close()

	news := getNews(srv.URL, "h4 a")

	if len(news) != 0 {
		t.Errorf("Expecting empty slice, received %v", news)
	}
}

func Test_GetNewsNon200(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer srv.Close()

	news := getNews(srv.URL, "h4 a")

	if len(news) != 0 {
		t.Errorf("Expecting empty slice, received %v", news)
	}
}
