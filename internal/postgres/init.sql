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

-- Insert demo users
INSERT INTO users (name, email, password_hash) VALUES
('John Doe', 'johndoe@example.com', '$2a$12$doE8XUHGFNy/l3bvn05F6eRpZUv3JOE40w2aIgbG6PZP55YpMXHAe'),
('Jane Doe', 'janedoe@example.com', '$2a$12$gleKw8KS0v7B.2LQ9WVjXuylOaSLHiLmkgQGkLidOR7J1yE2Iq/Gy');

-- Insert demo accounts
INSERT INTO accounts (user_id, account_number, account_name, balance) VALUES
((SELECT id FROM users WHERE email='johndoe@example.com'), '1111111111', 'Account 1', 500.00),
((SELECT id FROM users WHERE email='johndoe@example.com'), '2222222222', 'Account 2', 1000.00);

INSERT INTO accounts (user_id, account_number, account_name, balance) VALUES
((SELECT id FROM users WHERE email='janedoe@example.com'), '3333333333', 'Account 3', 750.00),
((SELECT id FROM users WHERE email='janedoe@example.com'), '4444444444', 'Account 4', 1500.00);
