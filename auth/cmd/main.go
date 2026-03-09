package main

import (
	"fmt"

	"github.com/IvanDrf/work-hunter/auth/internal/config"
)

func main() {
	config := config.LoadFromYAML()

	fmt.Println(config)
}
