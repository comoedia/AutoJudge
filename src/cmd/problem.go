/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"hoseo.dev/autojudge/src/Client"
	"hoseo.dev/autojudge/src/config"
	Log "hoseo.dev/autojudge/src/log"
	visual_viewport "hoseo.dev/autojudge/src/visual/viewport"
)

// problemCmd represents the problem command
var problemCmd = &cobra.Command{
	Use:   "problem",
	Short: "문제의 정보를 보여줍니다.",
	Long:  `메모리/시간 제한, 제출 언어, 문제 설명 등의 정보를 보여줍니다.`,
	Run: func(cmd *cobra.Command, args []string) {
		conf_problem := config.GetProblem()
		conf_endpoint := config.GetEndpoint()
		conf_submit := config.GetSubmit()

		isInvaildConf := (conf_problem.Title == "") || (conf_endpoint.Problem == "") || (conf_endpoint.Submissions == "") || (conf_endpoint.Submit == "") || (conf_problem.Description == "") || (conf_problem.Limit.Memory == "") || (conf_problem.Limit.Time == "")
		if isInvaildConf {
			Log.Error.Fatalln("환경 구성 파일의 정보가 부족합니다. 환경 구성을 다시 해주세요.")
		}

		client := Client.Get()

		contentPrefix := color.New(color.Bold).Sprint("1. 제한 사항 ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━") + "\n시간 제한: " + color.New(color.Bold).Sprint(conf_problem.Limit.Time) + "\n메모리 제한: " + color.New(color.Bold).Sprint(conf_problem.Limit.Memory) + "\n\n" + color.New(color.Bold).Sprint("2. 제출 언어 ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		langs := client.GetLangList(Client.Problem{
			Name:        conf_problem.Title,
			ProblemUrl:  conf_endpoint.Problem,
			ResultsUrl:  conf_endpoint.Submissions,
			SubmitUrl:   conf_endpoint.Submit,
			Description: conf_problem.Description,
			TimeLimit:   conf_problem.Limit.Time,
			MemoryLimit: conf_problem.Limit.Memory,
		})

		for _, lang := range langs {
			contentPrefix += "\n"
			if lang.Code == conf_submit.Lang.Index {
				contentPrefix += color.New(color.Bold).Add(color.FgHiCyan).Sprintf("* %s", lang.Name)
			} else {
				contentPrefix += "* " + lang.Name
			}
		}

		contentPrefix += "\n\n" + color.New(color.Bold).Sprint("3. 문제 설명 ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━") + "\n"

		visual_viewport.DisplayContent(conf_problem.Title, contentPrefix+conf_problem.Description)
	},
}

func init() {
	rootCmd.AddCommand(problemCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// problemCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// problemCmd.Flags().BoolP("stdout", "-o", false, "stdout으로 출력합니다.")
}
