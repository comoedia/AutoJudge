/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"hoseo.dev/autojudge/src/Client"
	Log "hoseo.dev/autojudge/src/log"
	visual_list "hoseo.dev/autojudge/src/visual/list"
)

// testserverCmd represents the testserver command
var testserverCmd = &cobra.Command{
	Use:   "testserver",
	Short: "서버에 접속해 데이터를 가져올 수 있는지 테스트합니다.",
	Long:  `서버에 접속해 데이터를 가져올 수 있는지 테스트합니다.`,
	Run: func(cmd *cobra.Command, args []string) {
		Log.Info.Println("서버 접속 테스트를 시작합니다.")

		client := Client.Get()
		Log.Info.Println("로그인에 성공했습니다!")

		classes := client.GetClasses()
		class := visual_list.GetSelctedClass(classes)
		Log.Verbose.Printf("Selected class: %s\n", class.Name)

		contests := client.GetContestList(class)
		contest := visual_list.GetSelectedContest(contests)
		Log.Verbose.Printf("Selected contest: %s\n", contest)

		problems := client.GetProblemList(contest)
		problem := visual_list.GetSelectedProblem(problems)
		Log.Verbose.Printf("Selected problem: %s\n", problem.Name)

		langs := client.GetLangList(problem)
		lang := visual_list.GetSelectedLanguage(langs)
		Log.Verbose.Printf("Selected language: %s (code: %d)\n", lang.Name, lang.Code)

		Log.Info.Println("데이터 가져오기 테스트를 성공했습니다.")
	},
}

func init() {
	rootCmd.AddCommand(testserverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// testserverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// testserverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
