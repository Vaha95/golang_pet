package DTO

import "time"

const (
	SHORTEN = "shorten"
	FOLLOW = "follow"
)

type BaseAuditItem struct {
	URL string `json:"url"`
	Action string `json:"action"`
	UserId *int `json:"user_id"`
	Ts string `json:"ts"`
}

func CrateBaseAuditItemShorten(url string, userId *int) BaseAuditItem {
	return  BaseAuditItem{
		URL: url,
		Action: SHORTEN,
		UserId: userId,
		Ts: time.Now().String(),
	}
}

func CrateBaseAuditItemFollow(url string, userId *int) BaseAuditItem {
	return  BaseAuditItem{
		URL: url,
		Action: FOLLOW,
		UserId: userId,
		Ts: time.Now().String(),
	}
}