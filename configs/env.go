// mongodb+srv://jongjate:<password>@cluster0.cgakqxe.mongodb.net/?retryWrites=true&w=majority

package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func EnvMySql() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Read the environment variables
	databaseURL := os.Getenv("DATABASE_URL")
	databaseName := os.Getenv("DATABASE_NAME")
	databaseUsername := os.Getenv("DATABASE_USERNAME")
	databasePassword := os.Getenv("DATABASE_PASSWORD")
	databasePort := os.Getenv("DATABASE_PORT")

	// You can use these variables as needed
	// For example, you might want to create a connection string
	// root:root@tcp(127.0.0.1:3306)/api_isoftel
	connectionString := databaseUsername + ":" + databasePassword + "@tcp(" + databaseURL + ":" + databasePort + ")/" + databaseName

	//print value
	// fmt.Println(connectionString)

	return connectionString
}
