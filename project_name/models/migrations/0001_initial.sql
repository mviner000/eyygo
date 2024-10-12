-- +migrate Up
CREATE TABLE IF NOT EXISTS auth_group (
 id INTEGER PRIMARY KEY AUTOINCREMENT,
 name text NOT NULL,
 CONSTRAINT uk_auth_group_name UNIQUE (name)
);

CREATE TABLE IF NOT EXISTS auth_permission (
 id INTEGER PRIMARY KEY AUTOINCREMENT,
 name text NOT NULL,
 content_type_id integer NOT NULL,
 codename text NOT NULL,
 unique_index text
);

CREATE TABLE IF NOT EXISTS auth_user (
 id INTEGER PRIMARY KEY AUTOINCREMENT,
 password text NOT NULL,
 last_login DATETIME,
 is_superuser numeric NOT NULL,
 username text NOT NULL,
 first_name text DEFAULT "",
 last_name text DEFAULT "",
 email text DEFAULT "",
 is_staff numeric NOT NULL,
 is_active numeric NOT NULL,
 date_joined DATETIME,
 group_id integer,
 user_permission_id integer,
 CONSTRAINT uk_auth_user_username UNIQUE (username)
);

CREATE TABLE IF NOT EXISTS eyygo_admin_log (
 id INTEGER PRIMARY KEY AUTOINCREMENT,
 action_time DATETIME,
 object_id text,
 object_repr text NOT NULL,
 action_flag integer NOT NULL,
 change_message text NOT NULL,
 content_type_id integer,
 user_id integer
);

CREATE TABLE IF NOT EXISTS eyygo_content_type (
 id INTEGER PRIMARY KEY AUTOINCREMENT,
 name text NOT NULL
);

CREATE TABLE IF NOT EXISTS eyygo_session (
 session_key TEXT,
 expire_date DATETIME,
 user_id integer NOT NULL,
 auth_token text NOT NULL
);

-- +migrate Down
DROP TABLE IF EXISTS eyygo_session;
DROP TABLE IF EXISTS eyygo_content_type;
DROP TABLE IF EXISTS eyygo_admin_log;
DROP TABLE IF EXISTS auth_user;
DROP TABLE IF EXISTS auth_permission;
DROP TABLE IF EXISTS auth_group;