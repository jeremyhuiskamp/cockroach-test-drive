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

		err = accts.Withdraw(50, *acct)
		if err == nil {
			t.Fatal("shouldn't be able to overdraft")
		}
		// TODO: assert on error type/content
		t.Logf("anticipated error from overdrafting: %#v", err)

		err = accts.Deposit(50, *acct)
		if err != nil {
			t.Fatal(err)
		}

		if balance := getBalance(acct.ID); balance != 50 {
			t.Fatalf("unexpected balance: %d", balance)
		}

		err = accts.Withdraw(10, *acct)
		if err != nil {
			t.Fatal(err)
		}

		if balance := getBalance(acct.ID); balance != 40 {
			t.Fatalf("unexpected balance: %d", balance)
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
