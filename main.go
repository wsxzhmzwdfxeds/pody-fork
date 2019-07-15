package main

import (
	"fmt"
	"os"
	"github.com/JulienBreux/pody/pkg"
	_ "github.com/JulienBreux/pody/pkg"
	"github.com/jroimartin/gocui"
	"log"

	//v1 "k8s.io/api/core/v1"
)

// Main or not main, that's the question^^
func main() {
	c := pkg.GetConfig()

	// Ask version: pody --version
	if c.AskVersion {
		fmt.Println(pkg.VersionFull())
		os.Exit(0)
	}

	// Ask Help: pody --help
	if c.AskHelp {
		fmt.Println(pkg.VersionFull())
		fmt.Println(pkg.HELP)
		os.Exit(0)
	}

	// Only used to check errors
	pkg.GetClientSet()

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Highlight = true
	g.SelFgColor = gocui.ColorGreen

	g.SetManagerFunc(uiLayout)

	if err := pkg.UiKey(g); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

// Define the UI layout
func uiLayout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	return pkg.ViewsTotal(g, maxX, maxY)
}
