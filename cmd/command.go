package cmd

import (
	"log"
	"os"

	"github.com/donghuynh99/ecommerce_api/database/migrations"
	"github.com/donghuynh99/ecommerce_api/database/seeders"
	"github.com/urfave/cli"
)

func InitCommand() {
	cmdApp := cli.NewApp()

	cmdApp.Commands = []cli.Command{
		{
			Name: "db:migrate",
			Action: func(c *cli.Context) error {
				migrations.RunMigration()

				return nil
			},
		},
		{
			Name: "generate:admin",
			Action: func(c *cli.Context) error {
				email := c.Args()[0]
				password := c.Args()[1]

				err := seeders.GenerateAdmin(email, password)

				if err != nil {
					return err
				}

				return nil
			},
		},
		{
			Name: "db:seed",
			Action: func(c *cli.Context) error {
				err := seeders.GenerateBaseCategory()

				if err != nil {
					return err
				}

				return nil
			},
		},
	}

	err := cmdApp.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
