package app

type App struct {
	r Repo
}

func New(r Repo) *App {
	return &App{r: r}
}
