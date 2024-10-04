CREATE TABLE `auth_user` (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `password` VARCHAR(128) NOT NULL,
    `last_login` DATETIME NULL,
    `is_superuser` TINYINT(1) NOT NULL,
    `username` VARCHAR(150) UNIQUE NOT NULL,
    `first_name` VARCHAR(150) DEFAULT '',
    `last_name` VARCHAR(150) DEFAULT '',
    `email` VARCHAR(254) DEFAULT '',
    `is_staff` TINYINT(1) NOT NULL,
    `is_active` TINYINT(1) NOT NULL,
    `date_joined` DATETIME NOT NULL,
    `groups_id` INT DEFAULT NULL,
    `user_permissions_id` INT DEFAULT NULL,
    FOREIGN KEY (`groups_id`) REFERENCES `auth_group`(`id`),
    FOREIGN KEY (`user_permissions_id`) REFERENCES `auth_permission`(`id`)
);

-- Create auth_group table
CREATE TABLE `auth_group` (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `name` VARCHAR(150) UNIQUE NOT NULL
);

-- Create auth_permission table
CREATE TABLE `auth_permission` (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `name` VARCHAR(255) NOT NULL,
    `content_type_id` INT NOT NULL,
    `codename` VARCHAR(100) NOT NULL,
    UNIQUE (`content_type_id`, `codename`),
    FOREIGN KEY (`content_type_id`) REFERENCES `django_content_type`(`id`)
);
