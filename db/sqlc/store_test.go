package db

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTransferTx(t *testing.T) {

	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	fmt.Println(">> before", account1.Balance, account2.Balance)
	//run n concurrent transfer transactions
	n := 5
	amount := int64(10)

	//channels are designed to connect concurrent go routines
	//channels allow routines to share data with each other without explicit locking
	errs := make(chan error)               // create a channel to detect errors in the below routine
	results := make(chan TransferTxResult) // create a channel to recieve Tx result data from the below routine

	for i := 0; i < n; i++ {
		go func() { // run each transfer in its own routine

			result, err := store.TransferTx(context.Background(), TransferTxParams{

				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})

			errs <- err       // send error back via channel to main routine
			results <- result // send result back via channel

		}()
	}

	//existed := make(map[int]bool)
	//after routines are done, check the results and errors
	for i := 0; i < n; i++ {

		err := <-errs // store errors in a variable in this main routine

		require.NoError(t, err)

		result := <-results

		require.NotEmpty(t, result)

		// Check the transfers
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransferByID(context.Background(), transfer.ID)
		require.NoError(t, err)

		//Check the entries

		// from entry
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntryByID(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		//to entry
		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntryByID(context.Background(), toEntry.ID)
		require.NoError(t, err)

		//todo : check accounts

		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		// todo : check account balances

		diff1 := account1.Balance - fromAccount.Balance //how much is being transferred out of account 1
		diff2 := toAccount.Balance - account2.Balance   //how much money is being received into account 2

		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0) // should be a positive number
		require.True(t, diff1%amount == 0)
		// Since its 5 concurrent transfers of amount:= 10
		//100 - 10 = 90  ---> 90 % 10 = 0
		//90 -  10 = 80  ---> 80 % 10 = 0
		//80 - 10 = 70   ---> 70 % 10 = 0

	}
}
