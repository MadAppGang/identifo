package localization

import (
	"fmt"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var DefaultLanguage = language.English

type Printer struct {
	printers   map[language.Tag]*message.Printer
	defaultTag language.Tag
}

func NewPrinter(defaultLocale string) (*Printer, error) {
	LoadDefaultCatalog()

	printers := make(map[language.Tag]*message.Printer)
	for _, t := range SupportedLangs {
		p := message.NewPrinter(t)
		printers[t] = p
	}

	l, err := language.Parse(defaultLocale)
	if err != nil {
		l = DefaultLanguage
	}

	// if we have no default printer, return error, and this error we could not localize
	_, ok := printers[l]
	if !ok {
		return nil, fmt.Errorf("No printer available for default locale: %v", l.String())
	}

	return &Printer{
		printers:   printers,
		defaultTag: l,
	}, nil
}

// String with language.
func (p *Printer) S(l language.Tag, s LocalizedString, params ...any) string {
	pp, ok := p.printers[l]
	if !ok {
		pp = p.printers[p.defaultTag]
	}

	return pp.Sprintf(string(s), params...)
}

// String default.
func (p *Printer) SD(s LocalizedString, params ...any) string {
	return p.S(p.defaultTag, s, params...)
}

// String with locale string in  BCP 47 format.
func (p *Printer) SL(localeOptions string, s LocalizedString, params ...any) string {
	matcher := language.NewMatcher(SupportedLangs)
	l, _ := language.MatchStrings(matcher, localeOptions)
	return p.S(l, s, params...)
}
