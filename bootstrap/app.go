package bootstrap

type Application struct {
	Cfg *Config
}

func App() (Application, error) {
	cfg := MustLoad()

	app := &Application{Cfg: cfg}
	return *app, nil
}
