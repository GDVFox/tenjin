CREATE TABLE IF NOT EXISTS `vote` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `employee_id` int(11) NOT NULL,
    `post_id` int(11) NULL,
    `comment_id` int(11) NULL,
    `delta` tinyint(1) NOT NULL,

    PRIMARY KEY (`id`),
    FOREIGN KEY (`employee_id`) REFERENCES `employee`(`person_id`),
    FOREIGN KEY (`post_id`) REFERENCES `employee_post`(`id`),
    FOREIGN KEY (`comment_id`) REFERENCES `comment`(`id`),
    UNIQUE (`employee_id`, `post_id`),
    UNIQUE (`employee_id`, `comment_id`),

    CONSTRAINT `check_delta_vote` CHECK (
        `delta` BETWEEN -1 AND 1
    ),

    CONSTRAINT `check_vote_post_comment` CHECK (
        CASE WHEN `post_id` IS NULL THEN 0 ELSE 1 END +
        CASE WHEN `comment_id` IS NULL THEN 0 ELSE 1 END = 1
    )
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
