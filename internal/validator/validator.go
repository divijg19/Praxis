package validator

var valid = map[string]bool{
	"cursor":    true,
	"buffer":    true,
	"composite": true,
}

func Exists(name string) bool {
	return valid[name]
}
