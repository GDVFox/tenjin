CREATE TABLE IF NOT EXISTS `employee_post` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `employee_id` int(11) NOT NULL,
    `text` text NOT NULL,
    `rating` int(11) NOT NULL DEFAULT 0,
    `status` enum('active', 'deleted') NOT NULL DEFAULT 'active',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    PRIMARY KEY (`id`),
    FOREIGN KEY (`employee_id`) REFERENCES `employee`(`person_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
