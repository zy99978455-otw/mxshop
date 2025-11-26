package config


type ServerConfig struct {
	Name string `mapstructure:"name" json:"name" yaml:"name"`
	MysqlInfo  MysqlConfig  `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	ConsulInfo ConsulConfig `mapstructure:"consul" json:"consul" yaml:"consul"`
}

// Mysql连接配置
type MysqlConfig struct {
	Host     string `mapstructure:"host" json:"host" yaml:"host"`
	Port     int    `mapstructure:"port" json:"port" yaml:"port"`
	Name     string `mapstructure:"db" json:"db" yaml:"db"`
	User     string `mapstructure:"user" json:"user" yaml:"user"`
	Password string `mapstructure:"password" json:"password" yaml:"password"`
}

// Consul连接配置
type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host" yaml:"host"`
	Port int    `mapstructure:"port" json:"port" yaml:"port"`
}

// Nacos连接配置
type NacosConfig struct {
	Host      string `mapstructure:"host"`
	Port      uint64 `mapstructure:"port"`
	Namespace string `mapstructure:"namespace"`
	User      string `mapstructure:"user"`
	Password  string `mapstructure:"password"`
	DataId    string `mapstructure:"dataid"`
	Group     string `mapstructure:"group"`
}