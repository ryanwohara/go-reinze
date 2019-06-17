// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runescape

import (
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_GetUsersOnline(t *testing.T) {
	getOsRsData := func(url string) []byte {
		return []byte("<p class='player-count'>There are currently 0 people playing!</p>")
	}

	msg := getUsersOnline(getOsRsData)

	if msg != "0" {
		t.Errorf("Expecting 0, received %s", msg)
	}

	getOsRsData = func(url string) []byte {
		return []byte("<p class='player-count'>There are currently 0 people playing OSRS!</p>")
	}

	msg = getUsersOnline(getOsRsData)

	if len(msg) > 0 {
		t.Errorf("Expecting length of 0, received %d", len(msg))
	}
}

func Test_GetOsRsPlayerData(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()
	Router().ServeHTTP(response, request)

	// srv := &http.Server{Addr: "localhost:12345"}
	// srv.ListenAndServe()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "hello world\n")
	})

	// body := string(getOsRsPlayerData("http://localhost:12345"))

	// if err := srv.Shutdown(context.Background()); err != nil {
	// 	panic(err) // failure/timeout shutting down the server gracefully
	// }

	if response.Code != 200 {
		t.Errorf("Expected 200 response, got %d", response.Code)
	}

	fmt.Println(response.Body)

	if response.Body.String() != "<p class='player-count'>There are currently 0 people playing!</p>" {
		t.Errorf("Expected `<p class='player-count'>There are currently 0 people playing!</p>` players, got `%s`", response.Body)
	}
}

func Router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", CreateEndpoint).Methods("GET")
	return router
}

func CreateEndpoint(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("<p class='player-count'>There are currently 0 people playing!</p>"))
}
