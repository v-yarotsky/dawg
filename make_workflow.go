package dawg

import (
	"archive/zip"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"

	"github.com/kardianos/osext"
)

func MakeWorkflowZIP(plist []byte, icons []ServiceIcon) (*bytes.Buffer, error) {
	fname, err := osext.Executable()
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)

	dawgFile, err := os.Open(fname)
	defer dawgFile.Close()

	iconReader, _ := gzip.NewReader(bytes.NewReader(iconPNG))

	type File struct {
		Name string
		Body io.Reader
	}
	files := []File{
		{"info.plist", bytes.NewReader(plist)},
		{"icon.png", iconReader},
		{"dawg", dawgFile},
	}

	for _, icon := range icons {
		// Workflow objects
		files = append(files, File{
			fmt.Sprintf("%s.png", icon.GUID),
			bytes.NewReader(icon.Data),
		})

		// XML results
		files = append(files, File{
			fmt.Sprintf("%s.png", icon.Service),
			bytes.NewReader(icon.Data),
		})
	}

	for _, file := range files {
		f, err := w.Create(file.Name)
		if err != nil {
			return nil, err
		}
		_, err = io.Copy(f, file.Body)
		if err != nil {
			return nil, err
		}
	}

	// Make sure to check the error on Close.
	if err = w.Close(); err != nil {
		return nil, err
	}
	return buf, nil
}

func MakeWorkflowPList(c Config) PList {
	objects := make(PArray, 0, len(c)+1)
	connections := make(PDict, len(c))
	uidata := make(PDict, len(c)+1)

	openURLObjectGUID := GUID()
	openURLObject := PDict{
		"config": PDict{
			"plusspaces": PBool(false),
			"url":        PString("{query}"),
			"utf8":       PBool(true),
		},
		"type":    PString("alfred.workflow.action.openurl"),
		"uid":     PString(openURLObjectGUID),
		"version": PInteger(0),
	}
	objects = append(objects, openURLObject)
	uidata[openURLObjectGUID] = PDict{
		"ypos": PReal(10),
	}

	ypos := 10
	for service, serviceConfig := range c {
		guid := serviceConfig.GUID
		obj := PDict{
			"config": PDict{
				"argumenttype":     PInteger(1),
				"escaping":         PInteger(102),
				"keyword":          PString(serviceConfig.Keyword),
				"queuedelaycustom": PBool(true),
				"queuedelaymode":   PInteger(0),
				"queuemode":        PInteger(1),
				"script":           PString(fmt.Sprintf("chmod +x ./dawg && ./dawg -s %s \"{query}\"", service)),
				"title":            PString(service),
				"type":             PInteger(0),
				"withspace":        PBool(true),
			},
			"type":    PString("alfred.workflow.input.scriptfilter"),
			"uid":     PString(guid),
			"version": PInteger(0),
		}
		objects = append(objects, obj)

		connections[guid] = PArray{
			PDict{
				"destinationuid":  PString(openURLObjectGUID),
				"modifiers":       PInteger(0),
				"modifiersubtext": PString(""),
			},
		}

		uidata[guid] = PDict{
			"ypos": PReal(ypos),
		}
		ypos += 120
	}

	plist := PList{
		"bundleid":    PString(""),
		"category":    PString("Productivity"),
		"connections": connections,
		"createdby":   PString("Vlad Yarotsky"),
		"description": PString("Serices at your fingertips :)"),
		"disabled":    PBool(false),
		"name":        PString("DAWG"),
		"objects":     objects,
		"readme":      PString(""),
		"uidata":      uidata,
		"webaddress":  PString("https://github.com/v-yarotsky/dawg"),
	}
	return plist
}
