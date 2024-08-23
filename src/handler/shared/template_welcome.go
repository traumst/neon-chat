package shared

import (
	"prplchat/src/model/app"
	t "prplchat/src/model/template"
)

func TemplateWelcome(user *app.User) (string, error) {
	var welcome t.WelcomeTemplate
	if user == nil {
		welcome = t.WelcomeTemplate{User: t.UserTemplate{}}
	} else {
		welcome = t.WelcomeTemplate{User: user.Template(0, 0, 0)}
	}
	return welcome.HTML()
}
