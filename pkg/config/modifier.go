package config

type Modifier struct {
	Type string         `json:"type" description:"name of the modifier"`
	Args map[string]any `json:"args" description:"modifier configuration"`
}

type LinePatchConfig struct {
	File            string  `mapstructure:"file" description:"the name of the file to be patched"`
	Line            int     `mapstructure:"line" description:"the line number in the file to be patched"`
	ReplaceTemplate *string `mapstructure:"template" description:"a special template to be used for patching the line"`
}
