CREATE TABLE IF NOT EXISTS `task_skill_requirement` (
    `task_id` int(11) NOT NULL,
    `skill_id` int(11) NOT NULL,
    `difficulty` enum('low', 'medium', 'hard', 'unsolvable') NOT NULL DEFAULT 'medium',

    PRIMARY KEY (`task_id`, `skill_id`),
    FOREIGN KEY (`task_id`) REFERENCES `task`(`post_id`),
    FOREIGN KEY (`skill_id`) REFERENCES `skill`(`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;