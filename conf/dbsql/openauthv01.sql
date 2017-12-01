DROP TABLE IF EXISTS `password`;
DROP TABLE IF EXISTS `domain`;
DROP TABLE IF EXISTS `client`;
DROP TABLE IF EXISTS `token`;
DROP TABLE IF EXISTS `project`;
DROP TABLE IF EXISTS `email`;
DROP TABLE IF EXISTS `phone`;
DROP TABLE IF EXISTS `auth`;
DROP TABLE IF EXISTS `function`;
DROP TABLE IF EXISTS `role`;
DROP TABLE IF EXISTS `group`;
DROP TABLE IF EXISTS `mapping`;
DROP TABLE IF EXISTS `user`;
DROP TABLE IF EXISTS `dbmanager`;


CREATE TABLE `user` (
`id` varchar(255) NOT NULL,
`name` varchar(255) NOT NULL,
`enabled` int(1) NOT NULL DEFAULT 0,
`last_active_time` int(64) NULL,
`extra` text NOT NULL DEFAULT '',
`domain_id` varchar(255) NOT NULL,
`password_id` int(16) NULL,
`create_at` int(64) NOT NULL,
`expires_active_days` int(16) NOT NUll,
PRIMARY KEY (`id`) 
)
ENGINE = InnoDB
AUTO_INCREMENT = 0
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_general_ci
COMMENT = '用户表';

CREATE TABLE `project` (
`id` varchar(255) NOT NULL,
`name` varchar(255) NOT NULL,
`description` text NOT NULL DEFAULT '',
`enabled` int(1) NOT NULL DEFAULT 0,
`domain_id` varchar(255) NOT NULL,
`extra` text NOT NULL DEFAULT '',
`create_at` int(64) NOT NULL,
PRIMARY KEY (`id`) 
)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_general_ci
COMMENT = '项目表';

CREATE TABLE `password` (
`id` int(11) NOT NULL AUTO_INCREMENT,
`password` varchar(255) NOT NULL,
`expires_at` int(64) NULL DEFAULT 0,
`create_at` int(64) NOT NULL DEFAULT 0,
`extra` text NOT NULL DEFAULT '',
`update_at` int(64) NOT NULL DEFAULT 0,
PRIMARY KEY (`id`) 
)
ENGINE = InnoDB
AUTO_INCREMENT = 1
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_general_ci
COMMENT = '用户密码表';

CREATE TABLE `token` (
`id` int(11) NOT NULL AUTO_INCREMENT,
`expire_at` int(64) NOT NULL DEFAULT 0,
`token` varchar(255) NOT NULL,
`type` int(2) NOT NULL,
`purpose` int(2) NULL,
`client_id` varchar(64) NULL,
`user_id` varchar(11) NULL,
`extra` text NOT NULL DEFAULT '',
`create_at` int(64) NOT NULL DEFAULT 0,
PRIMARY KEY (`id`) 
)
ENGINE = InnoDB
AUTO_INCREMENT = 1
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_general_ci
COMMENT = '用户访问令牌表';

CREATE TABLE `domain` (
`id` varchar(255) NOT NULL,
`name` varchar(255) NOT NULL UNIQUE,
`display_name` varchar(255) NOT NULL DEFAULT '',
`description` text NOT NULL DEFAULT '',
`enabled` int(1) NOT NULL,
`extra` text NOT NULL DEFAULT '',
`create_at` int(64) NOT NULL DEFAULT 0,
`update_at` int(64) NOT NULL DEFAULT 0,
PRIMARY KEY (`id`)
)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_general_ci
COMMENT = '用户域表';

CREATE TABLE `role` (
`id` varchar(255) NOT NULL,
`name` varchar(255) NOT NULL,
`domain_id` varchar(64) NULL,
`description` text NOT NULL DEFAULT '',
`extra` text NOT NULL DEFAULT '',
`create_at` int(64) NOT NULL DEFAULT 0,
PRIMARY KEY (`id`) 
)
ENGINE = InnoDB
AUTO_INCREMENT = 0
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_general_ci
COMMENT = '用户角色表';

