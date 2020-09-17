package main

import (
	"github.com/ronyv89/leedprojects/internal/routes"
)

func main() {
	r := routes.LPRouter()
	r.Run(":9000") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
