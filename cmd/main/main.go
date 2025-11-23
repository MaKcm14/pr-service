package main

import "github.com/MaKcm14/pr-service/internal/app"

func main() {
	s := app.NewService()
	_ = s.Start()
}
