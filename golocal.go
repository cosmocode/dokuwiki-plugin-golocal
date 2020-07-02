package main

import (
	"flag"
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/widget"
	setup "github.com/splitbrain/golocal/setup"
	"log"
	"os"
)

func main() {
	flagInstall := flag.Bool("install", false, "Install the protocol handler")
	flagUninstall := flag.Bool("uninstall", false, "Uninstall the protocol handler")
	flag.Parse()

	if *flagInstall {
		install(nil)
	} else if *flagUninstall {
		uninstall(nil)
	} else {
		gui()
	}
}

func gui() {
	application := app.New()

	label := "Hello"
	if len(os.Args) > 1 {
		label = os.Args[1]
	}

	w := application.NewWindow("Test")
	lblIntro := widget.NewLabel(label)
	btnInstall := widget.NewButton("Install", func() { install(w) })
	btnUninstall := widget.NewButton("Uninstall", func() { uninstall(w) })

	w.SetContent(
		widget.NewVBox(
			lblIntro,
			btnInstall,
			btnUninstall,
		),
	)

	w.Resize(fyne.NewSize(800, 400))
	w.ShowAndRun()
}

func install(window fyne.Window) {
	err := setup.LinuxInstall()
	errHandler(err, "Protocol handler installed", window)
}

func uninstall(window fyne.Window) {
	err := setup.LinuxUninstall()
	errHandler(err, "Protocol handler removed", window)
}

// Outputs either error or success message using the appropriate channel based on if
// window is available or nil
func errHandler(err error, success string, window fyne.Window) {
	if err == nil {
		if window == nil {
			log.Println(success)
		} else {
			dialog.ShowInformation("Success", success, window)
		}
	} else {
		if window == nil {
			log.Fatal(err)
		} else {
			dialog.ShowError(err, window)
		}
	}
}
