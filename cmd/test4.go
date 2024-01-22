package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/taylormonacelli/itsohio/test4"
)

var test4Cmd = &cobra.Command{
	Use:   "test4",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("test4 called")
		err := test4.Test4()
		if err != nil {
			fmt.Println("error running test4")
		}
	},
}

func init() {
	rootCmd.AddCommand(test4Cmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// test4Cmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// test4Cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	var (
		userCount int
		batchSize int
	)

	test4Cmd.Flags().IntVar(&userCount, "user-count", 10, "number of users to insert")
	test4Cmd.Flags().IntVar(&batchSize, "batch-size", 3, "sqlite batch size")

	test4Cmd.PreRun = func(cmd *cobra.Command, args []string) {
		err := viper.BindPFlag("user-count", test4Cmd.Flags().Lookup("user-count"))
		if err != nil {
			fmt.Println("error binding user-count flag")
			os.Exit(1)
		}

		err = viper.BindPFlag("batch-size", test4Cmd.Flags().Lookup("batch-size"))
		if err != nil {
			fmt.Println("error binding batch-size flag")
			os.Exit(1)
		}
	}
}
