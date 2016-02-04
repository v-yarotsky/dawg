package dawg

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/jtacoma/uritemplates"
)

type ServiceIcon struct {
	Service string
	GUID    string
	URL     string
	Data    []byte
}

func (i *ServiceIcon) fetch() error {
	res, err := http.Get(i.URL)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	iconData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	i.Data = iconData
	return nil
}

func DownloadServiceIcons(c Config) ([]ServiceIcon, error) {
	tpl, _ := uritemplates.Parse("https://logo.clearbit.com/{domain}?size=128&format=png")
	nIcons := len(c)
	ic := make(chan ServiceIcon)
	ec := make(chan error)

	icons := make([]ServiceIcon, 0, nIcons)
	for serviceName, serviceCfg := range c {
		url, _ := tpl.Expand(map[string]interface{}{"domain": serviceName})
		icon := ServiceIcon{
			Service: serviceName,
			GUID:    serviceCfg.GUID,
			URL:     url,
		}
		go func(i ServiceIcon) {
			err := i.fetch()
			if err != nil {
				ec <- err
			} else {
				ic <- i
			}
		}(icon)
	}

	for nIcons > 0 {
		select {
		case i := <-ic:
			icons = append(icons, i)
			nIcons--
		case err := <-ec:
			return []ServiceIcon{}, err
		case <-time.After(60 * time.Second):
			return []ServiceIcon{}, fmt.Errorf("Timed out while downloading service icons")
		}
	}
	return icons, nil
}
