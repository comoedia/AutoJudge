/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"hoseo.dev/autojudge/src/Client"
	"hoseo.dev/autojudge/src/config"
	Log "hoseo.dev/autojudge/src/log"
	"hoseo.dev/autojudge/src/util"
	visual_list "hoseo.dev/autojudge/src/visual/list"
)

// submitCmd represents the submit command
var submitCmd = &cobra.Command{
	Use:   "submit",
	Short: "소스코드를 제출합니다.",
	Long:  `소스코드를 제출합니다. autujudge submit ./source.c 형식으로 사용합니다.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		conf_submit := config.GetSubmit()

		// if config file is not set, check if the source code file exists
		if conf_submit.File == "" {
			Log.Verbose.Println("submitCmd > cannot find the source code file path in the config file.")

			// invalid arguments
			if len(args) == 0 || args[0] == "" {
				Log.Error.Println("소스코드 파일 경로를 입력해주세요.")
			}

			// set the source code file path to the config file
			conf_submit.File = args[0]

		} else if len(args) > 0 {
			Log.Verbose.Println("submitCmd > found the source code file path in the config file. but the argument is given.")

			if conf_submit.File != args[0] {
				Log.Verbose.Println("submitCmd > the source code file path in the config file will be overwritten by the argument.")

				// set the source code file path to the config file
				conf_submit.File = args[0]
			}
		}

		isExist := util.IsExistFile(conf_submit.File)
		if !isExist {
			Log.Error.Fatalln("소스코드 파일을 찾지 못했습니다. 경로를 다시 확인해주세요.")
		}

		config.SetSubmit(conf_submit)
		Log.Verbose.Println("submitCmd > set the source code file path to the config file.")

		Log.Info.Println("제출을 시작합니다.")

		filename := conf_submit.File
		if filename == "" {
			Log.Error.Println("소스코드 파일 경로를 입력해주세요. autujudge submit <소스코드 경로> 형식으로 입력하면 됩니다.")
			os.Exit(1)
		}

		sourcecode := util.GetTextFromFile(filename)

		client := Client.Get()
		Log.Info.Println("로그인에 성공했습니다!")

		conf_problem := config.GetProblem()
		conf_endpoint := config.GetEndpoint()

		var langcode int
		var problem Client.Problem

		// Check if config file is set, if not, use interactive mode
		if (conf_problem.Title == "") || (conf_endpoint.Problem == "") || (conf_endpoint.Submissions == "") || (conf_endpoint.Submit == "") || (conf_problem.Description == "") || (conf_problem.Limit.Memory == "") || (conf_problem.Limit.Time == "") || (conf_submit.Lang.Index == 0) {
			Log.Verbose.Printf("config file is not set. Use interactive mode.\n")
			Log.Info.Println("환경 구성 파일에서 필요한 정보를 일부를 찾지 못했습니다. 대화형 모드로 진행합니다.")

			classes := client.GetClasses()
			class := visual_list.GetSelctedClass(classes)
			Log.Verbose.Printf("Selected class: %s\n", class.Name)

			contests := client.GetContestList(class)
			contest := visual_list.GetSelectedContest(contests)
			Log.Verbose.Printf("Selected contest: %s\n", contest)

			problems := client.GetProblemList(contest)
			selctedproblem := visual_list.GetSelectedProblem(problems)
			Log.Verbose.Printf("Selected problem: %s\n", selctedproblem.Name)

			langs := client.GetLangList(selctedproblem)
			lang := visual_list.GetSelectedLanguage(langs)
			Log.Verbose.Printf("Selected language: %s (code: %d)\n", lang.Name, lang.Code)

			langcode = lang.Code
			problem = selctedproblem
		} else {
			langcode = conf_submit.Lang.Index
			problem = Client.Problem{
				Name:        conf_problem.Title,
				ProblemUrl:  conf_endpoint.Problem,
				ResultsUrl:  conf_endpoint.Submissions,
				SubmitUrl:   conf_endpoint.Submit,
				Description: conf_problem.Description,
				TimeLimit:   conf_problem.Limit.Time,
				MemoryLimit: conf_problem.Limit.Memory,
			}
		}

		Log.Info.Println("(1/3) 코드를 제출하고 있습니다...")
		client.SubmitSolution(problem, sourcecode, langcode)
		time.Sleep(1 * time.Second)

		Log.Info.Println("(2/3) 결과를 기다리는 중입니다...")
		result, score := client.GetResult(problem)

		Log.Info.Println("(3/3) 결과를 받았습니다!")
		fmt.Println("")
		if result == Client.RESULT_ACCEPT || result == Client.RESULT_PASS {
			color.New(color.Bold).Add(color.FgHiGreen).Println("맞았습니다!!")
			fmt.Println("모든 케이스를 통과하였습니다.")
		} else {
			switch result {
			case Client.RESULT_TIMELIMIT:
				fmt.Println(color.HiRedString("시간 초과"))
			case Client.RESULT_MEMLIMIT:
				fmt.Println(color.HiRedString("메모리 초과"))
			case Client.RESULT_COMPILE:
				fmt.Println(color.MagentaString("컴파일 에러"))
			case Client.RESULT_OUTPUTLIMIT:
				fmt.Println(color.HiRedString("출력 초과"))
			case Client.RESULT_RUNTIME:
				fmt.Println(color.HiRedString("런타임 에러"))
			case Client.RESULT_PRESENTATION:
				fmt.Println(color.HiRedString("출력 형식이 잘못되었습니다"))
			case Client.RESULT_EMPTYDATA:
				fmt.Println(color.YellowString("테스트 데이터가 없습니다"))
				fmt.Println("일반적으로 이 문제는 출제 오류일 수 있습니다.")
			case Client.RESULT_INVAILDCASE:
				fmt.Println(color.YellowString("테스트 케이스가 유효하지 않습니다"))
				fmt.Println("일반적으로 이 문제는 출제 오류일 수 있습니다.")
			}
			fmt.Printf("점수: %.1f/100\n", score)
		}
	},
}

func init() {
	rootCmd.AddCommand(submitCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// submitCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// submitCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
