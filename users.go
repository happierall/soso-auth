package auth

import "errors"

type Users struct {
	List []*User

	lastID int64
}

func (u *Users) Get(id int64) (*User, error) {
	for _, user := range u.List {
		if user.ID == id {
			return user, nil
		}
	}
	return nil, errors.New("User not found")
}

func (u *Users) Create(user *User) {
	u.lastID++
	user.ID = u.lastID

	u.List = append(u.List, user)
}

func (r *Users) Remove(id int64) {
	for key, item := range r.List {
		if item.ID == id {

			copy(r.List[key:], r.List[key+1:])
			r.List[len(r.List)-1] = nil
			r.List = r.List[:len(r.List)-1]

		}
	}
}

func (r *Users) Flush() {
	r.List = []*User{}
}
