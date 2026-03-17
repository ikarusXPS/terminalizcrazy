package tui

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewZoomState(t *testing.T) {
	zs := NewZoomState()

	assert.NotNil(t, zs)
	assert.False(t, zs.IsZoomed())
	assert.Equal(t, "", zs.GetZoomedPane())
}

func TestZoomStateSetZoomed(t *testing.T) {
	zs := NewZoomState()

	zs.SetZoomed(true, "pane-1")

	assert.True(t, zs.IsZoomed())
	assert.Equal(t, "pane-1", zs.GetZoomedPane())
}

func TestZoomStateSaveLayout(t *testing.T) {
	zs := NewZoomState()
	styles := DefaultStyles()

	pane1 := NewPane("pane-1", PaneTypeChat, "Test 1", 100, 50, styles)
	pane1.X = 0
	pane1.Y = 0

	pane2 := NewPane("pane-2", PaneTypeChat, "Test 2", 100, 50, styles)
	pane2.X = 100
	pane2.Y = 0

	panes := []*Pane{pane1, pane2}

	zs.SaveLayout(panes)

	snapshot := zs.GetSavedLayout()
	assert.NotNil(t, snapshot)
	assert.Len(t, snapshot.Panes, 2)

	assert.Equal(t, "pane-1", snapshot.Panes[0].ID)
	assert.Equal(t, 0, snapshot.Panes[0].X)
	assert.Equal(t, 100, snapshot.Panes[0].Width)

	assert.Equal(t, "pane-2", snapshot.Panes[1].ID)
	assert.Equal(t, 100, snapshot.Panes[1].X)
}

func TestZoomStateClearZoom(t *testing.T) {
	zs := NewZoomState()
	styles := DefaultStyles()

	pane := NewPane("pane-1", PaneTypeChat, "Test", 100, 50, styles)
	zs.SaveLayout([]*Pane{pane})
	zs.SetZoomed(true, "pane-1")

	assert.True(t, zs.IsZoomed())
	assert.NotNil(t, zs.GetSavedLayout())

	zs.ClearZoom()

	assert.False(t, zs.IsZoomed())
	assert.Equal(t, "", zs.GetZoomedPane())
	assert.Nil(t, zs.GetSavedLayout())
}

func TestNewZoomManager(t *testing.T) {
	pm := NewPaneManager(200, 100, DefaultStyles())
	zm := NewZoomManager(pm)

	assert.NotNil(t, zm)
	assert.False(t, zm.IsZoomed())
}

func TestZoomManagerZoomFocused(t *testing.T) {
	pm := NewPaneManager(200, 100, DefaultStyles())
	zm := NewZoomManager(pm)

	result := zm.ZoomFocused()

	assert.True(t, result)
	assert.True(t, zm.IsZoomed())
	assert.NotEmpty(t, zm.GetZoomedPaneID())
}

func TestZoomManagerZoomPane(t *testing.T) {
	pm := NewPaneManager(200, 100, DefaultStyles())
	zm := NewZoomManager(pm)

	focusedPane := pm.GetFocusedPane()
	result := zm.ZoomPane(focusedPane.ID)

	assert.True(t, result)
	assert.True(t, zm.IsZoomed())
	assert.Equal(t, focusedPane.ID, zm.GetZoomedPaneID())
}

func TestZoomManagerZoomPaneNotFound(t *testing.T) {
	pm := NewPaneManager(200, 100, DefaultStyles())
	zm := NewZoomManager(pm)

	result := zm.ZoomPane("nonexistent")

	assert.False(t, result)
	assert.False(t, zm.IsZoomed())
}

func TestZoomManagerRestore(t *testing.T) {
	pm := NewPaneManager(200, 100, DefaultStyles())
	zm := NewZoomManager(pm)

	zm.ZoomFocused()
	assert.True(t, zm.IsZoomed())

	result := zm.Restore()

	assert.True(t, result)
	assert.False(t, zm.IsZoomed())
}

func TestZoomManagerRestoreNotZoomed(t *testing.T) {
	pm := NewPaneManager(200, 100, DefaultStyles())
	zm := NewZoomManager(pm)

	result := zm.Restore()

	assert.False(t, result)
}

func TestZoomManagerToggleZoom(t *testing.T) {
	pm := NewPaneManager(200, 100, DefaultStyles())
	zm := NewZoomManager(pm)

	// Toggle on
	result := zm.ToggleZoom()
	assert.True(t, result)
	assert.True(t, zm.IsZoomed())

	// Toggle off
	result = zm.ToggleZoom()
	assert.True(t, result)
	assert.False(t, zm.IsZoomed())
}

func TestZoomManagerGetZoomStatus(t *testing.T) {
	pm := NewPaneManager(200, 100, DefaultStyles())
	zm := NewZoomManager(pm)

	assert.Equal(t, "", zm.GetZoomStatus())

	zm.ZoomFocused()

	assert.Equal(t, "[ZOOMED]", zm.GetZoomStatus())
}

func TestPaneManagerToggleZoom(t *testing.T) {
	pm := NewPaneManager(200, 100, DefaultStyles())

	result := pm.ToggleZoom()
	assert.True(t, result)
	assert.True(t, pm.IsZoomed())

	result = pm.ToggleZoom()
	assert.True(t, result)
	assert.False(t, pm.IsZoomed())
}

