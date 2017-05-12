package finesmith

// AssetConfig allows to copy single sources into multiple destinations in parallel
type AssetConfig struct {
	SourceDir       string   `json:"source"`
	DestinationDirs []string `json:"destinationDirs"`
}

// Config contains all required and optional flags to run finesmith
type Config struct {
	SourceDir       string      `json:"sourceBaseDir"`
	DestinationDir  string      `json:"destinationBaseDir"`
	TemplateDir     string      `json:"templateDir"`
	AssetsDir       AssetConfig `json:"assets,omitempty"`
	PrismicURL      string      `json:"prismicUrl"`
	PrismicToken    string      `json:"prismicToken"`
	EnablePreview   bool        `json:"enablePreview,omitempty"`
	EnableBuildHook bool        `json:"enableBuildHook,omitempty"`
	ServerPort      int         `json:"port"`
}
