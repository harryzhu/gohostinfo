/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// collectCmd represents the collect command
var collectCmd = &cobra.Command{
	Use:   "collect",
	Short: "gohostinfo collect --idc=us-east-01",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("collect is beginning ...")

		WalkOneByOne()
	},
}

func init() {
	rootCmd.AddCommand(collectCmd)
	collectCmd.Flags().StringVar(&IDC, "idc", "", "--idc=us-east-01")
	collectCmd.Flags().BoolVar(&WithDocker, "withdocker", false, "--withdocker=true|false")
	collectCmd.MarkFlagRequired("idc")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// collectCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// collectCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
