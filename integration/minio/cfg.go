package minio

type Config struct {
	AK       string `mapstructure:"ak"`
	SK       string `mapstructure:"sk"`
	Endpoint string `mapstructure:"endpoint"`
	Bucket   string `mapstructure:"bucket"`
	Token    string `mapstructure:"token"`
	UseSSL   bool   `mapstructure:"use_ssl"`
}
