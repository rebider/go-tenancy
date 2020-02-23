package auth

import (
	"time"

	"GoTenancy/config"
	"GoTenancy/config/bindatafs"
	"GoTenancy/config/db"
	"GoTenancy/libs/auth"
	"GoTenancy/libs/auth/authority"
	"GoTenancy/libs/auth/providers/facebook"
	"GoTenancy/libs/auth/providers/github"
	"GoTenancy/libs/auth/providers/google"
	"GoTenancy/libs/auth/providers/twitter"
	"GoTenancy/libs/auth_themes/clean"
	"GoTenancy/libs/render"
	"GoTenancy/models/users"
)

var (

	// Auth initialize Auth for Authentication
	Auth = clean.New(&auth.Config{
		DB:         db.DB,
		Mailer:     config.Mailer,
		Render:     render.New(&render.Config{AssetFileSystem: bindatafs.AssetFS.NameSpace("auth")}),
		UserModel:  users.User{},
		Redirector: auth.Redirector{RedirectBack: config.RedirectBack},
		ViewPaths:  append([]string{}, "GoTenancy/libs/auth_themes/clean/views"),
	})

	// Authority initialize Authority for Authorization
	Authority = authority.New(&authority.Config{
		Auth: Auth,
	})
)

func init() {
	//Auth.RegisterProvider(password.New(&password.Config{}))
	Auth.RegisterProvider(github.New(&config.Config.Github))
	Auth.RegisterProvider(google.New(&config.Config.Google))
	Auth.RegisterProvider(facebook.New(&config.Config.Facebook))
	Auth.RegisterProvider(twitter.New(&config.Config.Twitter))

	Authority.Register("logged_in_half_hour", authority.Rule{TimeoutSinceLastLogin: time.Minute * 30})

}
