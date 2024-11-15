package model

import (
	rms_users "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-users"
	uuid "github.com/satori/go.uuid"
)

// User represents one device's user
type User struct {
	// ID is a device/API key
	ID string `gorm:"primaryKey"`

	// Name is username
	Name *string

	// Info stores short comment about device
	Info string

	// TelegramUserId is associated Telegram user
	TelegramUserId *int32

	// Permissions
	Permissions int32

	// Domain
	Domain string
}

// GenerateID automation generation of user ID
func (u *User) GenerateID() {
	u.ID = uuid.NewV4().String()
}

var permissionsMap = map[rms_users.Permissions]int32{
	rms_users.Permissions_Search:             1,
	rms_users.Permissions_ConnectingToTheBot: 1 << 1,
	rms_users.Permissions_AccountManagement:  1 << 2,
	rms_users.Permissions_SendNotifications:  1 << 3,
	rms_users.Permissions_ListeningMusic:     1 << 4,
}

func (u *User) SetPermissions(perms []rms_users.Permissions) {
	u.Permissions = 0
	for _, p := range perms {
		u.Grant(p)
	}
}

func (u *User) IsAllowed(perm rms_users.Permissions) bool {
	return u.Permissions&permissionsMap[perm] != 0
}

func (u *User) Grant(perm rms_users.Permissions) {
	u.Permissions |= permissionsMap[perm]
}

func (u *User) GetPermissions() []rms_users.Permissions {
	var result []rms_users.Permissions
	for p, m := range permissionsMap {
		if u.Permissions&m != 0 {
			result = append(result, p)
		}
	}
	return result
}
