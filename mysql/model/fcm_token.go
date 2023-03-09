package model

type FcmToken struct {
	Id         int    `json:"id"`
	Platform   int    `json:"platform"`
	UserId     int    `json:"user_id"`
	Token      string `json:"token"`
	CreateTime int    `json:"create_time"`
	UpdateTime int    `json:"update_time"`
}

func (FcmToken) TableName() string {
	return "cs_fcm_token"
}
