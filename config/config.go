package config

type ConfigOptionError struct {
	msg string
}

func (e *ConfigOptionError) Error() string {
	return e.msg
}

var config = map[string]string{}

func SetOption(name, value string) {
	config[name] = value
}

func GetOption(name string) string {
	v := config[name]
	return v
}
