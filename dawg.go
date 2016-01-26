package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"os"
	"os/user"
	"sort"
	"strconv"
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

type AlfredOutputItems struct {
	XMLName xml.Name           `xml:"items"`
	Items   []AlfredOutputItem `xml:"item"`
}

func (s AlfredOutputItems) Len() int {
	return len(s.Items)
}
func (s AlfredOutputItems) Swap(i, j int) {
	s.Items[i], s.Items[j] = s.Items[j], s.Items[i]
}
func (s AlfredOutputItems) Less(i, j int) bool {
	return s.Items[i].pos < s.Items[j].pos
}

type AlfredOutputItem struct {
	XMLName      xml.Name `xml:"item"`
	UID          string   `xml:"uid,attr"`
	Valid        string   `xml:"valid,attr"`
	Autocomplete string   `xml:"autocomplete,attr"`
	Title        string   `xml:"title"`
	Arg          string   `xml:"arg"`
	pos          int
}

//
// <?xml version="1.0"?>
// <items>
//   <item uid="#{q dashboard_name}" valid="YES" autocomplete="#{q dashboard_name}">
//     <title>#{dashboard_name}</title>
//     <arg>https://app.datadoghq.com/screen/#{dashboard_id}</arg>
//   </item>
// </items>
//
func main() {
	var service string
	flag.StringVar(&service, "s", "", "Service name")
	flag.Parse()
	pattern := flag.Arg(0)

	c := mustReadConfig()
	serviceConfig, err := c.GetService(service)
	if err != nil {
		panic(err)
	}

	alfredOut := AlfredOutputItems{Items: make([]AlfredOutputItem, 0, 10)}
	for shortcut, _ := range serviceConfig.Substitutions {
		matchPos := strings.Index(shortcut, pattern)
		if matchPos == -1 {
			continue
		}
		url, err := serviceConfig.GetURL(shortcut)
		if err != nil {
			panic(err)
		}
		unquotedURL, _ := strconv.Unquote(url)
		alfredOut.Items = append(alfredOut.Items, AlfredOutputItem{
			UID:          shortcut,
			Valid:        "YES",
			Autocomplete: shortcut,
			Title:        shortcut,
			Arg:          unquotedURL,
			pos:          matchPos,
		})
	}
	sort.Sort(alfredOut)

	xmlout, err := xml.MarshalIndent(alfredOut, "", "  ")
	if err != nil {
		panic(err)
	}
	xmlWithHeader := xml.Header + string(xmlout)
	fmt.Println(xmlWithHeader)
}

func mustReadConfig() Config {
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
