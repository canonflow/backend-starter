CREATE TABLE IF NOT EXISTS role_permissions (
    `role_id` VARCHAR(128) NOT NULL,
    `permission_id` VARCHAR(128) NOT NULL,
    PRIMARY KEY (`role_id`, `permission_id`),
    FOREIGN KEY (`role_id`) REFERENCES roles(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`permission_id`) REFERENCES permissions(`id`) ON DELETE CASCADE
)