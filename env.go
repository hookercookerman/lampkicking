package lampkicking

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var environment string

func Getenv(name string) string {
	return os.Getenv(name)
}

func init() {
	environment = os.Getenv("GO_ENV")
	if environment == "" {
		environment = "development"
	}
	LoadEnvFiles()
}

func GetEnv() string {
	return environment
}

// was fighting a losing battle
func LoadEnvFiles() {
	file := fmt.Sprintf("./../../.%v_env", GetEnv())
	err := godotenv.Load(file)
	if err != nil {
		file := fmt.Sprintf(".%v_env", GetEnv())
		err := godotenv.Load(file)
		if err != nil {
			log.Fatalf("Error loading file: %v", file)
		}
	}
}
