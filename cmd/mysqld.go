/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/before80/go/entrance/mysqldEntrance"
	"github.com/spf13/cobra"
)

// mysqldCmd represents the mysqld command
var mysqldCmd = &cobra.Command{
	Use:   "mysqld",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		mysqldEntrance.Do(cmd)
	},
}

func init() {
	mysqldCmd.Flags().IntP("thread-num", "t", 3, "输入线程数")

	rootCmd.AddCommand(mysqldCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// mysqldCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// mysqldCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
