package Client

import (
	"strings"

	"hoseo.dev/autojudge/src/config"
	Log "hoseo.dev/autojudge/src/log"
	visual_textinput "hoseo.dev/autojudge/src/visual/textinput"
)

func Get() *JudgeClient {
	preparedCredential := getCredential()
	preparedCredentialTryCnt := 0

	var client JudgeClient
	isLogined := false
	for !isLogined {
		var credential visual_textinput.Credential

		if preparedCredential.ServerAddress != "" && preparedCredential.Username != "" && preparedCredential.Password != "" {
			// if prepated credential is not empty (success to get credential from autojudge.json file)
			credential = preparedCredential
			preparedCredentialTryCnt++

			// if prepared credential try count is over 3, clear prepared credential
			if preparedCredentialTryCnt > 3 {
				preparedCredential.ServerAddress = ""
				preparedCredential.Username = ""
				preparedCredential.Password = ""

				Log.Error.Println("환경 구성을 통한 로그인에 실패했습니다. 수동 로그인 모드로 접근합니다.")
				credential = visual_textinput.GetSelectedCredential()
			}
		} else {
			// if prepared credential is empty (fail to get credential from autojudge.json file)
			credential = visual_textinput.GetSelectedCredential()
		}

		client = JudgeClient{
			Host:     strings.TrimRight(credential.ServerAddress, "/"), // remove last slash
			Username: credential.Username,
			Password: credential.Password,
		}

		client.Init()
		isLogined = client.Login()

		if !isLogined {
			Log.Error.Println("로그인에 실패했습니다. 서버 주소, 아이디, 비밀번호를 다시 확인해주세요.")
		}
	}

	return &client
}

func getCredential() visual_textinput.Credential {
	endpoint := config.GetEndpoint()
	credential := config.GetCredentials()

	if endpoint.Host == "" || credential.Username == "" || credential.Password == "" {
		Log.Error.Println("환경 구성파일을 불러오는데 실패했습니다. 수동 로그인 모드로 접근합니다.")
		return visual_textinput.Credential{
			ServerAddress: "",
			Username:      "",
			Password:      "",
		}
	}

	Log.Info.Println("AutoJudge 환경 구성파일을 불러왔습니다.")

	return visual_textinput.Credential{
		ServerAddress: endpoint.Host,
		Username:      credential.Username,
		Password:      credential.Password,
	}
}
