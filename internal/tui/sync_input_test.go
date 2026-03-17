package tui

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewInputSynchronizer(t *testing.T) {
	is := NewInputSynchronizer()

	assert.NotNil(t, is)
	assert.False(t, is.IsBroadcastEnabled())
	assert.Empty(t, is.GetAllGroups())
}

func TestInputSynchronizerCreateGroup(t *testing.T) {
	is := NewInputSynchronizer()

	group := is.CreateGroup("group-1", "Test Group")

	assert.NotNil(t, group)
	assert.Equal(t, "group-1", group.ID)
	assert.Equal(t, "Test Group", group.Name)
	assert.True(t, group.Enabled)
	assert.Empty(t, group.PaneIDs)
}

func TestInputSynchronizerDeleteGroup(t *testing.T) {
	is := NewInputSynchronizer()

	is.CreateGroup("group-1", "Test Group")
	assert.Len(t, is.GetAllGroups(), 1)

	result := is.DeleteGroup("group-1")
	assert.True(t, result)
	assert.Empty(t, is.GetAllGroups())

	// Delete non-existent
	result = is.DeleteGroup("group-999")
	assert.False(t, result)
}

func TestInputSynchronizerGetGroup(t *testing.T) {
	is := NewInputSynchronizer()

	is.CreateGroup("group-1", "Test Group")

	group := is.GetGroup("group-1")
	assert.NotNil(t, group)
	assert.Equal(t, "group-1", group.ID)

	nilGroup := is.GetGroup("group-999")
	assert.Nil(t, nilGroup)
}

func TestInputSynchronizerAddPaneToGroup(t *testing.T) {
	is := NewInputSynchronizer()

	is.CreateGroup("group-1", "Test Group")

	result := is.AddPaneToGroup("group-1", "pane-1")
	assert.True(t, result)

	group := is.GetGroup("group-1")
	assert.Contains(t, group.PaneIDs, "pane-1")

	// Add to non-existent group
	result = is.AddPaneToGroup("group-999", "pane-2")
	assert.False(t, result)
}

func TestInputSynchronizerRemovePaneFromGroup(t *testing.T) {
	is := NewInputSynchronizer()

	is.CreateGroup("group-1", "Test Group")
	is.AddPaneToGroup("group-1", "pane-1")
	is.AddPaneToGroup("group-1", "pane-2")

	result := is.RemovePaneFromGroup("pane-1")
	assert.True(t, result)

	group := is.GetGroup("group-1")
	assert.NotContains(t, group.PaneIDs, "pane-1")
	assert.Contains(t, group.PaneIDs, "pane-2")

	// Remove non-existent
	result = is.RemovePaneFromGroup("pane-999")
	assert.False(t, result)
}

func TestInputSynchronizerGetPaneGroup(t *testing.T) {
	is := NewInputSynchronizer()

	is.CreateGroup("group-1", "Test Group")
	is.AddPaneToGroup("group-1", "pane-1")

	group := is.GetPaneGroup("pane-1")
	assert.NotNil(t, group)
	assert.Equal(t, "group-1", group.ID)

	nilGroup := is.GetPaneGroup("pane-999")
	assert.Nil(t, nilGroup)
}

func TestInputSynchronizerIsPaneInGroup(t *testing.T) {
	is := NewInputSynchronizer()

	is.CreateGroup("group-1", "Test Group")
	is.AddPaneToGroup("group-1", "pane-1")

	assert.True(t, is.IsPaneInGroup("pane-1"))
	assert.False(t, is.IsPaneInGroup("pane-999"))
}

func TestInputSynchronizerGetSyncTargets(t *testing.T) {
	is := NewInputSynchronizer()
	allPanes := []string{"pane-1", "pane-2", "pane-3", "pane-4"}

	is.CreateGroup("group-1", "Test Group")
	is.AddPaneToGroup("group-1", "pane-1")
	is.AddPaneToGroup("group-1", "pane-2")
	is.AddPaneToGroup("group-1", "pane-3")

	// From pane-1, should get pane-2 and pane-3
	targets := is.GetSyncTargets("pane-1", allPanes)
	assert.Len(t, targets, 2)
	assert.Contains(t, targets, "pane-2")
	assert.Contains(t, targets, "pane-3")

	// Pane not in group should have no targets
	targets = is.GetSyncTargets("pane-4", allPanes)
	assert.Empty(t, targets)
}

func TestInputSynchronizerBroadcast(t *testing.T) {
	is := NewInputSynchronizer()
	allPanes := []string{"pane-1", "pane-2", "pane-3"}

	assert.False(t, is.IsBroadcastEnabled())

	is.SetBroadcast(true)
	assert.True(t, is.IsBroadcastEnabled())

	// In broadcast mode, should get all other panes
	targets := is.GetSyncTargets("pane-1", allPanes)
	assert.Len(t, targets, 2)
	assert.Contains(t, targets, "pane-2")
	assert.Contains(t, targets, "pane-3")

	is.SetBroadcast(false)
	assert.False(t, is.IsBroadcastEnabled())
}

