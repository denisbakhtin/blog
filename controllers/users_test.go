package controllers

import (
	"github.com/denisbakhtin/blog/shared"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserIndex(t *testing.T) {
	req, _ := http.NewRequest("GET", "/admin/users", nil)
	w := httptest.NewRecorder()
	shared.Init()
	UserIndex(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("/admin/users didn't return %v\n", http.StatusOK)
	}
}
