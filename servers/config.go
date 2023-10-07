package servers

type Config struct {
	Host    string  `json:"host,omitempty"`
	Port    int     `json:"port,omitempty"`
	Context string  `json:"context,omitempty"`
	DistDir string  `json:"distDir,omitempty"`
	Options Options `json:"options"`
}

type Options struct {
	SuccessCode string `json:"successCode,omitempty"`
	ErrorCode   string `json:"errorCode,omitempty"`
}
