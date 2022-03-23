package configs

import (
	"time"

	"ty/car-prices-master/pkg/env"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var config = new(Config)

type Config struct {
	BasePath string `toml:"BasePath"`
	HomeUrl  string `toml:"HomeUrl"`
	MySQL    struct {
		Read struct {
			Addr string `toml:"addr"`
			User string `toml:"user"`
			Pass string `toml:"pass"`
			Name string `toml:"name"`
		} `toml:"read"`
		Write struct {
			Addr string `toml:"addr"`
			User string `toml:"user"`
			Pass string `toml:"pass"`
			Name string `toml:"name"`
		} `toml:"write"`
		Base struct {
			MaxOpenConn     int           `toml:"maxOpenConn"`
			MaxIdleConn     int           `toml:"maxIdleConn"`
			ConnMaxLifeTime time.Duration `toml:"connMaxLifeTime"`
		} `toml:"base"`
	} `toml:"mysql"`

	Redis struct {
		Addr         string `toml:"addr"`
		Pass         string `toml:"pass"`
		Db           int    `toml:"db"`
		MaxRetries   int    `toml:"maxRetries"`
		PoolSize     int    `toml:"poolSize"`
		MinIdleConns int    `toml:"minIdleConns"`
	} `toml:"redis"`

	Mongodb struct {
		Addr string   `toml:"address"`
		DBs  []string `toml:"dbs"`
	} `toml:"mongodb"`

	Mail struct {
		Host string `toml:"host"`
		Port int    `toml:"port"`
		User string `toml:"user"`
		Pass string `toml:"pass"`
		To   string `toml:"to"`
	} `toml:"mail"`

	JWT struct {
		Secret         string        `toml:"secret"`
		ExpireDuration time.Duration `toml:"expireDuration"`
	} `toml:"jwt"`

	URLToken struct {
		Secret         string        `toml:"secret"`
		ExpireDuration time.Duration `toml:"expireDuration"`
	} `toml:"urlToken"`

	HashIds struct {
		Secret string `toml:"secret"`
		Length int    `toml:"length"`
	} `toml:"hashids"`

	DataCenter struct {
		Host   string `toml:"host"`
		AppID  string `toml:"appid"`
		AppKey string `toml:"appkey"`
	} `toml:"datacenter"`
	DataCenterV2 struct {
		Host  string `toml:"host"`
		AppID int    `toml:"appid"`
	} `toml:"datacenterv2"`
	AccessToken struct {
		Addr   string `toml:"addr"`
		Host   string `toml:"host"`
		AppID  string `toml:"appid"`
		AppKey string `toml:"appkey"`
	} `toml:"accessToken"`
	KafkaES struct {
		Addr  []string `toml:"addr"`
		Topic string   `toml:"topic"`
	} `toml:"kafkaES"`
	KafkaSQM struct {
		Addr  []string `toml:"addr"`
		Topic string   `toml:"topic"`
	} `toml:"kafkaSQM"`
	H5UsrServiceUrl  string `toml:"H5UsrServiceUrl"`
	H5UsrInvoiceUrl  string `toml:"H5UsrInvoiceUrl"`
	YoupinOrderUrl   string `toml:"YoupinOrderUrl"`
	YoupinHomeUrl    string `toml:"YoupinHomeUrl"`
	MyAppointmentUrl string `toml:"MyAppointmentUrl"`
	Api2Host         string `toml:"Api2Host"`
	Api2Url          string `toml:"Api2Url"`
	Fwpthost         string `toml:"fwpthost"`
	WeatherH5Url     string `toml:"WeatherH5Url"`
	ZuulUrl          string `toml:"ZuulUrl"`
}

func init() {
	viper.SetConfigName(env.Active().Value() + "_configs")
	viper.SetConfigType("toml")
	viper.AddConfigPath("./configs")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := viper.Unmarshal(config); err != nil {
		panic(err)
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		if err := viper.Unmarshal(config); err != nil {
			panic(err)
		}
	})
}

func Get() Config {
	return *config
}
