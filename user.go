package auth

import "strconv"

type User struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`

	GithubID int    `json:"githubID"`
	GoogleID string `json:"googleID"`

	ServiceToken string `json:"-"`
}

func (u *User) StringID() string {
	return strconv.FormatInt(u.ID, 10)
}
