package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	mockdb "github.com/fsobh/simplebank/db/mock"
	db "github.com/fsobh/simplebank/db/sqlc"
	"github.com/fsobh/simplebank/util"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAccountAPI(t *testing.T) {

	account := randomAccount()

	//Creating a table driven test set here to get 100% coverage (don't account for happy scenario only)
	testCases := []struct {
		name          string
		accountID     int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {

				// I expect  the get account method to be called with any Context and the specified ID argument, called only once, expected to return the random account we made with nil error
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(account, nil)

			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				//check response code
				require.Equal(t, http.StatusOK, recorder.Code)

				//check the response body
				requireBodyMatchAccount(t, recorder.Body, account)

			},
		},
		{
			name:      "NotFound",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(db.Account{}, sql.ErrNoRows)

			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {

				//check response code
				require.Equal(t, http.StatusNotFound, recorder.Code)

			},
		},
		{
			name:      "InternalError",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {

				require.Equal(t, http.StatusInternalServerError, recorder.Code)

			},
		},
		{
			name:      "InvalidID",
			accountID: 0,
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {

				require.Equal(t, http.StatusBadRequest, recorder.Code)

			},
		},
	}

	for i := range testCases {

		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {

			ctrl := gomock.NewController(t) // create a new Mock Controller

			defer ctrl.Finish() // very important. It makes sure all functions that were expected to be called were called

			store := mockdb.NewMockStore(ctrl) // create new Mock store

			//build stubs :
			tc.buildStubs(store)

			//Decalre and initialize a New test server

			server := NewServer(store)
			//instead of starting up a real http server, we use the recorder feature of the HTTP package
			recorder := httptest.NewRecorder()

			//specify the path of the API we want to call
			url := fmt.Sprintf("/accounts/%d", tc.accountID)

			request, err := http.NewRequest(http.MethodGet, url, nil)

			require.NoError(t, err)

			// This will send our API request through the router and will store the response in the recorder
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})

	}

}

func randomAccount() db.Account {

	return db.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {

	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotAccount db.Account

	//unmarshal the body into the struct
	err = json.Unmarshal(data, &gotAccount)

	require.NoError(t, err)
	require.Equal(t, account, gotAccount) // make sure the accounts match

}
