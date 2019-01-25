CREATE TABLE `applications` (
`id` char(64) NOT NULL,
`name` varchar(64) NOT NULL DEFAULT '' COMMENT '应用的名称',
`user_id` char(64) NOT NULL DEFAULT '' COMMENT '应用持有者的用户ID',
`website` varchar(255) NOT NULL DEFAULT '' COMMENT '应用站点的URL',
`logo_image` varchar(255) NOT NULL DEFAULT '' COMMENT '应用的LOG',
`description` text NOT NULL COMMENT '应用的描述信息',
`create_at` int(64) UNSIGNED NOT NULL DEFAULT 0 COMMENT '应用的创建时间',
`update_at` int(64) UNSIGNED NOT NULL DEFAULT 0 COMMENT '应用更新时间',
`redirect_uri` varchar(128) NOT NULL DEFAULT '' COMMENT '当使用AuthCode时, 重定向的地址',
`client_id` char(64) NOT NULL DEFAULT '' COMMENT '该应用的客户端ID',
`client_secret` char(255) NOT NULL DEFAULT '' COMMENT '客户端凭证',
`locked` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否冻结, 冻结后禁止该APP访问',
`last_login_time` int(64) UNSIGNED NOT NULL DEFAULT 0 COMMENT '最近一次访问的时间',
`last_login_ip` varchar(255) NOT NULL DEFAULT '' COMMENT '最近一次访问的IP',
`login_failed_times` int(64) UNSIGNED NOT NULL DEFAULT 0 COMMENT '登录失败的次数',
`login_success_times` int(64) UNSIGNED NOT NULL DEFAULT 0 COMMENT '登录成功的次数',
`token_expire_time` int(64) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'token过期时间',
`extra` text NOT NULL COMMENT '预留字段',
PRIMARY KEY (`id`) ,
INDEX `client` (`client_id` ASC)
)
ENGINE = InnoDB
AUTO_INCREMENT = 0
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_general_ci
COMMENT = '用户应用'
ROW_FORMAT = dynamic;

CREATE TABLE `auth_codes` (
`code` varchar(128) NOT NULL COMMENT '授权码',
`user_id` char(64) NOT NULL DEFAULT '' COMMENT '用户ID',
`create_at` int(64) UNSIGNED NOT NULL DEFAULT 0 COMMENT '发放时间',
`expires_at` int(64) UNSIGNED NOT NULL DEFAULT 0 COMMENT '过期时间',
`extra` text NOT NULL COMMENT '预留字段',
PRIMARY KEY (`code`) 
)
ENGINE = InnoDB
AUTO_INCREMENT = 0
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_general_ci
COMMENT = 'OAuth2.0授权码'
ROW_FORMAT = dynamic;

CREATE TABLE `dbmanager` (
`id` int(11) NOT NULL AUTO_INCREMENT COMMENT '支持SQL计数',
`version` int(11) NOT NULL COMMENT 'SQL文件的版本',
`md5` int(11) NOT NULL COMMENT 'SQL文件内容的MD5消息摘要',
`create_at` int(64) UNSIGNED NOT NULL DEFAULT 0 COMMENT '执行数据库升级的时间',
`description` text NOT NULL COMMENT '该次schema变更说明',
PRIMARY KEY (`id`) 
)
ENGINE = InnoDB
AUTO_INCREMENT = 0
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_general_ci
COMMENT = '数据库版本变更管理表'
ROW_FORMAT = dynamic;

