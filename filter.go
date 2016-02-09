package dawg

import (
	"fmt"
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
			UID:          fmt.Sprintf("dawg:%s:%s", svc, shortcut),
			Autocomplete: shortcut,
			Title:        shortcut,
			Arg:          unquotedURL,
			Pos:          matchPos,
			Icon:         fmt.Sprintf("./%s.png", svc),
		})
	}
	sort.Sort(alfredOut)
	return alfredOut, nil
}
