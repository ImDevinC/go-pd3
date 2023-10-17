package app

import (
	"strings"

	"github.com/ImDevinC/go-pd3/internal/config"
	"github.com/ImDevinC/go-pd3/internal/models"
	"github.com/ImDevinC/go-pd3/internal/pd3"
	"github.com/gdamore/tcell/v2"
	_ "github.com/joho/godotenv/autoload"
	"github.com/rivo/tview"
)

type App struct {
	app           *tview.Application
	hideCompleted bool
	hideLocked    bool
	tui           tui
	filter        string
	Challenges    []models.PD3DataResponse
	username      string
	password      string
}

type tui struct {
	login  *tview.InputField
	filter *tview.InputField
	grid   *tview.Grid
	table  *tview.Table
	title  tview.Primitive
	footer *tview.Flex
}

func newPrimitive(text string) tview.Primitive {
	return tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText(text)
}

func (a *App) Run() error {
	a.hideCompleted = false
	a.hideLocked = false
	a.tui.title = newPrimitive("PD3 Challenges")
	a.tui.filter = tview.NewInputField().SetLabel("Filter: ").SetFieldWidth(0).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter || key == tcell.KeyEscape {
			a.refreshGrid(false, false)
		}
	}).SetChangedFunc(func(text string) {
		a.filter = strings.ToLower(text)
		a.populateTable()
	})
	a.tui.login = tview.NewInputField().SetLabel("Username: ").SetFieldWidth(0).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			if a.tui.login.GetLabel() == "Username: " {
				a.username = a.tui.login.GetText()
				a.tui.login.SetLabel("Password: ").SetText("").SetMaskCharacter('*')
			} else {
				a.password = a.tui.login.GetText()
				a.refreshGrid(false, false)
				a.tui.login.SetLabel("Username: ").SetText("").SetMaskCharacter(0)
				a.populateFooter("Retrieving challenges...")
				err := a.login()
				if err != nil {
					a.populateFooter(err.Error())
				} else {
					a.populateFooter("Challenges updated")
				}
			}
			return
		}
		if key == tcell.KeyEscape {
			a.refreshGrid(false, false)
		}

	})
	a.tui.footer = tview.NewFlex()
	a.populateFooter("")
	a.buildTable()
	a.populateTable()
	a.tui.grid = tview.NewGrid()
	a.app = tview.NewApplication().SetRoot(a.tui.grid, true).SetFocus(a.tui.grid)
	a.refreshGrid(false, false)
	return a.app.Run()
}

func (a *App) buildTable() {
	a.tui.table = tview.NewTable()
	a.populateTable()
	a.tui.table.SetFixed(1, 0)
	a.tui.table.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			a.app.Stop()
		}
	})
	a.tui.table.SetSelectable(true, false)
	a.tui.table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'c':
			a.hideCompleted = !a.hideCompleted
			a.populateTable()
		case '/':
			a.refreshGrid(true, false)
		case 'r':
			// a.refreshData()
			a.refreshGrid(false, true)
		case 'l':
			a.hideLocked = !a.hideLocked
			a.populateTable()
		}
		return event
	})
	a.tui.table.Select(1, 0)
}

