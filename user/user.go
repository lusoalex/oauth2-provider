package user

import "oauth2-provider/models"

type UserManager interface {
	MatchingCredentials(login, password string) (*models.User, bool)
}

type DefaultUserManager struct{}

func (*DefaultUserManager) MatchingCredentials(login, password string) (*models.User, bool) {
	switch login {
	case "alexandre":
		switch password {
		case "padawan":
			return &models.User{Name: "Jedi", Firstname: "Alexandre"}, true
		default:
			return nil, false
		}
	case "damien":
		switch password {
		case "master":
			return &models.User{Name: "Sith", Firstname: "Damien"}, true
		default:
			return nil, false
		}
	case "health_check":
		switch password {
		case "check_health":
			return &models.User{Name: "Check", Firstname: "Health"}, true
		default:
			return nil, false
		}
	default:
		return nil, false
	}
}
