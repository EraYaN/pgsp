package pgsp

import "charm.land/lipgloss/v2"

var (
	BoldStyle = func() lipgloss.Style {
		return lipgloss.NewStyle().
			Bold(true)
	}()
)
