package config

import (
	"io/ioutil"

	"github.com/SimonIshai/helloWorld/model"

	"github.com/SimonIshai/helloWorld/errors"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Settings struct {
		LogFolder string `yaml:"log_folder""`
		TenantID  string `yaml:"tenant_id""`
	} `yaml:"settings""`
	TestCases []model.TestCase `yaml:"test_cases""`
}

var cfg *Config

func Init(filePath string) (*Config, error) {

	if data, err := ioutil.ReadFile(filePath); err != nil {
		return nil, errors.WrapWithKind(err, errors.KindFileSystem, "read file")

	} else {
		//log.Println("config", string(data))
		cfg = &Config{}

		if err := yaml.Unmarshal(data, cfg); err != nil {
			return nil, errors.WrapWithKind(err, errors.KindParse, "parse config")

		} else {
			return cfg, nil
		}
	}
}

func GetConfig() (*Config, error) {
	if cfg == nil {
		return nil, errors.New(errors.KindConfig, "config is not initiated")
	}
	return cfg, nil
}
