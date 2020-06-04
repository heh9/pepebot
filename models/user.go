package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	ID                  *primitive.ObjectID   `bson:"_id, omitempty" json:"id, omitempty"`
	Fullname            string                `json:"fullname"  form:"fullname"  validate:"required"`
	Username            string                `json:"username"  form:"username"  validate:"required,unique=users"`
	Email               string                `json:"email"     form:"email"     validate:"required,email,unique=users"`
	Password            string                `json:"-"         form:"password"  validate:"required"`
	Avatar              string                `json:"avatar"`
	IsActive            bool                  `json:"is_active"`
	Token               string                `json:"token"`
	RoleId              *primitive.ObjectID   `json:"role_id"`
	LastLogin           time.Time             `json:"last_login"`
	JoinedAt            time.Time             `json:"joined_at"`
	UpdatedAt           time.Time             `json:"updated_at"`
}

type SearchUser struct {
	Query     string    `json:"query" form:"query" validate:"required"`
}

// Validate users's password
func (u *User) ValidatePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// Hash users's password with bcrypt
func (User) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return string(bytes), err
}

// set users password
func (u *User) SetPassword(password string) (err error) {
	u.Password, err = u.HashPassword(password)
	return
}
