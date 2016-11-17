package auth

import (
	"github.com/google/go-github/github"
	soso "github.com/happierall/soso-server"
	"golang.org/x/oauth2"
	githubConf "golang.org/x/oauth2/github"
)

// scopes: user:email
func UseGithubAuth(
	auth *Auth,
	clientID, clientSecret string, scopes []string) *githubAuth {

	g := &githubAuth{}

	g.Name = "github"
	g.Auth = auth
	g.ClientID = clientID
	g.ClientSecret = clientSecret
	g.Scopes = scopes
	g.RedirectURL = "http://localhost:4000/oauth/callback/" + g.Name

	g.Sessions = soso.NewSessionList()

	g.CallbackHandler = g.callbackGithub
	g.Endpoint = githubConf.Endpoint

	defer g.Handle()
	return g
}

type githubAuth struct {
	Base
}

func (g *githubAuth) callbackGithub(session soso.Session) {

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: g.Token.AccessToken},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	client := github.NewClient(tc)
	data, _, err := client.Users.Get("")
	if err != nil {
		Log.Error(err)
		return
	}

	name := *data.Name
	email := *data.Email

	g.registerUser(name, email, session)
}
