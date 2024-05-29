/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	Log "hoseo.dev/autojudge/src/log"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "현재 버전을 보여줍니다.",
	Long:  `현재 버전을 보여줍니다.`,
	Run: func(cmd *cobra.Command, args []string) {
		Log.Info.Println("Version: 1.3.0")
		Log.Info.Println("Changes: 문제 정보를 볼 수 있게 되었습니다.")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// versionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// versionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
