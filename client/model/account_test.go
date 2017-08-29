package model

import (
	"testing"

	"github.com/cockroachdb/cockroach-go/testserver"
	"github.com/mattes/migrate"
	"github.com/mattes/migrate/database/postgres"
	_ "github.com/mattes/migrate/source/file"
)

func TestDeposit(t *testing.T) {
	withTestDB(t, func(accts *AccountRepository) {
		acct, err := accts.OpenAccount()
		if err != nil {
			t.Fatal(err)
		}

		err = accts.Deposit(50, *acct)
		if err != nil {
			t.Fatal(err)
		}

		acct, err = accts.GetAccount(acct.ID)
		if err != nil {
			t.Fatal(err)
		}

		if acct.Balance() != 50 {
			t.Fatalf("unexpected balance: %d", acct.Balance())
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
