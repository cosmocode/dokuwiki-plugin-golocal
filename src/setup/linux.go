//go:build linux
// +build linux

package setup

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"strings"
)

const DESKTOPFILE = `
[Desktop Entry]
Type=Application
Name=Local Link Scheme Handler
Exec=%s %%u
StartupNotify=false
MimeType=%s;
`

// see https://unix.stackexchange.com/questions/497146/create-a-custom-url-protocol-handler
func Install() error {
	self, _ := os.Executable()
	desktopFile := desktopFile()
	schemeMimeType := fmt.Sprintf("x-scheme-handler/%s", PROTOCOL)
	desktopEntry := strings.TrimLeft(fmt.Sprintf(DESKTOPFILE, self, schemeMimeType), "\n")
	desktopFilePath, err1 := desktopFilePath()
	if err1 != nil {
		return err1
	}

	log.Print(desktopEntry)
	err2 := ioutil.WriteFile(desktopFilePath, []byte(desktopEntry), 0644)
	if err2 != nil {
		return err2
	}

	out, err3 := exec.Command("xdg-mime", "default", desktopFile, schemeMimeType).CombinedOutput()
	if err3 != nil {

		return fmt.Errorf("Failed to execute xdg-mime command.\n%s\n%s", err3.Error(), out)
	}

	return nil
}

func Uninstall() error {
	desktopFilePath, err0 := desktopFilePath()
	if err0 != nil {
		return err0
	}

	_, err1 := os.Stat(desktopFilePath)
	if os.IsNotExist(err1) {
		return errors.New("No handler found.")
	}

	err2 := os.Remove(desktopFilePath)
	if err2 != nil {
		return fmt.Errorf("Failed to remove desktop file.\n%s\n%s", desktopFilePath, err2.Error())
	}
	return nil
}

func Run(path string) error {
	out, err := exec.Command("xdg-open", "smb:"+path).CombinedOutput()
	if err != nil {
		return fmt.Errorf("Failed to execute xdg-open command.\n%s\n%s", err.Error(), out)
	}

	return nil
}

func desktopFile() string {
	return fmt.Sprintf("%s-handler.desktop", PROTOCOL)
}

func desktopFilePath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	location := fmt.Sprintf("%s/.local/share/applications/%s", usr.HomeDir, desktopFile())
	return location, nil
}
