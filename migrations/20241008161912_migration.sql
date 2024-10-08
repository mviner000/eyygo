-- +migrate Up
CREATE TABLE IF NOT EXISTS users (
 id integer PRIMARY KEY AUTOINCREMENT,
 username text NOT NULL,
 email text NOT NULL,
 password text NOT NULL,
 role_id integer NOT NULL,
 created_at TEXT,
 updated_at TEXT,
 CONSTRAINT fk_users_role_id FOREIGN KEY (role_id) REFERENCES roles(ID) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS posts (
 id integer PRIMARY KEY AUTOINCREMENT,
 user_id integer NOT NULL,
 content text NOT NULL,
 created_at TEXT,
 updated_at TEXT,
 CONSTRAINT fk_posts_user_id FOREIGN KEY (user_id) REFERENCES users(ID) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS followers (
 id integer PRIMARY KEY AUTOINCREMENT,
 user_id integer NOT NULL,
 follower_id integer NOT NULL,
 created_at TEXT,
 updated_at TEXT,
 CONSTRAINT fk_followers_user_id FOREIGN KEY (user_id) REFERENCES users(ID) ON DELETE CASCADE,
 CONSTRAINT fk_followers_follower_id FOREIGN KEY (follower_id) REFERENCES users(ID) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS roles (
 id integer PRIMARY KEY AUTOINCREMENT,
 name text NOT NULL
);

CREATE TABLE IF NOT EXISTS likes (
 id integer PRIMARY KEY AUTOINCREMENT,
 post_id integer NOT NULL,
 user_id integer NOT NULL,
 created_at TEXT,
 updated_at TEXT,
 CONSTRAINT fk_likes_post_id FOREIGN KEY (post_id) REFERENCES likes(ID) ON DELETE CASCADE,
 CONSTRAINT fk_likes_user_id FOREIGN KEY (user_id) REFERENCES likes(ID) ON DELETE CASCADE
);

-- +migrate Down
DROP TABLE IF EXISTS followers;
DROP TABLE IF EXISTS likes;
DROP TABLE IF EXISTS posts;
DROP TABLE IF EXISTS roles;
DROP TABLE IF EXISTS users;