CREATE TABLE `domains` (
`id` char(64) NOT NULL,
`name` varchar(128) NOT NULL DEFAULT '' COMMENT '域名称, 不允许重复',
`display_name` varchar(255) NOT NULL DEFAULT '' COMMENT '域显示名称, 用于界面进行展示, 一般为公司名称',
`logo_path` varchar(128) NOT NULL DEFAULT '' COMMENT '公司的LOGO',
`description` text NOT NULL COMMENT '域描述信息',
`enabled` tinyint(1) UNSIGNED NOT NULL DEFAULT 0 COMMENT '是否启用或者禁用改域(0, 1)',
`type` tinyint(1) UNSIGNED NOT NULL DEFAULT 0 COMMENT '域类型(1: 个人域, 2:  企业域, 3: 合作伙伴域)',
`create_at` int(64) UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建时间(时间戳)',
`update_at` int(64) UNSIGNED NOT NULL DEFAULT 0 COMMENT '更新时间(时间戳)',
`size` varchar(64) NOT NULL DEFAULT '' COMMENT '公司规模 ',
`location` varchar(64) NOT NULL DEFAULT '' COMMENT '公司所处位置, 比如成都',
`industry` varchar(64) NOT NULL DEFAULT '' COMMENT '公司所属行业',
`address` varchar(64) NOT NULL DEFAULT '' COMMENT '公司地址, 长度不能超过255',
`fax` varchar(128) NOT NULL DEFAULT '' COMMENT '公司传真',
`phone` varchar(128) NOT NULL DEFAULT '' COMMENT '公司固定电话',
`contacts_name` varchar(32) NOT NULL DEFAULT '' COMMENT '公司联系人姓名',
`contacts_title` varchar(32) NOT NULL DEFAULT '' COMMENT '联系人的职位',
`contacts_mobile` varchar(32) NOT NULL DEFAULT '' COMMENT '联系人电话',
`contacts_email` varchar(32) NOT NULL DEFAULT '' COMMENT '联系人邮箱',
`owner_id` char(64) NOT NULL DEFAULT '' COMMENT '域的持有者, 如果持有者注销了自己的账号, 该域也需要被回收',
`extra` text NOT NULL COMMENT '保留字段',
PRIMARY KEY (`id`) ,
UNIQUE INDEX `name` (`name` ASC)
)
ENGINE = InnoDB
AUTO_INCREMENT = 0
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_general_ci
COMMENT = '域(用户空间)'
ROW_FORMAT = dynamic;

CREATE TABLE `features` (
`id` char(64) NOT NULL,
`name` varchar(32) NOT NULL DEFAULT '' COMMENT '实例功能名称',
`tag` varchar(256) NOT NULL DEFAULT '' COMMENT 'HTTP 方法的名称',
`endpoint` varchar(256) NOT NULL DEFAULT '' COMMENT '功能的对应的URL',
`description` text NOT NULL COMMENT '功能描述',
`is_deleted` int(2) UNSIGNED NOT NULL DEFAULT 0 COMMENT '该功能是否已经被删除',
`when_deleted_version` varchar(255) NOT NULL DEFAULT '' COMMENT '在那个版本被删除',
`when_deleted_time` int(64) UNSIGNED NOT NULL DEFAULT 0 COMMENT '功能删除的时间',
`is_added` int(2) UNSIGNED NOT NULL DEFAULT 0 COMMENT '是否是新增的功能',
`when_added_version` varchar(255) NOT NULL DEFAULT '' COMMENT '那个版本新增的功能',
`when_added_time` int(64) UNSIGNED NOT NULL DEFAULT 0 COMMENT '功能注册的时间',
`service_id` char(64) NOT NULL DEFAULT '' COMMENT '该功能属于那个服务',
`extra` text NOT NULL,
PRIMARY KEY (`id`) 
)
ENGINE = InnoDB
AUTO_INCREMENT = 0
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_general_ci
COMMENT = '服务功能注册表'
ROW_FORMAT = dynamic;

