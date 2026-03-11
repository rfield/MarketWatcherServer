package notification

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/protobuf/types/known/timestamppb"
	"rjfield.com/backend/db"
	"rjfield.com/backend/generated/pb"
)

type NotificationServer struct {
	pb.UnimplementedNotificationServiceServer
}

// ListNotifications retrieves a list of notifications for the specified user.
func (s *NotificationServer) ListNotifications(ctx context.Context, req *pb.ListNotificationsRequest) (*pb.ListNotificationsReply, error) {
	log.Printf("ListNotifications() - received: %v", req.GetParent())

	user_id := db.UserIDFromResourceName(req.GetParent())

	notifications, err := db.ListNotifications(user_id)
	if err != nil {
		return nil, fmt.Errorf("unable to list notifications: %v", err)
	}
	log.Printf("ListNotifications() - found %d notifications for user ID: %s", len(notifications), user_id)
	return &pb.ListNotificationsReply{
		Notifications: notifications,
	}, nil
}

// UpdateNotification updates the status of a notification, such as marking it as read or archived.
func (s *NotificationServer) UpdateNotification(ctx context.Context, req *pb.UpdateNotificationRequest) (*pb.UpdateNotificationReply, error) {
	log.Printf("UpdateNotification() - received: %v", req.GetName())
	log.Printf("UpdateNotification() - STUB IMPLEMETATION ONLY ****")
	// Implementation for updating a notification
	return &pb.UpdateNotificationReply{
		Notification: &pb.Notification{
			Name:      req.GetName(),
			Type:      pb.Notification_ACCOUNT_ALERT,
			Title:     "Stock Price Alert",
			Message:   "Your stock price alert has been triggered.",
			Status:    pb.Notification_UNREAD,
			CreatedAt: timestamppb.Now(),
		},
	}, nil
}
