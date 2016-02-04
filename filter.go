package dawg

import (
	"sort"
	"strconv"
	"strings"
)

func Filter(c Config, svc, pat string) (AlfredOutput, error) {
	serviceConfig, err := c.GetService(svc)
	if err != nil {
		return AlfredOutput{}, err
	}

	alfredOut := make(AlfredOutput, 0, 10)
	for shortcut, _ := range serviceConfig.Substitutions {
		matchPos := strings.Index(shortcut, pat)
		if matchPos == -1 {
			continue
		}
		url, err := serviceConfig.GetURL(shortcut)
		if err != nil {
			return AlfredOutput{}, err
		}
		unquotedURL, _ := strconv.Unquote(url)
		alfredOut = append(alfredOut, AlfredOutputItem{
			UID:          "dawg:" + shortcut,
			Autocomplete: shortcut,
			Title:        shortcut,
			Arg:          unquotedURL,
			Pos:          matchPos,
		})
	}
	sort.Sort(alfredOut)
	return alfredOut, nil
}
