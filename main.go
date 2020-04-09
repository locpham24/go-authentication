package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/joho/godotenv"
	"github.com/locpham24/go-authentication/cmd"
	"github.com/micro/cli/v2"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	app := &cli.App{
		Name:  "auth",
		Usage: "Go Authentication Service",
		Action: func(c *cli.Context) error {
			fmt.Println("Welcome to Go Authentication Service")
			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "database",
				Aliases: []string{"db", "d"},
				EnvVars: []string{"DB_URI"},
			},
		},
		Before: func(c *cli.Context) error {
			db, err := gorm.Open("mysql", c.String("database"))
			if err != nil {
				panic(err)
			}
			c.App.Metadata["db"] = db
			return nil
		},
		Commands: []*cli.Command{&cmd.Migrate, &cmd.Start},
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
