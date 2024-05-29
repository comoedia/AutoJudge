package visual_list

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"hoseo.dev/autojudge/src/Client"
	Log "hoseo.dev/autojudge/src/log"
)

const listHeight = 14

var stateResult = ""

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type item string

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

type model struct {
	list     list.Model
	choice   string
	quitting bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choice = string(i)
			}
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.choice != "" {
		stateResult = m.choice
		return quitTextStyle.Render(fmt.Sprintf("✅ \"%s\"로 선택되었습니다.", m.choice))
	}
	if m.quitting {
		return quitTextStyle.Render("프로그램을 종료하여 선택되지 않았습니다.")
	}
	return "\n" + m.list.View()
}

func GetSelctedClass(classes []Client.ClassInfo) Client.ClassInfo {
	items := []list.Item{}

	for _, v := range classes {
		items = append(items, item(v.Name))
	}

	const defaultWidth = 20

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "목표로 할 수업을 선택하세요."
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	m := model{list: l}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	var selected Client.ClassInfo
	for _, v := range classes {
		if v.Name == stateResult {
			selected = v
		}
	}

	if selected.Name == "" {
		Log.Error.Fatalf("GetSelectedClass > failed to get selected class. (selected: %s)", stateResult)
	}

	stateResult = ""

	return selected
}

func GetSelectedContest(contests []Client.Contest) Client.Contest {
	items := []list.Item{}

	for _, v := range contests {
		items = append(items, item(v.Name))
	}

	const defaultWidth = 20

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "목표로 할 컨테스트를 선택하세요."
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	m := model{list: l}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	var selected Client.Contest
	for _, v := range contests {
		if v.Name == stateResult {
			selected = v
		}
	}

	if selected.Name == "" {
		Log.Error.Fatalf("GetSelectedContest > failed to get selected contest. (selected: %s)", stateResult)
	}

	stateResult = ""

	return selected
}

func GetSelectedProblem(problems []Client.Problem) Client.Problem {
	items := []list.Item{}

	for _, v := range problems {
		items = append(items, item(v.Name))
	}

	const defaultWidth = 20

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "목표로 할 문제를 선택하세요."
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	m := model{list: l}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	var selected Client.Problem
	for _, v := range problems {
		if v.Name == stateResult {
			selected = v
		}
	}

	if selected.Name == "" {
		Log.Error.Fatalf("GetSelectedProblem > failed to get selected problem. (selected: %s)", stateResult)
	}

	stateResult = ""

	return selected
}

func GetSelectedLanguage(languages []Client.LangInfo) Client.LangInfo {
	items := []list.Item{}

	for _, v := range languages {
		items = append(items, item(v.Name))
	}

	const defaultWidth = 20

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "제출할 소스코드의 언어를 선택하세요."
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	m := model{list: l}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	var selected Client.LangInfo
	for _, v := range languages {
		if v.Name == stateResult {
			selected = v
		}
	}

	if selected.Name == "" {
		Log.Error.Fatalf("GetSelectedLanguage > failed to get selected language. (selected: %s)", stateResult)
	}

	stateResult = ""

	return selected
}

func GetCreateBatYN() bool {
	items := []list.Item{
		item("네, 만들어 주세요!"),
		item("아뇨, 알아서 할게요!"),
	}

	const defaultWidth = 20

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "Windows 사용자를 위해 원클릭 배치파일을 생성할까요?"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	m := model{list: l}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	if stateResult == "네, 만들어 주세요!" {
		stateResult = ""
		return true
	}

	stateResult = ""
	return false
}

func GetRemoveOldConfYN() bool {
	items := []list.Item{
		item("네, 새로 만들래요!"),
		item("아뇨, 좀 더 생각해볼래요!"),
	}

	const defaultWidth = 20

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "이미 이 환경은 구성 되어있는 것 같습니다. 기존 환경 구성 파일을 삭제할까요?"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	m := model{list: l}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	if stateResult == "네, 새로 만들래요!" {
		stateResult = ""
		return true
	}

	stateResult = ""
	return false
}
