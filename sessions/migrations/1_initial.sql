CREATE TABLE eyygo_session (
    session_key CHAR(40) NOT NULL PRIMARY KEY,
    session_data TEXT NOT NULL,
    expire_date DATETIME NOT NULL
);

CREATE INDEX idx_expire_date ON eyygo_session (expire_date);