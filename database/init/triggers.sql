USE `tenjin`;

DROP TRIGGER IF EXISTS `person_delete_trigger`;
DELIMITER $$
CREATE TRIGGER `person_delete_trigger`
AFTER UPDATE
ON `person` FOR EACH ROW
BEGIN
    IF NEW.`status` = 'deleted' THEN
        UPDATE `works_in` SET `date_to` = CURRENT_TIMESTAMP WHERE `date_to` IS NULL AND `employee_id` = NEW.`id`;
    END IF;
END$$
DELIMITER ;

DROP TRIGGER IF EXISTS `employee_post_delete_trigger`;
DELIMITER $$
CREATE TRIGGER `employee_post_delete_trigger`
AFTER UPDATE
ON `employee_post` FOR EACH ROW
BEGIN
    IF new.`status` = 'deleted' THEN
        UPDATE `attachment` SET `status`='deleted' WHERE `post_id`=new.`id`;
        UPDATE `comment` SET `status`='deleted' WHERE `post_id`=new.`id`;
    END IF;
END$$
DELIMITER ;

DROP TRIGGER IF EXISTS `employee_post_insert_check_trigger`;
DELIMITER $$
CREATE TRIGGER `employee_post_insert_check_trigger`
BEFORE INSERT
ON `employee_post` FOR EACH ROW
BEGIN
    IF NEW.`employee_id` NOT IN (
        SELECT `p`.`id` FROM `employee` AS `e`
        JOIN `person` AS `p` ON `e`.`person_id`=`p`.`id`
        WHERE `p`.`status` = 'active' AND `e`.`person_id` = NEW.`employee_id`
    ) THEN signal sqlstate '45000' set message_text = 'Employee not found!';
    END IF;
END$$
DELIMITER ;

DROP TRIGGER IF EXISTS `solution_insert_check_trigger`;
DELIMITER $$
CREATE TRIGGER `solution_insert_check_trigger`
BEFORE INSERT
ON `solution` FOR EACH ROW
BEGIN
    IF NEW.`task_id` NOT IN (
        SELECT `t`.`post_id` FROM `task` AS `t`
        JOIN `employee_post` AS `p` ON `t`.`post_id`=`p`.`id`
        WHERE `p`.`status` = 'active' AND `t`.`post_id` = NEW.`task_id`
    ) THEN signal sqlstate '45000' set message_text = 'Task not found!';
    END IF;
END$$
DELIMITER ;

DROP TRIGGER IF EXISTS `comment_delete_trigger`;
DELIMITER $$
CREATE TRIGGER `comment_delete_trigger`
AFTER UPDATE
ON `comment` FOR EACH ROW
BEGIN
    IF NEW.`status` = 'deleted' THEN
        UPDATE `attachment` SET `status` = 'deleted' WHERE `comment_id` = NEW.`id`;
    END IF;
END$$
DELIMITER ;

DROP TRIGGER IF EXISTS `comment_insert_check_trigger`;
DELIMITER $$
CREATE TRIGGER `comment_insert_check_trigger`
BEFORE INSERT
ON `comment` FOR EACH ROW
BEGIN
    IF NEW.`employee_id` NOT IN (
        SELECT `p`.`id` FROM `employee` AS `e`
        JOIN `person` AS `p` ON `e`.`person_id`=`p`.`id`
        WHERE `p`.`status` = 'active' AND `e`.`person_id` = NEW.`employee_id`
    ) THEN signal sqlstate '45000' set message_text = 'Employee not found!';
    ELSEIF NEW.`post_id` NOT IN (
        SELECT `p`.`id` FROM `employee_post` AS `p`
        WHERE `status` = 'active' AND `id` = NEW.`post_id`
    ) THEN signal sqlstate '45000' set message_text = 'Post not found!';
    ELSEIF NEW.`parent` IS NOT NULL AND  NEW.`parent` NOT IN (
        SELECT `c`.`id` FROM `comment` AS `c`
        WHERE `status` = 'active' AND `c`.`post_id` = NEW.`post_id` AND `id` = NEW.`parent`
    ) THEN signal sqlstate '45000' set message_text = 'Parent not found!';
    END IF;
END$$
DELIMITER ;

DROP TRIGGER IF EXISTS `vacancy_close_trigger`;
DELIMITER $$
CREATE TRIGGER `vacancy_close_trigger`
AFTER UPDATE
ON `vacancy` FOR EACH ROW
BEGIN
    IF new.`status` = 'closed' THEN
        UPDATE `interview` SET `status`='canceled' WHERE `vacancy_id`=new.`id` AND `status`='waiting';
    END IF;
END$$
DELIMITER ;

DROP TRIGGER IF EXISTS `interview_insert_check_trigger`;
DELIMITER $$
CREATE TRIGGER `interview_insert_check_trigger`
BEFORE INSERT
ON `interview` FOR EACH ROW
BEGIN
    IF NEW.`person_id` NOT IN (
        SELECT `p`.`id` FROM `person` AS `p`
        WHERE `p`.`status` = 'active' AND `p`.`id` = NEW.`person_id`
    ) THEN signal sqlstate '45000' set message_text = 'Person not found!';
    ELSEIF NEW.`interviewer_id` IS NOT NULL AND NEW.`interviewer_id` NOT IN (
        SELECT `e`.`person_id` FROM `employee` AS `e`
        JOIN `person` AS `p` ON `e`.`person_id` = `p`.`id`
        WHERE `p`.`status` = 'active' AND `e`.`person_id` = NEW.`interviewer_id`
    ) THEN signal sqlstate '45000' set message_text = 'Interviewer not found!';
    ELSEIF NEW.`vacancy_id` NOT IN (
        SELECT `v`.`id` FROM `vacancy` AS `v`
        WHERE `v`.`status` = 'active' AND `v`.`id` = NEW.`vacancy_id`
    ) THEN signal sqlstate '45000' set message_text = 'Vacancy not found!';
    END IF;
END$$
DELIMITER ;

DROP TRIGGER IF EXISTS `vote_insert_trigger`;
DELIMITER $$
CREATE TRIGGER `vote_insert_trigger`
AFTER INSERT
ON `vote` FOR EACH ROW
BEGIN
    IF NEW.`post_id` IS NOT NULL THEN
        UPDATE `employee_post` SET `rating` = `rating` + NEW.`delta` WHERE `id` = NEW.`post_id`;
    END IF;
    IF NEW.`comment_id` IS NOT NULL THEN
        UPDATE `comment` SET `rating` = `rating` + NEW.`delta` WHERE `id` = NEW.`comment_id`;
    END IF;
END$$
DELIMITER ;

DROP TRIGGER IF EXISTS `vote_update_trigger`;
DELIMITER $$
CREATE TRIGGER `vote_update_trigger`
AFTER UPDATE
ON `vote` FOR EACH ROW
BEGIN
    IF NEW.`post_id` IS NOT NULL THEN
        UPDATE `employee_post` SET `rating` = `rating` + (NEW.`delta` - OLD.`delta`) WHERE `id` = NEW.`post_id`;
    END IF;
    IF NEW.`comment_id` IS NOT NULL THEN
        UPDATE `comment` SET `rating` = `rating`+ (NEW.`delta` - OLD.`delta`) WHERE `id` = NEW.`comment_id`;
    END IF;
END$$
DELIMITER ;
