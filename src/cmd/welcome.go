/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	Log "hoseo.dev/autojudge/src/log"
)

// welcomeCmd represents the welcome command
var welcomeCmd = &cobra.Command{
	Use:   "welcome",
	Short: "짧은 환영 메시지를 보여줍니다.",
	Long:  `짧은 환영 메시지를 보여줍니다.`,
	Run: func(cmd *cobra.Command, args []string) {
		Log.Info.Println("AutoJudge에 오신 것을 환영합니다!")
		Log.Info.Println("이 프로그램은 호서대학교 온라인 저지 시스템을 명령줄 도구로 관리할 수 있게 해주는 프로그램입니다.")
		Log.Info.Println("")
		Log.Info.Println("환경을 만들고자 하는 경로로 이동해서 autojudge setup을 입력하면 설정을 시작할 수 있습니다.")
		Log.Info.Println("그럼 이제 알고리즘의 굴레에 어서오세요!")
		Log.Info.Println("")
		Log.Info.Println("Made by hoseo.dev")
	},
}

func init() {
	rootCmd.AddCommand(welcomeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// welcomeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// welcomeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
