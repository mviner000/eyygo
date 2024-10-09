-- +migrate Up
CREATE TABLE IF NOT EXISTS roles (
 id INTEGER PRIMARY KEY AUTOINCREMENT,
 name text NOT NULL
);

CREATE TABLE IF NOT EXISTS accounts (
 id INTEGER PRIMARY KEY AUTOINCREMENT,
 username text NOT NULL,
 email text NOT NULL,
 password text NOT NULL,
 role_id integer NOT NULL,
 created_at TEXT,
 updated_at TEXT,
 CONSTRAINT fk_accounts_role_id FOREIGN KEY (role_id) REFERENCES roles(ID) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS posts (
 id INTEGER PRIMARY KEY AUTOINCREMENT,
 account_id integer NOT NULL,
 content text NOT NULL,
 created_at TEXT,
 updated_at TEXT,
 CONSTRAINT fk_posts_account_id FOREIGN KEY (account_id) REFERENCES accounts(ID) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS likes (
 id INTEGER PRIMARY KEY AUTOINCREMENT,
 post_id integer NOT NULL,
 account_id integer NOT NULL,
 created_at TEXT,
 updated_at TEXT,
 CONSTRAINT fk_likes_post_id FOREIGN KEY (post_id) REFERENCES posts(ID) ON DELETE CASCADE,
 CONSTRAINT fk_likes_account_id FOREIGN KEY (account_id) REFERENCES accounts(ID) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS comments (
 id INTEGER PRIMARY KEY AUTOINCREMENT,
 post_id integer NOT NULL,
 account_id integer NOT NULL,
 content text NOT NULL,
 created_at TEXT,
 updated_at TEXT,
 CONSTRAINT fk_comments_post_id FOREIGN KEY (post_id) REFERENCES posts(ID) ON DELETE CASCADE,
 CONSTRAINT fk_comments_account_id FOREIGN KEY (account_id) REFERENCES accounts(ID) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS followers (
 id INTEGER PRIMARY KEY AUTOINCREMENT,
 account_id integer NOT NULL,
 follower_id integer NOT NULL,
 created_at TEXT,
 updated_at TEXT,
 CONSTRAINT fk_followers_account_id FOREIGN KEY (account_id) REFERENCES accounts(ID) ON DELETE CASCADE,
 CONSTRAINT fk_followers_follower_id FOREIGN KEY (follower_id) REFERENCES accounts(ID) ON DELETE CASCADE
);

-- +migrate Down
DROP TABLE IF EXISTS followers;
DROP TABLE IF EXISTS comments;
DROP TABLE IF EXISTS likes;
DROP TABLE IF EXISTS posts;
DROP TABLE IF EXISTS accounts;
DROP TABLE IF EXISTS roles;