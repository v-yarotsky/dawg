package dawg

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"strings"

	"github.com/jtacoma/uritemplates"
)

type URITemplate uritemplates.UriTemplate

func (t *URITemplate) UnmarshalJSON(b []byte) error {
	tpl, err := uritemplates.Parse(string(b))
	if err != nil {
		return err
	}
	*t = URITemplate((*tpl))
	return nil
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
type Config map[string]ServiceConfig

func (c Config) GetService(service string) (ServiceConfig, error) {
	if s, ok := map[string]ServiceConfig(c)[service]; ok {
		return s, nil
	} else {
		return ServiceConfig{}, fmt.Errorf("service '%s' not found", service)
	}
}

type ServiceConfig struct {
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

	casted := uritemplates.UriTemplate(s.Template)
	if expanded, err := casted.Expand(shortcutVars); err != nil {
		return "", fmt.Errorf("could not expand URL template: %v", err)
	} else {
		return expanded, nil
	}
}

func MustReadConfig() Config {
	usr, _ := user.Current()
	dir := usr.HomeDir
	file, err := os.Open(strings.Join([]string{dir, ".dawg.json"}, "/"))
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var cfg Config
	dec := json.NewDecoder(file)
	if err := dec.Decode(&cfg); err != nil {
		panic(err)
	}
	return cfg
}
