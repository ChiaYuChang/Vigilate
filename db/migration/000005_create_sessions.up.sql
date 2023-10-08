CREATE TABLE sessions (
    token   text        PRIMARY KEY,
    data    bytea       NOT NULL,
    expiry  timestamp   NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX sessions_expiry_idx ON sessions (expiry);