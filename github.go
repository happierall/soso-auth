package auth

import (
	"github.com/google/go-github/github"
	soso "github.com/happierall/soso-server"
	"golang.org/x/oauth2"
	githubConf "golang.org/x/oauth2/github"
)

// scopes: user:email
// redirectURL can be empty ""
func UseGithubAuth(
	auth *Auth,
	clientID, clientSecret string, scopes []string, redirectURL string) *githubAuth {

	g := &githubAuth{}

	g.Name = "github"
	g.Auth = auth
	g.ClientID = clientID
	g.ClientSecret = clientSecret
	g.Scopes = scopes

	if redirectURL != "" {
		// RedirectURL = "http://localhost:4000/oauth/callback/github"
		g.RedirectURL = redirectURL + "/oauth/callback/" + g.Name
	}

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

	user := User{}

	user.GithubID = *data.ID

	if data.Name != nil {
		user.Name = *data.Name
	} else if data.Login != nil {
		user.Name = *data.Login
	}

	if data.Email != nil {
		user.Email = *data.Email
	}

	g.registerUser(user, session)
}
