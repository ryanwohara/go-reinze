// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runescape

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"
)

func Test_GetUsersOnline(t *testing.T) {
	srv := server()

	msg := getUsersOnline("127.0.0.1")

	if msg != "0" {
		t.Errorf("Expecting 0, received %s", msg)
	}

	msg = getUsersOnline("127.0.0.1")

	if len(msg) > 0 {
		t.Errorf("Expecting length of 0, received %d", len(msg))
	}

	serverShutdown(srv)
}

func Test_GetOsRsPlayerData(t *testing.T) {
	// request, _ := http.NewRequest("GET", "/", nil)
	// response := httptest.NewRecorder()
	// Router().ServeHTTP(response, request)

	srv := server()

	response := getUsersOnline("http://localhost:12345")

	// body := string(getOsRsPlayerData("http://localhost:12345"))

	// if err := srv.Shutdown(context.Background()); err != nil {
	// 	panic(err) // failure/timeout shutting down the server gracefully
	// }

	// if response.Code != 200 {
	// 	t.Errorf("Expected 200 response, got %d", response.Code)
	// }

	fmt.Println(response)

	if response != "<p class='player-count'>There are currently 0 people playing!</p>" {
		t.Errorf("Expected `<p class='player-count'>There are currently 0 people playing!</p>` players, got `%s`", response)
	}

	serverShutdown(srv)
}

func server() *http.Server {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "<p class='player-count'>There are currently 0 people playing!</p>\n")
	})

	srv := &http.Server{Addr: ":12345"}

	srv.ListenAndServe()

	return srv
}

func serverShutdown(srv *http.Server) {
	if err := srv.Shutdown(context.Background()); err != nil {
		//
	}
}
