CREATE SCHEMA `usr`;

USE `usr`;

DROP TABLE IF EXISTS `_user`;

CREATE TABLE `_user` (
    `id` int(10) unsigned zerofill NOT NULL AUTO_INCREMENT,
    `username`  varchar(16) CHARACTER SET latin1 NOT NULL,
    `email` varchar(50) CHARACTER SET latin1 NOT NULL,
    `nickname` varchar(40) CHARACTER SET utf8 NOT NULL,
    `pwd` varbinary(64) NOT NULL,
    `salt` varbinary(10) NOT NULL,
    `master` int(1) NOT NULL DEFAULT 0,
    PRIMARY KEY (`id`),
    UNIQUE KEY `username` (`username`),
    UNIQUE KEY `email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;