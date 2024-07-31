//go:build windows
// +build windows

package setup

import (
	"fmt"
	"golang.org/x/sys/windows/registry"
	"os"
	"os/exec"
	"strings"
)

// https://stackoverflow.com/a/3964401
func Install() error {
	var (
		err       error
		keyRoot   registry.Key
		keyOpener registry.Key
	)

	keyRoot, _, err = registry.CreateKey(registry.CLASSES_ROOT, PROTOCOL, registry.ALL_ACCESS)
	if err != nil {
		// if system wide install fails, use current user's registry
		keypath := fmt.Sprintf("Software\\Classes\\%s", PROTOCOL)
		keyRoot, _, err = registry.CreateKey(registry.CURRENT_USER, keypath, registry.ALL_ACCESS)
		if err != nil {
			return fmt.Errorf("failed to create root key. %s", err.Error())
		}
	}

	err = keyRoot.SetStringValue("", fmt.Sprintf("URL:%s Protocol", PROTOCOL))
	if err != nil {
		return err
	}

	err = keyRoot.SetStringValue("URL Protocol", "")
	if err != nil {
		return err
	}

	keyOpener, _, err = registry.CreateKey(keyRoot, "shell\\open\\command", registry.ALL_ACCESS)
	if err != nil {
		return err
	}

	err = keyOpener.SetStringValue("", fmt.Sprintf("\"%s\" \"%%1\"", os.Args[0]))
	if err != nil {
		return err
	}

	return nil
}

func Uninstall() error {
	var (
		err error
	)

	// try system wide key first
	_ = registry.DeleteKey(registry.CLASSES_ROOT, fmt.Sprintf("%s\\shell\\open\\command", PROTOCOL))
	_ = registry.DeleteKey(registry.CLASSES_ROOT, fmt.Sprintf("%s\\shell\\open", PROTOCOL))
	_ = registry.DeleteKey(registry.CLASSES_ROOT, fmt.Sprintf("%s\\shell", PROTOCOL))
	err = registry.DeleteKey(registry.CLASSES_ROOT, fmt.Sprintf("%s", PROTOCOL))
	if err == nil {
		return nil
	}

	// still here? try user key
	_ = registry.DeleteKey(registry.CURRENT_USER, fmt.Sprintf("Software\\Classes\\%s\\shell\\open\\command", PROTOCOL))
	_ = registry.DeleteKey(registry.CURRENT_USER, fmt.Sprintf("Software\\Classes\\%s\\shell\\open", PROTOCOL))
	_ = registry.DeleteKey(registry.CURRENT_USER, fmt.Sprintf("Software\\Classes\\%s\\shell", PROTOCOL))
	err = registry.DeleteKey(registry.CURRENT_USER, fmt.Sprintf("Software\\Classes\\%s", PROTOCOL))

	if err != nil {
		return fmt.Errorf("Handler removal failed.\n%s", err.Error())
	}

	return nil
}

func PreparePath(path string, isLetter bool) string {
	return strings.Replace(path, "/", "\\", -1)
}

func Run(path string) error {
	// note: the empty parameter is the title of the window and required even though it should be optional
	out, err := exec.Command("cmd", "/C", "start", "", "/B", path).CombinedOutput()

	if err != nil {
		return fmt.Errorf("Failed to execute command.\n%s\n%s", err.Error(), out)
	}

	return nil
}