CREATE TABLE `instances` (
`id` char(64) NOT NULL,
`name` varchar(32) NOT NULL DEFAULT '' COMMENT '服务实例的名称',
`address` varchar(32) NOT NULL DEFAULT '' COMMENT '服务实例的地址',
`service_type` varchar(32) NOT NULL DEFAULT '' COMMENT '实例所属的服务类型(controller, worker)',
`service_name` varchar(128) NOT NULL DEFAULT '' COMMENT '实例所属的服务名称(一个服务可以有多个实例)',
`version` varchar(255) NOT NULL DEFAULT '' COMMENT '实例的版本好',
`git_branch` varchar(255) NOT NULL DEFAULT '' COMMENT '实例构建时 对应的git 分支',
`git_commit` varchar(255) NOT NULL DEFAULT '' COMMENT '实例构建时 对应的git commit id',
`build_env` varchar(255) NOT NULL DEFAULT '' COMMENT '实例构建的环境(JDK1.8, GO1.10...)',
`build_at` varchar(255) NOT NULL DEFAULT '' COMMENT '实例的构建时间(字符串)',
`status` varchar(255) NOT NULL COMMENT '实例当前的状态(0: known, 1: online 2: offline)',
`online_at` int(64) UNSIGNED NOT NULL DEFAULT 0 COMMENT '实例上线时间(时间戳)',
`offline_at` int(64) UNSIGNED NOT NULL DEFAULT 0 COMMENT '示例下线时间(时间戳)',
`extra` text NOT NULL COMMENT '预留字段',
PRIMARY KEY (`id`) 
)
ENGINE = InnoDB
AUTO_INCREMENT = 0
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_general_ci
COMMENT = '服务实例, 属于服务发现部分'
ROW_FORMAT = dynamic;

CREATE TABLE `passwords` (
`user_id` char(64) NOT NULL DEFAULT '' COMMENT '用户ID',
`password` varchar(255) NOT NULL DEFAULT '' COMMENT 'hash过后的密码',
`expires_at` int(64) UNSIGNED NOT NULL DEFAULT 0 COMMENT '密码过期时间',
`create_at` int(64) UNSIGNED NOT NULL DEFAULT 0 COMMENT '密码创建时间',
`update_at` int(64) UNSIGNED NOT NULL DEFAULT 0,
`extra` text NOT NULL COMMENT '预留字段',
PRIMARY KEY (`user_id`) 
)
ENGINE = InnoDB
AUTO_INCREMENT = 0
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_general_ci
COMMENT = '用户密码'
ROW_FORMAT = dynamic;

CREATE TABLE `projects` (
`id` char(64) NOT NULL,
`name` varchar(128) NOT NULL DEFAULT '' COMMENT '项目名称',
`picture` varchar(128) NOT NULL DEFAULT '' COMMENT '项目图片',
`latitude` float(32,0) NOT NULL DEFAULT 0 COMMENT '项目所处的 经度',
`longitude` float(32,0) NOT NULL DEFAULT 0 COMMENT '项目所处的 维度',
`enabled` int(1) UNSIGNED NOT NULL DEFAULT 0 COMMENT '是否启用或者禁用此项目',
`owner_id` char(64) NOT NULL DEFAULT '' COMMENT ' 项目拥有者ID',
`description` text NOT NULL COMMENT '项目描述',
`domain_id` char(64) NOT NULL DEFAULT '' COMMENT '项目所属域',
`create_at` int(64) UNSIGNED NOT NULL DEFAULT 0 COMMENT '项目创建时间',
`update_at` int(64) UNSIGNED NOT NULL DEFAULT 0 COMMENT '最近一次访问时间',
`extra` text NOT NULL COMMENT '预留字段',
PRIMARY KEY (`id`) 
)
ENGINE = InnoDB
AUTO_INCREMENT = 0
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_general_ci
COMMENT = '资源容器'
ROW_FORMAT = dynamic;

CREATE TABLE `roles` (
`id` char(64) NOT NULL,
`name` varchar(64) NOT NULL DEFAULT '' COMMENT '名称',
`description` text NOT NULL COMMENT '描述',
`create_at` int(64) UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建时间',
`update_at` int(64) UNSIGNED NOT NULL DEFAULT 0 COMMENT '更新时间',
`extra` text NOT NULL COMMENT '预留字段',
PRIMARY KEY (`id`) ,
UNIQUE INDEX `name` (`name` ASC)
)
ENGINE = InnoDB
AUTO_INCREMENT = 0
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_general_ci
ROW_FORMAT = dynamic;

