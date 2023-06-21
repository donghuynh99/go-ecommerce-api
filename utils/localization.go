package utils

import (
	"encoding/json"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var localizer *i18n.Localizer

func InitLocalizer(langs ...string) {
	bundle := i18n.NewBundle(language.English)

	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	bundle.MustLoadMessageFile("./langs/en.json")

	bundle.MustLoadMessageFile("./langs/vi.json")

	localizer = i18n.NewLocalizer(bundle, langs...)
}

func Translation(messageId string, data map[string]interface{}, counter interface{}) string {
	if data == nil {
		data = make(map[string]interface{})
	}

	var pluralCount int

	if counter == nil {
		pluralCount = 0
	} else {
		pluralCount = counter.(int)
	}

	translation := localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID:    messageId,
		TemplateData: data,
		PluralCount:  pluralCount,
	})

	return translation
}
