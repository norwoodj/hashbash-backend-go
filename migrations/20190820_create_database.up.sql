CREATE DATABASE `hashbash`;

CREATE TABLE `rainbow_table` (
  `id` smallint(6) NOT NULL AUTO_INCREMENT,
  `name` varchar(64) NOT NULL,
  `numChains` bigint(20) NOT NULL,
  `chainLength` bigint(20) NOT NULL,
  `passwordLength` smallint(6) NOT NULL,
  `characterSet` varchar(256) NOT NULL,
  `hashFunction` varchar(16) NOT NULL,
  `finalChainCount` bigint(20) NOT NULL DEFAULT '0',
  `chainsGenerated` bigint(20) NOT NULL DEFAULT '0',
  `status` varchar(24) NOT NULL,
  `generateStarted` datetime DEFAULT NULL,
  `generateCompleted` datetime DEFAULT NULL,
  `created` datetime NOT NULL,
  `lastUpdated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`),
  KEY `numChains` (`numChains`),
  KEY `chainLength` (`chainLength`),
  KEY `passwordLength` (`passwordLength`),
  KEY `hashFunction` (`hashFunction`),
  KEY `finalChainCount` (`finalChainCount`)
);

CREATE TABLE `rainbow_chain` (
  `startPlaintext` varchar(32) NOT NULL,
  `endHash` varchar(128) NOT NULL,
  `rainbowTableId` smallint(6) NOT NULL,
  PRIMARY KEY (`rainbowTableId`,`endHash`),
  CONSTRAINT `rainbow_chain_ibfk_1` FOREIGN KEY (`rainbowTableId`) REFERENCES `rainbow_table` (`id`) ON DELETE CASCADE
);

CREATE TABLE `rainbow_table_search` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `rainbowTableId` smallint(6) NOT NULL,
  `hash` varchar(128) NOT NULL,
  `status` varchar(16) NOT NULL,
  `password` varchar(32) DEFAULT NULL,
  `searchStarted` datetime DEFAULT NULL,
  `searchCompleted` datetime DEFAULT NULL,
  `created` datetime NOT NULL,
  `lastUpdated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `rainbowTableId` (`rainbowTableId`,`hash`),
  KEY `rainbowTableId_2` (`rainbowTableId`,`status`),
  KEY `rainbowTableId_3` (`rainbowTableId`,`password`),
  CONSTRAINT `rainbow_table_search_ibfk_1` FOREIGN KEY (`rainbowTableId`) REFERENCES `rainbow_table` (`id`) ON DELETE CASCADE
);
