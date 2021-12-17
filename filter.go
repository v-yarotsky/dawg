package dawg

import (
	"fmt"
	"strconv"
)

func Filter(c Config, svc, pat string) (AlfredOutput, error) {
	serviceConfig, err := c.GetService(svc)
	if err != nil {
		return AlfredOutput{}, err
	}

	alfredOut := make(AlfredOutput, 0, 10)

	filteredShortcuts := FilterChoices(serviceConfig.Shortcuts(), pat)
	for _, shortcut := range filteredShortcuts {
		url, err := serviceConfig.GetURL(shortcut)
		if err != nil {
			return AlfredOutput{}, err
		}
		unquotedURL, _ := strconv.Unquote(url)
		alfredOut = append(alfredOut, AlfredOutputItem{
			UID:          fmt.Sprintf("dawg:%s:%s", svc, shortcut),
			Autocomplete: shortcut,
			Title:        shortcut,
			Arg:          unquotedURL,
			Icon:         fmt.Sprintf("./%s.png", svc),
		})
	}
	return alfredOut, nil
}