CREATE TABLE `role_feature_mappings` (
`feature_id` char(64) NOT NULL DEFAULT '',
`role_id` char(64) NOT NULL DEFAULT ''
)
ENGINE = InnoDB
AUTO_INCREMENT = 0
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_general_ci
ROW_FORMAT = dynamic;

CREATE TABLE `role_user_mappings` (
`user_id` char(64) NOT NULL DEFAULT '' COMMENT '用户',
`domain_id` char(64) NOT NULL DEFAULT '' COMMENT '域',
`role_id` char(64) NOT NULL DEFAULT '' COMMENT '角色'
)
ENGINE = InnoDB
AUTO_INCREMENT = 0
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_general_ci
COMMENT = '用于在同步域的角色对应表'
ROW_FORMAT = dynamic;

CREATE TABLE `services` (
`id` char(64) NOT NULL COMMENT '服务的ID',
`type` char(64) NOT NULL COMMENT '服务类型',
`name` varchar(255) NOT NULL DEFAULT '' COMMENT '服务名称',
`description` text NOT NULL COMMENT '服务的功能简介',
`enabled` int(2) UNSIGNED NOT NULL DEFAULT 0 COMMENT '服务状态',
`status` varchar(255) NOT NULL DEFAULT '' COMMENT '服务状态',
`status_update_at` int(64) UNSIGNED NOT NULL DEFAULT 0,
`current_version` varchar(128) NOT NULL DEFAULT '',
`upgrade_version` varchar(128) NOT NULL DEFAULT '',
`downgrade_version` varchar(128) NOT NULL DEFAULT '',
`create_at` int(64) UNSIGNED NOT NULL DEFAULT 0,
`update_at` int(64) UNSIGNED NOT NULL DEFAULT 0,
`client_id` char(128) NOT NULL DEFAULT '' COMMENT '客户端id',
`client_secret` char(255) NOT NULL DEFAULT '' COMMENT '客户端凭证',
`token_expire_time` int(64) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'token过期时间',
`extra` text NOT NULL,
PRIMARY KEY (`id`) ,
INDEX `name` (`name` ASC, `client_id` ASC)
)
ENGINE = InnoDB
AUTO_INCREMENT = 0
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_general_ci
COMMENT = '服务注册表'
ROW_FORMAT = dynamic;

CREATE TABLE `user_domain_mappings` (
`user_id` char(64) NOT NULL DEFAULT '',
`domain_id` char(64) NOT NULL DEFAULT '',
`join_at` int(64) UNSIGNED NOT NULL DEFAULT 0 COMMENT '加入时间',
`extra` text NOT NULL COMMENT '预留字段'
)
ENGINE = InnoDB
AUTO_INCREMENT = 0
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_general_ci
COMMENT = '用户第三方域'
ROW_FORMAT = dynamic;

CREATE TABLE `tokens` (
`access_token` varchar(128) NOT NULL COMMENT '访问令牌',
`refresh_token` varchar(128) NOT NULL DEFAULT '' COMMENT '刷新令牌, 当访问令牌过期后, 可以使用该token刷新一个新的token, 仅能使用一次',
`scope` varchar(255) NOT NULL DEFAULT '' COMMENT '授权的作用域',
`name` char(64) NOT NULL DEFAULT '' COMMENT '令牌的名称, 用于申请API访问内心的token',
`grant_type` char(64) NOT NULL DEFAULT '' COMMENT '授权类型(password, client, authcode, implement, upgrade, sdk)',
`token_type` char(64) NOT NULL DEFAULT '' COMMENT '令牌的类型(bearer,  jwt...)',
`create_at` int(64) UNSIGNED NOT NULL DEFAULT 0 COMMENT '令牌发放时间(时间戳)',
`expire_at` int(64) UNSIGNED NOT NULL DEFAULT 0 COMMENT '令牌过期时间(时间戳)',
`user_id` char(64) NOT NULL DEFAULT '' COMMENT '令牌的持有者ID',
`domain_id` char(64) NOT NULL DEFAULT '' COMMENT '令牌的作用域 - 域内可用',
`project_id` char(64) NOT NULL DEFAULT '' COMMENT '令牌的作用域 - 仅能作用在某一个项目内使用',
`application_id` char(64) NOT NULL DEFAULT '' COMMENT '如果token是发放给app就存在该值',
`service_id` char(64) NOT NULL DEFAULT '' COMMENT '服务名称, 仅用于当 client模式, 专门颁发给 服务注册使用',
`description` text NOT NULL COMMENT '令牌的描述信息',
`extra` text NOT NULL COMMENT '预留字段',
PRIMARY KEY (`access_token`) ,
UNIQUE INDEX `token` (`refresh_token` ASC)
)
ENGINE = InnoDB
AUTO_INCREMENT = 0
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_general_ci
COMMENT = '用户访问令牌'
ROW_FORMAT = dynamic;

