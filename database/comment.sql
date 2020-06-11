CREATE TABLE IF NOT EXISTS `comment` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `post_id` int(11) NOT NULL,
    `employee_id` int(11) NOT NULL,
    `parent` int(11) NULL DEFAULT NULL,
    `text` text NOT NULL,
    `rating` int(11) NOT NULL DEFAULT 0,
    `status` enum('active', 'deleted') NOT NULL DEFAULT 'active',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    PRIMARY KEY (`id`),
    FOREIGN KEY (`employee_id`) REFERENCES `employee`(`person_id`),
    FOREIGN KEY (`post_id`) REFERENCES `employee_post`(`id`),
    FOREIGN KEY (`parent`) REFERENCES `comment`(`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;