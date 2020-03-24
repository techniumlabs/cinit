package cmd

type CinitConfig struct {
	SecretProviders []string       `mapstructure:"providers"`
	Templates       []TemplateDefs `mapstructure:"templates"`
}

type TemplateDefs struct {
	Source string `mapstructure:"source"`
	Dest   string `mapstructure:"dest"`
}
