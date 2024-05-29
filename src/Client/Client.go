package Client

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	Log "hoseo.dev/autojudge/src/log"
	"hoseo.dev/autojudge/src/util"
)

type (
	submitScore  float64
	submitResult uint8
)

const (
	RESULT_UNKNOWN      submitResult = iota // Custom: Unknown Error
	RESULT_PASS         submitResult = iota // Pass
	RESULT_PENDING      submitResult = iota // Pending
	RESULT_ACCEPT       submitResult = iota // Accept
	RESULT_WRONG        submitResult = iota // Wrong Answer
	RESULT_TIMELIMIT    submitResult = iota // Time Limit
	RESULT_MEMLIMIT     submitResult = iota // Memory Limit
	RESULT_COMPILE      submitResult = iota // Compile Error
	RESULT_OUTPUTLIMIT  submitResult = iota // Output Limit
	RESULT_RUNTIME      submitResult = iota // Run-time Error
	RESULT_PRESENTATION submitResult = iota // Presentation Error
	RESULT_EMPTYDATA    submitResult = iota // Empty Test-data
	RESULT_INVAILDCASE  submitResult = iota // Invaild Case
	RESULT_REJECT       submitResult = iota // Reject
)

type ClassInfo struct {
	Name string
	Url  string
}

type Contest struct {
	Name string
	Url  string
}

type LangInfo struct {
	Code int
	Name string
}

type Problem struct {
	Name        string
	ProblemUrl  string
	ResultsUrl  string
	SubmitUrl   string
	Description string
	TimeLimit   string
	MemoryLimit string
}

// type SubmitResult struct {
// 	Name  string
// 	Score string
// }

type JudgeClient struct {
	Host     string
	Username string
	Password string
	client   *http.Client
}

func (t *JudgeClient) Init() {
	jar, err := cookiejar.New(nil)
	if err != nil {
		Log.Error.Fatalf("Login > preprocess: Failed to create a new cookie jar.\n\t%s", err)
	}

	t.client = &http.Client{
		Jar: jar,
	}
}

func (t *JudgeClient) Login() bool {
	Log.Verbose.Println("Login > phase 1: getting ci_session cookies")
	resource := "/index.php/auth/login"
	resp, err := t.client.Get(t.Host + resource)
	if err != nil {
		Log.Verbose.Printf("Login > failed to login (phase 1): %s", err)
		return false
	}

	// Check HTTP status code
	if resp.StatusCode != 200 {
		Log.Verbose.Printf("Login > failed to login (phase 1, status code: %d)", resp.StatusCode)
		return false
	}

	// prepare form data
	values := make(url.Values)
	values.Set("id", t.Username)
	values.Set("password", t.Password)

	Log.Verbose.Printf("credential length: %d", len(values.Encode()))

	Log.Verbose.Println("Login > phase 2: attempting to login")
	resource = "/index.php/auth/authentication?returnURL="
	resp, err = t.client.Post(t.Host+resource, "application/x-www-form-urlencoded", strings.NewReader(values.Encode()))
	if err != nil {
		Log.Verbose.Printf("Login > failed to login (phase 2): %s", err)
		return false
	}

	// Check HTTP status code
	if resp.StatusCode != 200 {
		Log.Verbose.Printf("Login > failed to login (phase 2, status code: %d)", resp.StatusCode)
		return false
	}

	Log.Verbose.Println("Login > phase 3: getting classes")
	resource = "/index.php/judge"
	resp, err = t.client.Get(t.Host + resource)
	if err != nil {
		Log.Verbose.Printf("Login > failed to login (phase 3): %s", err)
		return false
	}

	// Check HTTP status code
	if resp.StatusCode != 200 {
		Log.Verbose.Printf("Login > failed to login (phase 3, status code: %d)", resp.StatusCode)
		return false
	}

	html, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		Log.Verbose.Printf("Login > failed to login (phase 3): %s", err)
		return false
	}

	result := html.Find("body > div.container-fluid > div > div.col-md-6 > div > a.list-group-item.active")
	if result.Text() == "" {
		Log.Verbose.Printf("Login > failed to login (phase 3, cannot find selector)")
		return false
	}

	Log.Verbose.Println("Login > OK.")
	return true
}

func (t *JudgeClient) GetClasses() []ClassInfo {
	resource := "/index.php/judge"
	resp, err := t.client.Get(t.Host + resource)
	if err != nil {
		log.Fatal(err)
	}

	// Check HTTP status code
	if resp.StatusCode != 200 {
		Log.Error.Fatalf("getClasses > invaild response. (status code: %d)", resp.StatusCode)
	}

	html, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var classes []ClassInfo

	ClassList := html.Find("body > div.container-fluid > div > div.col-md-6 > div > a.list-group-item").Each(func(i int, s *goquery.Selection) {
		// drop title element
		if i == 0 {
			return
		}

		val, isExist := s.Attr("href")
		if !isExist {
			Log.Warn.Printf("GetClasses > %d, %s > cannot find href attribute. ignored.\n", i, s.Text())
			return
		}

		Log.Verbose.Printf("(%d) Name: %s, location: %s\n", i, s.Text(), val)
		classes = append(classes, ClassInfo{
			Name: s.Text(),
			Url:  val,
		})
	})

	if ClassList.Text() == "" {
		Log.Error.Fatalf("GetClasses > failed to login (phase 3, cannot find selector)")
	}

	return classes
}

