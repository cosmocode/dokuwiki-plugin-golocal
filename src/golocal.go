package main

import (
	"flag"
	"fmt"
	fyne "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/cosmocode/golocal/i18n"
	setup "github.com/cosmocode/golocal/setup"
	"log"
	"os"
	"regexp"
)

func main() {
	flag.Usage = usage
	flagInstall := flag.Bool("install", false, "Install the protocol handler")
	flagUninstall := flag.Bool("uninstall", false, "Uninstall the protocol handler")
	flag.Parse()

	i18n.Initialize()

	if *flagInstall {
		install(nil)
	} else if *flagUninstall {
		uninstall(nil)
	} else {
		_, window := guiInit()
		if len(os.Args) > 1 {
			go run(os.Args[1], window)
		} else {
			go guiInstaller(window)
		}

		// start the main loop
		window.ShowAndRun()
	}
}

func guiInit() (fyne.App, fyne.Window) {
	application := app.New()
	w := application.NewWindow(fmt.Sprintf("%s handler", setup.PROTOCOL))
	w.Resize(fyne.NewSize(800, 400))
	return application, w
}

func guiInstaller(window fyne.Window) {
	lblIntro := widget.NewLabel(i18n.T("intro", nil))
	lblIntro.Wrapping = fyne.TextWrapWord
	btnInstall := widget.NewButton(i18n.T("install", nil), func() { install(window) })
	btnUninstall := widget.NewButton(i18n.T("uninstall", nil), func() { uninstall(window) })

	window.SetContent(
		container.New(
			layout.NewVBoxLayout(),
			lblIntro,
			layout.NewSpacer(),
			btnInstall,
			btnUninstall,
		),
	)
}

func run(path string, window fyne.Window) {
	// remove protocol
	r, _ := regexp.Compile("^.*?://")
	path = r.ReplaceAllString(path, "")

	// FIXME decode URI, parse it maybe?

	window.SetContent(widget.NewLabel(path))

	err := setup.Run(path)
	errHandler(err, "", window)
	if err == nil {
		window.Close()
	}
}

func install(window fyne.Window) {
	err := setup.Install()
	errHandler(err, i18n.T("installed", nil), window)
}

func uninstall(window fyne.Window) {
	err := setup.Uninstall()
	errHandler(err, i18n.T("uninstalled", nil), window)
}

// Outputs either error or success message using the appropriate channel based on if
// window is available or nil
func errHandler(err error, success string, window fyne.Window) {
	if err == nil {
		if window == nil {
			log.Println(success)
		} else if success != "" {
			dialog.ShowInformation(i18n.T("success", nil), success, window)
		}
	} else {
		if window == nil {
			log.Fatal(err)
		} else {
			dialog.ShowError(err, window)
		}
	}
}

func usage() {
	fmt.Printf("Usage: %s %s://path \n", os.Args[0], setup.PROTOCOL)
	fmt.Println("  Protocol handling. Will try to open the given path locally.")
	fmt.Println()

	fmt.Printf("Usage: %s [OPTION]\n", os.Args[0])
	fmt.Println("  Install or uninstall the protocol handler")
	flag.PrintDefaults()
}
