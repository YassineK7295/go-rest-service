
### Database Creation ###
DROP DATABASE IF EXISTS membership_service;

CREATE DATABASE membership_service;

USE membership_service;

### user Table Creation ###
CREATE TABLE user(
	id INT NOT NULL AUTO_INCREMENT,
    first_name VARCHAR(32) NOT NULL,
    last_name VARCHAR(32) NOT NULL,
    user_id VARCHAR(64) NOT NULL,
    PRIMARY KEY (id)
);

ALTER TABLE user
ADD UNIQUE `uniq_user_id` (user_id);

### group Table Creation ###
CREATE TABLE `group`(
	id INT NOT NULL AUTO_INCREMENT,
    name VARCHAR(64) NOT NULL,
    PRIMARY KEY (id)
);

ALTER TABLE `group` ADD UNIQUE `uniq_name` (name);

### membership Table Creation ###
CREATE TABLE `membership`(
	id INT NOT NULL AUTO_INCREMENT,
    group_id INT NOT NULL,
    user_id INT NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (group_id) REFERENCES `group`(id),
    FOREIGN KEY (user_id) REFERENCES user(id)
);

ALTER TABLE `membership` ADD UNIQUE `uniq_group_id_user_id` (`group_id`, `user_id`);

### Store Procedures ###

DELIMITER //
CREATE PROCEDURE get_user(
	IN user_id VARCHAR(64)
)
BEGIN
	SELECT *
    FROM `user` AS U
    WHERE U.user_id = user_id;
END //

CREATE PROCEDURE get_user_membership(
	IN user_id int
)
BEGIN
	SELECT G.*
    FROM `membership` M
    INNER JOIN `group` G
		ON M.group_id = G.id
        AND M.user_id = user_id;
END //

CREATE PROCEDURE get_group(
	IN group_name VARCHAR(256)
)
BEGIN
	SELECT *
    FROM `group`
    WHERE name = group_name;
END //

CREATE PROCEDURE get_group_membership(
	IN group_id INT
)
BEGIN
	SELECT U.*
    FROM user U
    INNER JOIN membership M
		ON U.id = M.user_id
	INNER JOIN `group` G
		ON M.group_id = G.id
        AND M.group_id = group_id;
END //

CREATE PROCEDURE ins_user(
	IN first_name VARCHAR(32),
    IN last_name VARCHAR(32),
    IN user_id VARCHAR(64)
)
BEGIN

	INSERT INTO `user` (first_name, last_name, user_id)
    VALUES (first_name, last_name, user_id);

    SELECT U.id FROM `user` AS U WHERE U.user_id = user_id;

END //

