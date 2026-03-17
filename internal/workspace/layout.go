package workspace

import (
	"fmt"
)

// LayoutConfig holds configuration for layout calculations
type LayoutConfig struct {
	Width  int
	Height int
	Gap    int // Gap between panes
}

// LayoutResult contains the calculated positions for all panes
type LayoutResult struct {
	Positions []PanePosition
}

// PanePosition represents a calculated pane position
type PanePosition struct {
	ID     string
	X      int
	Y      int
	Width  int
	Height int
}

// CalculateLayout calculates pane positions based on layout type
func CalculateLayout(layout LayoutType, config LayoutConfig, paneCount int) (*LayoutResult, error) {
	if paneCount <= 0 {
		return &LayoutResult{Positions: []PanePosition{}}, nil
	}

	switch layout {
	case LayoutQuad:
		return calculateQuadLayout(config, paneCount)
	case LayoutTall:
		return calculateTallLayout(config, paneCount)
	case LayoutWide:
		return calculateWideLayout(config, paneCount)
	case LayoutStack:
		return calculateStackLayout(config, paneCount)
	case LayoutSingle:
		return calculateSingleLayout(config)
	default:
		return nil, ErrInvalidLayout
	}
}

// calculateQuadLayout calculates a 2x2 grid layout
// ┌─────────────────┬─────────────────┐
// │   Terminal 1    │   Terminal 2    │
// │   (main)        │   (ai/context)  │
// ├─────────────────┼─────────────────┤
// │   Terminal 3    │   Terminal 4    │
// │   (build/test)  │   (git/logs)    │
// └─────────────────┴─────────────────┘
func calculateQuadLayout(config LayoutConfig, paneCount int) (*LayoutResult, error) {
	positions := make([]PanePosition, 0, paneCount)

	halfWidth := (config.Width - config.Gap) / 2
	halfHeight := (config.Height - config.Gap) / 2

	// Top-left (pane 1)
	if paneCount >= 1 {
		positions = append(positions, PanePosition{
			ID:     "pane-1",
			X:      0,
			Y:      0,
			Width:  halfWidth,
			Height: halfHeight,
		})
	}

	// Top-right (pane 2)
	if paneCount >= 2 {
		positions = append(positions, PanePosition{
			ID:     "pane-2",
			X:      halfWidth + config.Gap,
			Y:      0,
			Width:  config.Width - halfWidth - config.Gap,
			Height: halfHeight,
		})
	}

	// Bottom-left (pane 3)
	if paneCount >= 3 {
		positions = append(positions, PanePosition{
			ID:     "pane-3",
			X:      0,
			Y:      halfHeight + config.Gap,
			Width:  halfWidth,
			Height: config.Height - halfHeight - config.Gap,
		})
	}

	// Bottom-right (pane 4)
	if paneCount >= 4 {
		positions = append(positions, PanePosition{
			ID:     "pane-4",
			X:      halfWidth + config.Gap,
			Y:      halfHeight + config.Gap,
			Width:  config.Width - halfWidth - config.Gap,
			Height: config.Height - halfHeight - config.Gap,
		})
	}

	return &LayoutResult{Positions: positions}, nil
}

// calculateTallLayout calculates a tall layout (1 main + 2 stacked side panes)
// ┌────────────────────┬──────────┐
// │                    │  Pane 2  │
// │     Main Pane      ├──────────┤
// │      (60%)         │  Pane 3  │
// │                    ├──────────┤
// │                    │  Pane 4  │
// └────────────────────┴──────────┘
func calculateTallLayout(config LayoutConfig, paneCount int) (*LayoutResult, error) {
	positions := make([]PanePosition, 0, paneCount)

	mainWidth := int(float64(config.Width) * 0.6)
	sideWidth := config.Width - mainWidth - config.Gap

	// Main pane (left, 60% width)
	if paneCount >= 1 {
		positions = append(positions, PanePosition{
			ID:     "pane-1",
			X:      0,
			Y:      0,
			Width:  mainWidth,
			Height: config.Height,
		})
	}

	// Calculate side pane heights
	sidePaneCount := paneCount - 1
	if sidePaneCount > 3 {
		sidePaneCount = 3
	}

	if sidePaneCount > 0 {
		sideHeight := (config.Height - config.Gap*(sidePaneCount-1)) / sidePaneCount
		y := 0

		for i := 0; i < sidePaneCount; i++ {
			height := sideHeight
			// Last pane gets remaining height
			if i == sidePaneCount-1 {
				height = config.Height - y
			}

			positions = append(positions, PanePosition{
				ID:     fmt.Sprintf("pane-%d", i+2),
				X:      mainWidth + config.Gap,
				Y:      y,
				Width:  sideWidth,
				Height: height,
			})

			y += height + config.Gap
		}
	}

	return &LayoutResult{Positions: positions}, nil
}

