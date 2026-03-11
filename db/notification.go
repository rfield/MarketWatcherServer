package db

import (
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
	"google.golang.org/protobuf/types/known/timestamppb"
	"rjfield.com/backend/generated/pb"
)

// listNotifications fetches all notifications for a given user ID.
func listNotifications(user_id string) ([]*pb.Notification, error) {
	log.Printf("listNotifications() - listing notifications for user ID: %s", user_id)
	// SQL statement with placeholders ($1, $2, etc. for Postgres)
	sqlStatement := `
		SELECT id, type, title, message, status, created_at FROM notifications WHERE user_id = $1`
	rows, err := db.Query(sqlStatement, user_id)
	if err != nil {
		return nil, fmt.Errorf("unable to query rows: %v", err)
	}
	defer rows.Close()

	var notifications []*pb.Notification

	for rows.Next() {
		var notification pb.Notification
		var notificationId string
		var timestamp time.Time
		err := rows.Scan(&notificationId, &notification.Type, &notification.Title, &notification.Message, &notification.Status, &timestamp)
		if err != nil {
			return nil, fmt.Errorf("unable to scan row: %v", err)
		}
		notification.Name = "users/" + user_id + "/notifications/" + notificationId
		notification.CreatedAt = timestamppb.New(timestamp)
		notifications = append(notifications, &notification)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %v", err)
	}

	log.Printf("listNotifications() - found %d notifications for user ID: %s", len(notifications), user_id)
	return notifications, nil
}

// updateNotification updates the status of a notification for a given user ID and notification ID.
func updateNotification(user_id string, notification_id string, status pb.Notification_NotificationStatus) (*pb.Notification, error) {
	log.Printf("updateNotification() - updating notification ID: %s for user ID: %s with status: %v", notification_id, user_id, status)
	// SQL statement with placeholders ($1, $2, etc. for Postgres)
	sqlStatement := `
		UPDATE notifications SET status = $1 WHERE id = $2 AND user_id = $3 RETURNING id, type, title, message, created_at`
	var notification pb.Notification
	err := db.QueryRow(sqlStatement, status, notification_id, user_id).Scan(&notification.Name, &notification.Type, &notification.Title, &notification.Message, &notification.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("unable to update notification: %v", err)
	}
	notification.Name = "users/" + user_id + "/notifications/" + notification_id
	log.Printf("updateNotification() - updated notification ID: %s for user ID: %s with status: %v", notification_id, user_id, status)
	return &notification, nil
}
