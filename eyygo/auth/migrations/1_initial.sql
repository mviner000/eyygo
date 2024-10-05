-- Create auth_user table
CREATE TABLE auth_user (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    password TEXT NOT NULL,
    last_login DATETIME NULL,
    is_superuser INTEGER NOT NULL CHECK (is_superuser IN (0, 1)),
    username TEXT UNIQUE NOT NULL,
    first_name TEXT DEFAULT '',
    last_name TEXT DEFAULT '',
    email TEXT DEFAULT '',
    is_staff INTEGER NOT NULL CHECK (is_staff IN (0, 1)),
    is_active INTEGER NOT NULL CHECK (is_active IN (0, 1)),
    date_joined DATETIME NOT NULL,
    groups_id INTEGER DEFAULT NULL,
    user_permissions_id INTEGER DEFAULT NULL,
    FOREIGN KEY (groups_id) REFERENCES auth_group(id),
    FOREIGN KEY (user_permissions_id) REFERENCES auth_permission(id)
);

-- Create auth_group table
CREATE TABLE auth_group (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL
);

-- Create auth_permission table
CREATE TABLE auth_permission (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    content_type_id INTEGER NOT NULL,
    codename TEXT NOT NULL,
    UNIQUE (content_type_id, codename),
    FOREIGN KEY (content_type_id) REFERENCES django_content_type(id)
);