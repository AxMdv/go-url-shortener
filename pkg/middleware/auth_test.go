package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSignupMiddleware(t *testing.T) {
	handler := SignUpMiddleware(testCall)
	srv := httptest.NewServer(handler)
	defer srv.Close()

	request, err := http.NewRequest(http.MethodGet, srv.URL, nil)
	require.NoError(t, err)

	t.Run("no cookie in request", func(t *testing.T) {
		resp, err := http.DefaultClient.Do(request)
		require.NoError(t, err)
		defer resp.Body.Close()
		for _, cookie := range resp.Cookies() {
			assert.Equal(t, "user_id", cookie.Name)
			assert.Equal(t, 200, resp.StatusCode)
		}
	})
	t.Run("valid cookie in request", func(t *testing.T) {
		cookie := &http.Cookie{
			Name:  "user_id",
			Value: "30316566613236612d656333362d366533362d383031362d303031353564623832353563f8c1885334e5c310884a9af35b2f189228001d67f322295417f17ff534a88b99",
		}
		request.AddCookie(cookie)
		resp, err := http.DefaultClient.Do(request)
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, 200, resp.StatusCode)
	})
	t.Run("invalid cookie in request", func(t *testing.T) {
		cookie := &http.Cookie{
			Name:  "user_id",
			Value: "3231652a1231231236333362d366533362d383031362d303031353564623832353563f8c1885334e5c310884a9af35b2f189228001d67f322295417f17ff534a88b99",
		}
		request.AddCookie(cookie)
		resp, err := http.DefaultClient.Do(request)
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, 200, resp.StatusCode)
	})
}

func TestValidateUserMiddleware(t *testing.T) {
	handler := ValidateUserMiddleware(testCall)
	srv := httptest.NewServer(handler)
	defer srv.Close()

	t.Run("no cookie in request", func(t *testing.T) {
		request, err := http.NewRequest(http.MethodGet, srv.URL, nil)
		require.NoError(t, err)
		resp, err := http.DefaultClient.Do(request)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	})
	t.Run("valid cookie in request", func(t *testing.T) {
		request, err := http.NewRequest(http.MethodGet, srv.URL, nil)
		require.NoError(t, err)
		cookie := &http.Cookie{
			Name:  "user_id",
			Value: "30316566613236612d656333362d366533362d383031362d303031353564623832353563f8c1885334e5c310884a9af35b2f189228001d67f322295417f17ff534a88b99",
		}
		request.AddCookie(cookie)
		resp, err := http.DefaultClient.Do(request)
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, 200, resp.StatusCode)
	})
	t.Run("invalid cookie in request", func(t *testing.T) {
		request, err := http.NewRequest(http.MethodGet, srv.URL, nil)
		require.NoError(t, err)
		cookie := &http.Cookie{
			Name:  "user_id",
			Value: "3231652asd3232366613236612d656333362d366533362d383031362d303031353564623832353563f8c1885334e5c310884a9af35b2f189228001d67f322295417f17ff534a88b99",
		}
		request.AddCookie(cookie)
		resp, err := http.DefaultClient.Do(request)
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

func testCall(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