func (a *App) populateTable() {
	a.tui.table.Clear()
	a.tui.table.SetCell(0, 1, tview.NewTableCell("Name").SetAlign(tview.AlignCenter).SetBackgroundColor(tcell.Color105))
	a.tui.table.SetCell(0, 2, tview.NewTableCell("Description").SetAlign(tview.AlignCenter).SetBackgroundColor(tcell.Color105))
	a.tui.table.SetCell(0, 3, tview.NewTableCell("Tags").SetAlign(tview.AlignCenter).SetBackgroundColor(tcell.Color105))
	count := 0
	for _, item := range a.Challenges {
		cellTextColor := tcell.ColorWhite
		status := ""
		isAchievement := false
		for _, t := range item.Challenge.Tags {
			if strings.ToLower(t) == "achievement" {
				isAchievement = true
				continue
			}
		}
		if isAchievement {
			continue
		}
		if strings.ToLower(item.Status) == "completed" {
			if a.hideCompleted {
				continue
			}
			cellTextColor = tcell.ColorGreen
		} else if strings.ToLower(item.Status) == "init" {
			if a.hideLocked {
				continue
			}
			cellTextColor = tcell.ColorGray
		}
		if a.filter != "" &&
			!strings.Contains(strings.ToLower(item.Challenge.Name), a.filter) &&
			!strings.Contains(strings.ToLower(item.Challenge.Description), a.filter) &&
			!strings.Contains(strings.ToLower(strings.Join(item.Challenge.Tags, ", ")), a.filter) {
			continue
		}
		count += 1
		a.tui.table.SetCell(count, 0, tview.NewTableCell(status).SetTextColor(cellTextColor).SetAlign(tview.AlignCenter))
		a.tui.table.SetCell(count, 1, tview.NewTableCell(item.Challenge.Name).SetTextColor(cellTextColor))
		a.tui.table.SetCell(count, 2, tview.NewTableCell(item.Challenge.Description).SetTextColor(cellTextColor))
		a.tui.table.SetCell(count, 3, tview.NewTableCell(strings.Join(item.Challenge.Tags, ", ")).SetTextColor(cellTextColor))
	}
}

func (a *App) refreshGrid(showFilter bool, showLogin bool) {
	rows := []int{1, 0, 1}
	if showFilter || showLogin {
		rows = []int{1, 1, 0, 1}
	}

	a.tui.grid.SetRows(rows...).
		SetColumns(1).
		SetBorders(true).
		AddItem(a.tui.title, 0, 0, 1, 3, 0, 0, false)

	if showFilter {
		a.tui.grid.
			AddItem(a.tui.filter, 1, 0, 1, 3, 0, 0, true).
			AddItem(a.tui.table, 2, 0, 1, 3, 0, 0, false).
			AddItem(a.tui.footer, 3, 0, 1, 3, 0, 0, false)
		a.app.SetFocus(a.tui.filter)
	} else if showLogin {
		a.tui.grid.
			AddItem(a.tui.login, 1, 0, 1, 3, 0, 0, true).
			AddItem(a.tui.table, 2, 0, 1, 3, 0, 0, false).
			AddItem(a.tui.footer, 3, 0, 1, 3, 0, 0, false)
		a.app.SetFocus(a.tui.login)
	} else {
		a.tui.grid.
			AddItem(a.tui.table, 1, 0, 1, 3, 0, 0, true).
			AddItem(a.tui.footer, 2, 0, 1, 3, 0, 0, false)
		a.app.SetFocus(a.tui.grid)
	}
}

func (a *App) populateFooter(errorMessage string) {
	a.tui.footer.Clear()
	a.tui.footer.AddItem(tview.NewTextView().SetText("<c> Toggle Completed").SetTextColor(tcell.ColorBlue), 0, 1, false)
	a.tui.footer.AddItem(tview.NewTextView().SetText("<l> Toggle Locked").SetTextColor(tcell.ColorBlue), 0, 1, false)
	a.tui.footer.AddItem(tview.NewTextView().SetText("< / > Filter").SetTextColor(tcell.ColorBlue), 0, 1, false)
	a.tui.footer.AddItem(tview.NewTextView().SetText("<r> Refresh Data").SetTextColor(tcell.ColorBlue), 0, 1, false)
	if errorMessage != "" {
		errorView := tview.NewTextView().SetText(errorMessage).SetTextColor(tcell.ColorRed)
		a.tui.footer.AddItem(errorView, 0, 1, false)
	}
}

func (a *App) login() error {
	client := pd3.New()
	err := client.Login(a.username, a.password)
	if err != nil {
		return err
	}
	challenges, err := client.GetChallenges()
	if err != nil {
		return err
	}
	a.Challenges = challenges
	a.populateTable()
	err = config.SaveChallenges(challenges)
	if err != nil {
		return err
	}
	return nil
}
