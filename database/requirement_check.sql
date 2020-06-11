CREATE TABLE IF NOT EXISTS `requirement_check` (
    `vacancy_id` int(11) NOT NULL,
    `interview_id` int(11) NOT NULL,
    `task_id` int(11) NOT NULL,
    `skill_id` int(11) NOT NULL,
    `comment` text NULL,
    `score` TINYINT(3) CHECK (score >= 0 AND score <= 100),

    PRIMARY KEY (`vacancy_id`, `interview_id`, `task_id`, `skill_id`),
    FOREIGN KEY (`interview_id`, `vacancy_id`) REFERENCES `interview`(`id`, `vacancy_id`),
    FOREIGN KEY (`vacancy_id`, `skill_id`) REFERENCES `vacancy_skill_requirement`(`vacancy_id`, `skill_id`),
    FOREIGN KEY (`task_id`, `skill_id`) REFERENCES `task_skill_requirement`(`task_id`, `skill_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;