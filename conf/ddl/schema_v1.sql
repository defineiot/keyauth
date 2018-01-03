DROP TABLE IF EXISTS `passwords`;
DROP TABLE IF EXISTS `domains`;
DROP TABLE IF EXISTS `applications`;
DROP TABLE IF EXISTS `clients`;
DROP TABLE IF EXISTS `tokens`;
DROP TABLE IF EXISTS `projects`;
DROP TABLE IF EXISTS `emails`;
DROP TABLE IF EXISTS `phones`;
DROP TABLE IF EXISTS `auth_codes`;
DROP TABLE IF EXISTS `features`;
DROP TABLE IF EXISTS `roles`;
DROP TABLE IF EXISTS `users_projects_mapping`;
DROP TABLE IF EXISTS `users`;
DROP TABLE IF EXISTS `roles_features_mapping`;
DROP TABLE IF EXISTS `services`;
DROP TABLE IF EXISTS `instances`;
DROP TABLE IF EXISTS `dbmanager`;


CREATE TABLE `domains` (
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

CREATE TABLE `projects` (
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

CREATE TABLE `users` (
`id` varchar(255) NOT NULL,
`name` varchar(255) NOT NULL,
`enabled` int(1) NOT NULL DEFAULT 0,
`last_active_time` int(64) NOT NULL DEFAULT 0,
`extra` text NOT NULL DEFAULT '',
`domain_id` varchar(255) NOT NULL,
`create_at` int(64) NOT NULL,
`expires_active_days` int(64) NOT NUll DEFAULT 0,
`default_project_id` varchar(255) NOT NULL DEFAULT '',
PRIMARY KEY (`id`) 
)
ENGINE = InnoDB
AUTO_INCREMENT = 0
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_general_ci
COMMENT = '用户表';

CREATE TABLE `emails` (
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

CREATE TABLE `phones` (
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

CREATE TABLE `passwords` (
`id` int(11) NOT NULL AUTO_INCREMENT,
`password` varchar(255) NOT NULL,
`expires_at` int(64) NULL DEFAULT 0,
`create_at` int(64) NOT NULL DEFAULT 0,
`extra` text NOT NULL DEFAULT '',
`update_at` int(64) NOT NULL DEFAULT 0,
`user_id` varchar(255) NOT NULL,
PRIMARY KEY (`id`) 
)
ENGINE = InnoDB
AUTO_INCREMENT = 1
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_general_ci
COMMENT = '用户密码表';

CREATE TABLE `users_projects_mapping` (
`id` int(11) NOT NULL AUTO_INCREMENT,
`user_id` varchar(255) NULL,
`project_id` varchar(255) NULL,
PRIMARY KEY (`id`) 
)
ENGINE = InnoDB
AUTO_INCREMENT = 1
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_general_ci
COMMENT = '关系映射表';

CREATE TABLE `applications` (
`id` varchar(255) NOT NULL,
`name` varchar(255) NOT NULL,
`user_id` varchar(64) NULL,
`website` varchar(255) NOT NULL DEFAULT '',
`logo_image` varchar(255) NOT NULL DEFAULT '',
`description` text NOT NULL DEFAULT '',
`create_at` int(64) NOT NULL DEFAULT 0,
`extra` text NOT NULL DEFAULT '',
PRIMARY KEY (`id`) 
)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_general_ci
COMMENT = '用户第三方应用';

CREATE TABLE `clients` (
`id` varchar(255) NOT NULL,
`secret` varchar(255) NOT NULL,
`type` varchar(64) NOT NULL,
`redirect_uri` varchar(255) NULL,
`application_id` varchar(128) NOT NULL DEFAULT '',
`service_id` varchar(128) NOT NULL DEFAULT '',
`extra` text NOT NULL DEFAULT '',
PRIMARY KEY (`id`) 
)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_general_ci
COMMENT = 'OAUTH2客户端凭证表';

CREATE TABLE `auth_codes` (
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

CREATE TABLE `tokens` (
`id` int(11) NOT NULL AUTO_INCREMENT,
`grant_type` varchar(64) NULL,
`access_token` varchar(255) NOT NULL,
`refresh_token` varchar(255) NOT NULL DEFAULT '',
`type` varchar(64) NOT NULL,
`create_at` int(64) NOT NULL DEFAULT 0,
`expire_at` int(64) NOT NULL DEFAULT 0,
`client_id` varchar(64) NULL,
`user_id` varchar(64) NULL,
`domain_id` varchar(255) NOT NULL DEFAULT '',
`project_id` varchar(255) NOT NULL DEFAULT '',
`extra` text NOT NULL DEFAULT '',
PRIMARY KEY (`id`) 
)
ENGINE = InnoDB
AUTO_INCREMENT = 1
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_general_ci
COMMENT = '用户访问令牌表';

CREATE TABLE `roles` (
`id` varchar(255) NOT NULL,
`name` varchar(255) NOT NULL,
`description` text NOT NULL DEFAULT '',
`create_at` int(64) NOT NULL DEFAULT 0,
`extra` text NOT NULL DEFAULT '',
PRIMARY KEY (`id`) 
)
ENGINE = InnoDB
AUTO_INCREMENT = 0
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_general_ci
COMMENT = '用户角色表';

CREATE TABLE `features` (
`id` int(11) NOT NULL AUTO_INCREMENT,
`name` varchar(255) NOT NULL,
`method` varchar(64) NOT NULL,
`endpoint` varchar(256) NOT NULL DEFAULT '',
`description` text NOT NULL DEFAULT '',
`is_deleted` int(2) NOT NULL DEFAULT 0,
`when_deleted_version` varchar(255) NOT NULL DEFAULT '',
`is_added` int(2) NOT NULL DEFAULT 0,
`when_added_version` varchar(255) NOT NULL DEFAULT '',
`create_at` int(64) NOT NULL DEFAULT 0,
`service_id` varchar(255) NOT NULL,
`extra` text NOT NULL DEFAULT '',
PRIMARY KEY (`id`) 
)
ENGINE = InnoDB
AUTO_INCREMENT = 1
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_general_ci
COMMENT = '服务功能表';

CREATE TABLE `roles_features_mapping` (
`id` int(11) NOT NULL AUTO_INCREMENT,
`feature_id` int(64) NOT NULL DEFAULT 0,
`role_id` varchar(128) NOT NULL DEFAULT '',
PRIMARY KEY (`id`) 
)
ENGINE = InnoDB
AUTO_INCREMENT = 1
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_general_ci
COMMENT = '角色与功能映射表';

CREATE TABLE `services` (
`id` varchar(128) NOT NULL AUTO_INCREMENT,
`name` varchar(255) NOT NULL,
`description` text NOT NULL DEFAULT '',
`enabled` int(2) NOT NULL DEFAULT 0,
`status` varchar(255) NOT NULL DEFAULT '',
`status_update_at` int(64) NOT NULL DEFAULT 0,
`version` varchar(255) NOT NULL DEFAULT '',
`create_at` int(64) NOT NULL DEFAULT 0,
`extra` text NOT NULL DEFAULT '',
PRIMARY KEY (`id`) 
)
ENGINE = InnoDB
AUTO_INCREMENT = 1
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_general_ci
COMMENT = '服务表';

CREATE TABLE `instances` (
`id` varchar(128) NOT NULL,
`name` varchar(255) NOT NULL,
`address` varchar(255) NOT NULL,
`update_at` int(64) NOT NULL DEFAULT 0,
`description` text NOT NULL DEFAULT '',
`features` text NOT NULL DEFAULT '',
`version` varchar(255) NOT NULL DEFAULT '',
`service_id` varchar(128) NOT NULL DEFAULT '',
`healthcheck_enabled` int(2) NOT NULL DEFAULT 0,
`healthcheck_endpoint` varchar(255) NOT NULL DEFAULT '',
`healthcheck_interval` int(64) NOT NULL DEFAULT 0,
`extra` text NOT NULL DEFAULT '',
PRIMARY KEY (`id`) 
)
ENGINE = InnoDB
AUTO_INCREMENT = 1
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_general_ci
COMMENT = '服务实例表';

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