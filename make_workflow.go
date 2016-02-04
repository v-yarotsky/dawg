package dawg

import (
	"archive/zip"
	"bytes"
	"crypto/rand"
	"fmt"
	"io"
	"os"

	"github.com/kardianos/osext"
)

func MakeWorkflowZIP(plist []byte) (*bytes.Buffer, error) {
	fname, err := osext.Executable()
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)

	dawgFile, err := os.Open(fname)
	defer dawgFile.Close()

	files := []struct {
		Name string
		Body io.Reader
	}{
		{"info.plist", bytes.NewReader(plist)},
		{"dawg", dawgFile},
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
	scriptFilterGUIDs := make([]string, 0, len(c))
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

	ypos := 130
	for service, serviceConfig := range c {
		guid := GUID()
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
		scriptFilterGUIDs = append(scriptFilterGUIDs, guid)

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
