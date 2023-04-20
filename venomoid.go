package venomoid

import (
	"github.com/spf13/viper"
	"os"
)

const (
	defaultErrorOnMissingFile = true
	defaultConfigLookup       = true
	defaultConfigType         = "yaml"
	defaultAutomaticEnv       = false
)

type ConfigBuilder struct {
	Name               string
	ConfigType         string
	Path               []string
	Defaults           map[string]interface{}
	ConfigFile         string
	ConfigLookup       bool
	AutomaticEnv       bool
	BindEnv            []string
	EnvPrefix          string
	ErrorOnMissingFile bool
}

func Config() *ConfigBuilder {
	return &ConfigBuilder{
		ConfigLookup:       defaultConfigLookup,
		ErrorOnMissingFile: defaultErrorOnMissingFile,
		ConfigType:         defaultConfigType,
		AutomaticEnv:       defaultAutomaticEnv,
	}
}

func (c *ConfigBuilder) Build(destStruct interface{}) error {
	if c.ConfigFile == "" && c.ConfigLookup == false && c.AutomaticEnv == false {
		return ErrorLookupAndFileMismatchAndAutomaticEnv
	}
	viper.SetConfigName(c.Name)
	viper.SetConfigType(c.ConfigType)

	if c.ConfigLookup {
		for _, path := range c.Path {
			viper.AddConfigPath(path)
		}
	}

	for key, value := range c.Defaults {
		viper.SetDefault(key, value)
	}

	if c.ConfigFile != "" {
		f, err := os.Open(c.ConfigFile)
		if err != nil && c.ErrorOnMissingFile {
			return &ErrorWrapper{
				InternalError: err,
				Label:         "error opening file",
			}
		}
		defer f.Close()

		if err := viper.ReadConfig(f); err != nil {
			// no need to handle viper.ConfigFileNotFoundError since os.Open takes care of that
			return &ErrorWrapper{
				InternalError: err,
				Label:         "could not read from config file",
			}

		}
	} else if c.ConfigLookup {
		if err := viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				if c.ErrorOnMissingFile {
					return ErrorMissingConfigFile
				}
			} else {
				return &ErrorWrapper{
					InternalError: err,
					Label:         "could not read from config file",
				}
			}
		}
	} else if c.AutomaticEnv {
		viper.AutomaticEnv()

		if c.EnvPrefix != "" {
			viper.SetEnvPrefix(c.EnvPrefix)
		}

		if len(c.BindEnv) != 0 {
			err := viper.BindEnv(c.BindEnv...)
			if err != nil {
				return err
			}
		}
	}

	return viper.Unmarshal(&destStruct)
}

func (c *ConfigBuilder) WithName(name string) *ConfigBuilder {
	c.Name = name
	return c
}

func (c *ConfigBuilder) WithType(fileType string) *ConfigBuilder {
	c.ConfigType = fileType
	return c
}

func (c *ConfigBuilder) WithPath(paths []string) *ConfigBuilder {
	c.Path = paths
	return c
}

func (c *ConfigBuilder) WithDefaults(defaults map[string]interface{}) *ConfigBuilder {
	c.Defaults = defaults
	return c
}

func (c *ConfigBuilder) WithFile(configFile string) *ConfigBuilder {
	c.ConfigFile = configFile
	return c
}

func (c *ConfigBuilder) WithConfigLookup(configLookup bool) *ConfigBuilder {
	c.ConfigLookup = configLookup
	return c
}

func (c *ConfigBuilder) WithErrorOnMissing(eom bool) *ConfigBuilder {
	c.ErrorOnMissingFile = eom
	return c
}

func (c *ConfigBuilder) WithAutomaticEnv(automaticEnv bool) *ConfigBuilder {
	c.AutomaticEnv = automaticEnv
	return c
}

func (c *ConfigBuilder) WithBindEnv(input ...string) *ConfigBuilder {
	for _, ip := range input {
		c.BindEnv = append(c.BindEnv, ip)
	}
	return c
}

func (c *ConfigBuilder) WithEnvPrefix(prefix string) *ConfigBuilder {
	c.EnvPrefix = prefix
	return c
}
