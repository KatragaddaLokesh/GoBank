package db

import (
	"context"
	"fmt"
	"testing"
	//"github.com/KatragaddaLokesh/Go_Bank/utils"
	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {

	store := NewStore(conn)
	acc1 := CreateRandomAccount(t)
	acc2 := CreateRandomAccount(t)
	fmt.Println(">> Before:", acc1.Balance, acc2.Balance)

	n := 10
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	// run n concurrent transfer transaction
	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: acc1.ID,
				ToAccountID:   acc2.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()
	}

	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		res := <-results
		require.NotEmpty(t, res)

		transfer := res.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, acc1.ID, transfer.FromAccountID)
		require.Equal(t, acc2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)

		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = testQueries.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		toEntry := res.ToEntry
		require.NotEmpty(t, toEntry)

		require.Equal(t, acc2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)

		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)

		require.NoError(t, err)

		fromEntry := res.FromEntry
		require.NotEmpty(t, fromEntry)

		require.Equal(t, acc1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)

		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)

		require.NoError(t, err)

		//check Balance
		FromAcc := res.FromAccount
		require.NotEmpty(t, FromAcc)
		require.Equal(t, acc1.ID, FromAcc.ID)

		ToAcc := res.ToAccount
		require.NotEmpty(t, ToAcc)
		require.Equal(t, acc2.ID, ToAcc.ID)
		//requIRE.EQUAL(T, ACC2.BALANCE, TOACC.BALANCE)

		fmt.Println(">> tx:", FromAcc.Balance, ToAcc.Balance)

		diff1 := acc1.Balance - FromAcc.Balance
		diff2 := ToAcc.Balance - acc2.Balance

		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	updateAcc, err := testQueries.GetAccount(context.Background(), acc1.ID)
	require.NoError(t, err)

	updateAcc1, err := testQueries.GetAccount(context.Background(), acc2.ID)
	require.NoError(t, err)

	fmt.Println(">> after:", updateAcc.Balance, updateAcc1.Balance)
	require.Equal(t, acc1.Balance-int64(n)*amount, updateAcc.Balance)
	require.Equal(t, acc2.Balance+int64(n)*amount, updateAcc1.Balance)
}

func TestTransferTxDeadLock(t *testing.T) {

	store := NewStore(conn)
	acc1 := CreateRandomAccount(t)
	acc2 := CreateRandomAccount(t)
	fmt.Println(">> Before:", acc1.Balance, acc2.Balance)

	n := 10
	amount := int64(10)

	errs := make(chan error)

	// run n concurrent transfer transaction
	for i := 0; i < n; i++ {
		FromAcc := acc1.ID
		ToAcc := acc2.ID

		if i%2 == 1 {
			FromAcc = acc2.ID
			ToAcc = acc1.ID
		}

		go func() {
			_, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: FromAcc,
				ToAccountID:   ToAcc,
				Amount:        amount,
			})

			errs <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

	}

	updateAcc, err := testQueries.GetAccount(context.Background(), acc1.ID)
	require.NoError(t, err)

	updateAcc1, err := testQueries.GetAccount(context.Background(), acc2.ID)
	require.NoError(t, err)

	fmt.Println(">> after:", updateAcc.Balance, updateAcc1.Balance)
	require.Equal(t, acc1.Balance, updateAcc.Balance)
	require.Equal(t, acc2.Balance, updateAcc1.Balance)
}
