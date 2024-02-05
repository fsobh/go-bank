package db

import (
	"context"
	"database/sql"
	"testing"
	"time"
	"github.com/fsobh/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T) Transfer {

	account1 := createRandomAccount(t) 
	account2 := createRandomAccount(t)

	arg:= CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID: account2.ID,
		Amount	   : util.RandomMoney(),
	}

	transfer,err := testQueries.CreateTransfer(context.Background(), arg)

	require.NoError(t,err)
	require.NotEmpty(t, transfer)

	require.Equal(t,  arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t,  arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t,  arg.Amount, transfer.Amount)

	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)

	return transfer
}
func createAccountTransfer(t *testing.T, fromAccountID int64, toAccountID int64 ) Transfer {

	arg := CreateTransferParams{}
	if fromAccountID == 0 && toAccountID == 0 {
		
		require.FailNow(t, "`From` and `To` account ID's not provided : `createAccountTransfer`")
	
	}else if fromAccountID > 0 && toAccountID == 0 {
		
		// create a from transfer using the from account ID and a random account for to
			
	account := createRandomAccount(t)

	arg = CreateTransferParams{
		FromAccountID: fromAccountID,
		ToAccountID: account.ID,
		Amount	   : util.RandomMoney(),
	}

	}else if fromAccountID == 0 && toAccountID > 0 {
		
		// create a to transfer using the from account ID and a random account for from
		account := createRandomAccount(t) 

		arg = CreateTransferParams{
			FromAccountID: account.ID,
			ToAccountID: toAccountID,
			Amount	   : util.RandomMoney(),
		}
		

	}else{

		arg = CreateTransferParams{
				FromAccountID: fromAccountID,
				ToAccountID: toAccountID,
				Amount	   : util.RandomMoney(),
			}
	
	}

	transfer,err := testQueries.CreateTransfer(context.Background(), arg)

		require.NoError(t,err)
		require.NotEmpty(t, transfer)

		require.Equal(t,  arg.FromAccountID, transfer.FromAccountID)
		require.Equal(t,  arg.ToAccountID, transfer.ToAccountID)
		require.Equal(t,  arg.Amount, transfer.Amount)

		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	createRandomTransfer(t)
}
func TestDeleteTransfer(t *testing.T){

	transfer := createRandomTransfer(t)
	err := testQueries.DeleteTransfer(context.Background(), transfer.ID)

	require.NoError(t,err)

	transfer2, err := testQueries.GetTransferByID(context.Background(),transfer.ID)

	require.Error(t,err)
	require.EqualError(t,err,sql.ErrNoRows.Error())
	require.Empty(t, transfer2)
}
func TestGetTransfer(t *testing.T){
	transfer1 := createRandomTransfer(t)
	transfer2, err := testQueries.GetTransferByID(context.Background(),transfer1.ID)

	require.NoError(t,err)
	require.NotEmpty(t,transfer2)



	require.Equal(t,  transfer1.ID, transfer2.ID)
	require.Equal(t,  transfer1.FromAccountID, transfer2.FromAccountID)
	require.Equal(t,  transfer1.ToAccountID, transfer2.ToAccountID)
	require.Equal(t,  transfer1.Amount, transfer2.Amount)

	require.WithinDuration(t,transfer1.CreatedAt, transfer2.CreatedAt, time.Second)
}

func TestListTransfers(t *testing.T){
	for i := 0; i < 10; i++ {
			createRandomTransfer(t)
	}

	// total 10 acc - this gets last 5
	arg := ListTransfersParams{
		Limit: 5,
		Offset: 5,
	}

	transfers, err := testQueries.ListTransfers(context.Background(),arg)
	
	require.NoError(t,err)
	require.Len(t,transfers,5)

	for _, transfer:= range transfers {

		require.NotEmpty(t,transfer)
	}
}

//-- name: GetAllTransfersByFromAccountID :many
func TestGetAllTransfersByFromAccountID(t *testing.T){

	account := createRandomAccount(t)

	for i := 0; i < 10; i++ {
		createAccountTransfer(t, account.ID, 0)
	}

	// total 10 acc - this gets last 5
	arg := GetAllTransfersByFromAccountIDParams{
		FromAccountID: account.ID,
		Limit: 5,
		Offset: 5,
	}

	transfers, err := testQueries.GetAllTransfersByFromAccountID(context.Background(),arg)
	
	require.NoError(t,err)
	require.Len(t,transfers,5)

	for _, transfer:= range transfers {

		require.NotEmpty(t,transfer)
	}
}

func TestGetAllTransfersByToAccountID(t *testing.T){

	account := createRandomAccount(t)

	for i := 0; i < 10; i++ {
		createAccountTransfer(t, 0, account.ID)
	}

	// total 10 acc - this gets last 5
	arg := GetAllTransfersByToAccountIDParams{
		ToAccountID : account.ID,
		Limit: 5,
		Offset: 5,
	}

	transfers, err := testQueries.GetAllTransfersByToAccountID(context.Background(),arg)
	
	require.NoError(t,err)
	require.Len(t,transfers,5)

	for _, transfer:= range transfers {

		require.NotEmpty(t,transfer)
	}
}

func TestGetAllTransfersByBetween(t *testing.T){

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	for i := 0; i < 10; i++ {
		createAccountTransfer(t, account1.ID, account2.ID)
	}

	// total 10 acc - this gets last 5
	arg := GetAllTransfersByBetweenParams{
		FromAccountID: account1.ID,
		ToAccountID : account2.ID,
		Limit: 5,
		Offset: 5,
	}

	transfers, err := testQueries.GetAllTransfersByBetween(context.Background(),arg)
	
	require.NoError(t,err)
	require.Len(t,transfers,5)

	for _, transfer:= range transfers {

		require.NotEmpty(t,transfer)
	}
}

