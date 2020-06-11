CREATE TABLE IF NOT EXISTS `vacancy_skill_requirement` (
    `vacancy_id` int(11) NOT NULL,
    `skill_id` int(11) NOT NULL,
    `difficulty` enum('low', 'medium', 'hard', 'unsolvable') NOT NULL DEFAULT 'medium',

    PRIMARY KEY (`vacancy_id`, `skill_id`),
    FOREIGN KEY (`vacancy_id`) REFERENCES `vacancy`(`id`),
    FOREIGN KEY (`skill_id`) REFERENCES `skill`(`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;