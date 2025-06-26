package obs

type Config struct {
	AK           string `mapstructure:"ak"`
	SK           string `mapstructure:"sk"`
	Endpoint     string `mapstructure:"endpoint"`
	Bucket       string `mapstructure:"bucket"`
	ACLOwnerID   string `mapstructure:"acl_owner_id"`
	ACLControlID string `mapstructure:"acl_control_id"`
}
