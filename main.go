package main

//TODO clean stuff up, show notifications for tag/title filtering with ctrl+f
//TODO create json file for configs for themes and keybindings
import (
	"fmt"
	"io"
	"os"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Achno/gocheat/internal/components"
	tlockstyles "github.com/Achno/gocheat/styles"
)

var FilterbyTag = false

var selectUserAscii = `

█▀▀ █ █ █▀▀ ▄▀█ ▀█▀ █▀ █ █ █▀▀ █▀▀ ▀█▀
█▄▄ █▀█ ██▄ █▀█  █  ▄█ █▀█ ██▄ ██▄  █
`

// impliments list.Item interface : FilterValue()
type SelectedItem struct {
	title string
	tag   string
}

// The value the fuzzy filter , filters by
func (item SelectedItem) FilterValue() string {
	if FilterbyTag {
		return item.tag
	} else {
		return item.title
	}
}

type SelectItemDelegate struct{}

func (delegate SelectItemDelegate) Height() int { return 3 }

// Spacing
func (delegate SelectItemDelegate) Spacing() int { return 0 }

// Update
func (d SelectItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

// Render
func (d SelectItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	item, ok := listItem.(SelectedItem)

	if !ok {
		return
	}

	// Decide the renderer based on focused index
	renderer := components.ListItemInactive
	if index == m.Index() {
		renderer = components.ListItemActive
	}

	// Render
	fmt.Fprint(w, renderer(65, string(item.title), string(item.tag)))
}

// Explanation: Keybindings that the list listens on.
//
// Impliments: Keymap interface
type selectItemKeyMap struct {
	Up     key.Binding
	Down   key.Binding
	Filter key.Binding
	Back   key.Binding
}

func (k selectItemKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Filter, k.Back}
}

func (k selectItemKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up},
		{k.Down},
		{k.Filter},
		{k.Back},
	}
}

// Initialize keybin
var selectItemKeys = selectItemKeyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Filter: key.NewBinding(
		key.WithKeys("filter", "/"),
		key.WithHelp("/", "filter the list"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "exit filtering"),
	),
}

// Model for the select screen
//
// Impliments the tea.Model interface : Init() Update() View()
type SelectItemScreen struct {

	// the List ui model
	listview list.Model
}

func InitSelectItemScreen() SelectItemScreen {
	return SelectItemScreen{
		listview: components.ListViewSimple(items, SelectItemDelegate{}, 65, min(12, len(items)*3)),
	}
}

func (screen SelectItemScreen) Init() tea.Cmd {
	return nil
}

func (screen SelectItemScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// List of cmds to send
	cmds := make([]tea.Cmd, 0)

	// Handle key presses
	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			return screen, tea.Quit

		case "ctrl+f":
			FilterbyTag = !FilterbyTag
		}

	}
	// Update listview
	screen.listview, cmd = screen.listview.Update(msg)
	cmds = append(cmds, cmd)

	// Return
	return screen, tea.Batch(cmds...)
}

// View
func (screen SelectItemScreen) View() string {
	// Set height
	screen.listview.SetHeight(min(12, len(items)*3))

	// List of items to render
	items := []string{
		tlockstyles.Title(selectUserAscii), "",
		// tlockstyles.Dimmed("Select a user to login as"), "",
		screen.listview.View(), "",
	}

	// Add paginator
	if screen.listview.Paginator.TotalPages > 1 {
		items = append(items, components.Paginator(screen.listview), "")
	}

	// Add help
	items = append(items, tlockstyles.HelpView(selectItemKeys))

	// Return
	joinedItems := lipgloss.JoinVertical(
		lipgloss.Center,
		items...,
	)

	// Place list in the center
	return lipgloss.Place(screen.listview.Width(), screen.listview.Height(), lipgloss.Center, lipgloss.Center, joinedItems)
}

var items = []list.Item{
	SelectedItem{title: "Maximize Window : meta+up", tag: "Kwin"},
	SelectedItem{title: "Minimize Window : meta+m", tag: "Kwin"},
	SelectedItem{title: "Rofi : fn+end", tag: "Rofi"},
	SelectedItem{title: "Take a screenshot  : f2", tag: "Flameshot"},
	SelectedItem{title: "Open the menu : f1", tag: "wlogout"},
	SelectedItem{title: "cube : meta + w", tag: "kwin"},
	SelectedItem{title: "resize Window : alt+k", tag: "Flameshot"},
	SelectedItem{title: "Lock windows in place : ctrl+alt", tag: "Kwin"},
}

func main() {

	// create the Screen where you select items
	model := InitSelectItemScreen()

	// Initialize all lipgloss styles based on the theme and accessed by 'Styles' variable
	tlockstyles.InitializeStyles(tlockstyles.InitTheme())

	if _, err := tea.NewProgram(model, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
