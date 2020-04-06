package config

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	SecretProviders []string       `mapstructure:"providers"`
	Templates       []TemplateDefs `mapstructure:"templates"`
}

type TemplateDefs struct {
	Source string `mapstructure:"source"`
	Dest   string `mapstructure:"dest"`
}

func Load(cfgFile string) (*Config, error) {
	var err error
	v := viper.New()
	if cfgFile != "" {
		// Use config file from the flag.
		v.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		v.AddConfigPath(".")
		v.AddConfigPath(home + "/")
		v.SetConfigName(".cinit")
	}

	v.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err = v.ReadInConfig(); err == nil {
		log.Infof("Using config file: %s", v.ConfigFileUsed())
		c := &Config{}
		err = v.Unmarshal(c)
		return c, err
	} else {
		log.Warnf("%s", err.Error())
		return nil, err
	}

}
