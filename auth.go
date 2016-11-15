package auth

import (
	"strconv"
	"time"

	"github.com/happierall/l"
	soso "github.com/happierall/soso-server"
)

var (
	Log = l.New()
)

func init() {
	Log.Prefix = l.Colorize("Auth ", l.Blue)
}

type Auth struct {
	Sign            string
	WithDefaultUser bool

	Router *soso.Engine
}

func New(sign string, WithDefaultUser bool, r *soso.Engine) *Auth {
	a := &Auth{sign, WithDefaultUser, r}

	if WithDefaultUser {
		initStore()
		a.Router.Middleware.Before(a.middleware)
	}

	return a
}

func (a *Auth) middleware(m *soso.Msg, start time.Time) {
	token, uid, err := ReadToken(m, a.Sign)
	if err != nil {
		Log.Debug("MiddlewareBefore, invalid token. ", err)
		return
	}

	if _, err := UsersData.Get(uid); err != nil {
		return
	}

	strID := strconv.FormatInt(uid, 10)
	m.User.ID = strID
	m.User.Token = token
	m.User.IsAuth = true

	// Register session
	soso.Sessions.Push(m.Session, strID)
}
