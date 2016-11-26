package auth

import (
	"github.com/happierall/l"
	soso "github.com/happierall/soso-server"
	"golang.org/x/oauth2"
	googleConf "golang.org/x/oauth2/google"
	"google.golang.org/api/plus/v1"
)

// scopes: email
func UseGoogleAuth(
	auth *Auth,
	clientID, clientSecret string, scopes []string, redirectURL string) *googleAuth {

	g := &googleAuth{}

	g.Name = "google"
	g.Auth = auth
	g.ClientID = clientID
	g.ClientSecret = clientSecret
	g.Scopes = scopes
	// g.RedirectURL = "http://localhost:4000/oauth/callback/google"
	g.RedirectURL = redirectURL + "/oauth/callback/" + g.Name

	g.Sessions = soso.NewSessionList()

	g.CallbackHandler = g.callbackGoogle
	g.Endpoint = googleConf.Endpoint

	defer g.Handle()
	return g
}

type googleAuth struct {
	Base
}

func (g *googleAuth) callbackGoogle(session soso.Session) {

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: g.Token.AccessToken},
	)
	client := oauth2.NewClient(oauth2.NoContext, ts)

	service, err := plus.New(client)
	if err != nil {
		l.Error(err, "Create Plus Client", 500)
		return
	}

	people := service.People.Get("me")
	data, err := people.Do()
	if err != nil {
		l.Error(err)
	}

	email := ""
	name := data.DisplayName

	for _, e := range data.Emails {
		if e.Type == "account" {
			email = e.Value
		}
	}

	g.registerUser(name, email, session)
}
