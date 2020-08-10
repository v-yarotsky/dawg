package dawg

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/jtacoma/uritemplates"
)

type URITemplate struct {
	*uritemplates.UriTemplate
	Raw string
}

func (t *URITemplate) UnmarshalJSON(b []byte) error {
	t.Raw = string(b)
	tpl, err := uritemplates.Parse(t.Raw)
	if err != nil {
		return err
	}
	t.UriTemplate = tpl
	return nil
}

func (t *URITemplate) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%q", t.Raw)), nil
}

//
// {
//    "myservice": {
//      "template": "https://foo/bar/{id}/{kind}",
//      "substitutions": {
//        "myshortcut": {
//          "id": "5000",
//          "kind": "funny"
//        }
//      }
//    }
// }
type Config map[string]*ServiceConfig

func (c Config) GetService(service string) (ServiceConfig, error) {
	if s, ok := c[service]; ok {
		return *s, nil
	} else {
		return ServiceConfig{}, fmt.Errorf("service '%s' not found", service)
	}
}

type ServiceConfig struct {
	GUID          string                            `json:"-"`
	Template      URITemplate                       `json:"template"`
	Keyword       string                            `json:"keyword"`
	Substitutions map[string]map[string]interface{} `json:"substitutions"`
}

func (s ServiceConfig) GetURL(shortcut string) (string, error) {
	var shortcutVars map[string]interface{}
	var ok bool
	if shortcutVars, ok = s.Substitutions[shortcut]; !ok {
		return "", fmt.Errorf("shortcut '%s' not found", shortcut)
	}

	if expanded, err := s.Template.Expand(shortcutVars); err != nil {
		return "", fmt.Errorf("could not expand URL template: %v", err)
	} else {
		return expanded, nil
	}
}

func ReadConfig(path string) (Config, error) {
	file, err := os.Open(path)

	if os.IsNotExist(err) {
		err = putSampleConfig(path)
		if err != nil {
			return Config{}, fmt.Errorf("could not create sample config file: %s", err)
		}
		file, err = os.Open(path)
	}

	if err != nil {
		return Config{}, fmt.Errorf("could not open config file: %s", err)
	}
	defer file.Close()

	var cfg Config
	dec := json.NewDecoder(file)
	if err := dec.Decode(&cfg); err != nil {
		return Config{}, fmt.Errorf("could not parse config file: %s", err)
	}

	for svc, _ := range cfg {
		cfg[svc].GUID = GUID() // randomly assign guids for alfred workflow objects
	}

	return cfg, nil
}

func putSampleConfig(path string) error {
	data, _ := json.MarshalIndent(sampleConfig(), "", "  ")
	return ioutil.WriteFile(path, data, 0644)
}

func GUID() string {
	u := make([]byte, 16)
	_, err := rand.Read(u)
	if err != nil {
		panic(err)
	}
	u[6] = (u[6] & 0x0f) | 0x40 // Version 4
	u[8] = (u[8] & 0x3f) | 0x80 // Variant is 10
	return fmt.Sprintf("%X-%X-%X-%X-%X", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
}
