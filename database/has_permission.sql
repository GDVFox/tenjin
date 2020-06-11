CREATE TABLE IF NOT EXISTS `has_permission` (
    `employee_id` int(11) NOT NULL,
    `permission_id` int(255) NOT NULL,
    `date_from` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `date_to` timestamp NULL DEFAULT NULL,

    PRIMARY KEY (`employee_id`, `permission_id`),
    FOREIGN KEY (`employee_id`) REFERENCES `employee`(`person_id`),
    FOREIGN KEY (`permission_id`) REFERENCES `permission`(`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;