// Package bridge reserves the Git compatibility seam for POC18.
//
// Intent: Git import/export/push/pull must share one conversion core later, but
// the first slice must not let Git remotes become the native authority model.
// Source: DI-dofoj; DI-harih
package bridge

// Operation names conventional Git bridge directions. These are compatibility
// adapter operations, not native PromiseGrid synchronization.
type Operation string

const (
	Import Operation = "import"
	Export Operation = "export"
	Pull   Operation = "pull"
	Push   Operation = "push"
)

// Mapping records the future conversion seam between Git labels and native CIDs.
type Mapping struct {
	Operation Operation `json:"operation"`
	GitLabel  string    `json:"git_label"`
	GridCID   string    `json:"grid_cid"`
	Outcome   string    `json:"outcome"`
}