func (t *JudgeClient) GetContestList(class ClassInfo) []Contest {
	resource := class.Url
	resp, err := t.client.Get(t.Host + resource)
	if err != nil {
		log.Fatal(err)
	}

	// Check HTTP status code
	if resp.StatusCode != 200 {
		Log.Error.Fatalf("GetProblemList > invaild response. (status code: %d)", resp.StatusCode)
	}

	html, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var contests []Contest

	html.Find("body > div.container-fluid > div > div:nth-child(1) > div > ul > li > a").Each(func(i int, s *goquery.Selection) {
		val, isExist := s.Attr("href")

		// if href attribute is not exist, ignore this contest
		if !isExist {
			Log.Warn.Printf("GetContestList > %d, %s > cannot find href attribute. ignored.\n", i, s.Text())
			return
		}

		Log.Verbose.Printf("GetContestList (at html.Find) > name: %s, url: %s", util.TrimString(s.Text()), val)

		contests = append(contests, Contest{
			Name: util.TrimString(s.Text()),
			Url:  val,
		})

		Log.Verbose.Printf("GetContestList > (%d) Name: %s, location: %s\n", i, util.TrimString(s.Text()), val)
	})

	return contests
}

func (t *JudgeClient) GetProblemList(contest Contest) []Problem {
	resource := contest.Url
	resp, err := t.client.Get(t.Host + resource)
	if err != nil {
		log.Fatal(err)
	}

	// Check HTTP status code
	if resp.StatusCode != 200 {
		Log.Error.Fatalf("GetProblemList > invaild response. (status code: %d)", resp.StatusCode)
	}

	html, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var problems []Problem

	html.Find("body > div.container-fluid > div > div:nth-child(5) > div > table > tbody > tr").Each(func(i int, s *goquery.Selection) {
		val, isExist := s.Find("td:nth-child(2) > a").Attr("href")

		// if href attribute is not exist, ignore this contest
		if !isExist {
			Log.Warn.Printf("GetProblemList > %d, %s > cannot find href attribute. ignored.\n", i, s.Text())
			return
		}

		// get problem content
		resource = val
		resp, err = t.client.Get(t.Host + resource)

		// parse problem content
		problemContent, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		// get problem information
		title := problemContent.Find("body > div.container-fluid > div > div:nth-child(5) > div.table-responsive > table > tbody > tr > td:nth-child(1)").Text()
		timelimit := problemContent.Find("body > div.container-fluid > div > div:nth-child(5) > div.table-responsive > table > tbody > tr > td:nth-child(2)").Text()
		memorylimit := problemContent.Find("body > div.container-fluid > div > div:nth-child(5) > div.table-responsive > table > tbody > tr > td:nth-child(3)").Text()
		description := problemContent.Find("body > div.container-fluid > div > div:nth-child(5) > pre").Text()

		// extract problem id
		// ex) /index.php/judge/contestprobleminfo/128/154/30 => [128, 154, 30] => 128/154/30
		spilted := strings.Split(val, "/")
		spilted = spilted[4:]
		id := strings.Join(spilted, "/")

		problems = append(problems, Problem{
			Name:        util.TrimString(title),
			ProblemUrl:  val,
			ResultsUrl:  "/index.php/judge/status/" + id, // need uid=<username>
			SubmitUrl:   "/index.php/judge/submit/" + id,
			Description: description,
			TimeLimit:   util.TrimString(timelimit),
			MemoryLimit: util.TrimString(memorylimit),
		})

		Log.Verbose.Printf("GetProblemList > title: %s\n", util.TrimString(title))
	})

	return problems
}

