-- account is a current snapshot of all account balances
CREATE TABLE account (
	id SERIAL PRIMARY KEY,
	amount INT CHECK (amount >= 0)
);

-- create a bunch of empty accounts to start with:
INSERT INTO account (amount) SELECT 0 FROM generate_series(1,500);

-- transfer is a log of all transfer between accounts
CREATE TABLE transfer (
	id SERIAL PRIMARY KEY,
	debtor INT NOT NULL REFERENCES account (id),
	creditor INT NOT NULL REFERENCES account (id),
	amount int CHECK (amount > 0)
);

-- deposit is a log of all money that enters from outside the system
CREATE TABLE deposit (
	id SERIAL PRIMARY KEY,
	account INT NOT NULL REFERENCES account (id),
	amount int CHECK (amount > 0)
);

-- withdrawal is a log of all money that leaves the system
CREATE TABLE withdrawal (
	id SERIAL PRIMARY KEY,
	account INT NOT NULL REFERENCES account (id),
	amount int CHECK (amount > 0)
);

-- TODO: create a view that calculates expected balance based on the log tables

