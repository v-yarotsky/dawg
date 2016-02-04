package dawg

import (
	"crypto/rand"
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
type Config map[string]*ServiceConfig

func (c Config) GetService(service string) (ServiceConfig, error) {
	if s, ok := map[string]*ServiceConfig(c)[service]; ok {
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

	for svc, _ := range cfg {
		cfg[svc].GUID = GUID() // randomly assign guids for alfred workflow objects
	}

	return cfg
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
