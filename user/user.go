package user

type User struct {
	login     string
	password  string
	Name      string
	Firstname string
}

func MatchingCredentials(login, password string) (*User, bool) {
	switch login {
	case "alexandre":
		switch password {
		case "padawan":
			return &User{Name: "Jedi", Firstname: "Alexandre"}, true
		default:
			return nil, false
		}
	case "damien":
		switch password {
		case "master":
			return &User{Name: "Sith", Firstname: "Damien"}, true
		default:
			return nil, false
		}
	default:
		return nil, false
	}
}
