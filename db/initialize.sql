CREATE USER `fss`@`%` IDENTIFIED BY 'password';
GRANT INSERT,SELECT,UPDATE,DELETE ON `fss_db`.* TO `fss`@`%`;

CREATE DATABASE IF NOT EXISTS `fss_db`;

CREATE TABLE IF NOT EXISTS `fss_db`.`users` (
    `id`               BIGINT UNSIGNED         NOT NULL,
    `username`         VARCHAR(32)             NOT NULL,
    `avatar`           VARCHAR(34)             NOT NULL,
    `locale`           VARCHAR(16),
    `deleted`          BOOLEAN                 NOT NULL,
    `created_at`       TIMESTAMP               NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `last_login`       TIMESTAMP               NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY ( `id` )
);

CREATE TABLE IF NOT EXISTS `fss_db`.`files` (
    `id`               BIGINT UNSIGNED         NOT NULL,
    `user_id`          VARCHAR(32)             NOT NULL,
    `name`             VARCHAR(128)            NOT NULL,
    `size`             BIGINT UNSIGNED         NOT NULL,
    `hash`             CHAR(60)                NOT NULL,
    `parent`           VARCHAR(60)
    `version`          TINYINT(16) UNSIGNED    NOT NULL,
    `deleted`          BOOLEAN                 NOT NULL,
    `created_at`       TIMESTAMP               NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `last_update`      TIMESTAMP               NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY ( `id` )
);

CREATE TABLE IF NOT EXISTS `fss_db`.`folders` (
    `id`               BIGINT UNSIGNED         NOT NULL,
    `user_id`          VARCHAR(32)             NOT NULL,
    `name`             VARCHAR(128)            NOT NULL,
    `parent`           VARCHAR(60)
    `deleted`          BOOLEAN                 NOT NULL,
    `created_at`       TIMESTAMP               NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `last_update`      TIMESTAMP               NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY ( `id` )
);