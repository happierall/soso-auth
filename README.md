# Auth for Soso-server (ouath2 github)

## Usage

```go
package main

import (
	"services"

	"github.com/happierall/l"
	mb "github.com/happierall/membase"
	auth "github.com/happierall/soso-auth"  
	soso "github.com/happierall/soso-server"
)

func init() {

}

func main() {

	soso.EnableDebug()

	Router := soso.Default()

	authObj := auth.New("super-key", true, Router)

	auth.UseGithubAuth(
		authObj,
		"clientid",
		"secretid",
		[]string{"user:email"},
	)

	Router.CREATE("product", func(m *soso.Msg) {
    if m.User.IsAuth {
      l.Print("create product")
    }
  })

	go Router.Run(4000)
	l.Print("Running app at localhost:4000")

	mb.Run()
}
```