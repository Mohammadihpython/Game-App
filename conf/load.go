package conf

import (
	"fmt"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"path/filepath"
	"strings"
)

func Load() Config {
	var k = koanf.New(".")

	//	 load default values using the confmap provider
	// A nested map can be loaded by setting the delimiter to an empty ""
	k.Load(confmap.Provider(defaultConfig, "."), nil)

	// Load Yaml config  and Merge it
	absPath, _ := filepath.Abs("../../config.yml")
	fmt.Println("Loading config from:", absPath)

	err := k.Load(file.Provider(absPath), yaml.Parser())
	if err != nil {
		panic(err)
	}
	var cfg Config

	k.Load(env.Provider("GAMEAPP_", ".", func(s string) string {
		str := strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, "GAMEAPP_")), "_", ".", -1)

		//fo multiword items such as "sign_key" that we shoud use like "GAMEAPP_AUTH_SIGN__KEY"
		// TODO Find a better solution if needed
		return strings.Replace(str, "..", "_", -1)

	}), nil)

	if err := k.Unmarshal("", &cfg); err != nil {
		panic(err)
	}
	fmt.Println(cfg)

	return cfg
}
