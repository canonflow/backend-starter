CREATE TABLE IF NOT EXISTS permissions (
    `id` VARCHAR(128) NOT NULL,
    `action_id` VARCHAR(128) NOT NULL,
    `resource` VARCHAR(100) NOT NULL,
    `description` VARCHAR(255),
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uq_permissions_resource_action` (`resource`, `action_id`),
    FOREIGN KEY (`action_id`) REFERENCES actions(`id`) ON DELETE RESTRICT
)