package config

import (
	"log"

	"gopkg.in/ini.v1"
)

type ConfigList struct {
	// web
	Port string
	// path
	StoragePath string
	RootCAPath  string
	SSLPath     string
}

var Config ConfigList

func init() {
	LoadConfig()
}

func LoadConfig() {
	cfg, err := ini.Load("C:\\Users\\MH\\go\\src\\golang_gRPC\\gRPC\\config.ini")

	if err != nil {
		log.Fatalln(err)
	}

	Config = ConfigList{
		// web
		Port: cfg.Section("web").Key("port").MustString("50051"),
		// path
		StoragePath: cfg.Section("path").Key("storage_path").String(),
		RootCAPath:  cfg.Section("path").Key("rootCA_path").String(),
		SSLPath:     cfg.Section("path").Key("ssl_path").String(),
	}
}
