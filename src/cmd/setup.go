/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"hoseo.dev/autojudge/src/Client"
	"hoseo.dev/autojudge/src/config"
	Log "hoseo.dev/autojudge/src/log"
	"hoseo.dev/autojudge/src/util"
	visual_list "hoseo.dev/autojudge/src/visual/list"
)

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "환경 구성을 시작합니다.",
	Long:  `autojudge 환경 구성을 시작합니다.`,
	Run: func(cmd *cobra.Command, args []string) {
		Log.Info.Println("AutoJudge 환경 구성을 시작합니다.")

		if util.IsExistFile("./autojudge.json") {
			Log.Info.Println("이미 존재하는 환경 구성을 찾았습니다.")
			isRemoveOldConf := visual_list.GetRemoveOldConfYN()

			if !isRemoveOldConf {
				Log.Info.Println("프로그램을 종료합니다.")
				os.Exit(0)
			}

			if err := os.Remove("./autojudge.json"); err != nil {
				Log.Error.Fatalf("기존 환경 구성을 삭제하는데 실패했습니다.")
			}
		}

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

		// if runtime.GOOS == "windows" {
		// 	Log.Info.Println("Windows 환경을 감지했습니다.")
		// 	isNeedBatchfile := visual_list.GetCreateBatYN()

		// 	if isNeedBatchfile {
		// 		// TODO: Create windows batchfile
		// 		Log.Info.Println("배치파일이 생성되었습니다.")
		// 	}
		// }

		config.New()
		config.SaveAll(config.FileRoot{
			Problem: config.Problem{
				Title:       problem.Name,
				Description: problem.Description,
				Limit: config.Limit{
					Time:   problem.TimeLimit,
					Memory: problem.MemoryLimit,
				},
			},
			Submit: config.Submit{
				Lang: config.Lang{
					Index: lang.Code,
					Str:   lang.Name,
				},
			},
			Credentials: config.Credentials{
				Username: client.Username,
				Password: client.Password,
			},
			Endpoint: config.Endpoint{
				Host: client.Host,
				Resources: config.Resources{
					Problem:     problem.ProblemUrl,
					Submissions: problem.ResultsUrl,
					Submit:      problem.SubmitUrl,
				},
			},
		})

		Log.Info.Println("환경 설정이 완료되었습니다. 이제 autojudge submit <소스코드 파일 경로>로 즉시 제출 할 수 있습니다.")
		Log.Info.Println("또한 autojudge setup로 설정을 변경할 수 있습니다.")
		Log.Info.Println("")
		Log.Info.Println("자세한 사용법은 autojudge help를 입력하여 확인하세요.")
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
