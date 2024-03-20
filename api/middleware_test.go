package api

import (
	"fmt"
	"github.com/fsobh/simplebank/token"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// #3.1
func addAuthorization(
	t *testing.T,
	request *http.Request,
	maker token.Maker,
	authorizationType string,
	username string,
	duration time.Duration,
) {

	tokenAuth, err := maker.CreateToken(username, duration)
	require.NoError(t, err)

	authorizationHeader := fmt.Sprintf("%s %s", authorizationType, tokenAuth)
	request.Header.Set(authorizationHeaderKey, authorizationHeader)
}

// #1
func TestAuthMiddleware(t *testing.T) {

	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, maker token.Maker) // this is to set up the auth header for the request
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		// #3
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, maker token.Maker) {
				// # 3.2
				addAuthorization(t, request, maker, authorizationTypeBearer, "user", time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},

		{
			name: "NoAuthorization",
			setupAuth: func(t *testing.T, request *http.Request, maker token.Maker) {
				//do nothing to test response when no auth header is added
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},

		{
			name: "UnsupportedAuthorization",
			setupAuth: func(t *testing.T, request *http.Request, maker token.Maker) {
				addAuthorization(t, request, maker, "unsupported", "user", time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},

		{
			name: "InvalidAuthorizationFormat",
			setupAuth: func(t *testing.T, request *http.Request, maker token.Maker) {
				addAuthorization(t, request, maker, "", "user", time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},

		{
			name: "ExpiredToken",
			setupAuth: func(t *testing.T, request *http.Request, maker token.Maker) {
				addAuthorization(t, request, maker, authorizationTypeBearer, "user", -time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	// #2 loop through the test case
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			server := newTestServer(t, nil)

			//api route for the sake of testing middleware
			authPath := "/auth"

			server.router.GET(
				authPath,                          // path
				authMiddleWare(server.tokenMaker), // middleware
				func(ctx *gin.Context) { // handler
					ctx.JSON(http.StatusOK, gin.H{})
				},
			)
			// create a new recorder to record the call
			recorder := httptest.NewRecorder()

			request, err := http.NewRequest(http.MethodGet, authPath, nil)
			require.NoError(t, err)

			// Then we add the authorization header to the request
			tc.setupAuth(t, request, server.tokenMaker)

			// serve the route
			server.router.ServeHTTP(recorder, request)

			// check the response
			tc.checkResponse(t, recorder)

		})
	}

}
