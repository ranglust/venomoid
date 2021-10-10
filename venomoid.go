package venomoid

import (
	"github.com/spf13/viper"
	"os"
)

const (
	defaultErrorOnMissingFile = true
	defaultConfigLookup       = true
	defaultConfigType         = "yaml"
)

type ConfigBuilder struct {
	name               string
	configType         string
	path               []string
	defaults           map[string]interface{}
	configFile         string
	configLookup       bool
	errorOnMissingFile bool
}

func Config() *ConfigBuilder {
	return &ConfigBuilder{
		configLookup:       defaultConfigLookup,
		errorOnMissingFile: defaultErrorOnMissingFile,
		configType:         defaultConfigType,
	}
}

func (c *ConfigBuilder) Build(destStruct interface{}) error {
	if c.configFile == "" && c.configLookup == false {
		return ErrorLookupAndFileMismatch
	}
	viper.SetConfigName(c.name)
	viper.SetConfigType(c.configType)

	if c.configLookup {
		for _, path := range c.path {
			viper.AddConfigPath(path)
		}
	}

	for key, value := range c.defaults {
		viper.SetDefault(key, value)
	}

	if c.configFile != "" {
		f, err := os.Open(c.configFile)
		if err != nil && c.errorOnMissingFile {
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
	} else if c.configLookup {
		if err := viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				if c.errorOnMissingFile {
					return ErrorMissingConfigFile
				}
			} else {
				return &ErrorWrapper{
					InternalError: err,
					Label:         "could not read from config file",
				}
			}
		}
	}

	return viper.Unmarshal(&destStruct)
}

func (c *ConfigBuilder) WithName(name string) *ConfigBuilder {
	c.name = name
	return c
}

func (c *ConfigBuilder) WithType(fileType string) *ConfigBuilder {
	c.configType = fileType
	return c
}

func (c *ConfigBuilder) WithPath(paths []string) *ConfigBuilder {
	c.path = paths
	return c
}

func (c *ConfigBuilder) WithDefaults(defaults map[string]interface{}) *ConfigBuilder {
	c.defaults = defaults
	return c
}

func (c *ConfigBuilder) WithFile(configFile string) *ConfigBuilder {
	c.configFile = configFile
	return c
}

func (c *ConfigBuilder) WithConfigLookup(configLookup bool) *ConfigBuilder {
	c.configLookup = configLookup
	return c
}

func (c *ConfigBuilder) WithErrorOnMissing(eom bool) *ConfigBuilder {
	c.errorOnMissingFile = eom
	return c
}
