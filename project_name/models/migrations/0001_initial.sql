-- +migrate Up
CREATE TABLE IF NOT EXISTS auth_groups (
 id INTEGER PRIMARY KEY AUTOINCREMENT,
 name text NOT NULL,
 CONSTRAINT uk_auth_groups_name UNIQUE (name)
);

CREATE TABLE IF NOT EXISTS auth_permissions (
 id INTEGER PRIMARY KEY AUTOINCREMENT,
 name text NOT NULL,
 content_type_id integer NOT NULL,
 codename text NOT NULL,
 unique_index text
);

CREATE TABLE IF NOT EXISTS auth_users (
 id INTEGER PRIMARY KEY AUTOINCREMENT,
 password text NOT NULL,
 last_login TEXT,
 is_superuser numeric NOT NULL,
 username text NOT NULL,
 first_name text DEFAULT "",
 last_name text DEFAULT "",
 email text DEFAULT "",
 is_staff numeric NOT NULL,
 is_active numeric NOT NULL,
 date_joined TEXT,
 groups_id integer DEFAULT null,
 user_permissions_id integer DEFAULT null,
 CONSTRAINT uk_auth_users_username UNIQUE (username)
);

CREATE TABLE IF NOT EXISTS eyygo_content_types (
 id INTEGER PRIMARY KEY AUTOINCREMENT,
 name text NOT NULL
);

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

CREATE TABLE IF NOT EXISTS followers (
 id INTEGER PRIMARY KEY AUTOINCREMENT,
 account_id integer NOT NULL,
 follower_id integer NOT NULL,
 created_at TEXT,
 updated_at TEXT,
 CONSTRAINT fk_followers_account_id FOREIGN KEY (account_id) REFERENCES accounts(ID) ON DELETE CASCADE,
 CONSTRAINT fk_followers_follower_id FOREIGN KEY (follower_id) REFERENCES accounts(ID) ON DELETE CASCADE
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

-- +migrate Down
DROP TABLE IF EXISTS comments;
DROP TABLE IF EXISTS likes;
DROP TABLE IF EXISTS posts;
DROP TABLE IF EXISTS followers;
DROP TABLE IF EXISTS accounts;
DROP TABLE IF EXISTS roles;
DROP TABLE IF EXISTS eyygo_content_types;
DROP TABLE IF EXISTS auth_users;
DROP TABLE IF EXISTS auth_permissions;
DROP TABLE IF EXISTS auth_groups;