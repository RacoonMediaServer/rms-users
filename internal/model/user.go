package model

import uuid "github.com/satori/go.uuid"

// User represents one device's user
type User struct {
	// ID is a device/API key
	ID string `gorm:"primaryKey"`

	// Info stores short comment about device
	Info string

	// Admin means the user have admin privileges
	Admin bool
}

// GenerateID automation generation of user ID
func (u *User) GenerateID() {
	u.ID = uuid.NewV4().String()
}
