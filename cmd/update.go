package main

import (
	tea "github.com/charmbracelet/bubbletea"

	"StatusApp/configs"
)

func tickUpdate(m BubbleTeaModel) (tea.Model, tea.Cmd) {
	m.Model.TickCounter++
	if m.Model.TickCounter%configs.SecondsBetweenAlternatingText == 0 {
		m.Model.AlternatingText = !m.Model.AlternatingText
	}
	if m.Model.TickCounter >= configs.SecondsBetweenRefresh {
		m.Model.TickCounter = 0
		return m, tea.Batch(fetchData(m), tickCmd())
	}
	return m, tickCmd()
}