func TestInputSynchronizerToggleBroadcast(t *testing.T) {
	is := NewInputSynchronizer()

	assert.False(t, is.IsBroadcastEnabled())

	result := is.ToggleBroadcast()
	assert.True(t, result)
	assert.True(t, is.IsBroadcastEnabled())

	result = is.ToggleBroadcast()
	assert.False(t, result)
	assert.False(t, is.IsBroadcastEnabled())
}

func TestInputSynchronizerEnableDisableGroup(t *testing.T) {
	is := NewInputSynchronizer()

	is.CreateGroup("group-1", "Test Group")

	// Initially enabled
	group := is.GetGroup("group-1")
	assert.True(t, group.Enabled)

	// Disable
	result := is.DisableGroup("group-1")
	assert.True(t, result)
	assert.False(t, is.GetGroup("group-1").Enabled)

	// Enable
	result = is.EnableGroup("group-1")
	assert.True(t, result)
	assert.True(t, is.GetGroup("group-1").Enabled)

	// Non-existent group
	result = is.EnableGroup("group-999")
	assert.False(t, result)
}

func TestInputSynchronizerToggleGroup(t *testing.T) {
	is := NewInputSynchronizer()

	is.CreateGroup("group-1", "Test Group")

	// Initially enabled
	result := is.ToggleGroup("group-1")
	assert.False(t, result)
	assert.False(t, is.GetGroup("group-1").Enabled)

	result = is.ToggleGroup("group-1")
	assert.True(t, result)
	assert.True(t, is.GetGroup("group-1").Enabled)

	// Non-existent group
	result = is.ToggleGroup("group-999")
	assert.False(t, result)
}

func TestInputSynchronizerGetAllGroups(t *testing.T) {
	is := NewInputSynchronizer()

	is.CreateGroup("group-1", "Test Group 1")
	is.CreateGroup("group-2", "Test Group 2")

	groups := is.GetAllGroups()
	assert.Len(t, groups, 2)
}

func TestInputSynchronizerClear(t *testing.T) {
	is := NewInputSynchronizer()

	is.CreateGroup("group-1", "Test Group")
	is.AddPaneToGroup("group-1", "pane-1")
	is.SetBroadcast(true)

	is.Clear()

	assert.Empty(t, is.GetAllGroups())
	assert.False(t, is.IsBroadcastEnabled())
	assert.False(t, is.IsPaneInGroup("pane-1"))
}

func TestInputSynchronizerGetStatus(t *testing.T) {
	is := NewInputSynchronizer()

	// No sync, no broadcast
	assert.Equal(t, "", is.GetStatus())

	// Broadcast mode
	is.SetBroadcast(true)
	assert.Equal(t, "[BROADCAST]", is.GetStatus())

	is.SetBroadcast(false)

	// Synced group
	is.CreateGroup("group-1", "Test Group")
	is.AddPaneToGroup("group-1", "pane-1")
	is.AddPaneToGroup("group-1", "pane-2")
	assert.Equal(t, "[SYNCED]", is.GetStatus())
}

func TestInputSynchronizerDisabledGroupNoTargets(t *testing.T) {
	is := NewInputSynchronizer()
	allPanes := []string{"pane-1", "pane-2", "pane-3"}

	is.CreateGroup("group-1", "Test Group")
	is.AddPaneToGroup("group-1", "pane-1")
	is.AddPaneToGroup("group-1", "pane-2")

	// Disable the group
	is.DisableGroup("group-1")

	// Should have no targets when group is disabled
	targets := is.GetSyncTargets("pane-1", allPanes)
	assert.Empty(t, targets)
}

func TestInputSynchronizerPaneMovesBetweenGroups(t *testing.T) {
	is := NewInputSynchronizer()

	is.CreateGroup("group-1", "Group 1")
	is.CreateGroup("group-2", "Group 2")

	is.AddPaneToGroup("group-1", "pane-1")
	assert.Equal(t, "group-1", is.GetPaneGroup("pane-1").ID)

	// Add to group-2 should remove from group-1
	is.AddPaneToGroup("group-2", "pane-1")
	assert.Equal(t, "group-2", is.GetPaneGroup("pane-1").ID)

	group1 := is.GetGroup("group-1")
	assert.NotContains(t, group1.PaneIDs, "pane-1")
}

func TestInputSynchronizerDeleteGroupRemovesPaneMappings(t *testing.T) {
	is := NewInputSynchronizer()

	is.CreateGroup("group-1", "Test Group")
	is.AddPaneToGroup("group-1", "pane-1")
	is.AddPaneToGroup("group-1", "pane-2")

	assert.True(t, is.IsPaneInGroup("pane-1"))
	assert.True(t, is.IsPaneInGroup("pane-2"))

	is.DeleteGroup("group-1")

	assert.False(t, is.IsPaneInGroup("pane-1"))
	assert.False(t, is.IsPaneInGroup("pane-2"))
}
