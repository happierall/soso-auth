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
	Log.Level = l.LevelInfo
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
	token, td, err := ReadToken(m, a.Sign)
	if err != nil {
		Log.Debug("MiddlewareBefore: ", err)
		return
	}

	if _, err := UsersData.Get(td.UID); err != nil {
		return
	}

	strID := strconv.FormatInt(td.UID, 10)
	m.User.ID = strID
	m.User.Token = token
	m.User.IsAuth = true
	m.User.IsAnonymous = td.IsAnonymous

	// Register session
	soso.Sessions.Push(m.Session, strID)
}

func EnableDebug() {
	Log.Level = l.LevelDebug
}

func DisableDebug() {
	Log.Level = l.LevelInfo
}
