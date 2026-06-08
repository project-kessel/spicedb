---
name: sync-upstream
description: Sync the SpiceDB fork from upstream. Determines the latest supported SpiceDB version from the operator's update graph, validates the upgrade path, and syncs if a new version is available.
disable-model-invocation: true
argument-hint: [--tag <tag>]
---

# Sync SpiceDB from Upstream

Syncs this fork from upstream [authzed/spicedb](https://github.com/authzed/spicedb)
by determining the latest supported version from the SpiceDB Operator's update graph.

## Arguments

- `--tag <tag>` — (optional) override the SpiceDB tag to sync to. If omitted,
  the skill determines the latest supported version from the operator's update graph.

## GitHub URLs

The skill uses these URLs to fetch operator data without needing the operator repo
locally:

- **Operator SYNC.md**: `https://raw.githubusercontent.com/project-kessel/spicedb-operator/refs/heads/main/SYNC.md`
- **Update graph**: `https://raw.githubusercontent.com/project-kessel/spicedb-operator/refs/heads/main/config/update-graph.yaml`

## Required Remotes

- `origin` — your personal fork (used for pushing branches)
- `kessel` — `project-kessel/spicedb` (the repo PRs target)
- `upstream` — `authzed/spicedb` (where upstream tags are fetched from)

## Process

### Step 1: Validate Prerequisites

Before starting the sync, validate that all required tools and configuration are
in place. Run all checks first, then report any issues together. For any issues
found, offer to fix them (with user confirmation). Do not ask the user to perform
manual steps unless they explicitly say they prefer to handle it themselves.

1. **Validate git remotes** by running `git remote -v` and checking:
   - `upstream` exists and points to `github.com/authzed/spicedb` (HTTP or SSH)
   - `kessel` exists and points to `github.com/project-kessel/spicedb` (HTTP or SSH)
   - `origin` exists (any URL is acceptable)

   If a remote is missing, offer to add it:
   ```bash
   git remote add upstream https://github.com/authzed/spicedb.git
   git remote add kessel https://github.com/project-kessel/spicedb.git
   ```
   If a remote exists but points to the wrong URL, offer to fix it:
   ```bash
   git remote set-url <name> <correct-url>
   ```

2. **Validate `gh` CLI is installed:**
   ```bash
   command -v gh
   ```
   If not installed, inform the user that `gh` is required for creating PRs and
   posting review diffs. Link to https://cli.github.com/ and ask how they want
   to proceed.

3. **Validate `gh` CLI is authenticated:**
   ```bash
   gh auth status
   ```
   If not authenticated, offer to run `gh auth login` for the user.

4. **Validate `gh` default repo is configured:**
   ```bash
   gh repo set-default --view
   ```
   The default repo should be `project-kessel/spicedb`. If it is not set or
   points to a different repo, offer to fix it:
   ```bash
   gh repo set-default project-kessel/spicedb
   ```

Only proceed to Step 2 once all prerequisites pass.

### Step 2: Determine the Target Tag

If `--tag` was provided, use that and skip to Step 3.

Otherwise, determine the latest supported SpiceDB version:

1. Check if `validate-upgrade-path` is in PATH:
   ```bash
   which validate-upgrade-path
   ```

2. **If the tool IS available**, use it with the GitHub URL to list versions:
   ```bash
   validate-upgrade-path \
     -g https://raw.githubusercontent.com/project-kessel/spicedb-operator/refs/heads/main/config/update-graph.yaml \
     --list-versions -d postgres
   ```
   The **first version listed** is the latest supported SpiceDB version.

3. **If the tool is NOT available**, ask the user:
   > `validate-upgrade-path` is not installed. Please provide either:
   > - The path to your spicedb-operator repo so I can build and run it, or
   > - The SpiceDB tag to sync to using `--tag`
   >
   > To install the tool: `cd <spicedb-operator-repo> && make install-validate-upgrade-path`

   If the user provides the operator repo path:
   ```bash
   cd <operator-repo> && make build-validate-upgrade-path
   bin/validate-upgrade-path --list-versions -d postgres
   ```

4. Read the current SpiceDB version from this repo's `SYNC.md`

5. If the latest version matches the current version, report that SpiceDB is
   already up to date and stop

6. Validate the upgrade path from current to latest:
   ```bash
   validate-upgrade-path \
     -g https://raw.githubusercontent.com/project-kessel/spicedb-operator/refs/heads/main/config/update-graph.yaml \
     <current-version> <latest-version>
   ```

7. If the upgrade is NOT supported, report the finding and stop. Ask the user
   how to proceed.

**Important**: Before proceeding, check the `go` directive in the upstream
`go.mod` at the target tag:
```bash
git fetch upstream --tags
git show tags/<tag>:go.mod | grep -E '^go '
```
If the Go version exceeds our go-toolset version, stop and report the issue.
The current go-toolset version constraint is documented in `README-redhat.md`.

### Step 3: Sync

1. `git fetch upstream --tags` (if not already done)
2. `git fetch kessel`
3. `git checkout -b sync-upstream-<tag> tags/<tag>`
4. `git checkout -b merge-upstream-<tag> kessel/main`
5. `git merge sync-upstream-<tag>`
6. Resolve any merge conflicts using the **Merge Action** column in the drift
   tracking table in `README-redhat.md`
7. For `go.mod` and `go.sum` conflicts, reset to upstream entirely:
   ```bash
   # Find all conflicting go.mod/go.sum/go.work files
   git diff --name-only --diff-filter=U --relative | grep -E "mod|sum|work"

   # Reset them to the upstream version
   git checkout sync-upstream-<tag> -- <file1> <file2> ...
   ```
8. For any file not listed in the table, accept the upstream (incoming) version

### Step 4: Update SYNC.md

Update `SYNC.md` with the new tag and the commit SHA of the upstream tag:
```bash
git rev-parse tags/<tag>
```
Commit the update.

### Step 5: Post-Merge Cleanup

1. Run `./scripts/redhat-diff.sh --stat`
2. Remove any **stale files** listed in the warning:
   ```bash
   git rm <file1> <file2> ...
   ```
3. Reset any **diverged files** listed in the warning:
   ```bash
   git checkout tags/<tag> -- <file1> <file2> ...
   ```
4. Commit the cleanup
5. Re-run `./scripts/redhat-diff.sh --stat` to confirm clean output

### Step 6: Push and Create PR

1. Push the branch to your fork:
   ```bash
   git push origin merge-upstream-<tag>
   ```
2. Create a PR against the kessel repo (**do not squash commits when merging**):
   ```bash
   gh pr create --repo project-kessel/spicedb --base main --title "..." --body "..."
   ```
3. Post the Red Hat diff on the PR:
   ```bash
   ./scripts/redhat-diff.sh --pr <pr-number>
   ```

### Step 7: Summary

Report:
- Old tag → new tag
- Upgrade path validation result (steps and migrations required)
- PR link
- Number of Red Hat changes, stale files cleaned, diverged files reset
- Remind the user to review the PR and check CI before merging
- Note: this PR should be merged BEFORE the operator PR, since the operator
  depends on the SpiceDB image being built and available
