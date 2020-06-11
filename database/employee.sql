CREATE TABLE IF NOT EXISTS `employee` (
    `person_id` int(11) NOT NULL,
    `password` blob NOT NULL,
    `email` varchar(320) NOT NULL,
    `hired_at` date NOT NULL,

    PRIMARY KEY (`person_id`),
    FOREIGN KEY (`person_id`) REFERENCES `person`(`id`),
    UNIQUE (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;