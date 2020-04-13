package cmd

import (
	"github.com/jinzhu/gorm"
	"github.com/locpham24/go-authentication/service"
	"github.com/micro/cli/v2"
)

var Client = cli.Command{
	Name:  "client",
	Usage: "Start the client",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "port",
			Aliases: []string{"p"},
		},
	},
	Action: func(c *cli.Context) error {
		db := c.App.Metadata["db"].(*gorm.DB)
		apiSvc := service.NewAPIService(db)
		apiSvc.Start(c)
		return nil
	},
}
