package dawg

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/jtacoma/uritemplates"
)

type serviceIcon struct {
	guid string
	data []byte
}

func DownloadServiceIcons(c Config) (map[string][]byte, error) {
	nIcons := len(c)
	ic := make(chan serviceIcon)
	ec := make(chan error)

	tpl, err := uritemplates.Parse("https://logo.clearbit.com/{domain}?size=128&format=png")
	if err != nil {
		return map[string][]byte{}, err
	}

	for serviceName, serviceCfg := range c {
		url, _ := tpl.Expand(map[string]interface{}{"domain": serviceName})
		go fetchServiceIcon(url, serviceCfg.GUID, ic, ec)
	}

	icons := make(map[string][]byte, nIcons)
	for nIcons > 0 {
		select {
		case i := <-ic:
			nIcons--
			icons[i.guid] = i.data
		case err := <-ec:
			return map[string][]byte{}, err
		case <-time.After(60 * time.Second):
			return map[string][]byte{}, fmt.Errorf("Timed out while downloading service icons")
		}
	}
	return icons, nil
}

func fetchServiceIcon(url, guid string, results chan serviceIcon, errc chan error) {
	res, err := http.Get(url)
	if err != nil {
		errc <- err
		return
	}
	iconData, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		errc <- err
		return
	}
	results <- serviceIcon{data: iconData, guid: guid}
}
