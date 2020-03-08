package controllers

import (
	"net/http"
)

// GET
// Authorization: 	Basic
// Params: 			None
// Body: 			None
// Check if a session is enabled
func SessionGet(w http.ResponseWriter, r *http.Request) {
	// todo
}

// GET
// Authorization: 	Basic + check of remote addr
// Params: 			duration
// Body: 			None
func SessionCreate(w http.ResponseWriter, r *http.Request) {
	// todo
}
