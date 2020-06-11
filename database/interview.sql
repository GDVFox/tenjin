CREATE TABLE IF NOT EXISTS `interview` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `vacancy_id` int(11) NOT NULL,
    `person_id` int(11) NOT NULL,
    `interviewer_id` int(11) NULL,
    `comment` text NULL,
    `planned_date` datetime NOT NULL,
    `status` enum('waiting', 'canceled', 'completed') NOT NULL DEFAULT 'waiting',
    `total_score` TINYINT(3) NULL CHECK (total_score >= 0 AND total_score <= 100),
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    PRIMARY KEY (`id`),
    FOREIGN KEY (`vacancy_id`) REFERENCES `vacancy`(`id`),
    FOREIGN KEY (`person_id`) REFERENCES `person`(`id`),
    FOREIGN KEY (`interviewer_id`) REFERENCES `employee`(`person_id`),

    UNIQUE KEY (`id`, `vacancy_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