CREATE TABLE `invitation_records` (
`code` char(64) NOT NULL DEFAULT '' COMMENT '邀请码',
`inviter` char(64) NOT NULL DEFAULT '' COMMENT '邀请人ID',
`invitee` char(64) NOT NULL DEFAULT '' COMMENT '被邀请人ID',
`invitee_domain` char(64) NOT NULL DEFAULT '' COMMENT '被邀请人域ID',
`invitee_roles` varchar(128) NOT NULL DEFAULT '' COMMENT '授予被邀请人的角色列表',
`invitation_time` int(64) NOT NULL DEFAULT 0 COMMENT '发放邀请码的时间',
`accept_time` int(64) UNSIGNED NOT NULL DEFAULT 0 COMMENT '接收邀请的时间',
`expire_time` int(64) UNSIGNED NOT NULL DEFAULT 0 COMMENT '邀请码过期时间',
`access_projects` varchar(255) NOT NULL DEFAULT '' COMMENT '允许访问的项目列表',
`extra` text NOT NULL COMMENT '预留字段',
PRIMARY KEY (`code`) 
)
ENGINE = InnoDB
AUTO_INCREMENT = 0
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_general_ci
COMMENT = '用户邀请记录'
ROW_FORMAT = dynamic;

CREATE TABLE `users` (
`id` char(64) NOT NULL,
`department` char(64) NOT NULL DEFAULT '' COMMENT '用户部门',
`account` char(64) NOT NULL DEFAULT '' COMMENT '用户名(用于登录系统)',
`mobile` char(11) NOT NULL DEFAULT '' COMMENT '主电话号码',
`email` char(64) NOT NULL DEFAULT '' COMMENT '主邮件地址',
`phone` char(20) NOT NULL DEFAULT '' COMMENT '用户座机',
`address` varchar(128) NOT NULL DEFAULT '' COMMENT '用户住址',
`real_name` varchar(128) NOT NULL DEFAULT '' COMMENT '真实姓名',
`nick_name` varchar(128) NOT NULL DEFAULT '' COMMENT '用户昵称',
`gender` char(1) NOT NULL DEFAULT '' COMMENT '用户的性别(m: 男性 male,  f:  女性 female)',
`avatar` varchar(128) NOT NULL DEFAULT '' COMMENT '用户头像',
`language` char(64) NOT NULL DEFAULT '' COMMENT '用户使用的语言',
`city` varchar(64) NOT NULL DEFAULT '' COMMENT '用户所在的城市',
`province` varchar(64) NOT NULL DEFAULT '用户所在的省',
`locked` int(1) UNSIGNED NOT NULL DEFAULT 0 COMMENT '是否启用/冻结该用户(0: 冻结, 1: 启用)',
`domain_id` char(64) NOT NULL DEFAULT '' COMMENT '用户所以域ID',
`create_at` int(64) UNSIGNED NOT NULL DEFAULT 0 COMMENT '用户创建时间',
`expires_active_days` int(64) UNSIGNED NOT NULL DEFAULT 0 COMMENT '用户多少天后未登录,  则用户过期',
`default_project_id` char(64) NOT NULL DEFAULT '' COMMENT '用户默认项目ID',
`extra` text NOT NULL,
PRIMARY KEY (`id`) ,
UNIQUE INDEX `user` (`email` ASC, `account` ASC, `mobile` ASC) USING BTREE
)
ENGINE = InnoDB
AUTO_INCREMENT = 0
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_general_ci
COMMENT = '用户'
ROW_FORMAT = dynamic;


