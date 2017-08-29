package model

import (
	"testing"

	"github.com/cockroachdb/cockroach-go/testserver"
	"github.com/mattes/migrate"
	"github.com/mattes/migrate/database/postgres"
	_ "github.com/mattes/migrate/source/file"
)

func TestDepositAndWithdraw(t *testing.T) {
	withTestDB(t, func(accts *AccountRepository) {
		getBalance := func(id int) uint {
			acct, err := accts.GetAccount(id)
			if err != nil {
				t.Fatal(err)
			}
			return acct.Balance()
		}

		acct, err := accts.OpenAccount()
		if err != nil {
			t.Fatal(err)
		}

		err = accts.Withdraw(50, acct)
		if err == nil {
			t.Fatal("shouldn't be able to overdraft")
		}
		// TODO: assert on error type/content
		// https://www.postgresql.org/docs/9.6/static/errcodes-appendix.html
		t.Logf("anticipated error from overdrafting: %#v", err)

		err = accts.Deposit(50, acct)
		if err != nil {
			t.Fatal(err)
		}

		if balance := getBalance(acct.ID); balance != 50 {
			t.Fatalf("unexpected balance: %d", balance)
		}

		err = accts.Withdraw(10, acct)
		if err != nil {
			t.Fatal(err)
		}

		if balance := getBalance(acct.ID); balance != 40 {
			t.Fatalf("unexpected balance: %d", balance)
		}
	})
}

func TestTransfer(t *testing.T) {
	withTestDB(t, func(accts *AccountRepository) {
		getBalance := func(id int) uint {
			acct, err := accts.GetAccount(id)
			if err != nil {
				t.Fatal(err)
			}
			return acct.Balance()
		}

		acct1, err := accts.OpenAccount()
		if err != nil {
			t.Fatal(err)
		}
		err = accts.Deposit(50, acct1)
		if err != nil {
			t.Fatal(err)
		}

		acct2, err := accts.OpenAccount()
		if err != nil {
			t.Fatal(err)
		}

		err = accts.Transfer(10, acct1, acct2)
		if err != nil {
			t.Fatal(err)
		}

		if balance := getBalance(acct1.ID); balance != 40 {
			t.Fatalf("unexpected balance: %d", balance)
		}
		if balance := getBalance(acct2.ID); balance != 10 {
			t.Fatalf("unexpected balance: %d", balance)
		}
	})
}

func TestFindRandomAccount(t *testing.T) {
	numAccts := 3
	numRandomQueries := 100
	withTestDB(t, func(accts *AccountRepository) {
		for i := 0; i < numAccts; i++ {
			_, err := accts.OpenAccount()
			if err != nil {
				t.Fatal(err)
			}
		}

		foundIDs := map[int]int{}
		for i := 0; i < 100; i++ {
			acct, err := accts.FindRandomAccount()
			if err != nil {
				t.Fatal(err)
			}
			foundIDs[acct.ID] += 1
		}

		if len(foundIDs) != numAccts {
			t.Fatalf("unexpected number of accounts found: %d, (%v)",
				len(foundIDs), foundIDs)
		}

		perfectRatio := float32(numRandomQueries) / float32(numAccts)
		maxExpectedPerAcct := perfectRatio * 1.5
		minExpectedPerAcct := perfectRatio / 1.5
		for _, count := range foundIDs {
			fcount := float32(count)
			if fcount > maxExpectedPerAcct || fcount < minExpectedPerAcct {
				t.Fatalf("random selection of accounts not random enough: %v",
					foundIDs)
			}
		}
	})
}

func withTestDB(t *testing.T, f func(*AccountRepository)) {
	db, close := testserver.NewDBForTestWithDatabase(t, "test")
	defer close()

	_, err := db.Exec("CREATE DATABASE IF NOT EXISTS test")
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec("SET DATABASE = test")
	if err != nil {
		t.Fatal(err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		t.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://../migrate",
		"postgres", driver)
	if err != nil {
		t.Fatal(err)
	}

	err = m.Up()
	if err != nil {
		t.Fatal(err)
	}

	f(NewAccountRepository(db))
}
