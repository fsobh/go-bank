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
		//txName := fmt.Sprintf("tx %d", i+1) //for fmt.logs
		go func() { // run each transfer in its own routine
			//ctx := context.WithValue(context.Background(), txKey, txName) //for fmt.logs
			result, err := store.TransferTx(context.Background(), TransferTxParams{

				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})

			errs <- err       // send error back via channel to main routine
			results <- result // send result back via channel

		}()
	}

	existed := make(map[int]bool)
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

		// Since its 5 concurrent transfers of amount:= 10
		fmt.Println(">> tx:", fromAccount.Balance, toAccount.Balance)
		diff1 := account1.Balance - fromAccount.Balance // *diff1* (initial account balance) - (final account balance)
		//100 - 90 = 10
		//90  - 80 = 20
		//80  - 70 = 30
		//...

		diff2 := toAccount.Balance - account2.Balance // *diff2* (final account balance) - (initial account balance)
		//60 - 50 = 10
		//70 - 50 = 20
		//80 - 50 = 30
		//...

		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)         // should be a positive number
		require.True(t, diff1%amount == 0) // n % n always = 0

		k := int(diff1 / amount)          // 10/10=1, 20/10=2, 30/10=3...50/10=5 (n=5)
		require.True(t, k >= 1 && k <= n) // so k will always be between 1 & 5
		require.NotContains(t, existed, k)
		existed[k] = true

	}
	// Check the final updated balances

	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println(">> after", updatedAccount1.Balance, updatedAccount2.Balance)
	require.Equal(t, account1.Balance-int64(n)*amount, updatedAccount1.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, updatedAccount2.Balance)

}

func TestTransferTxDeadLock(t *testing.T) {

	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	fmt.Println(">> before", account1.Balance, account2.Balance)

	n := 10
	amount := int64(10)

	errs := make(chan error)

	for i := 0; i < n; i++ {

		fromAccountID := account1.ID
		toAccountID := account2.ID

		//since we want half the transactions to be acc1 --> acc2 and other half acc2 -->acct1
		if i%2 == 1 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}

		go func() {

			_, err := store.TransferTx(context.Background(), TransferTxParams{

				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})

			errs <- err

		}()
	}

	for i := 0; i < n; i++ {

		err := <-errs // store errors in a variable in this main routine
		require.NoError(t, err)
	}
	// Check the final updated balances

	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println(">> after", updatedAccount1.Balance, updatedAccount2.Balance)
	require.Equal(t, account1.Balance, updatedAccount1.Balance)
	require.Equal(t, account2.Balance, updatedAccount2.Balance)

}
