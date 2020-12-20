package services

const (
	USER_SESSION_REDIS_PRE = "onbio_session:%s"
	USER_REDIRECT_URL      = "http://www.qq.com"
)

type SessionContent struct {
	UserName    string `json:"user_name"`
	UserAvatar  string `json:"user_avatar"`
	UserID      uint64 `json:"user_id"`
	UserLink    string `json:"user_link"`
	IsConfirmed int    `json:"is_confirmed"`
	Email       string `json:"email"`
	LoginTime   uint64 `json:"login_time"`
}
