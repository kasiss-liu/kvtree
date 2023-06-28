package useradmin

import (
	"time"

	"github.com/kasiss-liu/kvtree/apps/server/config"
)

type UserAuth struct {
	config.User
	NameSpaceAuth map[string]struct{}
}

func NewUserAuth(user *config.User) *UserAuth {
	ua := &UserAuth{}
	ua.Name = user.Name
	ua.Passwd = user.Passwd
	ua.Super = user.Super
	ua.Source = user.Source
	ua.NameSpaceAuth = make(map[string]struct{})
	for _, s := range ua.Source {
		ua.NameSpaceAuth[s] = struct{}{}
	}
	return ua
}

type UsersPermit struct {
	users map[string]*UserAuth
}

func NewUserPermit(users []*config.User) *UsersPermit {
	up := &UsersPermit{users: make(map[string]*UserAuth, len(users))}
	for _, user := range users {
		up.users[user.Name] = NewUserAuth(user)
	}
	return up
}

func (up *UsersPermit) Get(name string) *UserAuth {
	return up.users[name]
}

type UserLogin struct {
	Name    string
	Token   string
	Expired int
	Time    time.Time
}
