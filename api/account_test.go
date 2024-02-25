package api

import (
	"bytes"
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

	ctrl := gomock.NewController(t) // create a new Mock Controller

	defer ctrl.Finish() // very important. It makes sure all functions that were expected to be called were called

	store := mockdb.NewMockStore(ctrl) // create new Mock store

	//build stubs :

	// I expect  the get account method to be called with any Context and the specified ID argument, called only once, expected to return the random account we made with nil error
	store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(account, nil)

	//Decalre and initialize a New test server

	server := NewServer(store)
	//instead of starting up a real http server, we use the recorder feature of the HTTP package
	recorder := httptest.NewRecorder()

	//specify the path of the API we want to call
	url := fmt.Sprintf("/accounts/%d", account.ID)

	request, err := http.NewRequest(http.MethodGet, url, nil)

	require.NoError(t, err)

	// This will send our API request through the router and will store the response in the recorder
	server.router.ServeHTTP(recorder, request)

	//check response code
	require.Equal(t, http.StatusOK, recorder.Code)

	//check the response body
	requireBodyMatchAccount(t, recorder.Body, account)

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
