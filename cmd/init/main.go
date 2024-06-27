package main

import (
	"log"
	"os/exec"

	"github.com/cgalvisleon/elvis/console"
)

func main() {
	_, err := Command([]string{
		"go get github.com/joho/godotenv/autoload",
		"go get github.com/redis/go-redis/v9",
		"go get github.com/google/uuid",
		"go get github.com/nats-io/nats.go",
		"go get golang.org/x/crypto/bcrypt",
		"go get golang.org/x/exp/slices",
		"go get github.com/manifoldco/promptui",
		"go get github.com/schollz/progressbar/v3",
		"go get github.com/spf13/cobra",
		"go get github.com/golang-jwt/jwt/v4",
		"go get github.com/go-chi/chi/v5",
		"go get github.com/shirou/gopsutil/v3/mem",
		"go get github.com/lib/pq",
		"go get github.com/dimiro1/banner",
		"go get github.com/mattn/go-colorable",
		"go get github.com/rs/cors",
		"go get github.com/cgalvisleon/elvis",
	})

	if err != nil {
		console.Error(err)
	}
}

func Command(coms []string) ([][]byte, error) {
	var result [][]byte
	for _, com := range coms {
		out, err := exec.Command(com).Output()
		if err != nil {
			log.Fatal(err)
		}
		result = append(result, out)
	}

	return result, nil
}
