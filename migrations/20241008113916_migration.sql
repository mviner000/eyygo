-- +migrate Up
CREATE TABLE IF NOT EXISTS users (
 ID integer KEY AUTOINCREMENT PRIMARY,
 Username text NOT NULL UNIQUE,
 Email text UNIQUE NOT NULL,
 Password text NOT NULL,
 RoleID integer NOT NULL,
 CreatedAt datetime,
 UpdatedAt datetime
);

CREATE TABLE IF NOT EXISTS posts (
 ID integer PRIMARY KEY AUTOINCREMENT,
 Title text NOT NULL,
 Content text NOT NULL,
 UserID integer NOT NULL,
 CreatedAt datetime,
 UpdatedAt datetime
);

CREATE TABLE IF NOT EXISTS followers (
 ID integer AUTOINCREMENT PRIMARY KEY,
 FollowerUserID integer NOT NULL,
 FollowedUserID integer NOT NULL,
 CreatedAt datetime
);

CREATE TABLE IF NOT EXISTS roles (
 ID integer KEY AUTOINCREMENT PRIMARY,
 Name text NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS categories (
 ID integer PRIMARY KEY AUTOINCREMENT,
 Name text NOT NULL UNIQUE
);

-- +migrate Down
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS posts;
DROP TABLE IF EXISTS followers;
DROP TABLE IF EXISTS roles;
DROP TABLE IF EXISTS categories;