CREATE TABLE IF NOT EXISTS domains (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    provider TEXT NOT NULL,
    status INT NOT NULL,
    created_at TIMESTAMP DEFAULT now (),
    expires_at TIMESTAMP,
    is_deleted BOOL DEFAULT FALSE
);

CREATE UNIQUE INDEX idx_domains_name ON domains (name);

CREATE INDEX idx_domains_status ON domains (status);

CREATE TABLE IF NOT EXISTS mailboxes (
    id SERIAL PRIMARY KEY,
    address TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    domain_id INT REFERENCES domains (id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT now (),
    status INT NOT NULL,
    is_deleted BOOL DEFAULT FALSE
);

CREATE UNIQUE INDEX idx_mailboxes_address ON mailboxes (address);

CREATE INDEX idx_mailboxes_domain_id ON mailboxes (domain_id);

CREATE INDEX idx_mailboxes_status ON mailboxes (status);
