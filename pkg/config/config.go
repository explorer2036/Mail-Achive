package config

import (
	"Mail-Achive/pkg/model"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

// server fields
type serverStruct struct {
	ListenAddr    string   `yaml:"listen_addr"`
	ElasticURL    string   `yaml:"elastic_url"`
	FirebaseCreds string   `yaml:"firebase_creds"`
	DocumentName  string   `yaml:"document_name"`
	MatchFields   []string `yaml:"es_match_fields"`
}

// log fields
type logStruct struct {
	OutputLevel        string `yaml:"output_level"`
	OutputPath         string `yaml:"output_path"`
	RotationPath       string `yaml:"rotation_path"`
	RotationMaxSize    int    `yaml:"rotation_max_size"`
	RotationMaxAge     int    `yaml:"rotation_max_age"`
	RotationMaxBackups int    `yaml:"rotation_max_backups"`
	JSONEncoding       bool   `yaml:"json_encoding"`
}

// Config structure for server
type Config struct {
	Server serverStruct          `yaml:"server"`
	Users  map[string]model.User `yaml:"users"`
	Log    logStruct             `yaml:"log"`
}

// ParseYamlFile the config file
func ParseYamlFile(filename string, c *Config) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, c)
}
