use onbio;


CREATE TABLE `t_user` (
  `id`  bigint(20)  NOT NULL AUTO_INCREMENT COMMENT 'id',
  `user_name` varchar(255) NOT NULL DEFAULT '' COMMENT '用户名',
  `user_pwd` varchar(255) NOT NULL DEFAULT '' COMMENT '用户密码',
  `user_avatar` varchar(255) NOT NULL DEFAULT '' COMMENT '用户头像',
  `user_type` varchar(255) NOT NULL DEFAULT '' COMMENT '用户类型',
  `user_src` int(11) NOT NULL DEFAULT '1' COMMENT '用户来源 1 自己注册 2 第三方登录',
  `user_extra` varchar(1024) NOT NULL DEFAULT '' COMMENT '保留字段',
  `user_link` varchar(25) NOT NULL DEFAULT '' COMMENT '用户个性链接',
  `is_confirmed` tinyint(2) NOT NULL DEFAULT '0' COMMENT '是否通过邮箱认证',
  `email` varchar(25) NOT NULL DEFAULT '' COMMENT '用户邮箱',
  `operator` varchar(255) NOT NULL DEFAULT '' COMMENT '操作人',
  `use_flag` tinyint(2) NOT NULL DEFAULT '0' COMMENT '是否有效',
  `create_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `last_updated_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '最后更新时间',
  PRIMARY KEY (`id`)
) ENGINE=INNODB AUTO_INCREMENT=10000 DEFAULT CHARSET=utf8 COMMENT '用户表';

CREATE TABLE `t_user_link` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'id',
  `user_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '用户id',
  `link_url` varchar(255) NOT NULL DEFAULT '' COMMENT '用户链接',
  `link_desc` varchar(2048) NOT NULL DEFAULT '' COMMENT '内容简述',
  `user_img` varchar(255) NOT NULL DEFAULT '' COMMENT '链接首图',
  `operator` varchar(255) NOT NULL DEFAULT '' COMMENT '操作人',
  `use_flag` tinyint(2) NOT NULL DEFAULT '0' COMMENT '是否有效',
  `create_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `last_updated_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '最后更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_user` (`user_id`)
) ENGINE=INNODB AUTO_INCREMENT=1  DEFAULT CHARSET=utf8 COMMENT='用户链接表';


  CREATE TABLE `t_login_log` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'id',
    `user_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '用户id',
    `user_name` varchar(255) NOT NULL DEFAULT '' COMMENT '用户名',
    `login_platform` varchar(128) NOT NULL DEFAULT 'pc' COMMENT '登录渠道 pc/ios/android',
    `session_id` varchar(255) NOT NULL DEFAULT '' COMMENT '登录的sessionid',
    `logout_type` int(11) NOT NULL DEFAULT '0' COMMENT '登出原因',
    `operator` varchar(255) NOT NULL DEFAULT '' COMMENT '操作人',
    `use_flag` tinyint(2) NOT NULL DEFAULT '0' COMMENT '是否有效',
    `create_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '创建时间',
    `last_updated_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '最后更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_user` (`user_id`)
  ) ENGINE=INNODB AUTO_INCREMENT=1  DEFAULT CHARSET=utf8 COMMENT='用户登录记录表';


