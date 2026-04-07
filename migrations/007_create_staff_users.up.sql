CREATE TABLE staff_users (
    id SERIAL PRIMARY KEY,
    account_id INTEGER NOT NULL UNIQUE,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    position TEXT,
    branch_id INTEGER,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_account
        FOREIGN KEY(account_id)
        REFERENCES accounts(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_branch
        FOREIGN KEY(branch_id)
        REFERENCES branches(id)
        ON DELETE SET NULL
);