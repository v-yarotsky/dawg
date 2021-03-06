package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/v-yarotsky/dawg"
)

const configPath = "./dawg.json"

func main() {
	service := flag.String("service", "", "Service name")
	updateWorkflow := flag.Bool("update", false, "Update the Alfred Workflow")
	flag.Parse()

	c, err := dawg.ReadConfig(configPath)
	handleError(err)

	switch true {
	case *updateWorkflow:
		handleError(updateAlfredWorkflow(c))
	case *service != "":
		pattern := flag.Arg(0)
		handleError(printAlfredXML(c, *service, pattern))
	default:
		flag.Usage()
		os.Exit(1)
	}
}

func printAlfredXML(c dawg.Config, service, pattern string) error {
	alfredOut, err := dawg.Filter(c, service, pattern)
	if err != nil {
		return fmt.Errorf("could not filter through config: %v", err)
	}

	rawXML, err := alfredOut.MakeXML()
	if err != nil {
		return fmt.Errorf("could not prepare xml output for Alfred: %v", err)
	}
	fmt.Println(string(rawXML))
	return nil
}

func updateAlfredWorkflow(c dawg.Config) error {
	plist := dawg.MakeWorkflowPList(c)
	icons, err := dawg.DownloadServiceIcons(c)
	if err != nil {
		return fmt.Errorf("could not download service icons: %v", err)
	}
	plistData, _ := plist.PListWithHeader()

	if err = dawg.SelfUpdate(plistData, icons); err != nil {
		return fmt.Errorf("could not self update: %v", err)
	}
	return nil
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
