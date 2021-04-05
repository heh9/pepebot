package config

import (
	"io/ioutil"

	"github.com/hashicorp/hcl"
)

type ConfMap struct {
	Discord DiscordConfigMap `hcl:"discord,block"`
	Steam   SteamConfigMap   `hcl:"steam,block"`
	DB      DBConfig         `hcl:"db,block"`
	Sounds  SoundsConfig     `hcl:"sounds,block"`
}

type DBConfig struct {
	Type string `hcl:"type"`
	Path string `hcl:"path"`
}

type SoundsConfig struct {
	Type  string   `hcl:"type"`
	Path  string   `hcl:"path"`
	Runes []string `hcl:"runes"`
	Win   []string `hcl:"win"`
	Loss  []string `hcl:"loss"`
}

type DiscordConfigMap struct {
	Token string `hcl:"token"`
}

type SteamConfigMap struct {
	WebApiToken string `hcl:"web_api_token"`
}

var Map = new(ConfMap)

func LoadFile(filename string) (err error) {
	d, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	obj, err := hcl.Parse(string(d))
	if err != nil {
		return err
	}
	// Build up the result
	if err := hcl.DecodeObject(&Map, obj); err != nil {
		return err
	}
	return
}
