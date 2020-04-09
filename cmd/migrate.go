package cmd

import (
	"github.com/jinzhu/gorm"
	"github.com/locpham24/go-authentication/model"
	"github.com/micro/cli/v2"
)

var Migrate = cli.Command{
	Name:  "migrate",
	Usage: "Migrate schema to database",
	Action: func(c *cli.Context) error {
		db := c.App.Metadata["db"].(*gorm.DB)
		db.DropTableIfExists(&model.User{})
		db.AutoMigrate(&model.User{})
		return nil
	},
}
