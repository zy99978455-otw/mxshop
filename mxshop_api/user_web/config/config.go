package config

type UserSrvConfig struct {
	Host string `mapstructure:"host" json:"host" yaml:"host"`
	Port int    `mapstructure:"port" json:"port" yaml:"port"`
	Name string `mapstructure:"name" json:"name" yaml:"name"`
}

type JWTConfig struct {
	SigningKey string `mapstructure:"key" json:"key"  yaml:"key"`
}

// type AliSmsConfig struct {
// 	ApiKey     string `mapstructure:"key" json:"key"`
// 	ApiSecrect string `mapstructure:"secrect" json:"secrect"`
// }

type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host" yaml:"host"`
	Port int    `mapstructure:"port" json:"port" yaml:"port"`
}

// type RedisConfig struct {
// 	Host   string `mapstructure:"host" json:"host"`
// 	Port   int    `mapstructure:"port" json:"port"`
// 	Expire int    `mapstructure:"expire" json:"expire"`
// }

type ServerConfig struct {
	Name        string        `mapstructure:"name" json:"name" yaml:"name"`
	Host        string        `mapstructure:"host" json:"host" yaml:"host"`
	Tags        []string      `mapstructure:"tags" json:"tags" yaml:"tags"`
	Port        int           `mapstructure:"port" json:"port" yaml:"port"`
	UserSrvInfo UserSrvConfig `mapstructure:"user_srv" json:"user_srv" yaml:"user_srv"`
	JWTInfo     JWTConfig     `mapstructure:"jwt" json:"jwt" yaml:"jwt"`
	ConsulInfo  ConsulConfig  `mapstructure:"consul" json:"consul" yaml:"consul"`
}

type NacosConfig struct {
	Host      string `mapstructure:"host"`
	Port      uint64    `mapstructure:"port"`
	Namespace string `mapstructure:"namespace"`
	User      string `mapstructure:"user"`
	Password  string `mapstructure:"password"`
	DataId    string `mapstructure:"dataid"`
	Group     string `mapstructure:"group"`
}
