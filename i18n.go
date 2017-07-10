package app

import (
	"github.com/gowww/i18n"
	"golang.org/x/text/language"
)

var confI18n struct {
	Locales  i18n.Locales
	Fallback language.Tag
	Parsers  []i18n.Parser
}

// Localize sets app locales with fallback and client locale parsers (order is mandatory and default are ParseFormValue, ParseAcceptLanguage).
func Localize(locs i18n.Locales, fallback language.Tag, parsers ...i18n.Parser) {
	if len(confI18n.Parsers) > 0 {
		panic("app: locales set multiple times")
	}
	if len(parsers) == 0 {
		parsers = []i18n.Parser{i18n.ParseFormValue, i18n.ParseAcceptLanguage}
	}
	confI18n.Fallback = fallback
	confI18n.Locales = locs
	confI18n.Parsers = parsers
}
