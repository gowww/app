package app

import (
	"net/http"

	"github.com/gowww/i18n"
	"golang.org/x/text/language"
)

const (
	// TnPlaceholder is the placeholder replaced by n in a translation, when using the Tn function.
	TnPlaceholder = "{{.n}}"
)

var (
	confI18n *configurationI18n

	// ParseFormValue parses the LocaleFieldName form value.
	ParseFormValue = i18n.ParseFormValue
	// ParseAcceptLanguage parses the Accept-Language header.
	ParseAcceptLanguage = i18n.ParseAcceptLanguage
	// ParseCookie parses the LocaleFieldName cookie.
	ParseCookie = i18n.ParseCookie
)

type configurationI18n struct {
	Locales  Locales
	Fallback language.Tag
	Parsers  []Parser
}

// Locales is a map of locales and their translations.
type Locales map[language.Tag]Translations

// Translations is a map of translations associated to keys.
type Translations map[string]string

// A Parser is a funcion that returns a list of accepted languages, most preferred first.
type Parser func(*http.Request) []language.Tag

// Localize sets app locales with fallback and client locale parsers (order is mandatory and default are ParseFormValue, ParseAcceptLanguage).
func Localize(locs Locales, fallback language.Tag, parsers ...Parser) {
	if confI18n != nil {
		panic("app: locales set multiple times")
	}
	if len(parsers) == 0 {
		parsers = []Parser{ParseFormValue, ParseAcceptLanguage}
	}
	confI18n = &configurationI18n{
		Fallback: fallback,
		Locales:  locs,
		Parsers:  parsers,
	}
}
