package tui

import (
	"sync"
)

// SyncGroup represents a group of panes that receive synchronized input
type SyncGroup struct {
	ID      string
	Name    string
	PaneIDs []string
	Enabled bool
}

// InputSynchronizer manages input synchronization across panes
type InputSynchronizer struct {
	groups      map[string]*SyncGroup
	paneToGroup map[string]string // Maps pane ID to group ID
	broadcast   bool              // Broadcast to all panes
	mu          sync.RWMutex
}

// NewInputSynchronizer creates a new input synchronizer
func NewInputSynchronizer() *InputSynchronizer {
	return &InputSynchronizer{
		groups:      make(map[string]*SyncGroup),
		paneToGroup: make(map[string]string),
		broadcast:   false,
	}
}

// CreateGroup creates a new sync group
func (is *InputSynchronizer) CreateGroup(id, name string) *SyncGroup {
	is.mu.Lock()
	defer is.mu.Unlock()

	group := &SyncGroup{
		ID:      id,
		Name:    name,
		PaneIDs: make([]string, 0),
		Enabled: true,
	}
	is.groups[id] = group
	return group
}

// DeleteGroup removes a sync group
func (is *InputSynchronizer) DeleteGroup(id string) bool {
	is.mu.Lock()
	defer is.mu.Unlock()

	group, ok := is.groups[id]
	if !ok {
		return false
	}

	// Remove pane associations
	for _, paneID := range group.PaneIDs {
		delete(is.paneToGroup, paneID)
	}

	delete(is.groups, id)
	return true
}

// GetGroup returns a sync group by ID
func (is *InputSynchronizer) GetGroup(id string) *SyncGroup {
	is.mu.RLock()
	defer is.mu.RUnlock()
	return is.groups[id]
}

// AddPaneToGroup adds a pane to a sync group
func (is *InputSynchronizer) AddPaneToGroup(groupID, paneID string) bool {
	is.mu.Lock()
	defer is.mu.Unlock()

	group, ok := is.groups[groupID]
	if !ok {
		return false
	}

	// Remove from existing group if any
	if existingGroupID, exists := is.paneToGroup[paneID]; exists {
		is.removePaneFromGroupUnsafe(existingGroupID, paneID)
	}

	group.PaneIDs = append(group.PaneIDs, paneID)
	is.paneToGroup[paneID] = groupID
	return true
}

// RemovePaneFromGroup removes a pane from its sync group
func (is *InputSynchronizer) RemovePaneFromGroup(paneID string) bool {
	is.mu.Lock()
	defer is.mu.Unlock()

	groupID, ok := is.paneToGroup[paneID]
	if !ok {
		return false
	}

	return is.removePaneFromGroupUnsafe(groupID, paneID)
}

// removePaneFromGroupUnsafe removes a pane from a group without locking
func (is *InputSynchronizer) removePaneFromGroupUnsafe(groupID, paneID string) bool {
	group, ok := is.groups[groupID]
	if !ok {
		return false
	}

	for i, id := range group.PaneIDs {
		if id == paneID {
			group.PaneIDs = append(group.PaneIDs[:i], group.PaneIDs[i+1:]...)
			delete(is.paneToGroup, paneID)
			return true
		}
	}

	return false
}

// GetPaneGroup returns the group a pane belongs to
func (is *InputSynchronizer) GetPaneGroup(paneID string) *SyncGroup {
	is.mu.RLock()
	defer is.mu.RUnlock()

	groupID, ok := is.paneToGroup[paneID]
	if !ok {
		return nil
	}

	return is.groups[groupID]
}

// IsPaneInGroup checks if a pane is in any sync group
func (is *InputSynchronizer) IsPaneInGroup(paneID string) bool {
	is.mu.RLock()
	defer is.mu.RUnlock()
	_, ok := is.paneToGroup[paneID]
	return ok
}

// GetSyncTargets returns the pane IDs that should receive input from a source pane
func (is *InputSynchronizer) GetSyncTargets(sourcePaneID string, allPaneIDs []string) []string {
	is.mu.RLock()
	defer is.mu.RUnlock()

	// If broadcast mode, return all panes except source
	if is.broadcast {
		targets := make([]string, 0)
		for _, id := range allPaneIDs {
			if id != sourcePaneID {
				targets = append(targets, id)
			}
		}
		return targets
	}

	// Check if pane is in a sync group
	groupID, ok := is.paneToGroup[sourcePaneID]
	if !ok {
		return nil
	}

	group := is.groups[groupID]
	if group == nil || !group.Enabled {
		return nil
	}

	// Return all other panes in the group
	targets := make([]string, 0)
	for _, id := range group.PaneIDs {
		if id != sourcePaneID {
			targets = append(targets, id)
		}
	}

	return targets
}

// SetBroadcast enables or disables broadcast mode
func (is *InputSynchronizer) SetBroadcast(enabled bool) {
	is.mu.Lock()
	defer is.mu.Unlock()
	is.broadcast = enabled
}

// IsBroadcastEnabled returns whether broadcast mode is enabled
func (is *InputSynchronizer) IsBroadcastEnabled() bool {
	is.mu.RLock()
	defer is.mu.RUnlock()
	return is.broadcast
}

// ToggleBroadcast toggles broadcast mode
func (is *InputSynchronizer) ToggleBroadcast() bool {
	is.mu.Lock()
	defer is.mu.Unlock()
	is.broadcast = !is.broadcast
	return is.broadcast
}

// EnableGroup enables a sync group
func (is *InputSynchronizer) EnableGroup(id string) bool {
	is.mu.Lock()
	defer is.mu.Unlock()

	group, ok := is.groups[id]
	if !ok {
		return false
	}

	group.Enabled = true
	return true
}

// DisableGroup disables a sync group
func (is *InputSynchronizer) DisableGroup(id string) bool {
	is.mu.Lock()
	defer is.mu.Unlock()

	group, ok := is.groups[id]
	if !ok {
		return false
	}

	group.Enabled = false
	return true
}

// ToggleGroup toggles a sync group's enabled state
func (is *InputSynchronizer) ToggleGroup(id string) bool {
	is.mu.Lock()
	defer is.mu.Unlock()

	group, ok := is.groups[id]
	if !ok {
		return false
	}

	group.Enabled = !group.Enabled
	return group.Enabled
}

// GetAllGroups returns all sync groups
func (is *InputSynchronizer) GetAllGroups() []*SyncGroup {
	is.mu.RLock()
	defer is.mu.RUnlock()

	groups := make([]*SyncGroup, 0, len(is.groups))
	for _, g := range is.groups {
		groups = append(groups, g)
	}
	return groups
}

// Clear removes all sync groups and resets broadcast mode
func (is *InputSynchronizer) Clear() {
	is.mu.Lock()
	defer is.mu.Unlock()

	is.groups = make(map[string]*SyncGroup)
	is.paneToGroup = make(map[string]string)
	is.broadcast = false
}

// GetStatus returns a status string for display
func (is *InputSynchronizer) GetStatus() string {
	is.mu.RLock()
	defer is.mu.RUnlock()

	if is.broadcast {
		return "[BROADCAST]"
	}

	if len(is.groups) > 0 {
		for _, g := range is.groups {
			if g.Enabled && len(g.PaneIDs) > 1 {
				return "[SYNCED]"
			}
		}
	}

	return ""
}
