CREATE TABLE IF NOT EXISTS user_roles(
    `user_id` BIGINT NOT NULL,
    `role_id` VARCHAR(128) NOT NULL,
    PRIMARY KEY (`user_id`, `role_id`),
    FOREIGN KEY (`user_id`) REFERENCES users(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`role_id`) REFERENCES roles(`id`) ON DELETE CASCADE
)