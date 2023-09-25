package main

import (
	"embed"
	"log"

	"github.com/ImDevinC/go-pd3/internal/app"
	"github.com/ImDevinC/go-pd3/internal/config"
	_ "github.com/joho/godotenv/autoload"
)

//go:embed default.json
var defaultChallenges embed.FS

func main() {
	challenges, err := config.LoadSaved()
	if err != nil {
		panic(err)
	}
	app := app.App{
		Challenges: challenges,
	}
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
