package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/taylormonacelli/itsohio/test2"
)

// test2Cmd represents the test2 command
var test2Cmd = &cobra.Command{
	Use:   "test2",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := test2.Test2()
		if err != nil {
			fmt.Println("error running test2:", err)
		}
	},
}

var (
	userCount int
	batchSize int
)

func init() {
	rootCmd.AddCommand(test2Cmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// test2Cmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// test2Cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	test2Cmd.Flags().IntVar(&userCount, "user-count", 50_000, "number of users to insert")
	test2Cmd.Flags().IntVar(&batchSize, "batch-size", 8_000, "sqlite batch size")

	test2Cmd.PreRun = func(cmd *cobra.Command, args []string) {
		err := viper.BindPFlag("user-count", cmd.Flags().Lookup("user-count"))
		if err != nil {
			fmt.Println("error binding user-count flag")
			os.Exit(1)
		}

		err = viper.BindPFlag("batch-size", test2Cmd.Flags().Lookup("batch-size"))
		if err != nil {
			fmt.Println("error binding batch-size flag")
			os.Exit(1)
		}
	}
}
