/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
)

// collectCmd represents the collect command
var collectCmd = &cobra.Command{
	Use:   "collect",
	Short: "gohostinfo collect --idc=us-east-01",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		WalkOneByOne()
	},
}

func init() {
	rootCmd.AddCommand(collectCmd)
	collectCmd.Flags().StringVar(&IDC, "idc", "", "--idc=us-east-01")
	collectCmd.Flags().StringVar(&File, "file", "gohostinfo.json", "--file=gohostinfo.json")
	collectCmd.Flags().StringVar(&Group, "group", "", "--group=")
	collectCmd.Flags().StringVar(&Tags, "tags", "", "--tags=\"tag1;tag2;tag3\"")

	collectCmd.Flags().BoolVar(&WithDocker, "withdocker", false, "--withdocker=true|false")

	collectCmd.MarkFlagRequired("idc")

}