// calculateWideLayout calculates a wide layout (1 top + 2 bottom panes)
// ┌─────────────────────────────────┐
// │           Top Pane              │
// │            (60%)                │
// ├─────────────────┬───────────────┤
// │    Pane 2       │    Pane 3     │
// └─────────────────┴───────────────┘
func calculateWideLayout(config LayoutConfig, paneCount int) (*LayoutResult, error) {
	positions := make([]PanePosition, 0, paneCount)

	topHeight := int(float64(config.Height) * 0.6)
	bottomHeight := config.Height - topHeight - config.Gap

	// Top pane (60% height)
	if paneCount >= 1 {
		positions = append(positions, PanePosition{
			ID:     "pane-1",
			X:      0,
			Y:      0,
			Width:  config.Width,
			Height: topHeight,
		})
	}

	// Calculate bottom pane widths
	bottomPaneCount := paneCount - 1
	if bottomPaneCount > 3 {
		bottomPaneCount = 3
	}

	if bottomPaneCount > 0 {
		bottomWidth := (config.Width - config.Gap*(bottomPaneCount-1)) / bottomPaneCount
		x := 0

		for i := 0; i < bottomPaneCount; i++ {
			width := bottomWidth
			// Last pane gets remaining width
			if i == bottomPaneCount-1 {
				width = config.Width - x
			}

			positions = append(positions, PanePosition{
				ID:     fmt.Sprintf("pane-%d", i+2),
				X:      x,
				Y:      topHeight + config.Gap,
				Width:  width,
				Height: bottomHeight,
			})

			x += width + config.Gap
		}
	}

	return &LayoutResult{Positions: positions}, nil
}

// calculateStackLayout calculates a vertical stack layout
// ┌─────────────────────────────────┐
// │           Pane 1                │
// ├─────────────────────────────────┤
// │           Pane 2                │
// ├─────────────────────────────────┤
// │           Pane 3                │
// ├─────────────────────────────────┤
// │           Pane 4                │
// └─────────────────────────────────┘
func calculateStackLayout(config LayoutConfig, paneCount int) (*LayoutResult, error) {
	positions := make([]PanePosition, 0, paneCount)

	if paneCount <= 0 {
		return &LayoutResult{Positions: positions}, nil
	}

	paneHeight := (config.Height - config.Gap*(paneCount-1)) / paneCount
	y := 0

	for i := 0; i < paneCount; i++ {
		height := paneHeight
		// Last pane gets remaining height
		if i == paneCount-1 {
			height = config.Height - y
		}

		positions = append(positions, PanePosition{
			ID:     fmt.Sprintf("pane-%d", i+1),
			X:      0,
			Y:      y,
			Width:  config.Width,
			Height: height,
		})

		y += height + config.Gap
	}

	return &LayoutResult{Positions: positions}, nil
}

// calculateSingleLayout calculates a single pane layout
func calculateSingleLayout(config LayoutConfig) (*LayoutResult, error) {
	return &LayoutResult{
		Positions: []PanePosition{
			{
				ID:     "pane-1",
				X:      0,
				Y:      0,
				Width:  config.Width,
				Height: config.Height,
			},
		},
	}, nil
}

// DefaultPanesForLayout returns the default pane configuration for a layout
func DefaultPanesForLayout(layout LayoutType) []PaneState {
	switch layout {
	case LayoutQuad:
		return []PaneState{
			{ID: "pane-1", Type: PaneTypeChat, Title: "Main"},
			{ID: "pane-2", Type: PaneTypeTerminal, Title: "Terminal"},
			{ID: "pane-3", Type: PaneTypeOutput, Title: "Output"},
			{ID: "pane-4", Type: PaneTypeHistory, Title: "History"},
		}
	case LayoutTall:
		return []PaneState{
			{ID: "pane-1", Type: PaneTypeChat, Title: "Main"},
			{ID: "pane-2", Type: PaneTypeTerminal, Title: "Terminal"},
			{ID: "pane-3", Type: PaneTypeOutput, Title: "Output"},
		}
	case LayoutWide:
		return []PaneState{
			{ID: "pane-1", Type: PaneTypeChat, Title: "Main"},
			{ID: "pane-2", Type: PaneTypeTerminal, Title: "Terminal"},
			{ID: "pane-3", Type: PaneTypeOutput, Title: "Output"},
		}
	case LayoutStack:
		return []PaneState{
			{ID: "pane-1", Type: PaneTypeChat, Title: "Chat"},
			{ID: "pane-2", Type: PaneTypeTerminal, Title: "Terminal"},
			{ID: "pane-3", Type: PaneTypeOutput, Title: "Output"},
			{ID: "pane-4", Type: PaneTypeHistory, Title: "History"},
		}
	case LayoutSingle:
		return []PaneState{
			{ID: "pane-1", Type: PaneTypeChat, Title: "Main"},
		}
	default:
		return []PaneState{
			{ID: "pane-1", Type: PaneTypeChat, Title: "Main"},
		}
	}
}

// ApplyLayoutToWorkspace updates pane positions based on layout
func ApplyLayoutToWorkspace(w *Workspace, width, height int) error {
	config := LayoutConfig{
		Width:  width,
		Height: height,
		Gap:    1, // 1 character gap between panes
	}

	result, err := CalculateLayout(w.Layout, config, len(w.Panes))
	if err != nil {
		return err
	}

	// Update pane positions
	for i := range w.Panes {
		if i < len(result.Positions) {
			pos := result.Positions[i]
			w.Panes[i].X = pos.X
			w.Panes[i].Y = pos.Y
			w.Panes[i].Width = pos.Width
			w.Panes[i].Height = pos.Height
		}
	}

	return nil
}