func (t *JudgeClient) GetLangList(problem Problem) []LangInfo {
	resource := problem.ProblemUrl
	resp, err := t.client.Get(t.Host + resource)
	if err != nil {
		Log.Error.Fatalln("GetLangList > failed to getting problem info.")
	}

	// Check HTTP status code
	if resp.StatusCode != 200 {
		Log.Error.Fatalf("GetLangList > invaild response. (status code: %d)\n", resp.StatusCode)
	}

	html, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	editorUrl, isExist := html.Find("body > div.container-fluid > div > div:nth-child(5) > div:nth-child(3) > a").Attr("href")
	if !isExist {
		Log.Error.Fatalln("GetLangList > cannot find editor url.")
	}

	resource = editorUrl
	resp, err = t.client.Get(t.Host + resource)
	if err != nil {
		Log.Error.Fatalln("GetLangList > failed to getting language list.")
	}

	// Check HTTP status code
	if resp.StatusCode != 200 {
		Log.Error.Fatalf("GetLangList > invaild response. (status code: %d)\n", resp.StatusCode)
	}

	html, err = goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var langlist []LangInfo
	html.Find("#submit_form > label").Each(func(i int, s *goquery.Selection) {
		classname, _ := s.Attr("class")
		if classname != "radio-inline" {
			Log.Verbose.Printf("GetLangList > (%d) this element is not radio button. ignored.\n", i)
			return
		}

		langCode, isExist := s.Find("input").Attr("value")
		if !isExist {
			Log.Verbose.Printf("GetLangList > (%d) cannot find langcode. ignored.\n", i)
			return
		}

		i, err = strconv.Atoi(langCode)
		if err != nil {
			Log.Verbose.Printf("GetLangList > %d > cannot convert value to integer. ignored.\n", i)
			return
		}

		langlist = append(langlist, LangInfo{
			Code: i,
			Name: util.TrimString(s.Text()),
		})

		Log.Verbose.Printf("GetLangList > %s (code: %d)\n", util.TrimString(s.Text()), i)
	})

	if len(langlist) == 0 {
		Log.Error.Fatalln("GetLangList > failed to get language list. (empty)")
	}

	Log.Verbose.Printf("GetLangList > find %d langcode(s)\n", len(langlist))

	return langlist
}

func (t *JudgeClient) SubmitSolution(problem Problem, sourceCode string, langCode int) {
	// prepare form data
	values := make(url.Values)
	values.Set("lang", fmt.Sprint(langCode))
	values.Set("real_source_code", sourceCode)

	// submit form
	resource := problem.SubmitUrl
	resp, err := t.client.Post(t.Host+resource, "application/x-www-form-urlencoded", strings.NewReader(values.Encode()))
	if err != nil {
		Log.Error.Fatalf("SubmitProblem > failed to submit solution. (problem: %s)", problem.Name)
	}

	// Check HTTP status code
	if resp.StatusCode != 200 {
		Log.Error.Fatalf("SubmitProblem > invaild response. (status code: %d)", resp.StatusCode)
	}
}

func (t *JudgeClient) GetResult(problem Problem) (submitResult, submitScore) {
	// prepare form data
	values := make(url.Values)
	values.Set("uid", t.Username)

	result := RESULT_UNKNOWN
	score := submitScore(-1.0)
	for i := 0; i < 10; i++ {
		resource := problem.ResultsUrl
		resp, err := t.client.Get(t.Host + resource + "?" + values.Encode())
		if err != nil {
			Log.Error.Fatalf("GetResult > failed to get result. (problem: %s)", problem.Name)
		}

		// Check HTTP status code
		if resp.StatusCode != 200 {
			Log.Error.Fatalf("GetResult > invaild response. (status code: %d)", resp.StatusCode)
		}

		html, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		html.Find("#result-tab > tbody > tr:nth-child(1)").Each(func(i int, s *goquery.Selection) {
			// get message and score
			resultStr := s.Find("td:nth-child(4) > btn").Text()
			scoreStr := s.Find("td:nth-child(5) > span").Text()

			if resultStr != "" {
				Log.Verbose.Printf("GetResult > result message: %s\n", resultStr)
				switch resultStr {
				case "Pass":
					result = RESULT_PASS
				case "Pending":
					Log.Verbose.Println("GetResult > judge is pending. waiting...")
				case "Accept":
					result = RESULT_ACCEPT
				case "Wrong Answer":
					result = RESULT_WRONG
				case "Time Limit":
					result = RESULT_TIMELIMIT
				case "Memory Limit":
					result = RESULT_MEMLIMIT
				case "Compile Error":
					result = RESULT_COMPILE
				case "Output Limit":
					result = RESULT_OUTPUTLIMIT
				case "Run-time Error":
					result = RESULT_RUNTIME
				case "Presentation Error":
					result = RESULT_PRESENTATION
				case "Empty Test-data":
					result = RESULT_EMPTYDATA
				case "Invaild Case":
					result = RESULT_INVAILDCASE
				case "Reject":
					result = RESULT_REJECT
				}
			}

			if scoreStr != "" {
				Log.Verbose.Printf("GetResult > raw score: %s\n", scoreStr)
				scoreFloat, err := strconv.ParseFloat(scoreStr, 64)
				if err != nil {
					Log.Warn.Printf("GetResult > invaild score. maybe judge is not finished yet.")
				}

				if scoreFloat >= 0 {
					score = submitScore(math.Floor(scoreFloat))
				}
			}
		})

		if score >= 0 {
			break
		}

		time.Sleep(1 * time.Second)
	}

	if score < 0 {
		Log.Error.Fatalf("여러번 점수 확인을 시도했으나, 점수를 확인할 수 없습니다.")
	}

	return result, score
}
