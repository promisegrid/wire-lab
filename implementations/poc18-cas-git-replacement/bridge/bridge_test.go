package bridge

import (
	"os"
	"path/filepath"
	"testing"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"

	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/store"
)

// TestGitBridgeRoundtripAndRemoteFlow proves that all four conventional Git
// bridge operations use the same conversion seam instead of separate fixtures.
func TestGitBridgeRoundtripAndRemoteFlow(t *testing.T) {
	root := t.TempDir()
	sourceRepositoryPath := filepath.Join(root, "source-git")
	createSourceRepository(t, sourceRepositoryPath)

	aliceCAS, aliceErr := store.Open(filepath.Join(root, "alice-cas"))
	if aliceErr != nil {
		t.Fatalf("Open(alice) error = %v", aliceErr)
	}
	adapter := NewAdapter(aliceCAS, "alice", "bob")
	importResult, importErr := adapter.ImportRepository(sourceRepositoryPath)
	if importErr != nil {
		t.Fatalf("ImportRepository() error = %v", importErr)
	}
	if importResult.HeadSnapshotCID == "" || importResult.MappingMessageCID == "" {
		t.Fatalf("import result missing snapshot or mapping: %#v", importResult)
	}
	if importResult.Counts["git_commit"] != 2 {
		t.Fatalf("imported git commits = %d, want 2", importResult.Counts["git_commit"])
	}

	headSnapshotCID, parseErr := store.ParseCIDText(importResult.HeadSnapshotCID)
	if parseErr != nil {
		t.Fatalf("ParseCIDText(head snapshot) error = %v", parseErr)
	}
	exportRepositoryPath := filepath.Join(root, "export-git")
	exportResult, exportErr := adapter.ExportSnapshot(headSnapshotCID, exportRepositoryPath)
	if exportErr != nil {
		t.Fatalf("ExportSnapshot() error = %v", exportErr)
	}
	if exportResult.HeadGitHash == "" || exportResult.MappingMessageCID == "" {
		t.Fatalf("export result missing git hash or mapping: %#v", exportResult)
	}
	assertExportedFiles(t, exportRepositoryPath)
	assertCommitCount(t, exportRepositoryPath, 2)

	remoteRepositoryPath := filepath.Join(root, "remote.git")
	createBareMainRepository(t, remoteRepositoryPath)
	pushResult, pushErr := adapter.PushSnapshot(headSnapshotCID, remoteRepositoryPath, filepath.Join(root, "push-worktree"))
	if pushErr != nil {
		t.Fatalf("PushSnapshot() error = %v", pushErr)
	}
	if pushResult.HeadGitHash == "" {
		t.Fatalf("push result missing git hash: %#v", pushResult)
	}

	bobCAS, bobErr := store.Open(filepath.Join(root, "bob-cas"))
	if bobErr != nil {
		t.Fatalf("Open(bob) error = %v", bobErr)
	}
	pullResult, pullErr := NewAdapter(bobCAS, "bob", "alice").PullRepository(remoteRepositoryPath, filepath.Join(root, "pull-worktree"))
	if pullErr != nil {
		t.Fatalf("PullRepository() error = %v", pullErr)
	}
	if pullResult.HeadSnapshotCID == "" || pullResult.MappingMessageCID == "" {
		t.Fatalf("pull result missing snapshot or mapping: %#v", pullResult)
	}
	pulledSnapshotCID, pulledErr := store.ParseCIDText(pullResult.HeadSnapshotCID)
	if pulledErr != nil {
		t.Fatalf("ParseCIDText(pulled snapshot) error = %v", pulledErr)
	}
	if !bobCAS.Has(pulledSnapshotCID) {
		t.Fatalf("Bob CAS missing pulled snapshot %s", pullResult.HeadSnapshotCID)
	}
}

