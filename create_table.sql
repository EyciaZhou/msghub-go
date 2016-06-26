drop table picref;
drop table plugin_conf;
drop table topic;
drop table msg;
drop table author;
drop table pic_task_queue;

CREATE DATABASE usr;

CREATE TABLE `_user` (
  `id` int(10) unsigned zerofill NOT NULL AUTO_INCREMENT,
  `username` varchar(16) NOT NULL,
  `email` varchar(50) NOT NULL,
  `nickname` varchar(40) CHARACTER SET utf8 NOT NULL,
  `pwd` varbinary(64) NOT NULL,
  `salt` varbinary(10) NOT NULL,
  `master` int(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `username` (`username`),
  UNIQUE KEY `email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

CREATE database msghub;

CREATE TABLE `pic_task_queue` (
  `id` int(10) unsigned zerofill NOT NULL AUTO_INCREMENT,
  `url` varchar(2083) NOT NULL,
  `status` tinyint(3) unsigned NOT NULL,
  `owner` varchar(40) NOT NULL,
  `time` timestamp NULL DEFAULT NULL,
  `trytimes` tinyint(4) unsigned NOT NULL DEFAULT '0',
  `nodetype` varchar(20) DEFAULT NULL,
  `nodenum` varchar(20) DEFAULT NULL,
  `ext` varchar(20) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `urlIndex` (`url`),
  KEY `queue` (`owner`,`status`,`time`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

CREATE TABLE `author` (
  `id` varchar(64) CHARACTER SET utf8 NOT NULL,
  `coverImg` int(10) unsigned zerofill DEFAULT NULL,
  `name` varchar(30) CHARACTER SET utf8mb4 NOT NULL,
  PRIMARY KEY (`id`),
  KEY `coverImg2pid_author` (`coverImg`),
  CONSTRAINT `coverImg2pid_author` FOREIGN KEY (`coverImg`) REFERENCES `pic_task_queue` (`id`) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

CREATE TABLE `msg` (
  `id` int(10) unsigned zerofill NOT NULL AUTO_INCREMENT,
  `SnapTime` int(10) unsigned NOT NULL,
  `PubTime` int(10) unsigned NOT NULL,
  `SourceURL` varchar(2083) NOT NULL,
  `Body` mediumtext CHARACTER SET utf8mb4 NOT NULL,
  `Title` text CHARACTER SET utf8mb4 NOT NULL,
  `SubTitle` text CHARACTER SET utf8mb4 NOT NULL,
  `CoverImg` int(10) unsigned zerofill DEFAULT NULL,
  `ViewType` tinyint(4) NOT NULL,
  `AuthorId` varchar(64) CHARACTER SET utf8 DEFAULT NULL,
  `Tag` varchar(30) CHARACTER SET utf8mb4 NOT NULL,
  `Topic` varchar(64) CHARACTER SET utf8 DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `SourceURL` (`SourceURL`,`Topic`),
  KEY `ma` (`AuthorId`),
  KEY `SnapTime` (`SnapTime`,`id`) USING BTREE,
  KEY `CoverImg` (`CoverImg`),
  KEY `Topic` (`Topic`,`SnapTime`,`id`) USING BTREE,
  CONSTRAINT `authorid2aid` FOREIGN KEY (`AuthorId`) REFERENCES `author` (`id`) ON DELETE SET NULL ON UPDATE CASCADE,
  CONSTRAINT `coverimg2pid_msg` FOREIGN KEY (`CoverImg`) REFERENCES `pic_task_queue` (`id`) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

CREATE TABLE `plugin_conf` (
  `id` varchar(64) NOT NULL,
  `plugin_type` varchar(20) NOT NULL,
  `config` text CHARACTER SET utf8mb4 NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

CREATE TABLE `topic` (
  `topic_hash` varchar(64) NOT NULL,
  `title` text CHARACTER SET utf8mb4 NOT NULL,
  `uid` int(10) unsigned zerofill NOT NULL,
  UNIQUE KEY `uid_topic_hash` (`topic_hash`,`uid`),
  KEY `uid` (`uid`),
  KEY `topic_hash` (`topic_hash`),
  CONSTRAINT `topic_plugin_conf_id` FOREIGN KEY (`topic_hash`) REFERENCES `plugin_conf` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `topic_uid` FOREIGN KEY (`uid`) REFERENCES `usr`.`_user` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

CREATE TABLE `picref` (
  `Ref` varchar(50) DEFAULT NULL,
  `Description` text CHARACTER SET utf8mb4 NOT NULL,
  `Pixes` varchar(10) DEFAULT NULL,
  `pid` int(10) unsigned zerofill NOT NULL,
  `mid` int(10) unsigned zerofill NOT NULL,
  UNIQUE KEY `mp` (`mid`,`pid`),
  KEY `forPid` (`pid`),
  KEY `forMid` (`mid`),
  CONSTRAINT `forMid` FOREIGN KEY (`mid`) REFERENCES `msg` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `forPid` FOREIGN KEY (`pid`) REFERENCES `pic_task_queue` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=latin1;