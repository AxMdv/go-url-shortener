package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithLogging(t *testing.T) {
	handler := WithLogging(testCall)
	srv := httptest.NewServer(handler)
	defer srv.Close()

	t.Run("no cookie in request", func(t *testing.T) {
		request, err := http.NewRequest(http.MethodGet, srv.URL, nil)
		require.NoError(t, err)
		resp, err := http.DefaultClient.Do(request)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, 200, resp.StatusCode)

	})
}
