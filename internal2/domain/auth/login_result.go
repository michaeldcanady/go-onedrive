package auth

import "time"

type LoginResult struct {
	AccessToken string
	RecordSaved bool
	ExpiresOn   time.Time
	Username    string
	Profile     string
}
