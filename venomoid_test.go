package venomoid

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

type testConfig struct {
	KeyString string `yaml:"keystring"`
	KeyBool   bool   `yaml:"keybool"`
	KeyInt    int    `yaml:"keyint"`
}

func TestConfigBuilder_BuildWithDefaults(t *testing.T) {
	c := Config()
	assert.Equal(t, c.ErrorOnMissingFile, true, "ErrorOnMissingFile default is not true")
	assert.Equal(t, c.ConfigLookup, true, "ConfigLookup default is not true")
	assert.Equal(t, c.ConfigType, "yaml", "ConfigType default is not yaml")
}

func TestConfigBuilder_BuildNoLookupOrFile(t *testing.T) {
	c := Config()
	config := &testConfig{}

	err := c.WithName("test-config").WithConfigLookup(false).WithFile("").Build(config)
	assert.Equal(t, ErrorLookupAndFileMismatchAndAutomaticEnv, err, "unexpected error")
}

func TestConfigBuilder_BuildWithConfigfileBadFile(t *testing.T) {
	tempFile, _ := ioutil.TempFile("/tmp", "venomoid_temp")
	defer os.Remove(tempFile.Name())

	output := []byte("JUNKInFile@#$%@@:dfadfn\n\\-n\nqefrkf%%W@$\n")
	_, _ = tempFile.Write(output)

	config := &testConfig{}

	c := Config()
	err := c.WithName("test-config").
		WithFile(tempFile.Name()).
		WithType("yaml").
		WithErrorOnMissing(true).
		WithConfigLookup(false).
		Build(config)

	assert.Equal(t,
		"could not read from config file",
		err.(*ErrorWrapper).Label,
		"wrong error message")

}

func TestConfigBuilder_BuildWithLookupBadFile(t *testing.T) {
	tempDir, _ := ioutil.TempDir("/tmp", "venmoid_tempd")
	defer os.RemoveAll(tempDir)

	output := []byte("dkv;jdf;j$@#%@$dfdafas\\-fvcfden\ndnslds: sdkjs\\nff--\\n")
	fileName := fmt.Sprintf("%s/%s", tempDir, "testconfig.yaml")
	_ = os.WriteFile(fileName, output, 0644)
	defer os.Remove(fileName)

	paths := []string{tempDir}

	config := &testConfig{}

	c := Config()
	err := c.WithName("testconfig").
		WithPath(paths).
		WithType("yaml").
		WithErrorOnMissing(true).
		WithConfigLookup(true).
		Build(config)

	assert.Equal(t,
		"could not read from config file",
		err.(*ErrorWrapper).Label,
		"wrong error message")

}

func TestConfigBuilder_BuildWithLookupMissingFile(t *testing.T) {
	tempDir, _ := ioutil.TempDir("/tmp", "venmoid_tempd")
	defer os.RemoveAll(tempDir)

	paths := []string{tempDir}
	config := &testConfig{}

	c := Config()
	err := c.WithName("testconfig").
		WithPath(paths).
		WithType("yaml").
		WithErrorOnMissing(true).
		WithConfigLookup(true).
		Build(config)

	assert.Equal(t, ErrorMissingConfigFile, err, "unexpected error")
}

func TestConfigBuilder_BuildWithLookup(t *testing.T) {
	tempDir, _ := ioutil.TempDir("/tmp", "venmoid_tempd")
	defer os.RemoveAll(tempDir)

	output := []byte("---\nkeystring: \"string\"\nkeybool: true\n")
	fileName := fmt.Sprintf("%s/%s", tempDir, "testconfig.yaml")
	_ = os.WriteFile(fileName, output, 0644)
	defer os.Remove(fileName)

	paths := []string{tempDir}
	defaults := map[string]interface{}{
		"keyint": 5,
	}

	config := &testConfig{}

	c := Config()
	err := c.WithName("testconfig").
		WithPath(paths).
		WithType("yaml").
		WithErrorOnMissing(true).
		WithConfigLookup(true).
		WithDefaults(defaults).
		Build(config)

	assert.NoError(t, err, "did not expect an error")
	assert.Equal(t, true, config.KeyBool, "boolean key mismatch")
	assert.Equal(t, "string", config.KeyString, "string key mismatch")
	assert.Equal(t, 5, config.KeyInt, "int key mismatch. default did not load")

}

func TestConfigBuilder_BuildWithLookupMultiple(t *testing.T) {
	tempDir, _ := ioutil.TempDir("/tmp", "venmoid_tempd*")
	defer os.RemoveAll(tempDir)

	output1 := []byte("---\nkeystring: \"string\"\n")
	output2 := []byte("---\nkeybool: true\n")

	fileName1 := filepath.Join(tempDir, "testconfig1.yaml")
	fileName2 := filepath.Join(tempDir, "testconfig2.yaml")

	_ = os.WriteFile(fileName1, output1, 0644)
	defer os.Remove(fileName1)

	_ = os.WriteFile(fileName2, output2, 0644)
	defer os.Remove(fileName2)

	paths := []string{tempDir}
	defaults := map[string]interface{}{
		"keyint": 5,
	}

	config := &testConfig{}

	c := Config()
	err := c.WithName("testconfig1").WithName("testconfig2").
		WithPath(paths).
		WithType("yaml").
		WithErrorOnMissing(true).
		WithConfigLookup(true).
		WithDefaults(defaults).
		Build(config)

	assert.NoError(t, err, "did not expect an error")
	assert.Equal(t, true, config.KeyBool, "boolean key mismatch")
	assert.Equal(t, "string", config.KeyString, "string key mismatch")
	assert.Equal(t, 5, config.KeyInt, "int key mismatch. default did not load")
}

