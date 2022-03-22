CREATE TABLE `qr` (
	`id` bigint unsigned NOT NULL AUTO_INCREMENT,
	`qr_code` varchar(255)  NOT NULL DEFAULT '' COMMENT '验证码',
	`auth_count` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '验证次数',
	`create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
	`update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	PRIMARY KEY (`id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4;