CREATE TABLE `application` (
`id` varchar(64) NOT NULL,
`name` varchar(255) NOT NULL,
`user_id` varchar(64) NULL,
`app_key` varchar(255) NOT NULL,
`app_secret` varchar(255) NOT NULL,
`description` text NOT NULL DEFAULT '',
`type` int(1) NOT NULL,
`redirect_uri` varchar(255) NULL,
`extra` text NOT NULL DEFAULT '',
`create_at` int(64) NOT NULL DEFAULT 0,
PRIMARY KEY (`id`) 
)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_general_ci
COMMENT = '服务开发者凭证表';

CREATE TABLE `auth` (
`id` int(11) NOT NULL AUTO_INCREMENT,
`code` varchar(255) NOT NULL,
`user_id` varchar(255) NOT NULL,
`used` int(1) NOT NULL,
`create_at` int(64) NOT NULL DEFAULT 0,
`expires_at` int(64) NOT NULL DEFAULT 0,
PRIMARY KEY (`id`) 
)
ENGINE = InnoDB
AUTO_INCREMENT = 1
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_general_ci
COMMENT = 'OAuth2授权码表';

CREATE TABLE `group` (
`id` varchar(255) NOT NULL,
`name` varchar(255) NOT NULL,
`domain_id` varchar(255) NOT NULL,
`description` varchar(255) NOT NULL DEFAULT '',
`extra` text NOT NULL DEFAULT '',
`create_at` int(64) NOT NULL DEFAULT 0,
PRIMARY KEY (`id`) 
)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_general_ci
COMMENT = '用户组表';

CREATE TABLE `function` (
`id` int(11) NOT NULL AUTO_INCREMENT,
`name` varchar(255) NOT NULL,
`service` varchar(255) NOT NULL,
`extra` text NOT NULL DEFAULT '',
`create_at` int(64) NOT NULL DEFAULT 0,
PRIMARY KEY (`id`) 
)
ENGINE = InnoDB
AUTO_INCREMENT = 1
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_general_ci
COMMENT = '服务功能表';

CREATE TABLE `email` (
`id` int(11) NOT NULL AUTO_INCREMENT,
`address` varchar(255) NOT NULL,
`primary` int(1) NULL,
`description` text NOT NULL DEFAULT '',
`user_id` varchar(255) NOT NULL,
PRIMARY KEY (`id`) 
)
ENGINE = InnoDB
AUTO_INCREMENT = 1
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_general_ci
COMMENT = '用户邮件表';

CREATE TABLE `phone` (
`id` int(11) NOT NULL AUTO_INCREMENT,
`numbers` int(32) NULL,
`primary` int(1) NULL,
`description` varchar(255) NOT NULL DEFAULT '',
`user_id` varchar(255) NOT NULL,
PRIMARY KEY (`id`) 
)
ENGINE = InnoDB
AUTO_INCREMENT = 1
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_general_ci
COMMENT = '用户电话号码表';

CREATE TABLE `mapping` (
`id` int(11) NOT NULL AUTO_INCREMENT,
`user_id` varchar(255) NULL,
`project_Id` varchar(255) NULL,
`group_id` varchar(255) NULL,
`role_id` varchar(255) NULL,
`function_id` int(11) NULL,
PRIMARY KEY (`id`) 
)
ENGINE = InnoDB
AUTO_INCREMENT = 1
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_general_ci
COMMENT = '关系映射表';

CREATE TABLE `dbmanager` (
`id` int(11) NOT NULL AUTO_INCREMENT,
`version` int(11) NOT NULL,
`description` text NOT NULL DEFAULT '',
PRIMARY KEY (`id`) 
)
ENGINE = InnoDB
AUTO_INCREMENT = 1
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_general_ci
COMMENT = '数据库SQL版本管理表';

/*
初始化成功
*/
INSERT INTO `dbmanager` (`version`, `description`) VALUES ('1', '初始版本');