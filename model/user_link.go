package model

import (
	"onbio/logger"
	"time"

	"go.uber.org/zap"
)

const (
	LinkTableName = "t_user_link"
)

/***

CREATE TABLE `t_user_link` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'id',
  `user_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '用户id',
  `link_url` varchar(255) NOT NULL DEFAULT '' COMMENT '用户链接',
  `link_desc` varchar(2048) NOT NULL DEFAULT '' COMMENT '内容简述',
  `link_img` varchar(255) NOT NULL DEFAULT '' COMMENT '链接首图',
  `operator` varchar(255) NOT NULL DEFAULT '' COMMENT '操作人',
  `use_flag` tinyint(2) NOT NULL DEFAULT '0' COMMENT '是否有效',
  `create_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `last_updated_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '最后更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_user` (`user_id`)
) ENGINE=INNODB AUTO_INCREMENT=1  DEFAULT CHARSET=utf8 COMMENT='用户链接表';
**/
type Link struct {
	ID              uint64 `gorm:"primaryKey"  json:"id"`
	UserID          uint64 `gorm:"column:user_id" json:"user_id"`
	LinkUrl         string `gorm:"column:link_url" json:"link_url"`
	LinkDesc        string `gorm:"column:link_desc" json:"link_desc"`
	LinkImg         string `gorm:"column:link_img" json:"link_img"`
	Operator        string `gorm:"column:operator" json:"operator"`
	UseFlag         int    `gorm:"column:use_flag" json:"use_flag"`
	CreateTime      uint64 `gorm:"column:create_time" json:"create_time"`
	LastUpdatedTime uint64 `gorm:"column:last_updated_time" json:"last_updated_time"`
}

func CreateLink(userID uint64, linkUrl, linkDesc, linkImg string) (err error) {

	newLink := Link{
		UserID:          userID,
		LinkUrl:         linkUrl,
		LinkDesc:        linkDesc,
		LinkImg:         linkImg,
		CreateTime:      uint64(time.Now().Unix()),
		LastUpdatedTime: uint64(time.Now().Unix()),
	}

	db := getMysqlConn().Table(UserTableName)
	db = db.Create(&newLink)
	if db.Error != nil {
		logger.Error("CreateLink::Find error: %s", zap.Error(db.Error))
		return db.Error
	}
	return
}

//更新链接记录
func UpdateLinkByID(linkID uint64, link Link) (err error) {
	return
}


//获取用户链接列表
func GetUserLinkList(userID uint64, page, pageSize int) (linkList []*Link, err error) {
  return 
}
