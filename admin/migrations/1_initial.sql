CREATE TABLE eyygo_admin_log (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    action_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    object_id TEXT,
    object_repr VARCHAR(200) NOT NULL,
    action_flag INTEGER NOT NULL,
    change_message TEXT NOT NULL,
    content_type_id INTEGER,
    user_id INTEGER NOT NULL,

    FOREIGN KEY (content_type_id) REFERENCES eyygo_content_type (id) ON DELETE SET NULL,
    FOREIGN KEY (user_id) REFERENCES auth_user (id) ON DELETE CASCADE
);