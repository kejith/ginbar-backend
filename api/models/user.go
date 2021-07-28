package models

import (
	"time"

	"github.com/kejith/ginbar-backend/com/kejith/ginbar-backend/mysql/db"
)

// PublicUserJSON is a representation in JSON for the Database Object User
// it only stores values that are viable for public access
type PublicUserJSON struct {
	ID        int32     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
}

// Populate fills the struct from db.User
func (u *PublicUserJSON) Populate(user db.User) {
	u.ID = user.ID
	u.CreatedAt = user.CreatedAt
	u.UpdatedAt = user.UpdatedAt
	u.Name = user.Name
}
