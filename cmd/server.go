package cmd

import (
	"github.com/jinzhu/gorm"
	"github.com/locpham24/go-authentication/service"
	"github.com/micro/cli/v2"
)

var Server = cli.Command{
	Name:  "server",
	Usage: "Start the server",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "port",
			Aliases: []string{"p"},
		},
	},
	Action: func(c *cli.Context) error {
		db := c.App.Metadata["db"].(*gorm.DB)
		authSvc := service.NewAuthService(db)
		authSvc.Start(c)
		return nil
	},
}
