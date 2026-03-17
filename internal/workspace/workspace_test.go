package workspace

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewWorkspace(t *testing.T) {
	w := NewWorkspace("ws-1", "Test Workspace", LayoutQuad)

	assert.Equal(t, "ws-1", w.ID)
	assert.Equal(t, "Test Workspace", w.Name)
	assert.Equal(t, LayoutQuad, w.Layout)
	assert.Empty(t, w.Panes)
	assert.False(t, w.CreatedAt.IsZero())
	assert.False(t, w.UpdatedAt.IsZero())
}

func TestWorkspaceAddPane(t *testing.T) {
	w := NewWorkspace("ws-1", "Test", LayoutQuad)

	pane := PaneState{
		ID:    "pane-1",
		Type:  PaneTypeChat,
		Title: "Main",
	}

	w.AddPane(pane)

	assert.Len(t, w.Panes, 1)
	assert.Equal(t, "pane-1", w.Panes[0].ID)
}

func TestWorkspaceRemovePane(t *testing.T) {
	w := NewWorkspace("ws-1", "Test", LayoutQuad)
	w.AddPane(PaneState{ID: "pane-1", Type: PaneTypeChat})
	w.AddPane(PaneState{ID: "pane-2", Type: PaneTypeTerminal})

	result := w.RemovePane("pane-1")
	assert.True(t, result)
	assert.Len(t, w.Panes, 1)
	assert.Equal(t, "pane-2", w.Panes[0].ID)

	// Remove non-existent pane
	result = w.RemovePane("pane-999")
	assert.False(t, result)
}

func TestWorkspaceGetPane(t *testing.T) {
	w := NewWorkspace("ws-1", "Test", LayoutQuad)
	w.AddPane(PaneState{ID: "pane-1", Type: PaneTypeChat, Title: "Main"})

	pane := w.GetPane("pane-1")
	assert.NotNil(t, pane)
	assert.Equal(t, "Main", pane.Title)

	nilPane := w.GetPane("pane-999")
	assert.Nil(t, nilPane)
}

func TestWorkspaceUpdatePane(t *testing.T) {
	w := NewWorkspace("ws-1", "Test", LayoutQuad)
	w.AddPane(PaneState{ID: "pane-1", Type: PaneTypeChat, Title: "Main"})

	updated := PaneState{ID: "pane-1", Type: PaneTypeChat, Title: "Updated"}
	result := w.UpdatePane(updated)

	assert.True(t, result)
	assert.Equal(t, "Updated", w.Panes[0].Title)

	// Update non-existent pane
	result = w.UpdatePane(PaneState{ID: "pane-999"})
	assert.False(t, result)
}

func TestWorkspaceClone(t *testing.T) {
	w := NewWorkspace("ws-1", "Original", LayoutQuad)
	w.AddPane(PaneState{ID: "pane-1", Type: PaneTypeChat})

	clone := w.Clone()

	assert.Equal(t, w.ID, clone.ID)
	assert.Equal(t, w.Name, clone.Name)
	assert.Len(t, clone.Panes, 1)

	// Modify clone
	clone.Name = "Modified"
	assert.NotEqual(t, w.Name, clone.Name)
}

func TestWorkspaceToJSON(t *testing.T) {
	w := NewWorkspace("ws-1", "Test", LayoutQuad)
	w.AddPane(PaneState{ID: "pane-1", Type: PaneTypeChat})

	jsonStr, err := w.ToJSON()
	require.NoError(t, err)
	assert.Contains(t, jsonStr, "ws-1")
	assert.Contains(t, jsonStr, "Test")
}

func TestFromJSON(t *testing.T) {
	original := NewWorkspace("ws-1", "Test", LayoutQuad)
	original.AddPane(PaneState{ID: "pane-1", Type: PaneTypeChat})

	jsonStr, err := original.ToJSON()
	require.NoError(t, err)

	restored, err := FromJSON(jsonStr)
	require.NoError(t, err)

	assert.Equal(t, original.ID, restored.ID)
	assert.Equal(t, original.Name, restored.Name)
	assert.Equal(t, original.Layout, restored.Layout)
}

