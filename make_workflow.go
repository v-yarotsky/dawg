package dawg

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

func SelfUpdate(plist []byte, icons []ServiceIcon) (err error) {
	root := "."

	plistFilename := path.Join(root, "info.plist")
	if _, err = os.Stat(plistFilename); os.IsNotExist(err) {
		return fmt.Errorf("PWD is not alfred workflow root. Cowardly refusing to proceed")
	}

	if err = ioutil.WriteFile(plistFilename, plist, 0644); err != nil {
		return
	}

	fs, _ := ioutil.ReadDir(root)
	for _, f := range fs {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".png") {
			os.Remove(f.Name())
		}
	}

	if err = ioutil.WriteFile("icon.png", iconPNG, 0644); err != nil {
		return
	}

	for _, icon := range icons {
		// Workflow objects
		filename := path.Join(root, fmt.Sprintf("%s.png", icon.GUID))
		if err = ioutil.WriteFile(filename, icon.Data, 0644); err != nil {
			return
		}

		// XML results
		filename = path.Join(root, fmt.Sprintf("%s.png", icon.Service))
		if err = ioutil.WriteFile(filename, icon.Data, 0644); err != nil {
			return
		}
	}
	return nil
}

func MakeWorkflowPList(c Config) PList {
	objects := make(PArray, 0)
	connections := make(PDict)
	uidata := make(PDict)

	selfUpdateKeywordObjectGUID := GUID()
	selfUpdateKeywordObject := PDict{
		"config": PDict{
			"argumenttype": PInteger(2),
			"keyword":      PString("dawg edit"),
			"text":         PString("Edit DAWG configuration"),
			"withspace":    PBool(false),
		},
		"type":    PString("alfred.workflow.input.keyword"),
		"uid":     PString(selfUpdateKeywordObjectGUID),
		"version": PInteger(0),
	}
	objects = append(objects, selfUpdateKeywordObject)
	uidata[selfUpdateKeywordObjectGUID] = PDict{
		"ypos": PReal(10),
	}

	selfUpdateScriptObjectGUID := GUID()
	selfUpdateScriptObject := PDict{
		"config": PDict{
			"concurrently": PBool(false),
			"escaping":     PInteger(102),
			"script": PString(`set -e
trap 'echo Failed to update DAWG. Make sure there are no errors in config' ERR
chmod +x ./dawg
open -n -W "$(./dawg -config)"
./dawg -update 2>&1
echo 'Updated!'`),
			"type": PInteger(0),
		},
		"type":    PString("alfred.workflow.action.script"),
		"uid":     PString(selfUpdateScriptObjectGUID),
		"version": PInteger(0),
	}
	objects = append(objects, selfUpdateScriptObject)
	uidata[selfUpdateScriptObjectGUID] = PDict{
		"ypos": PReal(10),
	}
	connections[selfUpdateKeywordObjectGUID] = PArray{
		PDict{
			"destinationuid":  PString(selfUpdateScriptObjectGUID),
			"modifiers":       PInteger(0),
			"modifiersubtext": PString(""),
		},
	}

	selfUpdateNotificationObjectGUID := GUID()
	selfUpdateNotificationObject := PDict{
		"config": PDict{
			"lastpathcomponent":        PBool(false),
			"onlyshowifquerypopulated": PBool(true),
			"output":                   PInteger(0),
			"removeextension":          PBool(false),
			"sticky":                   PBool(false),
			"text":                     PString("{query}"),
			"title":                    PString("DAWG"),
		},
		"type":    PString("alfred.workflow.output.notification"),
		"uid":     PString(selfUpdateNotificationObjectGUID),
		"version": PInteger(0),
	}
	objects = append(objects, selfUpdateNotificationObject)
	uidata[selfUpdateNotificationObjectGUID] = PDict{
		"ypos": PReal(10),
	}
	connections[selfUpdateScriptObjectGUID] = PArray{
		PDict{
			"destinationuid":  PString(selfUpdateNotificationObjectGUID),
			"modifiers":       PInteger(0),
			"modifiersubtext": PString(""),
		},
	}

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
		"ypos": PReal(140),
	}

	ypos := 140
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
				"script":           PString(fmt.Sprintf("chmod +x ./dawg && ./dawg -service %s \"{query}\"", service)),
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
		"bundleid":    PString(BundleID),
		"category":    PString("Productivity"),
		"connections": connections,
		"createdby":   PString("Vlad Yarotsky"),
		"description": PString("Services at your fingertips :)"),
		"disabled":    PBool(false),
		"name":        PString("DAWG"),
		"objects":     objects,
		"readme":      PString(""),
		"uidata":      uidata,
		"webaddress":  PString("https://github.com/v-yarotsky/dawg"),
	}
	return plist
}
