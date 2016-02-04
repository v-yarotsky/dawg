package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/v-yarotsky/dawg"
)

func main() {
	service := flag.String("s", "", "Service name")
	makeWorkflow := flag.Bool("generate", false, "Generate Alfred Workflow")
	flag.Parse()

	c := dawg.MustReadConfig()

	switch true {
	case *makeWorkflow:
		handleError(makeAlfredWorkflow(c))
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

func makeAlfredWorkflow(c dawg.Config) error {
	plist := dawg.MakeWorkflowPList(c)
	out, _ := plist.PListWithHeader()

	zipfile, err := dawg.MakeWorkflowZIP(out)
	if err != nil {
		return fmt.Errorf("could not make a workflow zip file: %v", err)
	}

	f, err := os.Create("DAWG.alfredworkflow")
	if err != nil {
		return fmt.Errorf("could not open workflow zip file for writing: %v", err)
	}

	_, err = io.Copy(f, zipfile)
	if err != nil {
		return fmt.Errorf("could not write to the workflow zip file: %v", err)
	}

	err = f.Close()
	if err != nil {
		return fmt.Errorf("could not close the workflow zip file: %v", err)
	}
	return nil
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
