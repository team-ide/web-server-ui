package servers

type Config struct {
	Host    string  `json:"host,omitempty" yaml:"host,omitempty"`
	Port    int     `json:"port,omitempty" yaml:"port,omitempty"`
	Context string  `json:"context,omitempty" yaml:"context,omitempty"`
	DistDir string  `json:"distDir,omitempty" yaml:"distDir,omitempty"`
	Options Options `json:"options,omitempty" yaml:"options,omitempty"`
}

type Options struct {
	SuccessCode string `json:"successCode,omitempty" yaml:"successCode,omitempty"`
	ErrorCode   string `json:"errorCode,omitempty" yaml:"errorCode,omitempty"`
}
