-- migrations/001_init.sql
CREATE TABLE accounts (
    id BIGSERIAL PRIMARY KEY,
    number VARCHAR(50) UNIQUE NOT NULL,
    balance DECIMAL(15,2) DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE transactions (
    id BIGSERIAL PRIMARY KEY,
    debit_account_id BIGINT REFERENCES accounts(id),
    credit_account_id BIGINT REFERENCES accounts(id),
    amount DECIMAL(15,2) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);