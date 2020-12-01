CREATE USER `fss`@`%` IDENTIFIED BY 'password';
GRANT INSERT,SELECT,UPDATE,DELETE ON `fss_db`.* TO `fss`@`%`;

CREATE DATABASE IF NOT EXISTS `fss_db`;

CREATE TABLE IF NOT EXISTS `fss_db`.`users` (
    `user_id`          BIGINT UNSIGNED         NOT NULL,
    `password`         VARCHAR(60)             NOT NULL,
    `locale`           VARCHAR(16),
    `role`             BOOLEAN                 NOT NULL DEFAULT 0,
    `created_at`       TIMESTAMP               NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `last_login`       TIMESTAMP               NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `deleted_at`       TIMESTAMP,
 
    PRIMARY KEY ( `user_id` )
);

CREATE TABLE IF NOT EXISTS `fss_db`.`user_profiles` (
    `user_id`          BIGINT UNSIGNED         NOT NULL,
    `name`             VARCHAR(32)             NOT NULL,
    `avatar`           VARCHAR(34)             NOT NULL,
    `email`            VARCHAR(128)            NOT NULL,

    PRIMARY KEY ( `user_id` )
);

CREATE TABLE IF NOT EXISTS `fss_db`.`file_types` (
    `filetype_id`      INT UNSIGNED            NOT NULL,
    `name`             VARCHAR(128)            NOT NULL,

    PRIMARY KEY ( `filetype_id` )
);

CREATE TABLE IF NOT EXISTS `fss_db`.`files` (
    `file_id`          BIGINT UNSIGNED         NOT NULL,
    `name`             VARCHAR(128)            NOT NULL,
    `user_id`          VARCHAR(32)             NOT NULL,
    `filetype_id`      INT UNSIGNED            NOT NULL,
    `size`             BIGINT UNSIGNED,
    `hash`             CHAR(60),
    `parent`           VARCHAR(60),
    `status`           SMALLINT UNSIGNED       NOT NULL DEFAULT 0,
    `version`          TINYINT(16) UNSIGNED    NOT NULL,
    `created_at`       TIMESTAMP               NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`       TIMESTAMP               NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `deleted_at`       TIMESTAMP,

    PRIMARY KEY ( `file_id` )
);

CREATE TABLE IF NOT EXISTS `fss_db`.`shared` (
    `shared_id`        BIGINT UNSIGNED         NOT NULL,
    `file_id`          BIGINT UNSIGNED         NOT NULL,
    `permission`       TINYINT(3) UNSIGNED     NOT NULL,
    `shared_user`      VARCHAR(32),
    `shared_at`        TIMESTAMP               NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`       TIMESTAMP               NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `deleted_at`       TIMESTAMP,

    PRIMARY KEY ( `shared_id` )
);

BULK INSERT `fss_db`.`shared`
FROM '/docker-entrypoint-initdb.d/filetype.csv'
WITH (
    FIRSTROW = 2,
    FIELDTERMINATOR = ",",
    ROWTERMINATOR = "\n"
);

LOAD DATA INFILE "/docker-entrypoint-initdb.d/filetype.csv"
    INTO TABLE `fss_db`.`file_types`
    FIELDS TERMINATED BY ','
    LINES TERMINATED BY '\n'
    IGNORE 1 LINES;