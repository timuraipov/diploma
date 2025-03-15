package bootstrap

type Application struct {
	Cfg *Config
}

func App() (Application, error) {

	cfg, err := MustLoad()
	if err != nil {
		return Application{}, err
	}

	app := &Application{Cfg: cfg}
	return *app, nil
}