CREATE PROCEDURE ins_membership(
	IN user_id INT,
    # comma delimited list of group names
    IN group_names TEXT
)
BEGIN
	### dynamic sql used to insert multiple rows with one call
	### this can lead to a sql injection
	### limiting user permissions and validating input values should provide more security
	SET @sql_stmt = CONCAT('
    INSERT INTO `membership` (group_id, user_id)
    SELECT G.id, ', user_id, '
    FROM `group` AS G
    WHERE `name` IN (', group_names, ');');

    PREPARE stmt FROM @sql_stmt;
    EXECUTE stmt;
    DEALLOCATE PREPARE stmt;

END //

CREATE PROCEDURE del_user(
    IN user_id VARCHAR(64)
)
BEGIN
    DECLARE id INT;
    
	DECLARE EXIT HANDLER FOR SQLEXCEPTION
    BEGIN
		ROLLBACK;
		RESIGNAL;
	END;
    
    SET id = (SELECT U.id FROM `user` AS U WHERE U.user_id = user_id);
    IF id IS NULL THEN
        SIGNAL SQLSTATE '45000'
            SET MESSAGE_TEXT = 'user does not exist', MYSQL_ERRNO = 3000;
	END IF;
    
    START TRANSACTION;
		DELETE M
		FROM membership M
		WHERE M.user_id = id;

		DELETE U
		FROM `user` AS U
		WHERE U.id = id;
	COMMIT;
END //

CREATE PROCEDURE upd_user(
	IN first_name VARCHAR(32),
    IN last_name VARCHAR(32),
    IN user_id VARCHAR(32)
)
BEGIN
	DECLARE id INT;
    SET id = (SELECT U.id FROM `user` AS U WHERE U.user_id = user_id);
    
    IF id IS NULL THEN
        SIGNAL SQLSTATE '45000'
            SET MESSAGE_TEXT = 'user does not exist', MYSQL_ERRNO = 3000;
	END IF;
    
    UPDATE `user` AS U
    SET U.first_name = first_name,
		U.last_name = last_name
    WHERE U.id = id;
    
    SELECT id;
END //

CREATE PROCEDURE upd_membership(
	IN user_id INT,
    # comma delimited list of groups names
    IN group_names TEXT
)
BEGIN
	DECLARE EXIT HANDLER FOR SQLEXCEPTION
    BEGIN
		ROLLBACK;
		RESIGNAL;
	END;
    
    IF (SELECT U.id FROM `user` AS U WHERE U.id = user_id) IS NULL THEN
        SIGNAL SQLSTATE '45000'
            SET MESSAGE_TEXT = 'user does not exist', MYSQL_ERRNO = 3000;
	END IF;
    
    START TRANSACTION;
		SET @sql_stmt = CONCAT( '
		DELETE M
		FROM membership AS M
		WHERE M.user_id = ', user_id, ';');
		
		PREPARE stmt FROM @sql_stmt;
		EXECUTE stmt;
		DEALLOCATE PREPARE stmt;
		
		SET @sql_stmt = CONCAT('INSERT INTO `membership` (group_id, user_id)
		SELECT id, ', user_id, '
		FROM `group`
		WHERE `name` IN (', group_names, ');');

		PREPARE stmt FROM @sql_stmt;
		EXECUTE stmt;
		DEALLOCATE PREPARE stmt;
	COMMIT;
END //

CREATE PROCEDURE ins_group(
	IN group_name VARCHAR(256)
)
BEGIN
	INSERT INTO `group` (name)
    VALUES (group_name);
END //

CREATE PROCEDURE upd_group_membership(
	IN group_name VARCHAR(256),
    # Comma delimited list
    IN user_ids TEXT
)
BEGIN
	DECLARE group_id INT;
    
	DECLARE EXIT HANDLER FOR SQLEXCEPTION
    BEGIN
		ROLLBACK;
		RESIGNAL;
	END;
    
    SET group_id = (SELECT G.id FROM `group` AS G WHERE G.name = group_name);
    
    IF group_id IS NULL THEN
        SIGNAL SQLSTATE '45000'
            SET MESSAGE_TEXT = 'group does not exist', MYSQL_ERRNO = 3000;
	END IF;
    
    START TRANSACTION;
		SET @sql_stmt = CONCAT( '
		DELETE M
		FROM membership AS M
		WHERE M.group_id = ', group_id, ';');
		
		PREPARE stmt FROM @sql_stmt;
		EXECUTE stmt;
		DEALLOCATE PREPARE stmt;
		
		SET @sql_stmt = CONCAT('INSERT INTO membership (group_id, user_id)
		SELECT ', group_id, ', id
		FROM `user`
		WHERE user_id IN (', user_ids, ');');
        
		PREPARE stmt FROM @sql_stmt;
		EXECUTE stmt;
		DEALLOCATE PREPARE stmt;
	COMMIT;
END //

CREATE PROCEDURE del_group(
	IN group_name VARCHAR(256)
)
BEGIN
	DECLARE group_id INT;
    
	DECLARE EXIT HANDLER FOR SQLEXCEPTION
    BEGIN
		ROLLBACK;
		RESIGNAL;
	END;
    
    SET group_id = (SELECT G.id FROM `group` AS G WHERE G.name = group_name);
    
    IF group_id IS NULL THEN
        SIGNAL SQLSTATE '45000'
            SET MESSAGE_TEXT = 'group does not exist', MYSQL_ERRNO = 3000;
	END IF;
    
    START TRANSACTION;
		DELETE M
        FROM membership M
        WHERE M.group_id = group_id;
    
		DELETE
		FROM `group`
		WHERE id = group_id;
	COMMIT;
END //

DELIMITER ;

GRANT ALL ON *.* TO 'username'@'%';

FLUSH PRIVILEGES;