package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/taylormonacelli/itsohio/test5"
)

var test5Cmd = &cobra.Command{
	Use:   "test5",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("test5 called")
		err := test5.Test5()
		if err != nil {
			fmt.Println("error running test5")
		}
	},
}

func init() {
	rootCmd.AddCommand(test5Cmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// test5Cmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// test5Cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	var (
		userCount    int
		batchSize    int
		gormLogLevel string
	)

	test5Cmd.Flags().IntVar(&userCount, "user-count", 10, "number of users to insert")
	test5Cmd.Flags().IntVar(&batchSize, "batch-size", 3, "sqlite batch size")
	test5Cmd.Flags().StringVar(&gormLogLevel, "gorm-log-level", "silent", "gorm log level (silent, warn, error, info)")

	test5Cmd.PreRun = func(cmd *cobra.Command, args []string) {
		err := viper.BindPFlag("user-count", test5Cmd.Flags().Lookup("user-count"))
		if err != nil {
			fmt.Println("error binding user-count flag")
			os.Exit(1)
		}

		err = viper.BindPFlag("batch-size", test5Cmd.Flags().Lookup("batch-size"))
		if err != nil {
			fmt.Println("error binding batch-size flag")
			os.Exit(1)
		}

		err = viper.BindPFlag("gorm-log-level", test5Cmd.Flags().Lookup("gorm-log-level"))
		if err != nil {
			fmt.Println("error binding gorm-log-level flag")
			os.Exit(1)
		}
	}
}
