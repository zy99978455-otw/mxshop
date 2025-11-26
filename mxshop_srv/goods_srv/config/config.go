package config

// 定义mysql连接的配置
type MysqlConfig struct{
	Host string `mapstructure:"host" json:"host" yaml:"host"`
	Port int    `mapstructure:"port" json:"port" yaml:"port"`
	Name string `mapstructure:"db" json:"db" yaml:"name"`
	User string `mapstructure:"user" json:"user" yaml:"user"`
	Password string `mapstructure:"password" json:"password" yaml:"password"`
}

// 定义consul的配置
type ConsulConfig struct{
	Host string `mapstructure:"host" json:"host" yaml:"host"`
	Port int    `mapstructure:"port" json:"port" yaml:"port"`
}

// 定义ElasticSearch的配置信息
type EsConfig struct{
	Host string `mapstructure:"host" json:"host" yaml:"host"`
	Port int    `mapstructure:"port" json:"port" yaml:"port"`
}

type ServerConfig struct{
	Name string `mapstructure:"name" json:"name" yaml:"name"`
	Host string `mapstructure:"host" json:"host" yaml:"host"`
	Tags []string `mapstructure:"tags" json:"tags" yaml:"tags"`
	MysqlInfo MysqlConfig `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	ConsulInfo ConsulConfig `mapstructure:"consul" json:"consul" yaml:"consul"`
	EsInfo EsConfig `mapstructure:"es" json:"es" yaml:"es"`
}

// 定义nacos的配置信息
type NacosConfig struct {
	Host      string `mapstructure:"host"`
	Port      uint64    `mapstructure:"port"`
	Namespace string `mapstructure:"namespace"`
	User      string `mapstructure:"user"`
	Password  string `mapstructure:"password"`
	DataId    string `mapstructure:"dataid"`
	Group     string `mapstructure:"group"`
}