package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	DbConnect          string `default:"sampleuser:samplepassword@db/sampledb" desc:"Db connection string"`
	Listen             string `default:"0.0.0.0:8080"`
	AdminEmail         string `default:"admin@example.com"`
	AdminPass          string `default:"adminpass"`
	TokenSecret        string `default:"samplesecret" desc:"secret to generate access JWT tokens"`
	RefreshTokenSecret string `default:"samplerefreshsecret" desc:"secret to generate refresh JWT tokens"`
	LogLevel           string `default:"debug"`
	PrettyLog          bool   `default:"true"`
}

var C Config

func InitConfig(appPrefix string) {
	envconfig.MustProcess(appPrefix, &C)
}
