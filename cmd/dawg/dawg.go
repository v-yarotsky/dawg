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
		makeAlfredWorkflow(c)
	case *service != "":
		pattern := flag.Arg(0)
		printAlfredXML(c, *service, pattern)
	default:
		flag.Usage()
		os.Exit(1)
	}
}

func printAlfredXML(c dawg.Config, service, pattern string) {
	alfredOut, err := dawg.Filter(c, service, pattern)
	if err != nil {
		log.Fatal(err)
	}

	rawXML, err := alfredOut.MakeXML()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(rawXML))
}

func makeAlfredWorkflow(c dawg.Config) {
	plist := dawg.MakeWorkflowPList(c)
	out, _ := plist.PListWithHeader()

	zipfile, err := dawg.MakeWorkflowZIP(out)
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create("DAWG.alfredworkflow")
	if err != nil {
		log.Fatal(err)
	}

	_, err = io.Copy(f, zipfile)
	if err != nil {
		log.Fatal(err)
	}

	err = f.Close()
	if err != nil {
		log.Fatal(err)
	}
}
