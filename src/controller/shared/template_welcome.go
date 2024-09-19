package shared

import (
	"neon-chat/src/model/app"
	"neon-chat/src/model/template"
)

func TemplateWelcome(user *app.User) (string, error) {
	var welcome template.WelcomeTemplate
	if user == nil {
		welcome = template.WelcomeTemplate{User: template.UserTemplate{}}
	} else {
		welcome = template.WelcomeTemplate{User: user.Template(0, 0, 0)}
	}
	return welcome.HTML()
}
