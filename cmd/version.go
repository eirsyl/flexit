package cmd

import (
	"fmt"

	"github.com/eirsyl/flexit/app"
	"github.com/spf13/cobra"
)

func CreateVersionCmd(a app.App, version string) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print application version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("%s version %s \n", a.GetName(), version)
		},
	}
}
