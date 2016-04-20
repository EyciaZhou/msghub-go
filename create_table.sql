CREATE SCHEMA `msghub`;

USE `msghub`;

DROP TABLE  IF EXISTS `picref`;
DROP TABLE  IF EXISTS `msg`;
DROP TABLE  IF EXISTS `pic_task_queue`;
DROP TABLE  IF EXISTS `topic`;
DROP TABLE  IF EXISTS `author`;

CREATE TABLE `topic` (
  `id` varchar(30) CHARACTER SET utf8 NOT NULL,
  `Title` text CHARACTER SET utf8 NOT NULL,
  `LastModify` int(10) unsigned NOT NULL,
  PRIMARY KEY (`id`),
  KEY `LastModify` (`LastModify`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

CREATE TABLE `pic_task_queue` (
  `id` int(10) unsigned zerofill NOT NULL AUTO_INCREMENT,
  `url` varchar(2083) NOT NULL,
  `status` tinyint(3) unsigned NOT NULL,
  `owner` int(10) unsigned NOT NULL,
  `time` timestamp NULL DEFAULT NULL,
  `trytimes` tinyint(4) unsigned NOT NULL DEFAULT '0',
  `nodenum` int(10) unsigned DEFAULT NULL,
  `ext` varchar(5) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `urlIndex` (`url`),
  KEY `queue` (`owner`,`status`,`time`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

CREATE TABLE `author` (
  `id` varchar(30) CHARACTER SET utf8 NOT NULL,
  `coverImg` int(10) unsigned zerofill DEFAULT NULL,
  `name` varchar(30) CHARACTER SET utf8 NOT NULL,
  PRIMARY KEY (`id`),
  CONSTRAINT `coverImg2pid` FOREIGN KEY (`coverImg`) REFERENCES `pic_task_queue` (`id`) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

CREATE TABLE `msg` (
  `id` int(10) unsigned zerofill NOT NULL AUTO_INCREMENT,
  `SnapTime` int(10) unsigned NOT NULL,
  `PubTime` int(10) unsigned NOT NULL,
  `SourceURL` varchar(2083) NOT NULL,
  `Body` MEDIUMTEXT CHARACTER SET utf8 NOT NULL,
  `Title` text CHARACTER SET utf8 NOT NULL,
  `SubTitle` text CHARACTER SET utf8 NOT NULL,
  `CoverImg` int(10) unsigned zerofill DEFAULT NULL,
  `ViewType` tinyint(4) NOT NULL,
  `Frm` varchar(30) CHARACTER SET utf8 NOT NULL,
  `Tag` varchar(30) CHARACTER SET utf8 NOT NULL,
  `Topic` varchar(30) CHARACTER SET utf8 DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `SourceURL` (`SourceURL`),
  KEY `SnapTime` (`SnapTime`),
  KEY `Topic` (`Topic`),
  CONSTRAINT `coverimg2pid` FOREIGN KEY (`CoverImg`) REFERENCES `pic_task_queue` (`id`) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

CREATE TABLE `picref` (
  `Ref` varchar(50) DEFAULT NULL,
  `Description` text CHARACTER SET utf8 NOT NULL,
  `Pixes` varchar(10) DEFAULT NULL,
  `pid` int(10) unsigned zerofill NOT NULL,
  `mid` int(10) unsigned zerofill NOT NULL,
  unique key `mp` (`mid`, `pid`),
  CONSTRAINT `forMid` FOREIGN KEY (`mid`) REFERENCES `msg` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `forPid` FOREIGN KEY (`pid`) REFERENCES `pic_task_queue` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
