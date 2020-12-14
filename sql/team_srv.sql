use db_team;


CREATE TABLE `svip_team` (
  `team_id`  bigint(20)  NOT NULL AUTO_INCREMENT COMMENT 'id',
  `team_name` varchar(255) NOT NULL DEFAULT '' COMMENT '队伍名称',
  `game_id` varchar(255) NOT NULL DEFAULT '' COMMENT '游戏ID',
  `game_name` varchar(255) NOT NULL DEFAULT '' COMMENT '游戏名称',
  `team_leader` varchar(255) NOT NULL DEFAULT '' COMMENT '队长用户id',
  `team_status` int(11) NOT NULL DEFAULT '1' COMMENT '队伍状态,1 是正常，2 是已删除',
  `team_property` varchar(1024) NOT NULL DEFAULT '' COMMENT '队伍属性值,保留字段',
  `bussiness_id` varchar(25) NOT NULL DEFAULT '' COMMENT '业务ID，区分不同业务',
  `team_del_time` bigint NOT NULL DEFAULT '0' COMMENT '队伍删除时间',
  `operator` varchar(255) NOT NULL DEFAULT '' COMMENT '操作人',
  `create_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `last_updated_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '最后更新时间',
  PRIMARY KEY (`team_id`)
) ENGINE=INNODB AUTO_INCREMENT=10000 DEFAULT CHARSET=utf8 COMMENT '组队表';


CREATE TABLE `svip_team_member` (
  `user_id` varchar(255) NOT NULL DEFAULT '' COMMENT '用户id',
  `user_name` varchar(255) NOT NULL  DEFAULT '' COMMENT '用户名',
  `team_id` bigint(20)  NOT NULL DEFAULT '0' COMMENT '队伍id',
  `team_name` varchar(255) NOT NULL DEFAULT '' COMMENT '队伍名称',
  `game_id` varchar(255) NOT NULL DEFAULT '' COMMENT '游戏ID',
  `game_name` varchar(255) NOT NULL DEFAULT '' COMMENT '游戏名称',
  `user_avatar` varchar(255) NOT NULL COMMENT '用户头像',
  `user_team_status`  int(11) NOT NULL DEFAULT '1' COMMENT '队内状态 1 正常 2 已离队',
  `user_team_role` int(11) NOT NULL DEFAULT '0' COMMENT '角色:队长/队员',
  `user_join_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '入队时间',
  `user_leave_reason`  int(11) NOT NULL DEFAULT '0' COMMENT '离队原因：1-队员主动离队；2-队长踢出；3-客服解散队伍',
  `user_leave_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '离队时间',
  `create_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `last_updated_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '最后更新时间',
  UNIQUE `uni_team` (`user_id`,`team_id`)
) ENGINE=InnoDB AUTO_INCREMENT=10000  DEFAULT CHARSET=utf8 COMMENT='组队成员表';


CREATE TABLE `svip_team_log` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'id',
  `user_id` varchar(255) NOT NULL DEFAULT '' COMMENT '用户id',
  `user_name` varchar(255) NOT NULL DEFAULT '' COMMENT '用户名',
  `team_id` int(11) NOT NULL DEFAULT '0' COMMENT '队伍ID',
  `team_name` varchar(255) NOT NULL DEFAULT '' COMMENT '队伍名称',
  `game_id` int(11)  NOT NULL DEFAULT '0' COMMENT '游戏ID',
  `game_name` varchar(255) NOT NULL DEFAULT '' COMMENT '游戏名称',
  `type` int(11) NOT NULL DEFAULT '0' COMMENT '类型',
  `content` varchar(255) NOT NULL DEFAULT '' COMMENT '日志内容',
  `operator` varchar(255) NOT NULL DEFAULT '' COMMENT '操作者(运营操作时为rtx)',
  `create_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `last_updated_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '最后更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_user` (`user_id`),
  KEY `idx_team` (`team_id`)
) ENGINE=INNODB AUTO_INCREMENT=10000  DEFAULT CHARSET=utf8 COMMENT='队伍操作日志表'; 
