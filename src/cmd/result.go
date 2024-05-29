/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"hoseo.dev/autojudge/src/Client"
	"hoseo.dev/autojudge/src/config"
	Log "hoseo.dev/autojudge/src/log"
)

// resultCmd represents the result command
var resultCmd = &cobra.Command{
	Use:   "result",
	Short: "최근 답안의 제출 결과를 가져옵니다.",
	Long:  `최근 답안의 제출 결과를 가져옵니다.`,
	Run: func(cmd *cobra.Command, args []string) {
		conf_problem := config.GetProblem()
		conf_endpoint := config.GetEndpoint()

		client := Client.Get()
		Log.Info.Println("로그인에 성공했습니다!")

		isInvaildConf := (conf_problem.Title == "") || (conf_endpoint.Problem == "") || (conf_endpoint.Submissions == "") || (conf_endpoint.Submit == "") || (conf_problem.Description == "") || (conf_problem.Limit.Memory == "") || (conf_problem.Limit.Time == "")
		if isInvaildConf {
			Log.Error.Fatalln("환경 구성 파일의 정보가 부족합니다. 환경 구성을 다시 해주세요.")
		}

		Log.Info.Println("(1/2) 결과를 기다리는 중입니다...")
		result, score := client.GetResult(Client.Problem{
			Name:        conf_problem.Title,
			ProblemUrl:  conf_endpoint.Problem,
			ResultsUrl:  conf_endpoint.Submissions,
			SubmitUrl:   conf_endpoint.Submit,
			Description: conf_problem.Description,
			TimeLimit:   conf_problem.Limit.Time,
			MemoryLimit: conf_problem.Limit.Memory,
		})

		Log.Info.Println("(2/2) 결과를 받았습니다!")
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
			default:
				fmt.Println(color.HiRedString("알 수 없는 오류"))
				fmt.Println("프로그램에 정의되지 않은 오류입니다.")
			}
			fmt.Printf("점수: %.1f/100\n", score)
		}
	},
}

func init() {
	rootCmd.AddCommand(resultCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// resultCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// resultCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