func TestConfigBuilder_BuildWithConfigFile(t *testing.T) {
	tempFile, _ := ioutil.TempFile("/tmp", "venomoid_tempf")
	defer os.Remove(tempFile.Name())

	output := []byte("---\nkeystring: \"string\"\nkeybool: true\n")
	_, _ = tempFile.Write(output)

	defaults := map[string]interface{}{
		"keyint": 5,
	}

	config := &testConfig{}

	c := Config()
	err := c.WithName("test-config").
		WithFile(tempFile.Name()).
		WithType("yaml").
		WithErrorOnMissing(true).
		WithConfigLookup(false).
		WithDefaults(defaults).
		Build(config)

	assert.NoError(t, err, "did not expect an error")
	assert.Equal(t, true, config.KeyBool, "boolean key mismatch")
	assert.Equal(t, "string", config.KeyString, "string key mismatch")
	assert.Equal(t, 5, config.KeyInt, "int key mismatch. default did not load")
}

func TestConfigBuilder_BuildWithMultipleConfigFiles(t *testing.T) {
	tempDir, _ := ioutil.TempDir("/tmp", "venmoid_tempd*")
	defer os.RemoveAll(tempDir)

	output1 := []byte("---\nkeystring: \"string\"\n")
	output2 := []byte("---\nkeybool: true\n")

	fileName1 := filepath.Join(tempDir, "testconfig1.yaml")
	fileName2 := filepath.Join(tempDir, "testconfig2.yaml")

	_ = os.WriteFile(fileName1, output1, 0644)
	defer os.Remove(fileName1)

	_ = os.WriteFile(fileName2, output2, 0644)
	defer os.Remove(fileName2)

	defaults := map[string]interface{}{
		"keyint": 5,
	}

	config := &testConfig{}

	c := Config()
	err := c.WithName("test-config").
		WithFile(fileName1).WithFile(fileName2).
		WithType("yaml").
		WithErrorOnMissing(true).
		WithConfigLookup(false).
		WithDefaults(defaults).
		Build(config)

	assert.NoError(t, err, "did not expect an error")
	assert.Equal(t, true, config.KeyBool, "boolean key mismatch")
	assert.Equal(t, "string", config.KeyString, "string key mismatch")
	assert.Equal(t, 5, config.KeyInt, "int key mismatch. default did not load")
}

func TestConfigBuilder_BuildWithMissingConfigFileFail(t *testing.T) {
	config := &testConfig{}

	c := Config()
	err := c.WithName("test-config").
		WithFile("some not existing file").
		WithType("yaml").
		WithErrorOnMissing(true).
		WithConfigLookup(false).
		Build(config)

	assert.Equal(t, "error opening file", err.(*ErrorWrapper).Label, "unexpect error received")

}
func TestConfigBuilder_BuildWithMissingConfigFileOk(t *testing.T) {
	config := &testConfig{}

	c := Config()
	err := c.WithName("test-config").
		WithFile("some not existing file").
		WithType("yaml").
		WithErrorOnMissing(false).
		WithConfigLookup(false).
		Build(config)

	assert.NoError(t, err, "unexpected error")
}

func TestConfigBuilder_BuildWithMissingAutomaticEnvError(t *testing.T) {
	config := &testConfig{}

	c := Config()
	err := c.WithName("test-config").
		WithType("yaml").
		WithErrorOnMissing(false).
		WithConfigLookup(false).
		WithAutomaticEnv(false).
		Build(config)

	assert.Error(t, err, ErrorLookupAndFileMismatchAndAutomaticEnv)
}

func TestConfigBuilder_WithAutomaticEnv(t *testing.T) {
	config := &testConfig{}
	_ = os.Setenv("KEYSTRING", "key_value")

	c := Config()
	err := c.WithName("test-config").
		WithType("yaml").
		WithErrorOnMissing(false).
		WithConfigLookup(false).
		WithAutomaticEnv(true).
		Build(config)

	assert.Equal(t, "key_value", viper.Get("keyString"))
	assert.NotEqual(t, "key_value", config.KeyString)
	assert.NoError(t, err)
}

func TestConfigBuilder_WithBindEnv(t *testing.T) {
	config := &testConfig{}
	_ = os.Setenv("KEYSTRING", "key_value")

	c := Config()
	err := c.WithName("test-config").
		WithType("yaml").
		WithErrorOnMissing(false).
		WithConfigLookup(false).
		WithAutomaticEnv(true).
		WithBindEnv("keystring").
		Build(config)

	assert.Equal(t, "key_value", viper.Get("keyString"))
	assert.Equal(t, "key_value", config.KeyString)
	assert.NoError(t, err)
}

func TestConfigBuilder_WithEnvPrefix(t *testing.T) {
	config := &testConfig{}
	_ = os.Setenv("PRE_KEYSTRING", "key_value")

	c := Config()
	err := c.WithName("test-config").
		WithType("yaml").
		WithErrorOnMissing(false).
		WithConfigLookup(false).
		WithAutomaticEnv(true).
		WithBindEnv("keystring").
		WithEnvPrefix("pre").
		Build(config)

	assert.Equal(t, "key_value", viper.Get("keyString"))
	assert.Equal(t, "key_value", config.KeyString)
	assert.NoError(t, err)
}
