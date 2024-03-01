package metric

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	Register(mux)

	testServer := httptest.NewServer(mux)
	resp, err := http.Get(testServer.URL + "/metrics")
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
