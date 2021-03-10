package config

// App ...
type App struct {
	ApiKey       string `yaml:"api_key"`
	TemplatePath string `yaml:"template_path"`
}

// New ...
func New() App {
	return App{
		TemplatePath: "./",
	}
}
