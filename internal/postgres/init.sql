CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE accounts (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    account_name VARCHAR(20) NOT NULL,
    account_number VARCHAR(20) UNIQUE NOT NULL,
    balance NUMERIC(15,2) NOT NULL DEFAULT 0.00,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE SET NULL,
    from_account_id INT REFERENCES accounts(id) ON DELETE SET NULL,
    to_account_id INT REFERENCES accounts(id) ON DELETE SET NULL,
    amount NUMERIC(15,2) NOT NULL CHECK (amount > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Insert demo user
INSERT INTO users (name, email, password_hash) VALUES (
    'John Doe',
    'johndoe@example.com',
    '$2a$12$v6Bg98bCJU8LtQRSAsnj5u6EsCJP3zWLtd0N0Vu1Z/WIaZKboeqtG'
);

-- Insert two demo accounts
INSERT INTO accounts (user_id, account_number, account_name, balance) VALUES
((SELECT id FROM users WHERE email='johndoe@example.com'), '1111111111', 'Demo Account 1', 1000.00),
((SELECT id FROM users WHERE email='johndoe@example.com'), '2222222222', 'Demo Account 2', 500.00);
