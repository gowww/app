/*
Package app provides a full featured framework for any web apps.
*/
package app

import (
	"log"
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
	ParseFormValue Parser = i18n.ParseFormValue
	// ParseAcceptLanguage parses the Accept-Language header.
	ParseAcceptLanguage Parser = i18n.ParseAcceptLanguage
	// ParseCookie parses the LocaleFieldName cookie.
	ParseCookie Parser = i18n.ParseCookie
)

// Locales is a map of locales and their translations.
type Locales map[language.Tag]Translations

// Translations is a map of translations associated to keys.
type Translations map[string]string

// A Parser is a funcion that returns a list of accepted languages, most preferred first.
type Parser i18n.Parser

// Localize sets app locales with fallback and client locale parsers (order is mandatory and default are ParseFormValue, ParseAcceptLanguage).
func Localize(locs Locales, fallback language.Tag, parsers ...Parser) {
	if confI18n != nil {
		log.Fatal("app: locales can be set only once")
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

type configurationI18n struct {
	Locales  Locales
	Fallback language.Tag
	Parsers  []Parser
}

func (conf *configurationI18n) handleI18n(handler *http.Handler) {
	if confI18n == nil {
		return
	}
	ll := make(i18n.Locales)
	for lang, trans := range conf.Locales {
		ll[lang] = i18n.Translations(trans)
	}
	var pp []i18n.Parser
	for _, parser := range conf.Parsers {
		pp = append(pp, i18n.Parser(parser))
	}
	*handler = i18n.Handle(*handler, ll, conf.Fallback, pp...)
}
