package upyun

type Config struct {
	Bucket    string            `mapstructure:"bucket"`
	Operator  string            `mapstructure:"operator"`
	Password  string            `mapstructure:"password"`
	Hosts     map[string]string `mapstructure:"hosts"`
	UserAgent string            `mapstructure:"user_agent"`
}
