CREATE SCHEMA `usth`;

USE `usth`;

DROP TABLE IF EXISTS `_reply`;
DROP TABLE IF EXISTS `info`;
DROP TABLE IF EXISTS `diggs`;

CREATE TABLE `_reply` (
  `id` int(10) unsigned zerofill NOT NULL AUTO_INCREMENT,
  `_time` int(10) unsigned NOT NULL,
  `author_name` varchar(10) CHARACTER SET utf8 NOT NULL,
  `stu_id` varchar(10) NOT NULL,
  `content` varchar(250) CHARACTER SET utf8 NOT NULL,
  `class_name` varchar(50) CHARACTER SET utf8 NOT NULL,
  `digg` int(10) unsigned NOT NULL,
  `refid` varchar(30),
  `ref_author_id` varchar(30),
  `ref_author` varchar(10) CHARACTER SET utf8,
  `ref_content` varchar(250) CHARACTER SET utf8,

  PRIMARY KEY (`id`),
  KEY `cn` (`class_name`, `_time`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

CREATE TABLE `diggs` (
  `reply_id` int(10) unsigned zerofill NOT NULL,
  `stu_id` varchar(10) NOT NULL,

  UNIQUE KEY `reply_id__stu_id` (`reply_id`, `stu_id`),
  CONSTRAINT `reply_id2` FOREIGN KEY (`reply_id`) REFERENCES `_reply` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

CREATE TABLE `info` (
  `stu_id` varchar(10) NOT NULL,
  `pwd` varchar(20) NOT NULL,
  `author_name` varchar(10) CHARACTER SET utf8 NOT NULL,

  PRIMARY KEY (`stu_id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
