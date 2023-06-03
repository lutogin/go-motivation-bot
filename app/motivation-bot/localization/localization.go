package localization

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

// Declare bundle as a global variable
var bundle *i18n.Bundle

func init() {
	// Initialize bundle and load message files
	bundle = i18n.NewBundle(language.English)
	bundle.MustLoadMessageFile("./localization/translations/en.json")
	bundle.MustLoadMessageFile("./localization/translations/ru.json")
}

func Tr(messageID, lang string) string {
	parsedLang := getLangCode(lang)
	localizer := i18n.NewLocalizer(bundle, parsedLang)
	return localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: messageID,
		},
		MessageID: messageID,
	})
}

func getLangCode(lang string) string {
	switch lang {
	case "be":
	case "ru":
		return "ru"
	default:
		return "en"
	}

	return "en"
}
