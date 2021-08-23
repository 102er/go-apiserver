package configer

type Config struct {
	Version string
	Env     string
	Log     LogCfg
}

type LogCfg struct {
	Level        string
	Filename     string
	MaxSize      int
	MaxAge       int
	MaxBackups   int
	PrintConsole bool
}
