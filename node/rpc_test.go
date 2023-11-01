package node

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPingRoute(t *testing.T) {
	app := NewServer()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	app.Server.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	// assert.Equal(t, )

}