func TestPaneManagerZoomFocused(t *testing.T) {
	pm := NewPaneManager(200, 100, DefaultStyles())

	result := pm.ZoomFocused()
	assert.True(t, result)
	assert.True(t, pm.IsZoomed())
	assert.NotEmpty(t, pm.GetZoomedPaneID())
}

func TestPaneManagerRestoreZoom(t *testing.T) {
	pm := NewPaneManager(200, 100, DefaultStyles())

	pm.ZoomFocused()
	assert.True(t, pm.IsZoomed())

	result := pm.RestoreZoom()
	assert.True(t, result)
	assert.False(t, pm.IsZoomed())
}

func TestPaneManagerToggleFloating(t *testing.T) {
	pm := NewPaneManager(200, 100, DefaultStyles())
	focused := pm.GetFocusedPane()

	assert.False(t, pm.IsFloating(focused.ID))

	result := pm.ToggleFloating()
	assert.True(t, result)
	assert.True(t, pm.IsFloating(focused.ID))

	result = pm.ToggleFloating()
	assert.True(t, result)
	assert.False(t, pm.IsFloating(focused.ID))
}

func TestPaneManagerMakeFloating(t *testing.T) {
	pm := NewPaneManager(200, 100, DefaultStyles())
	focused := pm.GetFocusedPane()

	fp := pm.MakeFloating(focused.ID)

	assert.NotNil(t, fp)
	assert.True(t, pm.IsFloating(focused.ID))
}

func TestPaneManagerDockFloating(t *testing.T) {
	pm := NewPaneManager(200, 100, DefaultStyles())
	focused := pm.GetFocusedPane()

	pm.MakeFloating(focused.ID)
	assert.True(t, pm.IsFloating(focused.ID))

	result := pm.DockFloating(focused.ID)
	assert.True(t, result)
	assert.False(t, pm.IsFloating(focused.ID))
}

func TestPaneManagerGetFloatingPanes(t *testing.T) {
	pm := NewPaneManager(200, 100, DefaultStyles())

	// Split to get multiple panes
	pm.SplitVertical(PaneTypeTerminal, "Terminal")

	panes := pm.GetAllPanes()
	for _, p := range panes {
		pm.MakeFloating(p.ID)
	}

	floating := pm.GetFloatingPanes()
	assert.Len(t, floating, 2)
}

func TestPaneManagerToggleBroadcast(t *testing.T) {
	pm := NewPaneManager(200, 100, DefaultStyles())

	assert.False(t, pm.IsBroadcastEnabled())

	result := pm.ToggleBroadcast()
	assert.True(t, result)
	assert.True(t, pm.IsBroadcastEnabled())
}

func TestPaneManagerCreateSyncGroup(t *testing.T) {
	pm := NewPaneManager(200, 100, DefaultStyles())

	group := pm.CreateSyncGroup("test-group", "Test Group")

	assert.NotNil(t, group)
	assert.Equal(t, "test-group", group.ID)
}

func TestPaneManagerAddPaneToSyncGroup(t *testing.T) {
	pm := NewPaneManager(200, 100, DefaultStyles())
	focused := pm.GetFocusedPane()

	pm.CreateSyncGroup("test-group", "Test Group")
	result := pm.AddPaneToSyncGroup("test-group", focused.ID)

	assert.True(t, result)
}

func TestPaneManagerRemovePaneFromSyncGroup(t *testing.T) {
	pm := NewPaneManager(200, 100, DefaultStyles())
	focused := pm.GetFocusedPane()

	pm.CreateSyncGroup("test-group", "Test Group")
	pm.AddPaneToSyncGroup("test-group", focused.ID)

	result := pm.RemovePaneFromSyncGroup(focused.ID)
	assert.True(t, result)
}

func TestPaneManagerGetSyncTargets(t *testing.T) {
	pm := NewPaneManager(200, 100, DefaultStyles())

	// Split to get multiple panes
	pm.SplitVertical(PaneTypeTerminal, "Terminal")
	panes := pm.GetAllPanes()

	// Create sync group with both panes
	pm.CreateSyncGroup("test-group", "Test Group")
	for _, p := range panes {
		pm.AddPaneToSyncGroup("test-group", p.ID)
	}

	targets := pm.GetSyncTargets(panes[0].ID)
	assert.Len(t, targets, 1)
	assert.Equal(t, panes[1].ID, targets[0])
}

func TestPaneManagerGetInputSyncStatus(t *testing.T) {
	pm := NewPaneManager(200, 100, DefaultStyles())

	// No sync
	assert.Equal(t, "", pm.GetInputSyncStatus())

	// Broadcast
	pm.ToggleBroadcast()
	assert.Equal(t, "[BROADCAST]", pm.GetInputSyncStatus())
}

func TestPaneManagerGetStatusLine(t *testing.T) {
	pm := NewPaneManager(200, 100, DefaultStyles())

	// Empty status
	status := pm.GetStatusLine()
	assert.Equal(t, "", status)

	// Zoom
	pm.ZoomFocused()
	status = pm.GetStatusLine()
	assert.Contains(t, status, "ZOOMED")

	pm.RestoreZoom()

	// Floating
	focused := pm.GetFocusedPane()
	pm.MakeFloating(focused.ID)
	status = pm.GetStatusLine()
	assert.Contains(t, status, "FLOATING")
}
