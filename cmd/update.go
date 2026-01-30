package main

import (
	tea "github.com/charmbracelet/bubbletea"

	"StatusApp/configs"
)

func tickUpdate(m *BubbleTeaModel) (tea.Model, tea.Cmd) {
	m.TickCounter++
	if m.TickCounter%configs.SecondsBetweenAlternatingText == 0 {
		m.AlternatingText = !m.AlternatingText
	}
	if m.TickCounter >= configs.SecondsBetweenRefresh {
		m.TickCounter = 0
		return m, tea.Batch(fetchData(m), tickCmd())
	}
	return m, tickCmd()
}
