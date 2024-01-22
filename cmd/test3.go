/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/taylormonacelli/itsohio/test3"
)

var test3Cmd = &cobra.Command{
	Use:   "test3",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("test3 called")
		err := test3.Test3()
		if err != nil {
			fmt.Println("error running test3")
		}
	},
}

func init() {
	rootCmd.AddCommand(test3Cmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// test3Cmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// test3Cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	var (
		userCount int
		batchSize int
	)

	test3Cmd.Flags().IntVar(&userCount, "user-count", 10, "number of users to insert")
	test3Cmd.Flags().IntVar(&batchSize, "batch-size", 3, "sqlite batch size")

	test3Cmd.PreRun = func(cmd *cobra.Command, args []string) {
		err := viper.BindPFlag("user-count", test3Cmd.Flags().Lookup("user-count"))
		if err != nil {
			fmt.Println("error binding user-count flag")
			os.Exit(1)
		}

		err = viper.BindPFlag("batch-size", test3Cmd.Flags().Lookup("batch-size"))
		if err != nil {
			fmt.Println("error binding batch-size flag")
			os.Exit(1)
		}
	}
}
