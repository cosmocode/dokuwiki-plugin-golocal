//go:build darwin

package setup

import (
	"fmt"
	homedir "github.com/mitchellh/go-homedir"
	"howett.net/plist"
	"os"
	"os/exec"
)

const launchServicesPlist = "~/Library/Preferences/com.apple.LaunchServices/com.apple.launchservices.secure.plist"

// see https://unix.stackexchange.com/questions/497146/create-a-custom-url-protocol-handler
func Install() error {
	self, _ := os.Executable()
	removeHandler(PROTOCOL)
	err := addHandler(PROTOCOL, self)
	if err != nil {
		return err
	}
	return nil
}

func Uninstall() error {
	_, err := getHandler(PROTOCOL)
	if err != nil {
		return err
	}

	err = removeHandler(PROTOCOL)
	if err != nil {
		return err
	}
	return nil
}

func Run(path string) error {
	out, err := exec.Command("open", "smb:"+path).CombinedOutput()
	if err != nil {
		return fmt.Errorf("Failed to execute open command.\n%s\n%s", err.Error(), out)
	}

	return nil
}

func readPlist() (map[string]interface{}, error) {
	plistPath, err := homedir.Expand(launchServicesPlist)
	if err != nil {
		return nil, err
	}

	content, err := os.ReadFile(plistPath)
	if err != nil {
		return nil, err
	}

	var prefs interface{}
	if _, err = plist.Unmarshal(content, &prefs); err != nil {
		return nil, err
	}

	prefsMap, ok := prefs.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unable to convert plist to map")
	}

	return prefsMap, nil
}

func writePlist(prefsMap map[string]interface{}) error {
	plistPath, err := homedir.Expand(launchServicesPlist)
	if err != nil {
		return err
	}

	updatedContent, err := plist.Marshal(prefsMap, plist.AutomaticFormat)
	if err != nil {
		return err
	}

	err = os.WriteFile(plistPath, updatedContent, 0644)
	if err != nil {
		return err
	}

	return nil
}

func addHandler(protocol, path string) error {
	prefsMap, err := readPlist()
	if err != nil {
		return err
	}

	// Get the LSHandlers slice
	lsHandlers, ok := prefsMap["LSHandlers"].([]interface{})
	if !ok {
		return fmt.Errorf("unable to get LSHandlers")
	}

	// Append a new handler to the LSHandlers slice
	lsHandlers = append(lsHandlers, map[string]string{
		"LSHandlerRoleAll":   path,
		"LSHandlerURLScheme": protocol,
	})

	// Update the LSHandlers slice in the prefsMap
	prefsMap["LSHandlers"] = lsHandlers

	err = writePlist(prefsMap)
	if err != nil {
		return err
	}

	return nil
}

func removeHandler(protocol string) error {
	prefsMap, err := readPlist()
	if err != nil {
		return err
	}

	// Get the LSHandlers slice
	lsHandlers, ok := prefsMap["LSHandlers"].([]interface{})
	if !ok {
		return fmt.Errorf("unable to get LSHandlers")
	}

	// Create a new slice to hold the handlers that do not match the given protocol
	newLSHandlers := make([]interface{}, 0)

	// Iterate over the LSHandlers slice and add the handler to the new slice if it does not match the given protocol
	for _, handler := range lsHandlers {
		handlerMap, ok := handler.(map[string]interface{})
		if !ok {
			return fmt.Errorf("unable to convert handler to map")
		}

		if handlerMap["LSHandlerURLScheme"] != protocol {
			newLSHandlers = append(newLSHandlers, handler)
		}
	}

	// Update the LSHandlers slice in the prefsMap with the new slice
	prefsMap["LSHandlers"] = newLSHandlers

	err = writePlist(prefsMap)
	if err != nil {
		return err
	}

	return nil
}

func getHandler(protocol string) (string, error) {
	prefsMap, err := readPlist()
	if err != nil {
		return "", err
	}

	// Get the LSHandlers slice
	lsHandlers, ok := prefsMap["LSHandlers"].([]interface{})
	if !ok {
		return "", fmt.Errorf("unable to get LSHandlers")
	}

	// Iterate over the LSHandlers slice and find the handler that matches the given protocol
	for _, handler := range lsHandlers {
		handlerMap, ok := handler.(map[string]interface{})
		if !ok {
			return "", fmt.Errorf("unable to convert handler to map")
		}

		if handlerMap["LSHandlerURLScheme"] == protocol {
			return handlerMap["LSHandlerRoleAll"].(string), nil
		}
	}

	return "", fmt.Errorf("handler for protocol %s not found", protocol)
}
