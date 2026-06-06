//go:build no_ui

package main

import "net/http"

func registerUIRoute(*http.ServeMux) {
	// No-op if the UI shouldn't be bundled
}
