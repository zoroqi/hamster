package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zoroqi/rubbish/zettel/index"
)

var rootCmd = &cobra.Command{
	Use:   "zettel tools",
	Short: "zettel tools https://zettelkasten.de/introduction/zh/",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}

func init() {
	rootCmd.AddCommand(
		index.ZettelIndex(),
	)
}

func main() {
	Execute()
}
