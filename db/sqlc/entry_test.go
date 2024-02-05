package db

import (
	"context"
	"database/sql"
	"testing"
	"time"
	"github.com/fsobh/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T) Entry {

	account1 := createRandomAccount(t) 

	arg:= CreateEntryParams{
		AccountID  : account1.ID,
		Amount	   : util.RandomMoney(),
	}

	entry,err := testQueries.CreateEntry(context.Background(), arg)

	require.NoError(t,err)
	require.NotEmpty(t, entry)

	require.Equal(t,  arg.AccountID, entry.AccountID)
	require.Equal(t,  arg.Amount, entry.Amount)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry
}

func createAccountEntry(t *testing.T,  accountID int64) Entry {

	arg:= CreateEntryParams{
		AccountID  : accountID,
		Amount	   : util.RandomMoney(),
	}

	entry,err := testQueries.CreateEntry(context.Background(), arg)

	require.NoError(t,err)
	require.NotEmpty(t, entry)

	require.Equal(t,  arg.AccountID, entry.AccountID)
	require.Equal(t,  arg.Amount, entry.Amount)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry
}
func TestCreateEntry(t *testing.T) {
	createRandomEntry(t)
}
func TestDeleteEntry(t *testing.T){
	entry1 := createRandomEntry(t)
	err := testQueries.DeleteEntry(context.Background(), entry1.ID)
	require.NoError(t,err)
	entry2, err := testQueries.GetEntryByID(context.Background(),entry1.ID)
	require.Error(t,err)
	require.EqualError(t,err,sql.ErrNoRows.Error())
	require.Empty(t, entry2)
}

func TestListEntriesByAccountID(t *testing.T){

	account1 := createRandomAccount(t)

	for i := 0; i < 10; i++ {
		//so we have a variety of entries to test against
		createAccountEntry(t,account1.ID)
		createRandomEntry(t)
	}

	// total 10 acc - this gets last 5
	arg := GetAllEntriesByAccountIDParams{
		AccountID: account1.ID,
		Limit: 5,
		Offset: 5,
	}
	entries, err := testQueries.GetAllEntriesByAccountID(context.Background(),arg)
	
	require.NoError(t,err)
	require.Len(t,entries,5)

	for _, entry:= range entries {

		require.NotEmpty(t,entry)
	}

}

func TestGetEntry(t *testing.T){
	entry1 := createRandomEntry(t)
	entry2, err := testQueries.GetEntryByID(context.Background(),entry1.ID)

	require.NoError(t,err)
	require.NotEmpty(t,entry2)



	require.Equal(t,  entry1.ID, entry2.ID)
	require.Equal(t,  entry1.AccountID, entry2.AccountID)
	require.Equal(t,  entry1.Amount, entry2.Amount)

	require.WithinDuration(t,entry1.CreatedAt, entry2.CreatedAt, time.Second)
}

func TestListEntries(t *testing.T){

	for i := 0; i < 10; i++ {
		//so we have a variety of entries to test against
		createRandomEntry(t)
	}

	// total 10 acc - this gets last 5
	arg := ListEntryParams{
		Limit: 5,
		Offset: 5,
	}
	entries, err := testQueries.ListEntry(context.Background(),arg)
	
	require.NoError(t,err)
	require.Len(t,entries,5)

	for _, entry:= range entries {

		require.NotEmpty(t,entry)
	}
}


