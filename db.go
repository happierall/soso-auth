package auth

import "github.com/happierall/membase"

var (
	MB        = membase.New()
	UsersData = Users{}
)

func init() {
	// membase.EnableDebug()
	MB.StoreFolder = "./authmem/"

	MB.ListenUnique(&UsersData, "auth_users")

	Log.Logf("Users count %d", len(UsersData.List))

	go MB.Run()
}
