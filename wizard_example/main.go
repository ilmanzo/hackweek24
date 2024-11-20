package main

// POC for a "wizard-style" application that asks some information to the user

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// main purpose of this projects is to learn and explore the Go import rules and directory structure

func main() {
	p := tea.NewProgram(NewWizard())
	m, err := p.Run() // The returned model is of type tea.Model, which is an interface.
	if err != nil {
		fmt.Printf("Whops, there's been an error: %v", err)
		os.Exit(1)
	}
	// The returned model is of type tea.Model, which is an interface.
	// You need to use type assertion (m, ok := model.(wizardModel)) to convert it to your specific Model type.
	wizard, ok := m.(wizardModel)
	if !ok {
		fmt.Println("Unexpected model type!")
		os.Exit(1)
	}
	fmt.Println("Input Wizard finished.")
	fmt.Printf("User name = %s\n", wizard.nickname)
}

type Step int

const (
	StepWelcome Step = iota
	StepNickname
	StepLikesCats
	StepSummary
)

// our data model
type wizardModel struct {
	currentStep Step // where are we in the wizard ?
	nickname    string
	likesCats   bool
	input       textinput.Model
}

// this function inits our model with default data and returns a model
func NewWizard() wizardModel {
	ti := textinput.New()
	ti.Focus()
	ti.CharLimit = 20
	ti.Width = 20
	return wizardModel{
		currentStep: StepWelcome,
		nickname:    "",
		likesCats:   false,
		input:       ti,
	}
}

func (w wizardModel) View() string {
	s := ""
	switch w.currentStep {
	case StepWelcome:
		s = "Welcome to the support wizard!\nPress enter to continue."
	case StepNickname:
		s = "Please enter your nickname:"
		s += w.input.View()
	case StepLikesCats:
		s = "Do you like kittens 🐱 ?"
		s += w.input.View()
	case StepSummary:
		s = "You entered these informations on the survey:\n"
		s += "Nickname: " + w.nickname + " "
		if w.likesCats {
			s += "Likes cats 🐱"
		} else {
			s += "Does not like cats 🚫"
		}
	default:
		return "unknown step" // should never reach this
	}
	return s
}

func (w wizardModel) Init() tea.Cmd {
	return textinput.Blink
}

// the Update function receives a message and updates the model accordingly
func (w wizardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		//on Enter press, go -> to the next wizard step
		case tea.KeyCtrlC, tea.KeyEsc:
			return w, tea.Quit
		case tea.KeyEnter:
			if w.currentStep == StepWelcome {
				w.currentStep = StepNickname
				return w, nil
			} else if w.currentStep == StepNickname {
				// store the entered value inside the model and advance the wizard
				w.nickname = w.input.Value()
				w.currentStep = StepLikesCats
				w.input.Reset() // Clear the text input
				w.input.SetSuggestions([]string{"Yes", "No"})
			} else if w.currentStep == StepLikesCats {
				answer := w.input.Value()
				if len(answer) > 0 && (answer[0] == 'Y' || answer[0] == 'y') {
					w.likesCats = true
				}
				if len(answer) > 0 {
					w.currentStep = StepSummary
				}
			} else if w.currentStep == StepSummary {
				return w, tea.Quit
			}
		}
	case tea.WindowSizeMsg:
		w.input.Focus()
	}
	// delegate other events to the textInput field
	w.input, cmd = w.input.Update(msg)
	return w, cmd
}
