package workspace

import "errors"

var (
	// ErrWorkspaceNotFound is returned when a workspace is not found
	ErrWorkspaceNotFound = errors.New("workspace not found")

	// ErrInvalidWorkspaceID is returned when a workspace ID is invalid
	ErrInvalidWorkspaceID = errors.New("invalid workspace ID")

	// ErrInvalidWorkspaceName is returned when a workspace name is invalid
	ErrInvalidWorkspaceName = errors.New("invalid workspace name")

	// ErrInvalidLayout is returned when a layout type is invalid
	ErrInvalidLayout = errors.New("invalid layout type")

	// ErrPaneNotFound is returned when a pane is not found
	ErrPaneNotFound = errors.New("pane not found")

	// ErrMaxWorkspacesReached is returned when the maximum number of workspaces is reached
	ErrMaxWorkspacesReached = errors.New("maximum number of workspaces reached")

	// ErrWorkspaceAlreadyExists is returned when trying to create a workspace that already exists
	ErrWorkspaceAlreadyExists = errors.New("workspace already exists")

	// ErrCannotDeleteLastWorkspace is returned when trying to delete the last workspace
	ErrCannotDeleteLastWorkspace = errors.New("cannot delete the last workspace")
)
