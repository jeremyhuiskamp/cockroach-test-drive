// Package model inappropriately mixes the domain and
// persistence models.
package model

import (
	"database/sql"
	"errors"

	"github.com/cockroachdb/cockroach-go/crdb"
	"github.com/jmoiron/sqlx"
)

type Account struct {
	ID     int
	Amount uint
}

func (a Account) Balance() uint {
	return a.Amount
}

type AccountRepository struct {
	db *sqlx.DB
}

func NewAccountRepository(db *sql.DB) *AccountRepository {
	return &AccountRepository{sqlx.NewDb(db, "postgres")}
}

func (r AccountRepository) OpenAccount() (Account, error) {
	var a Account
	err := r.db.Get(&a, "insert into account (amount) values (0) returning *")
	return a, err
}

func (r AccountRepository) GetAccount(id int) (Account, error) {
	var a Account
	err := r.db.Get(&a, "select * from account where id = $1", id)
	return a, err
}

func (r AccountRepository) FindRandomAccount() (Account, error) {
	var a Account
	// Will not scale beyond tables that can easily be held in memory:
	err := r.db.Get(&a, "select * from account order by random() limit 1")
	return a, err
}

func (r AccountRepository) Deposit(amount uint, acct Account) error {
	return crdb.ExecuteTx(r.db.DB, func(tx *sql.Tx) error {
		_, err := tx.Exec("insert into deposit (account, amount) values ($1, $2)",
			acct.ID, amount)
		// TODO: cleaner handling for accounts that don't exist?
		if err != nil {
			return err
		}

		_, err = tx.Exec("update account set amount = amount + $1 where id = $2",
			amount, acct.ID)
		return err
	})
}

func (r AccountRepository) Withdraw(amount uint, acct Account) error {
	return crdb.ExecuteTx(r.db.DB, func(tx *sql.Tx) error {
		_, err := tx.Exec("insert into withdrawal (account, amount) values ($1, $2)",
			acct.ID, amount)
		// TODO: cleaner handling for accounts that don't exist?
		if err != nil {
			return err
		}

		_, err = tx.Exec("update account set amount = amount - $1 where id = $2",
			amount, acct.ID)
		// TODO: cleaner handling for overdraft?
		return err
	})
}

func (r AccountRepository) Transfer(amount uint, from, to Account) error {
	return crdb.ExecuteTx(r.db.DB, func(tx *sql.Tx) error {
		_, err := tx.Exec("insert into transfer (debtor, creditor, amount) values ($1, $2, $3)",
			from.ID, to.ID, amount)
		// TODO: cleaner handling for accounts that don't exist?
		if err != nil {
			return err
		}

		_, err = tx.Exec("update account set amount = amount - $1 where id = $2",
			amount, from.ID)
		// TODO: cleaner handling for overdraft?
		if err != nil {
			return err
		}

		_, err = tx.Exec("update account set amount = amount + $1 where id = $2",
			amount, to.ID)
		return err
	})
}

func (r AccountRepository) Validate(acct Account) error {
	// select sum of withdrawals, deposits, debts and credits
	// select current amount
	// compare to make sure they're correct
	return errors.New("not yet implemented")
}