func TestWorkspaceValidate(t *testing.T) {
	tests := []struct {
		name    string
		ws      *Workspace
		wantErr error
	}{
		{
			name:    "valid workspace",
			ws:      NewWorkspace("ws-1", "Test", LayoutQuad),
			wantErr: nil,
		},
		{
			name: "missing ID",
			ws: &Workspace{
				Name:   "Test",
				Layout: LayoutQuad,
			},
			wantErr: ErrInvalidWorkspaceID,
		},
		{
			name: "missing name",
			ws: &Workspace{
				ID:     "ws-1",
				Layout: LayoutQuad,
			},
			wantErr: ErrInvalidWorkspaceName,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.ws.Validate()
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestLayoutDescription(t *testing.T) {
	tests := []struct {
		layout LayoutType
		want   string
	}{
		{LayoutQuad, "2x2 grid layout"},
		{LayoutTall, "Main pane with side stack"},
		{LayoutWide, "Top pane with bottom row"},
		{LayoutStack, "Vertical stack"},
		{LayoutSingle, "Single pane"},
		{LayoutCustom, "Custom layout"},
		{LayoutType("unknown"), "Unknown layout"},
	}

	for _, tt := range tests {
		t.Run(string(tt.layout), func(t *testing.T) {
			assert.Equal(t, tt.want, tt.layout.Description())
		})
	}
}

func TestAvailableLayouts(t *testing.T) {
	layouts := AvailableLayouts()

	assert.Contains(t, layouts, LayoutQuad)
	assert.Contains(t, layouts, LayoutTall)
	assert.Contains(t, layouts, LayoutWide)
	assert.Contains(t, layouts, LayoutStack)
	assert.Contains(t, layouts, LayoutSingle)
}

func TestCalculateLayout(t *testing.T) {
	config := LayoutConfig{
		Width:  100,
		Height: 50,
		Gap:    1,
	}

	tests := []struct {
		name      string
		layout    LayoutType
		paneCount int
		wantLen   int
	}{
		{"quad with 4 panes", LayoutQuad, 4, 4},
		{"quad with 2 panes", LayoutQuad, 2, 2},
		{"tall with 3 panes", LayoutTall, 3, 3},
		{"wide with 3 panes", LayoutWide, 3, 3},
		{"stack with 4 panes", LayoutStack, 4, 4},
		{"single", LayoutSingle, 1, 1},
		{"empty", LayoutQuad, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := CalculateLayout(tt.layout, config, tt.paneCount)
			require.NoError(t, err)
			assert.Len(t, result.Positions, tt.wantLen)
		})
	}
}

func TestCalculateQuadLayout(t *testing.T) {
	config := LayoutConfig{
		Width:  100,
		Height: 50,
		Gap:    0,
	}

	result, err := calculateQuadLayout(config, 4)
	require.NoError(t, err)
	assert.Len(t, result.Positions, 4)

	// Top-left
	assert.Equal(t, 0, result.Positions[0].X)
	assert.Equal(t, 0, result.Positions[0].Y)

	// Top-right
	assert.Equal(t, 50, result.Positions[1].X)
	assert.Equal(t, 0, result.Positions[1].Y)

	// Bottom-left
	assert.Equal(t, 0, result.Positions[2].X)
	assert.Equal(t, 25, result.Positions[2].Y)

	// Bottom-right
	assert.Equal(t, 50, result.Positions[3].X)
	assert.Equal(t, 25, result.Positions[3].Y)
}

func TestCalculateTallLayout(t *testing.T) {
	config := LayoutConfig{
		Width:  100,
		Height: 50,
		Gap:    0,
	}

	result, err := calculateTallLayout(config, 3)
	require.NoError(t, err)
	assert.Len(t, result.Positions, 3)

	// Main pane should be 60% width
	assert.Equal(t, 60, result.Positions[0].Width)
	assert.Equal(t, 50, result.Positions[0].Height)
}

func TestCalculateWideLayout(t *testing.T) {
	config := LayoutConfig{
		Width:  100,
		Height: 50,
		Gap:    0,
	}

	result, err := calculateWideLayout(config, 3)
	require.NoError(t, err)
	assert.Len(t, result.Positions, 3)

	// Top pane should be 60% height
	assert.Equal(t, 100, result.Positions[0].Width)
	assert.Equal(t, 30, result.Positions[0].Height)
}

func TestCalculateStackLayout(t *testing.T) {
	config := LayoutConfig{
		Width:  100,
		Height: 40,
		Gap:    0,
	}

	result, err := calculateStackLayout(config, 4)
	require.NoError(t, err)
	assert.Len(t, result.Positions, 4)

	// Each pane should be 10 height (40 / 4)
	for _, pos := range result.Positions {
		assert.Equal(t, 100, pos.Width)
		assert.Equal(t, 10, pos.Height)
	}
}

func TestDefaultPanesForLayout(t *testing.T) {
	tests := []struct {
		layout   LayoutType
		expected int
	}{
		{LayoutQuad, 4},
		{LayoutTall, 3},
		{LayoutWide, 3},
		{LayoutStack, 4},
		{LayoutSingle, 1},
	}

	for _, tt := range tests {
		t.Run(string(tt.layout), func(t *testing.T) {
			panes := DefaultPanesForLayout(tt.layout)
			assert.Len(t, panes, tt.expected)
		})
	}
}

func TestApplyLayoutToWorkspace(t *testing.T) {
	w := NewWorkspace("ws-1", "Test", LayoutQuad)
	w.Panes = DefaultPanesForLayout(LayoutQuad)

	err := ApplyLayoutToWorkspace(w, 100, 50)
	require.NoError(t, err)

	// Verify positions were set
	for _, pane := range w.Panes {
		assert.True(t, pane.Width > 0)
		assert.True(t, pane.Height > 0)
	}
}

func TestNewManager(t *testing.T) {
	storage := NewInMemoryStorage()
	manager, err := NewManager(ManagerConfig{
		Storage: storage,
		Width:   100,
		Height:  50,
	})
	require.NoError(t, err)

	// Should have default workspace
	assert.Equal(t, 1, manager.GetWorkspaceCount())
	assert.NotNil(t, manager.GetActiveWorkspace())
}

func TestManagerCreateWorkspace(t *testing.T) {
	storage := NewInMemoryStorage()
	manager, err := NewManager(ManagerConfig{Storage: storage})
	require.NoError(t, err)

	ws, err := manager.CreateWorkspace("Test", LayoutTall)
	require.NoError(t, err)

	assert.Equal(t, "Test", ws.Name)
	assert.Equal(t, LayoutTall, ws.Layout)
	assert.Equal(t, 2, manager.GetWorkspaceCount())
}

func TestManagerSwitchWorkspace(t *testing.T) {
	storage := NewInMemoryStorage()
	manager, err := NewManager(ManagerConfig{Storage: storage})
	require.NoError(t, err)

	ws2, err := manager.CreateWorkspace("Second", LayoutTall)
	require.NoError(t, err)

	err = manager.SwitchWorkspace(ws2.ID)
	require.NoError(t, err)

	assert.Equal(t, ws2.ID, manager.GetActiveWorkspace().ID)
}

func TestManagerDeleteWorkspace(t *testing.T) {
	storage := NewInMemoryStorage()
	manager, err := NewManager(ManagerConfig{Storage: storage})
	require.NoError(t, err)

	ws2, err := manager.CreateWorkspace("Second", LayoutTall)
	require.NoError(t, err)

	err = manager.DeleteWorkspace(ws2.ID)
	require.NoError(t, err)

	assert.Equal(t, 1, manager.GetWorkspaceCount())
}

func TestManagerCannotDeleteLastWorkspace(t *testing.T) {
	storage := NewInMemoryStorage()
	manager, err := NewManager(ManagerConfig{Storage: storage})
	require.NoError(t, err)

	ws := manager.GetActiveWorkspace()
	err = manager.DeleteWorkspace(ws.ID)

	assert.Equal(t, ErrCannotDeleteLastWorkspace, err)
}

func TestManagerRenameWorkspace(t *testing.T) {
	storage := NewInMemoryStorage()
	manager, err := NewManager(ManagerConfig{Storage: storage})
	require.NoError(t, err)

	ws := manager.GetActiveWorkspace()
	err = manager.RenameWorkspace(ws.ID, "Renamed")
	require.NoError(t, err)

	assert.Equal(t, "Renamed", manager.GetActiveWorkspace().Name)
}

func TestManagerSetLayout(t *testing.T) {
	storage := NewInMemoryStorage()
	manager, err := NewManager(ManagerConfig{Storage: storage, Width: 100, Height: 50})
	require.NoError(t, err)

	ws := manager.GetActiveWorkspace()
	err = manager.SetLayout(ws.ID, LayoutTall)
	require.NoError(t, err)

	assert.Equal(t, LayoutTall, manager.GetActiveWorkspace().Layout)
}

func TestManagerDuplicateWorkspace(t *testing.T) {
	storage := NewInMemoryStorage()
	manager, err := NewManager(ManagerConfig{Storage: storage})
	require.NoError(t, err)

	ws := manager.GetActiveWorkspace()
	duplicate, err := manager.DuplicateWorkspace(ws.ID)
	require.NoError(t, err)

	assert.NotEqual(t, ws.ID, duplicate.ID)
	assert.Contains(t, duplicate.Name, "copy")
	assert.Equal(t, 2, manager.GetWorkspaceCount())
}

func TestManagerMaxWorkspaces(t *testing.T) {
	storage := NewInMemoryStorage()
	manager, err := NewManager(ManagerConfig{Storage: storage})
	require.NoError(t, err)

	// Create max workspaces - 1 (we already have 1)
	for i := 1; i < MaxWorkspaces; i++ {
		_, err := manager.CreateWorkspace("Test", LayoutQuad)
		require.NoError(t, err)
	}

	// This should fail
	_, err = manager.CreateWorkspace("TooMany", LayoutQuad)
	assert.Equal(t, ErrMaxWorkspacesReached, err)
}

func TestManagerOnChange(t *testing.T) {
	storage := NewInMemoryStorage()
	manager, err := NewManager(ManagerConfig{Storage: storage, Width: 100, Height: 50})
	require.NoError(t, err)

	var changedWs *Workspace
	manager.OnChange(func(w *Workspace) {
		changedWs = w
	})

	ws := manager.GetActiveWorkspace()
	err = manager.SetLayout(ws.ID, LayoutTall)
	require.NoError(t, err)

	assert.NotNil(t, changedWs)
	assert.Equal(t, LayoutTall, changedWs.Layout)
}

func TestManagerOnSwitch(t *testing.T) {
	storage := NewInMemoryStorage()
	manager, err := NewManager(ManagerConfig{Storage: storage})
	require.NoError(t, err)

	ws2, err := manager.CreateWorkspace("Second", LayoutTall)
	require.NoError(t, err)

	var oldWs, newWs *Workspace
	manager.OnSwitch(func(old, new *Workspace) {
		oldWs = old
		newWs = new
	})

	err = manager.SwitchWorkspace(ws2.ID)
	require.NoError(t, err)

	assert.NotNil(t, oldWs)
	assert.NotNil(t, newWs)
	assert.Equal(t, ws2.ID, newWs.ID)
}

func TestInMemoryStorage(t *testing.T) {
	storage := NewInMemoryStorage()

	ws := NewWorkspace("ws-1", "Test", LayoutQuad)
	ws.AddPane(PaneState{ID: "pane-1", Type: PaneTypeChat})

	// Save
	err := storage.SaveWorkspace(ws)
	require.NoError(t, err)

	// Get
	retrieved, err := storage.GetWorkspace("ws-1")
	require.NoError(t, err)
	assert.Equal(t, ws.ID, retrieved.ID)

	// List
	list, err := storage.ListWorkspaces()
	require.NoError(t, err)
	assert.Len(t, list, 1)

	// Delete
	err = storage.DeleteWorkspace("ws-1")
	require.NoError(t, err)

	_, err = storage.GetWorkspace("ws-1")
	assert.Equal(t, ErrWorkspaceNotFound, err)
}

func TestWorkspaceGetFloatingPanes(t *testing.T) {
	w := NewWorkspace("ws-1", "Test", LayoutQuad)
	w.AddPane(PaneState{ID: "pane-1", Floating: true})
	w.AddPane(PaneState{ID: "pane-2", Floating: false})
	w.AddPane(PaneState{ID: "pane-3", Floating: true})

	floating := w.GetFloatingPanes()
	assert.Len(t, floating, 2)
}

func TestWorkspaceGetDockedPanes(t *testing.T) {
	w := NewWorkspace("ws-1", "Test", LayoutQuad)
	w.AddPane(PaneState{ID: "pane-1", Floating: true})
	w.AddPane(PaneState{ID: "pane-2", Floating: false})
	w.AddPane(PaneState{ID: "pane-3", Floating: true})

	docked := w.GetDockedPanes()
	assert.Len(t, docked, 1)
}

func TestManagerSwitchByIndex(t *testing.T) {
	storage := NewInMemoryStorage()
	manager, err := NewManager(ManagerConfig{Storage: storage})
	require.NoError(t, err)

	ws2, err := manager.CreateWorkspace("Second", LayoutTall)
	require.NoError(t, err)

	err = manager.SwitchWorkspaceByIndex(2)
	require.NoError(t, err)

	assert.Equal(t, ws2.ID, manager.GetActiveWorkspace().ID)

	// Invalid index
	err = manager.SwitchWorkspaceByIndex(99)
	assert.Equal(t, ErrWorkspaceNotFound, err)
}

func TestManagerGetActiveIndex(t *testing.T) {
	storage := NewInMemoryStorage()
	manager, err := NewManager(ManagerConfig{Storage: storage})
	require.NoError(t, err)

	// Get the initial active workspace
	initialWs := manager.GetActiveWorkspace()
	initialIndex := manager.GetActiveIndex()
	assert.Equal(t, 1, initialIndex) // Should be index 1 (first workspace)

	// Create a new workspace (small delay to ensure different timestamps)
	time.Sleep(time.Millisecond * 10)
	ws2, err := manager.CreateWorkspace("Second", LayoutTall)
	require.NoError(t, err)

	// Initial workspace should still be active
	assert.Equal(t, initialWs.ID, manager.GetActiveWorkspace().ID)

	// Switch to ws2
	err = manager.SwitchWorkspace(ws2.ID)
	require.NoError(t, err)

	// The new workspace index depends on creation order
	newIndex := manager.GetActiveIndex()
	assert.True(t, newIndex >= 1 && newIndex <= 2, "Active index should be 1 or 2")
	assert.Equal(t, ws2.ID, manager.GetActiveWorkspace().ID)
}

func TestWorkspaceSetActivePane(t *testing.T) {
	w := NewWorkspace("ws-1", "Test", LayoutQuad)
	w.AddPane(PaneState{ID: "pane-1"})
	w.AddPane(PaneState{ID: "pane-2"})

	beforeUpdate := w.UpdatedAt
	time.Sleep(time.Millisecond)

	w.SetActivePane("pane-2")

	assert.Equal(t, "pane-2", w.ActivePane)
	assert.True(t, w.UpdatedAt.After(beforeUpdate))
}

func TestCalculateLayoutInvalidType(t *testing.T) {
	config := LayoutConfig{Width: 100, Height: 50, Gap: 1}
	_, err := CalculateLayout(LayoutType("invalid"), config, 4)
	assert.Equal(t, ErrInvalidLayout, err)
}

func TestManagerAddPaneToWorkspace(t *testing.T) {
	storage := NewInMemoryStorage()
	manager, err := NewManager(ManagerConfig{Storage: storage, Width: 100, Height: 50})
	require.NoError(t, err)

	ws := manager.GetActiveWorkspace()
	initialCount := len(ws.Panes)

	err = manager.AddPaneToWorkspace(ws.ID, PaneState{
		ID:    "new-pane",
		Type:  PaneTypeOutput,
		Title: "New Output",
	})
	require.NoError(t, err)

	updated := manager.GetWorkspace(ws.ID)
	assert.Len(t, updated.Panes, initialCount+1)
}

func TestManagerRemovePaneFromWorkspace(t *testing.T) {
	storage := NewInMemoryStorage()
	manager, err := NewManager(ManagerConfig{Storage: storage, Width: 100, Height: 50})
	require.NoError(t, err)

	ws := manager.GetActiveWorkspace()
	require.True(t, len(ws.Panes) > 0)

	paneToRemove := ws.Panes[0].ID
	err = manager.RemovePaneFromWorkspace(ws.ID, paneToRemove)
	require.NoError(t, err)

	updated := manager.GetWorkspace(ws.ID)
	for _, pane := range updated.Panes {
		assert.NotEqual(t, paneToRemove, pane.ID)
	}
}
