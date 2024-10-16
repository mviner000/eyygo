CREATE TABLE eyygo_session (
    session_key CHAR(40) NOT NULL PRIMARY KEY,
    expire_date DATETIME NOT NULL
);

CREATE INDEX idx_expire_date ON eyygo_session (expire_date);

ALTER TABLE eyygo_session ADD COLUMN user_id INTEGER NOT NULL;
ALTER TABLE eyygo_session ADD COLUMN auth_token TEXT NOT NULL;