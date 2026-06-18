CREATE TABLE IF NOT EXISTS loans (
    id          SERIAL PRIMARY KEY,
    user_id     INT NOT NULL,
    book_id     INT NOT NULL,
    loaned_at   TIMESTAMP NOT NULL DEFAULT NOW(),
    returned_at TIMESTAMP,
    status      VARCHAR(20) NOT NULL DEFAULT 'active',
    CONSTRAINT chk_status CHECK (status IN ('active', 'returned'))
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_loans_active_user_book
    ON loans(user_id, book_id)
    WHERE status = 'active';

CREATE INDEX IF NOT EXISTS idx_loans_user_id ON loans(user_id);