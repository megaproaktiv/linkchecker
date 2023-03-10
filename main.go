package main

import (
	"archive/zip"
	"fmt"
	"net/http"
	"os"

	"github.com/antchfx/xmlquery"
    "github.com/fatih/color"
	"github.com/alecthomas/kong"

)

var CLI struct {
	Check struct {
	  Paths []string `arg:"" name:"path" help:"Paths to powerpoint pptx file." type:"path"`
	} `cmd:"" help:"check links in an pptx file for reachability."`
}  

//const fileName = "testdata/someuris.pptx"

func main() {
	var presentationPath string
    // Open the PowerPoint file and parse the XML data
	ctx := kong.Parse(&CLI, 
		kong.Name("linkechecker"),
		kong.Description("Check powerpoint links reachability."),
		)
	switch ctx.Command() {
	case "check <path>":		
		presentationPath = CLI.Check.Paths[0]
	default:
	  panic(ctx.Command())
	}

    pptx, err := zip.OpenReader(presentationPath)
	if err != nil {
		panic(err)
	}
	defer pptx.Close()

	// Find the presentation.xml file inside the PowerPoint file
	var presentation *zip.File
	for _, file := range pptx.File {
		// if file.Name == "ppt/presentation.xml" {
		if file.Name == "ppt/slides/_rels/slide2.xml.rels" {
			presentation = file
			break
		}
	}
	if presentation == nil {
		panic("presentation.xml not found in PowerPoint file")
	}

	// Read the XML data from the presentation.xml file
	xmlData, err := presentation.Open()
	if err != nil {
		panic(err)
	}
	defer xmlData.Close()

    doc, err := xmlquery.Parse(xmlData)
	if err != nil {
		panic(err)
	}
    

	// Find all Relationship elements using XPath
    elements := xmlquery.Find(doc, "//Relationship")

    // Filter the Relationship elements to only include those with https targets
    for _, element := range elements {
        target := element.SelectAttr("Target")
        if target != "" && startsWithHTTPS(target) {
            fmt.Printf("Found https target: %s\n", target)
            if webIsReachable(target) {
                color.Green("☑️ Target is reachable")
            } else {
                color.Red("❌ Target is not reachable")
            }
        }
    }

}


func startsWithHTTPS(s string) bool {
	return len(s) >= 5 && s[:5] == "https"
}

func webIsReachable(web string) bool {
    response, errors := http.Get(web)

    if errors != nil {
        _, netErrors := http.Get("https://www.google.com")

        if netErrors != nil {
            fmt.Fprintf(os.Stderr, "no internet\n")
            os.Exit(1)
        }

        return false
    }

    if response.StatusCode == 200 {
        return true
    }

    return false
}
