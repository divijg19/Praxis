package validator

var valid = map[string]bool{
	"cursor": true,
}

func Exists(name string) bool {
	return valid[name]
}
