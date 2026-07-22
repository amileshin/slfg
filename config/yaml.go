package config

import "gopkg.in/yaml.v3"

type yamlFile struct {
	Configuration yamlConfig `yaml:"configuration"`
}

type yamlConfig struct {
	Appenders []yamlAppender `yaml:"appender"`
	Root      yamlRoot       `yaml:"root"`
}

type yamlAppender struct {
	Name    string       `yaml:"name"`
	Class   string       `yaml:"class"`
	Target  string       `yaml:"target"`
	File    string       `yaml:"file"`
	Append  *bool        `yaml:"append"`
	Filters []yamlFilter `yaml:"filter"`
}

type yamlFilter struct {
	Class string `yaml:"class"`
	Level string `yaml:"level"`
}

type yamlRoot struct {
	Level        string           `yaml:"level"`
	AppenderRefs []yamlAppenderRef `yaml:"appender-ref"`
}

type yamlAppenderRef struct {
	Ref string `yaml:"ref"`
}

func parseYAML(data []byte) (*Configuration, error) {
	var raw yamlFile
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	cfg := &Configuration{
		Root: RootConfig{
			Level: raw.Configuration.Root.Level,
		},
	}

	for _, ra := range raw.Configuration.Root.AppenderRefs {
		cfg.Root.AppenderRefs = append(cfg.Root.AppenderRefs, ra.Ref)
	}

	for _, ya := range raw.Configuration.Appenders {
		ac := AppenderConfig{
			Name:   ya.Name,
			Class:  ya.Class,
			Target: ya.Target,
			File:   ya.File,
			Append: ya.Append == nil || *ya.Append,
		}
		for _, yf := range ya.Filters {
			ac.Filters = append(ac.Filters, FilterConfig{
				Class: yf.Class,
				Level: yf.Level,
			})
		}
		cfg.Appenders = append(cfg.Appenders, ac)
	}

	return cfg, nil
}
