package venomoid

import (
	"errors"
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
	Name               []string
	ConfigType         string
	Path               []string
	Defaults           map[string]interface{}
	ConfigFiles        []string
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
	if len(c.ConfigFiles) == 0 && c.ConfigLookup == false && c.AutomaticEnv == false {
		return ErrorLookupAndFileMismatchAndAutomaticEnv
	}

	viper.SetConfigType(c.ConfigType)

	if c.ConfigLookup {
		for _, path := range c.Path {
			viper.AddConfigPath(path)
		}
	}

	for key, value := range c.Defaults {
		viper.SetDefault(key, value)
	}

	if len(c.ConfigFiles) > 0 {
		for i, configFile := range c.ConfigFiles {
			f, err := os.Open(configFile)
			if err != nil && c.ErrorOnMissingFile {
				return &ErrorWrapper{
					InternalError: err,
					Label:         "error opening file",
				}
			}

			if i > 0 {
				if err := viper.MergeConfig(f); err != nil {
					_ = f.Close()
					return &ErrorWrapper{
						InternalError: err,
						Label:         "could not read from config file",
					}
				}
				_ = f.Close()
				continue
			}

			if err := viper.ReadConfig(f); err != nil {
				_ = f.Close()
				// no need to handle viper.ConfigFileNotFoundError since os.Open takes care of that
				return &ErrorWrapper{
					InternalError: err,
					Label:         "could not read from config file",
				}
			}
			_ = f.Close()
		}
	} else if c.ConfigLookup {
		for i, name := range c.Name {
			viper.SetConfigName(name)

			if i > 0 {
				if err := viper.MergeInConfig(); err != nil {
					var configFileNotFoundError viper.ConfigFileNotFoundError
					if errors.As(err, &configFileNotFoundError) {
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
				continue
			}

			if err := viper.ReadInConfig(); err != nil {
				var configFileNotFoundError viper.ConfigFileNotFoundError
				if errors.As(err, &configFileNotFoundError) {
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
		}
	}

	if c.AutomaticEnv {
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
	if name == "" {
		return c
	}

	for _, cName := range c.Name {
		if cName == name {
			return c
		}
	}
	c.Name = append(c.Name, name)
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
	if configFile == "" {
		return c
	}

	for _, cFile := range c.ConfigFiles {
		if cFile == configFile {
			return c
		}
	}
	c.ConfigFiles = append(c.ConfigFiles, configFile)
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
