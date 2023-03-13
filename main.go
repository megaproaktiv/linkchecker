package main

import (
	"archive/zip"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/antchfx/xmlquery"
	"github.com/fatih/color"
)

var CLI struct {
	Check struct {
	  Internal bool `help:"Deactivate internet availability check for internal checks. True = no check" `

	  Paths []string `arg:"" name:"path" help:"Paths to powerpoint pptx file." type:"path"`
	} `cmd:"" help:"check links in an pptx file for reachability."`
}  

func main() {
	var presentationPath string
    // Open the PowerPoint file and parse the XML data
	ctx := kong.Parse(&CLI, 
		kong.Name("linkchecker"),
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
    internal := CLI.Check.Internal
	// Find the presentation.xml file inside the PowerPoint file
	var presentation *zip.File
	for _, file := range pptx.File {
		// if file.Name == "ppt/presentation.xml" {
		// if strings.HasPrefix(file.Name, == "ppt/slides/_rels/slide2.xml.rels" {
		if strings.HasPrefix(file.Name, "ppt/slides/_rels/slide") {
			parts := strings.Split(file.Name, "/")
			slideParts := strings.Split(parts[3], ".")
			slide := slideParts[0]
			presentation = file
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
					// fmt.Printf("Found https target: %s in slide %v\n", target,slide)
					if webIsReachable(target,internal) {
						color.Green("☑️ Target %v is reachable in slide: %v \n", target, slide)
					} else {
						color.Red("❌ Target %v is not reachable in slide: %v \n", target, slide)
					}
				}
			}
		}
	}
	if presentation == nil {
		panic("No slides not found in PowerPoint file")
	}

}


func startsWithHTTPS(s string) bool {
	return len(s) >= 5 && s[:5] == "https"
}

func webIsReachable(web string, internal bool) bool {
    response, errors := http.Get(web)

    if errors != nil {
		if !internal {
			_, netErrors := http.Get("https://www.google.com")
			
			if netErrors != nil {
				fmt.Fprintf(os.Stderr, "no internet\n")
				os.Exit(1)
			}
			
			return false
		}
		return false
    }

    if response.StatusCode == 200 {
        return true
    }

    return false
}
