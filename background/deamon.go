package background

import (
	"github.com/go-macaron/cache"
	"github.com/go-macaron/captcha"
	"github.com/go-macaron/csrf"
	"github.com/go-macaron/session"
	"gopkg.in/macaron.v1"
)

func daemon() {
	m := macaron.Classic()
	m.Use(session.Sessioner(session.Options{
		Provider:       "redis",
		ProviderConfig: "addr=127.0.0.1:6379",
	}))
	m.Use(macaron.Renderer())
	m.Use(csrf.Csrfer())
	m.Use(cache.Cacher())
	m.Use(captcha.Captchaer())

	m.Run()
}
