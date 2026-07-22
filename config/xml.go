package config

import "encoding/xml"

type xmlConfiguration struct {
	XMLName   xml.Name      `xml:"configuration"`
	Appenders []xmlAppender `xml:"appender"`
	Root      xmlRoot       `xml:"root"`
}

type xmlAppender struct {
	Name    string      `xml:"name,attr"`
	Class   string      `xml:"class,attr"`
	Target  string      `xml:"target"`
	File    string      `xml:"file"`
	Append  string      `xml:"append"`
	Filters []xmlFilter `xml:"filter"`
}

type xmlFilter struct {
	Class string `xml:"class,attr"`
	Level string `xml:"level"`
}

type xmlRoot struct {
	Level        string           `xml:"level,attr"`
	AppenderRefs []xmlAppenderRef `xml:"appender-ref"`
}

type xmlAppenderRef struct {
	Ref string `xml:"ref,attr"`
}

func parseXML(data []byte) (*Configuration, error) {
	var raw xmlConfiguration
	if err := xml.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	cfg := &Configuration{
		Root: RootConfig{
			Level: raw.Root.Level,
		},
	}

	for _, ra := range raw.Root.AppenderRefs {
		cfg.Root.AppenderRefs = append(cfg.Root.AppenderRefs, ra.Ref)
	}

	for _, xa := range raw.Appenders {
		ac := AppenderConfig{
			Name:   xa.Name,
			Class:  xa.Class,
			Target: xa.Target,
			File:   xa.File,
			Append: xa.Append != "false",
		}
		for _, xf := range xa.Filters {
			ac.Filters = append(ac.Filters, FilterConfig{
				Class: xf.Class,
				Level: xf.Level,
			})
		}
		cfg.Appenders = append(cfg.Appenders, ac)
	}

	return cfg, nil
}
