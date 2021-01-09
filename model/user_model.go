package model

import (
	"errors"
	"onbio/logger"
	"onbio/mysql"
	"onbio/utils"
	"time"

	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

const (
	UserTableName = "t_user"
)

/***

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


**/
type User struct {
	ID              uint64 `gorm:"primaryKey"  json:"id"`
	UserName        string `gorm:"column:user_name" json:"user_name"`
	UserPwd         string `gorm:"column:user_pwd" json:"user_pwd"`
	UserAvatar      string `gorm:"column:user_avatar" json:"user_avatar"`
	UserType        string `gorm:"column:user_type" json:"user_type"`
	UserSrc         int    `gorm:"column:user_src" json:"user_src"`
	UserExtra       string `gorm:"column:user_extra" json:"user_extra"`
	UserLink        string `gorm:"column:user_link" json:"user_link"`
	IsConfirmed     int    `gorm:"column:is_confirmed" json:"is_confirmed"`
	Email           string `gorm:"column:email" json:"email"`
	Operator        string `gorm:"column:operator" json:"operator"`
	UseFlag         int    `gorm:"column:use_flag" json:"use_flag"`
	CreateTime      uint64 `gorm:"column:create_time" json:"create_time"`
	LastUpdatedTime uint64 `gorm:"column:last_updated_time" json:"last_updated_time"`
}

func CreateUser(userName, userAvatar, userPwd, email string) error {

	//md5sum
	userPwd, _ = utils.Md5Sum(userPwd)

	newUser := User{
		UserName:        userName,
		UserPwd:         userPwd,
		Email:           email,
		UserAvatar:      userAvatar,
		CreateTime:      uint64(time.Now().Unix()),
		LastUpdatedTime: uint64(time.Now().Unix()),
	}

	db := getMysqlConn().Table(UserTableName)
	db = db.Create(&newUser)
	if db.Error != nil {
		logger.Error("CreateUser::Find error: %s", zap.Error(db.Error))
		return db.Error
	}
	return nil
}

func IsEmailExisted(email string) (err error, isExisted bool) {

	db := getMysqlConn().Table(UserTableName)
	if len(email) != 0 {
		db = db.Where("email = ?", email)
	}
	var count int
	err = db.Count(&count).Error
	if err != nil {
		logger.Error("get user from db failed ")
		return db.Error, true
	}
	if count > 0 {
		return nil, true
	}
	return nil, false
}

func IsUserExisted(userName string) (err error, isExisted bool) {

	db := getMysqlConn().Table(UserTableName)
	if len(userName) != 0 {
		db = db.Where("user_name = ?", userName)
	}
	var count int
	err = db.Count(&count).Error
	if err != nil {
		logger.Error("get user from db failed ")
		return db.Error, true
	}
	if count > 0 {
		return nil, true
	}
	return nil, false
}

func UpdateUserInfoByID(userID uint64, info User) (err error) {
	user := User{
		ID: userID,
	}

	updates := map[string]interface{}{}

	if info.IsConfirmed != 0 {
		updates["is_confirmed"] = info.IsConfirmed
	}
	if info.Email != "" {
		updates["email"] = info.Email
	}

	if info.UserPwd != "" {
		updates["user_pwd"] = info.UserPwd
	}

	if info.UserAvatar != "" {
		updates["user_avatar"] = info.UserAvatar
	}

	updates["last_updated_time"] = uint64(time.Now().Unix())

	db := getMysqlConn().Table(UserTableName)
	err = db.Model(&user).Updates(updates).Error
	if err != nil {
		logger.Error("update user info ", zap.Any("model", user), zap.Any("updates", updates))
		return
	}
	return
}

func GetUserInfo(userEmail, userName string, userID uint64) (err error, user User) {
	db := getMysqlConn().Table(UserTableName)
	if len(userName) != 0 {
		db = db.Where("user_name = ?", userName)
	}

	if len(userEmail) != 0 {
		db = db.Where("email = ?", userEmail)
	}

	if userID != 0 {
		db = db.Where("id = ?", userID)
	}

	err = db.Scan(&user).Error
	if err != nil {
		logger.Error("get user info from db failed ")
		return
	}

	if user.UserName == "" {
		logger.Error("invalid user ")
		err = errors.New("invalid user")
	}

	return

}

func CheckUserPwd(userName, userPwd string) (err error, user User) {

	//md5sum
	userPwd, _ = utils.Md5Sum(userPwd)
	db := getMysqlConn().Table(UserTableName)
	if len(userName) != 0 {
		db = db.Where("user_name = ?", userName)
	}

	err = db.Scan(&user).Error
	if err != nil {
		logger.Error("get user info from db failed ")
		return
	}

	if user.UserName == "" {
		logger.Error("invalid user ")
		err = errors.New("invalid user")
	}

	if user.UserPwd != userPwd {
		err = errors.New("pwd incorrect")
		return
	}
	return
}

func getMysqlConn() *gorm.DB {
	return mysql.GetDBConn("teamDB")
}
