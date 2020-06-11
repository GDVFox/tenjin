CREATE TABLE IF NOT EXISTS `works_in` (
    `employee_id` int(11) NOT NULL,
    `department_id` int(11) NOT NULL,
    `appointement_id` int(11) NOT NULL,
    `date_from` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `date_to` timestamp NULL DEFAULT NULL,

    PRIMARY KEY (`employee_id`, `department_id`, `appointement_id`),
    FOREIGN KEY (`employee_id`) REFERENCES `employee`(`person_id`),
    FOREIGN KEY (`department_id`) REFERENCES `department`(`id`),
    FOREIGN KEY (`appointement_id`) REFERENCES `appointement`(`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;