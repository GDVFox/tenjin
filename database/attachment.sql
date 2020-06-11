CREATE TABLE IF NOT EXISTS `attachment` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `post_id` int(11) NULL,
    `comment_id` int(11) NULL,
    `uri` varchar(512) NOT NULL,
    `status` enum('active', 'deleted') NOT NULL DEFAULT 'active',
    `loaded_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (`id`),
    FOREIGN KEY (`post_id`) REFERENCES `employee_post`(`id`),
    FOREIGN KEY (`comment_id`) REFERENCES `comment`(`id`),
    UNIQUE (`uri`, `post_id`),
    UNIQUE (`uri`, `comment_id`),

    CONSTRAINT `check_attachment_post_comment` CHECK (
        CASE WHEN `post_id` IS NULL THEN 0 ELSE 1 END +
        CASE WHEN `comment_id` IS NULL THEN 0 ELSE 1 END = 1
    )
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
