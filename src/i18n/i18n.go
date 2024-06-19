package i18n

import (
	"embed"
	"github.com/jeandeaual/go-locale"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v2"
	"log"
	"strings"
)

//go:embed *.yaml
var localeFiles embed.FS

var localizer *i18n.Localizer

// Initialize initializes the i18n bundle and localizer
func Initialize() {
	lang := getPreferredLanguage()

	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)

	// Load English translations first
	if err := loadTranslations(bundle, "en.yaml"); err != nil {
		log.Fatalf("Failed to load default English translation: %s", err)
	}

	// Try to load the specified language translations
	if lang != "en" {
		if err := loadTranslations(bundle, lang+".yaml"); err != nil {
			// Try loading the base language if the specific locale fails
			baseLang := strings.Split(lang, "-")[0]
			if baseLang != lang { // Only attempt if baseLang is different from lang
				log.Printf("Failed to load %s translation, trying base language %s: %s", lang, baseLang, err)
				if err := loadTranslations(bundle, baseLang+".yaml"); err != nil {
					log.Printf("Failed to load base language %s translation, falling back to English: %s", baseLang, err)
				} else {
					lang = baseLang
				}
			} else {
				log.Printf("Failed to load %s translation, falling back to English: %s", lang, err)
			}
		}
	}

	// Create a localizer for the specified language
	localizer = i18n.NewLocalizer(bundle, lang)
}

// loadTranslations loads translation files into the bundle
func loadTranslations(bundle *i18n.Bundle, filename string) error {
	data, err := localeFiles.ReadFile(filename)
	if err != nil {
		return err
	}

	// Must use bundle.ParseMessageFileBytes instead of LoadMessageFile for embedded data
	_, err = bundle.ParseMessageFileBytes(data, filename)
	if err != nil {
		return err
	}
	return nil
}

// T translates a message ID to the current language
func T(messageID string, templateData map[string]string) string {
	result, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: templateData,
	})
	if err != nil {
		log.Printf("Failed to localize message: %s", err)
		return messageID + "*"
	}
	return result
}

// getPreferredLanguage determines the user's preferred language using go-locale
func getPreferredLanguage() string {
	lang, err := locale.GetLocale()
	if err != nil {
		log.Printf("Failed to get preferred language, defaulting to English: %s", err)
		return "en"
	}

	// Normalize the language code
	lang = strings.Split(lang, ".")[0] // Remove encoding part
	lang = strings.ReplaceAll(lang, "_", "-")

	return lang
}
