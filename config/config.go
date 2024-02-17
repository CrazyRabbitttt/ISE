package config

type Config struct {
	Addr  string `yaml:"addr"`
	Shard string `yaml:"shard"` // shard number
	Auth  string `yaml:"auth"`  // 权限认证

}