func createSourceRepository(t *testing.T, repositoryPath string) {
	t.Helper()
	repository, initErr := git.PlainInit(repositoryPath, false)
	if initErr != nil {
		t.Fatalf("PlainInit(source) error = %v", initErr)
	}
	if err := repository.Storer.SetReference(plumbing.NewSymbolicReference(plumbing.HEAD, plumbing.NewBranchReferenceName("main"))); err != nil {
		t.Fatalf("SetReference(HEAD) error = %v", err)
	}
	writeFile(t, repositoryPath, "README.md", "hello from Git\n")
	firstCommit := commitAll(t, repository, "first commit")
	writeFile(t, repositoryPath, filepath.Join("docs", "app.txt"), "second Git file\n")
	secondCommit := commitAll(t, repository, "second commit")
	if firstCommit == secondCommit {
		t.Fatalf("expected distinct commits")
	}
}

func createBareMainRepository(t *testing.T, repositoryPath string) {
	t.Helper()
	repository, initErr := git.PlainInit(repositoryPath, true)
	if initErr != nil {
		t.Fatalf("PlainInit(remote) error = %v", initErr)
	}
	if err := repository.Storer.SetReference(plumbing.NewSymbolicReference(plumbing.HEAD, plumbing.NewBranchReferenceName("main"))); err != nil {
		t.Fatalf("SetReference(remote HEAD) error = %v", err)
	}
}

func writeFile(t *testing.T, root string, relPath string, content string) {
	t.Helper()
	path := filepath.Join(root, relPath)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("MkdirAll(%s) error = %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile(%s) error = %v", path, err)
	}
}

func commitAll(t *testing.T, repository *git.Repository, message string) plumbing.Hash {
	t.Helper()
	worktree, worktreeErr := repository.Worktree()
	if worktreeErr != nil {
		t.Fatalf("Worktree() error = %v", worktreeErr)
	}
	if err := worktree.AddWithOptions(&git.AddOptions{All: true}); err != nil {
		t.Fatalf("AddWithOptions() error = %v", err)
	}
	hash, commitErr := worktree.Commit(message, &git.CommitOptions{
		Author:    bridgeSignature("alice"),
		Committer: bridgeSignature("alice"),
	})
	if commitErr != nil {
		t.Fatalf("Commit(%s) error = %v", message, commitErr)
	}
	return hash
}

func assertExportedFiles(t *testing.T, repositoryPath string) {
	t.Helper()
	readme, readErr := os.ReadFile(filepath.Join(repositoryPath, "README.md"))
	if readErr != nil {
		t.Fatalf("ReadFile(README) error = %v", readErr)
	}
	if string(readme) != "hello from Git\n" {
		t.Fatalf("README content = %q", string(readme))
	}
	app, appErr := os.ReadFile(filepath.Join(repositoryPath, "docs", "app.txt"))
	if appErr != nil {
		t.Fatalf("ReadFile(app) error = %v", appErr)
	}
	if string(app) != "second Git file\n" {
		t.Fatalf("app content = %q", string(app))
	}
}

func assertCommitCount(t *testing.T, repositoryPath string, want int) {
	t.Helper()
	repository, openErr := git.PlainOpen(repositoryPath)
	if openErr != nil {
		t.Fatalf("PlainOpen(export) error = %v", openErr)
	}
	head, headErr := repository.Head()
	if headErr != nil {
		t.Fatalf("Head() error = %v", headErr)
	}
	seen := map[plumbing.Hash]bool{}
	var walk func(plumbing.Hash)
	walk = func(hash plumbing.Hash) {
		if seen[hash] {
			return
		}
		seen[hash] = true
		commit, commitErr := repository.CommitObject(hash)
		if commitErr != nil {
			t.Fatalf("CommitObject(%s) error = %v", hash, commitErr)
		}
		for _, parentHash := range commit.ParentHashes {
			walk(parentHash)
		}
	}
	walk(head.Hash())
	if len(seen) != want {
		t.Fatalf("commit count = %d, want %d", len(seen), want)
	}
}
