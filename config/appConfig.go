package config

import (
	"os"
	"io/ioutil"
	yaml "gopkg.in/yaml.v2"
)

type AppConfig struct {
	ClusterURL      string `yaml:"clusterUrl"`
	Token           string `yaml:"token"`
	KafkaURL      string `yaml:"kakfaURL"`
	Interval        int    `yaml:"interval"`
	MetricResources []struct {
		Name string `yaml:"name"`
		Path string `yaml:"path"`
	} `yaml:"metricResources"`
}

func Config(configFilePath string) (AppConfig, error) {

	var SetConfig AppConfig
	path, _ := os.Getwd()
	reqYamlConfig, err := os.Open(path + configFilePath)
	defer reqYamlConfig.Close()
	if err != nil {
		return SetConfig, err
	}
	buf, err := ioutil.ReadAll(reqYamlConfig)
	if err != nil {
		return SetConfig, err
	}
	yaml.Unmarshal(buf, &SetConfig)
	return SetConfig, err
}
