package dawg

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"net/http"
	"time"

	"golang.org/x/image/draw"

	"github.com/jtacoma/uritemplates"
	"github.com/nfnt/resize"
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
		return fmt.Errorf("could not fetch logo from %s: %s", i.URL, err)
	}
	defer res.Body.Close()

	m, _, err := image.Decode(res.Body)
	if err != nil {
		return fmt.Errorf("could not decode logo from %s: %s", i.URL, err)
	}

	m, err = generateLogoThumbnail(m)
	if err != nil {
		return fmt.Errorf("could not generate thumbnail for %s: %s", i.URL, err)
	}

	var out bytes.Buffer
	if err = png.Encode(&out, m); err != nil {
		return fmt.Errorf("could not encode thumbnail for logo %s: %s", i.URL, err)
	}

	i.Data = out.Bytes()

	return nil
}

func generateLogoThumbnail(m image.Image) (image.Image, error) {
	dst := image.NewRGBA(image.Rect(0, 0, 128, 128))
	resized := resize.Resize(128, 0, m, resize.Bicubic) // resize to width 128, preserve aspect ratio
	p := image.Point{
		(dst.Bounds().Dx() - resized.Bounds().Dx()) / 2.0,
		(dst.Bounds().Dy() - resized.Bounds().Dy()) / 2.0,
	}
	draw.Copy(dst, p, resized, resized.Bounds(), draw.Src, nil)
	return dst, nil
}

func DownloadServiceIcons(c Config) ([]ServiceIcon, error) {
	tpl, _ := uritemplates.Parse("https://logo.clearbit.com/{domain}?size=128")
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
