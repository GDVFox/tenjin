CREATE TABLE IF NOT EXISTS `task` (
    `post_id` int(11) NOT NULL,
    `title` varchar(255) NOT NULL,

    PRIMARY KEY (`post_id`),
    FOREIGN KEY (`post_id`) REFERENCES `employee_post`(`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;