package api

import (
	"net/http"
	"strconv"

	db "github.com/danilshap/domains-generator/internal/db/sqlc"
	"github.com/danilshap/domains-generator/internal/models/view"
	"github.com/danilshap/domains-generator/internal/views/components/notifications"
	"github.com/danilshap/domains-generator/internal/views/layouts"
	"github.com/go-chi/chi/v5"
)

func (s *Server) handleListNotifications(w http.ResponseWriter, r *http.Request) {
	page := 1
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	user, err := getCurrentUser(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	const pageSize = 10
	offset := (page - 1) * pageSize

	notificationData, err := s.store.GetNotifications(r.Context(), db.GetNotificationsParams{
		UserID: user.UserID,
		Limit:  int32(pageSize),
		Offset: int32(offset),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	totalCount, err := s.store.GetNotificationsCount(r.Context(), user.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	unreadCount, err := s.store.GetUnreadNotificationsCount(r.Context(), user.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	notificationViews := make([]view.NotificationView, len(notificationData))
	for i, n := range notificationData {
		notificationViews[i] = view.NotificationView{
			ID:        int32(n.ID),
			Title:     n.Title,
			Message:   n.Message,
			Type:      string(n.Type),
			IsRead:    n.ReadAt.Valid,
			CreatedAt: n.CreatedAt,
		}
	}

	totalPages := (totalCount + int64(pageSize) - 1) / int64(pageSize)

	data := view.NotificationListData{
		Notifications: notificationViews,
		UnreadCount:   int32(unreadCount),
		CurrentPage:   int32(page),
		TotalPages:    int32(totalPages),
		PageSize:      int32(pageSize),
	}

	if r.Header.Get("HX-Request") == "true" {
		notifications.Page(data).Render(r.Context(), w)
		return
	}

	component := layouts.Base(notifications.Page(data))
	component.Render(r.Context(), w)
}

func (s *Server) handleMarkNotificationRead(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.store.MarkNotificationRead(r.Context(), int64(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleMarkAllNotificationsRead(w http.ResponseWriter, r *http.Request) {
	user, err := getCurrentUser(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	err = s.store.MarkAllNotificationsRead(r.Context(), user.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.handleListNotifications(w, r)
}