CREATE TABLE `login_records` (
`id` char(64) NOT NULL COMMENT '通过hash({user_id}.{application_id})获得',
`ip` char(64) NOT NULL DEFAULT '' COMMENT '最近一次登录的IP(user password认证用户)',
`login` int(64) UNSIGNED NOT NULL DEFAULT 0 COMMENT '最近一次登录时间(user password认证用户)',
`logout` int(64) UNSIGNED NOT NULL DEFAULT 0 COMMENT '用户最近一次退出系统的时间',
`failed` int(64) UNSIGNED NOT NULL DEFAULT 0 COMMENT '用户连续登录失败的次数',
`success` int(64) UNSIGNED NOT NULL DEFAULT 0 COMMENT '成功登陆的次数',
`user_id` char(64) NOT NULL DEFAULT '' COMMENT '令牌的持有者ID',
`application_id` char(64) NOT NULL DEFAULT '' COMMENT '用户通过那个应用登录的',
`grant_type` char(64) NOT NULL DEFAULT '' COMMENT '授权类型(password, client, authcode, implement, upgrade, sdk)',
`extra` text NOT NULL,
PRIMARY KEY (`id`) 
)
ENGINE = InnoDB
AUTO_INCREMENT = 0
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_general_ci
COMMENT = '用户登录记录'
ROW_FORMAT = dynamic;


CREATE TABLE `user_project_mappings` (
`user_id` char(64) NULL DEFAULT '',
`project_id` char(64) NULL DEFAULT ''
)
ENGINE = InnoDB
AUTO_INCREMENT = 0
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_general_ci
ROW_FORMAT = dynamic;

CREATE TABLE `verification_code` (
`code` char(8) NOT NULL COMMENT '验证码',
`purpose` tinyint(2) UNSIGNED NOT NULL DEFAULT 0 COMMENT '验证码的用途  (1: 注册码,  2: 密码找回,  3:  登录)',
`sending_mode` tinyint(2) UNSIGNED NOT NULL DEFAULT 0 COMMENT '验证码的发送方式 (1: 邮箱, 2: 短信)',
`sending_target` varchar(128) NOT NULL DEFAULT '' COMMENT '发送者的具体值和sending_mode对应',
`create_at` int(64) UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建时间',
`expire_at` int(64) UNSIGNED NOT NULL DEFAULT 0 COMMENT '过期时间',
`status` int(64) UNSIGNED NOT NULL DEFAULT 0,
`extra` text NOT NULL COMMENT '预留字段',
PRIMARY KEY (`code`) 
)
ENGINE = InnoDB
AUTO_INCREMENT = 0
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_general_ci
COMMENT = '验证码'
ROW_FORMAT = dynamic;

CREATE TABLE `departments` (
`id` char(64) NOT NULL,
`name` char(64) NOT NULL COMMENT '部门名称',
`number` char(64) NOT NULL COMMENT '部门编号',
`parent` char(64) NOT NULL DEFAULT '' COMMENT '上级部门ID',
`grade` tinyint(3) NOT NULL DEFAULT 0 COMMENT '第几层',
`path` text NOT NULL DEFAULT '' COMMENT '具体路径',
`manager` char(64) NOT NULL DEFAULT '' COMMENT '部门负责人',
`domain_id` char(64) NOT NULL DEFAULT '' COMMENT '所属的域ID',
`create_at` int(64) UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建时间',
`extra` text NOT NULL COMMENT '预留字段',
PRIMARY KEY (`id`),
UNIQUE INDEX `dept` (`number` ASC) 
)
ENGINE = InnoDB
AUTO_INCREMENT = 0
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_general_ci
COMMENT = '用户部门'
ROW_FORMAT = dynamic;

