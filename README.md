# Auth for Soso-server (Warning! Not maintained)

## Usage with store
```go
package main

import (
	"services"

	"github.com/happierall/l"
	mb "github.com/happierall/membase"
	auth "github.com/happierall/soso-auth"  
	soso "github.com/happierall/soso-server"
)

func main() {
	Router := soso.Default()

	auth.EnableDebug()
	authObj := auth.New("super-key", true, Router)

	auth.UseGithubAuth(
		authObj,
		"clientid",
		"secretid",
		[]string{"user:email"},
		"http://localhost:4000", // Then in github settings set: "http://localhost:4000/oauth/callback/github" 
	)

	Router.CREATE("product", func(m *soso.Msg) {
		if m.User.IsAuth {
			l.Print("create product")
		}
	})

	go Router.Run(4000)

	mb.Run() // Need for store
}
```

## Usage with custom store and user
```go
package main

import (
	"strconv"
	"time"

	"soso-auth"

	"github.com/happierall/l"
	soso "github.com/happierall/soso-server"
)

var lastID int64 = 0

type MyUser struct {
	ID    int64
	Name  string
	Email string
}

var Users = []*MyUser{}

// 1. Run and make request from client: soso.request("auth", "github")
// and you will get auth_url, open it in iframe
// 2. Wait onSuccess handler on server
// 3. Send token to user from success handler
func main() {
	Router := soso.Default()

	authObj := auth.New("super_duper_key", false, Router)

	githubAuth := auth.UseGithubAuth(
		authObj,
		"clientID",
		"clientSecret",
		[]string{"user:email"},
		"http://localhost:4000", // Then in github settings set: "http://localhost:4000/oauth/callback/github"
	)

	githubAuth.OnSuccess = func(userData *auth.User, session soso.Session) {

		user := &MyUser{
			Name:  userData.Name,
			Email: userData.Email,
		}

		// Try auth current
		for _, u := range Users {
			if u.Email == user.Email {
				successResponce(u, authObj.Sign, session)
				return
			}
		}

		// Or register
		lastID++
		user.ID = lastID
		Users = append(Users, user)

		successResponce(user, authObj.Sign, session)
	}

	// Read token from every request (other.token)
	Router.Middleware.Before(func(m *soso.Msg, start time.Time) {
		token, tokenData, err := auth.ReadToken(m, authObj.Sign)
		if err != nil {
			return
		}

		for _, u := range Users {
			if u.ID == tokenData.UID {

				strID := strconv.FormatInt(tokenData.UID, 10)
				m.User.ID = strID
				m.User.Token = token
				m.User.IsAuth = true
				m.User.IsAnonymous = tokenData.IsAnonymous

				// Register session
				soso.Sessions.Push(m.Session, strID)

			}
		}

	})

	Router.SEARCH("user", func(m *soso.Msg) {
		if m.User.IsAuth {

			uid, _ := strconv.Atoi(m.User.ID)

			for _, user := range Users {
				if user.ID == int64(uid) {

					l.Logf("User email: %s", user.Email)

				}
			}
		}
	})

	l.Print("Running app at localhost:4000")
	Router.Run(4000)
}

func successResponce(user *MyUser, sign string, session soso.Session) {
	authToken := auth.CreateToken(map[string]interface{}{
		"uid":         user.ID,
		"isAnonymous": false,
	}, sign)
	soso.SendMsg("auth", "SUCCESS", session, map[string]interface{}{
		"token": authToken,
		"user":  user,
	})
}
```

### License
[MIT](http://opensource.org/licenses/MIT)
