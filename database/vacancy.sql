CREATE TABLE IF NOT EXISTS `vacancy` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `department_id` int(11) NOT NULL,
    `appointement_id` int(11) NOT NULL,
    `description` text NULL,
    `priority` enum('low', 'medium', 'high') NOT NULL DEFAULT 'medium',
    `status` enum('active', 'paused', 'closed') NOT NULL DEFAULT 'active',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    PRIMARY KEY (`id`),
    FOREIGN KEY (`department_id`) REFERENCES `department`(`id`),
    FOREIGN KEY (`appointement_id`) REFERENCES `appointement`(`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;