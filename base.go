package auth

import (
	"context"
	"net/http"

	soso "github.com/happierall/soso-server"
	"github.com/rs/xid"
	"golang.org/x/oauth2"
)

type Base struct {
	Auth *Auth

	Name         string
	ClientID     string
	ClientSecret string
	Scopes       []string
	Endpoint     oauth2.Endpoint
	RedirectURL  string

	Sessions soso.SessionList

	OnSuccess func(user *User, session soso.Session, authType string)
	OnError   func(error, soso.Session)

	CallbackHandler func(soso.Session)

	Token *oauth2.Token
}

func (b *Base) handler(m *soso.Msg) {
	m.Success(map[string]interface{}{
		"url": b.authUrl(m),
	})
}

func (b *Base) Handle() {
	if b.Name == "" {
		Log.Crit("Base.Name can't be empty")
	}
	b.Auth.Router.Handle("auth", b.Name, b.handler)
	http.HandleFunc("/oauth/callback/"+b.Name, b.Callback)

	// Listen Default Sessions,
	// because they have event onClose in router
	soso.Sessions.OnClose(b.Sessions.OnCloseExecute)
}

func (b *Base) conf() *oauth2.Config {
	conf := &oauth2.Config{
		ClientID:     b.ClientID,
		ClientSecret: b.ClientSecret,
		Scopes:       b.Scopes,
		Endpoint:     b.Endpoint,
	}

	if b.RedirectURL != "" {
		conf.RedirectURL = b.RedirectURL
	}

	return conf
}

func (b *Base) authUrl(m *soso.Msg) string {
	uid := xid.New().String()
	b.Sessions.Push(m.Session, uid)

	return b.conf().AuthCodeURL(uid, oauth2.AccessTypeOffline)
}

func (b *Base) Callback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	uid := r.URL.Query().Get("state")

	if code != "" && uid != "" {
		ctx := context.Background()

		var session soso.Session
		sessions := b.Sessions.Get(uid)

		if len(sessions) > 0 {
			session = sessions[0]
		}

		token, err := b.conf().Exchange(ctx, code)
		if err != nil {
			Log.Error(err, session.ID(), uid)
			b.OnError(err, session)
			return
		}

		b.Token = token

		if b.CallbackHandler != nil {

			b.CallbackHandler(session)

		}

	}
}

func (b *Base) registerUser(user User, session soso.Session) {

	// Run handler onSuccess and exit
	if !b.Auth.WithDefaultUser && b.OnSuccess != nil {
		b.OnSuccess(&user, session, b.Name)
		return
	}

	// or run default mechanic to save user
	for _, u := range UsersData.List {
		if (user.Email != "" && u.Email == user.Email) ||
			(user.GithubID != 0 && u.GithubID == user.GithubID) ||
			(user.GoogleID != "" && u.GoogleID == user.GoogleID) {

			u.ServiceToken = b.Token.AccessToken

			if b.OnSuccess != nil {
				b.OnSuccess(u, session, b.Name)
			}

			b.successResponse(u, session)
			return
		}
	}

	UsersData.Create(&user)

	if b.OnSuccess != nil {
		b.OnSuccess(&user, session, b.Name)
	}

	b.successResponse(&user, session)
}

func (b *Base) successResponse(user *User, session soso.Session) {
	authToken := CreateToken(map[string]interface{}{
		"UID":         user.ID,
		"IsAnonymous": false,
	}, b.Auth.Sign)

	soso.SendMsg("auth", "SUCCESS", session, map[string]interface{}{
		"token": authToken,
		"user":  user,
	})
}
