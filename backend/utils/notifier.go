package utils

import (
	"fmt"
	"careerconnect/backend/models"
)

// Mock notification system (replace with Twilio/Gupshup API)
func SendNotification(user models.Candidate, job models.Job, daysLeft int) {
	msg := fmt.Sprintf("🔔 Reminder: %s closes in %d days. Apply soon!", job.Title, daysLeft)
	fmt.Printf("Sending SMS to %s: %s\n", user.Phone, msg)
}
