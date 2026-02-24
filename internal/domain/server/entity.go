package server

import "time"

type Server struct {
	ID              string
	UserID          string
	Name            string
	Host            string
	Port            int
	Username        string
	SSHKeyEncrypted string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
