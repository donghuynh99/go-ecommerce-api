package cmd

import (
	"github.com/donghuynh99/ecommerce_api/database/migrations"
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
	}
}
