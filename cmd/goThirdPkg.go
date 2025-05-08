/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/before80/go/entrance/goThirdPkgEntrance"
	"github.com/spf13/cobra"
)

// goThirdPkgCmd represents the goThirdPkg command
var goThirdPkgCmd = &cobra.Command{
	Use:   "goThirdPkg",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		goThirdPkgEntrance.Do(cmd)
	},
}

func init() {
	goThirdPkgIndexCmd.Flags().IntP("thread-num", "t", 3, "输入线程数")
	rootCmd.AddCommand(goThirdPkgCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:

	// goThirdPkgCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// goThirdPkgCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
