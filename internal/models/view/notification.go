package view

import "time"

type NotificationView struct {
	ID        int32
	Title     string
	Message   string
	Type      string // success, error, warning, info
	IsRead    bool
	CreatedAt time.Time
}

type NotificationListData struct {
	Notifications []NotificationView
	UnreadCount   int32
	CurrentPage   int32
	TotalPages    int32
	PageSize      int32
}
