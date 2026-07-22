package config

// Configuration — модель конфигурации logback.
type Configuration struct {
	Appenders []AppenderConfig
	Root      RootConfig
}

type AppenderConfig struct {
	Name    string
	Class   string
	Target  string // ConsoleAppender: "System.out" / "System.err"
	File    string // FileAppender: путь к файлу
	Append  bool   // FileAppender: дозапись (по умолчанию true)
	Filters []FilterConfig
}

type FilterConfig struct {
	Class string
	Level string
}

type RootConfig struct {
	Level        string
	AppenderRefs []string
}
