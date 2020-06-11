USE `mysql`;

CREATE DATABASE IF NOT EXISTS `tenjin`;
USE `tenjin`;

CREATE TABLE IF NOT EXISTS `department` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `name` varchar(255) NOT NULL,

    PRIMARY KEY (`id`),
    UNIQUE (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `appointement` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `name` varchar(255) NOT NULL,

    PRIMARY KEY (`id`),
    UNIQUE (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `person` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `photo_uri` text NULL,
    `first_name` varchar(128) NOT NULL,
    `last_name` varchar(128) NOT NULL,
    `status` enum('active', 'blocked', 'deleted') NOT NULL DEFAULT 'active',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `employee` (
    `person_id` int(11) NOT NULL,
    `password` blob NOT NULL,
    `email` varchar(320) NOT NULL,
    `hired_at` date NOT NULL,

    PRIMARY KEY (`person_id`),
    FOREIGN KEY (`person_id`) REFERENCES `person`(`id`),
    UNIQUE (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `employee_post` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `employee_id` int(11) NOT NULL,
    `text` text NOT NULL,
    `status` enum('active', 'deleted') NOT NULL DEFAULT 'active',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    PRIMARY KEY (`id`),
    FOREIGN KEY (`employee_id`) REFERENCES `employee`(`person_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


CREATE TABLE IF NOT EXISTS `comment` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `post_id` int(11) NOT NULL,
    `employee_id` int(11) NOT NULL,
    `parent` int(11) NULL DEFAULT NULL,
    `text` text NOT NULL,
    `status` enum('active', 'deleted') NOT NULL DEFAULT 'active',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    PRIMARY KEY (`id`),
    FOREIGN KEY (`employee_id`) REFERENCES `employee`(`person_id`),
    FOREIGN KEY (`post_id`) REFERENCES `employee_post`(`id`),
    FOREIGN KEY (`parent`) REFERENCES `comment`(`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `task` (
    `post_id` int(11) NOT NULL,
    `title` varchar(255) NOT NULL,

    PRIMARY KEY (`post_id`),
    FOREIGN KEY (`post_id`) REFERENCES `employee_post`(`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `solution` (
    `post_id` int(11) NOT NULL,
    `task_id` int(11) NOT NULL,
    `is_approved` boolean NOT NULL DEFAULT false,

    PRIMARY KEY (`post_id`),
    FOREIGN KEY (`post_id`) REFERENCES `employee_post`(`id`),
    FOREIGN KEY (`task_id`) REFERENCES `task`(`post_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


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


CREATE TABLE IF NOT EXISTS `skill` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `name` varchar(255) NOT NULL,

    PRIMARY KEY (`id`),
    UNIQUE (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `permission` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `name` varchar(255) NOT NULL,
    `description` text NOT NULL,

    PRIMARY KEY (`id`),
    UNIQUE (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `vacancy` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `department_id` int(11) NOT NULL,
    `appointement_id` int(11) NOT NULL,
    `description` text NULL,
    `priority` enum('low', 'medium', 'high') NOT NULL DEFAULT 'medium',
    `status` enum('active', 'paused', 'closed') NOT NULL DEFAULT 'active',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    PRIMARY KEY (`id`),
    FOREIGN KEY (`department_id`) REFERENCES `department`(`id`),
    FOREIGN KEY (`appointement_id`) REFERENCES `appointement`(`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `interview` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `vacancy_id` int(11) NOT NULL,
    `person_id` int(11) NOT NULL,
    `comment` text NULL,
    `planned_date` datetime NOT NULL,
    `status` enum('waiting', 'canceled', 'completed') NOT NULL DEFAULT 'waiting',
    `total_score` TINYINT(3) CHECK (total_score >= 0 AND total_score <= 100),
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    PRIMARY KEY (`id`),
    FOREIGN KEY (`vacancy_id`) REFERENCES `vacancy`(`id`),
    FOREIGN KEY (`person_id`) REFERENCES `person`(`id`),

    UNIQUE KEY (`id`, `vacancy_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `task_skill_requirement` (
    `task_id` int(11) NOT NULL,
    `skill_id` int(11) NOT NULL,
    `difficulty` enum('low', 'medium', 'hard', 'unsolvable') NOT NULL DEFAULT 'medium',

    PRIMARY KEY (`task_id`, `skill_id`),
    FOREIGN KEY (`task_id`) REFERENCES `task`(`post_id`),
    FOREIGN KEY (`skill_id`) REFERENCES `skill`(`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `vacancy_skill_requirement` (
    `vacancy_id` int(11) NOT NULL,
    `skill_id` int(11) NOT NULL,
    `difficulty` enum('low', 'medium', 'hard', 'unsolvable') NOT NULL DEFAULT 'medium',

    PRIMARY KEY (`vacancy_id`, `skill_id`),
    FOREIGN KEY (`vacancy_id`) REFERENCES `vacancy`(`id`),
    FOREIGN KEY (`skill_id`) REFERENCES `skill`(`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

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

CREATE TABLE IF NOT EXISTS `has_permission` (
    `employee_id` int(11) NOT NULL,
    `permission_id` int(255) NOT NULL,
    `date_from` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `date_to` timestamp NULL DEFAULT NULL,

    PRIMARY KEY (`employee_id`, `permission_id`),
    FOREIGN KEY (`employee_id`) REFERENCES `employee`(`person_id`),
    FOREIGN KEY (`permission_id`) REFERENCES `permission`(`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

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

