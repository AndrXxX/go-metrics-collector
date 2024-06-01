package main

import (
	"github.com/AndrXxX/go-metrics-collector/internal/server"
)

func main() {
	if err := server.Run(); err != nil {
		panic(err)
	}
}
