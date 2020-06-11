CREATE TABLE IF NOT EXISTS `solution` (
    `post_id` int(11) NOT NULL,
    `task_id` int(11) NOT NULL,
    `is_approved` boolean NOT NULL DEFAULT false,

    PRIMARY KEY (`post_id`),
    FOREIGN KEY (`post_id`) REFERENCES `employee_post`(`id`),
    FOREIGN KEY (`task_id`) REFERENCES `task`(`post_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
