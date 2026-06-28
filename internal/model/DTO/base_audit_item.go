package DTO

import "time"

// Audit action type constants.
const (
	// SHORTEN denotes the action of creating a short URL.
	SHORTEN = "shorten"
	// FOLLOW denotes the action of following a short URL to its target.
	FOLLOW = "follow"
)

// BaseAuditItem is the payload sent to the audit log on URL operations.
type BaseAuditItem struct {
	URL    string `json:"url"`
	Action string `json:"action"`
	UserId *int   `json:"user_id"`
	Ts     int64  `json:"ts"`
}

// CreateBaseAuditItemShorten creates an audit item for a shorten action.
func CreateBaseAuditItemShorten(url string, userId *int) BaseAuditItem {
	return BaseAuditItem{
		URL:    url,
		Action: SHORTEN,
		UserId: userId,
		Ts:     time.Now().Unix(),
	}
}

// CreateBaseAuditItemFollow creates an audit item for a follow action.
func CreateBaseAuditItemFollow(url string, userId *int) BaseAuditItem {
	return BaseAuditItem{
		URL:    url,
		Action: FOLLOW,
		UserId: userId,
		Ts:     time.Now().Unix(),
	}
}
