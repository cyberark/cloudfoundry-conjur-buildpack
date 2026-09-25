---
name: ai-pr-review
description: Automated audit of code changes before opening a Pull Request against this repo (cloudfoundry-conjur-buildpack, the CyberArk Secrets Manager buildpack for Cloud Foundry, bash + Go). Analyzes the diff against the base branch, checks compliance with CyberArk/Conjur conventions used here (architecture, security, quality), and generates a PR_REVIEW_REPORT.md ready to paste into the PR description. Use before every commit that finalizes a task, or when the user writes "@ai-pr-review run analysis and generate report".
---

# AI PR Review

Automated audit of code changes before opening a Pull Request. The skill analyzes the diff, checks compliance with the CyberArk/Conjur conventions used in this repo (`cloudfoundry-conjur-buildpack`), and generates a report ready to paste into the PR description.

## Workflow

### 1. Context Extraction

- Run `git diff` against the base branch (usually `main`).
- Identify the changed components:
  - `bin/supply` (+ `.bat`/`.ps1` variants) — the CF supply-script entrypoint
  - `lib/0001_retrieve-secrets.sh` (+ `.bat`) — the `profile.d` script that fetches and exports secrets at app start
  - `lib/install_go.sh` — Go toolchain install-at-supply-time logic
  - `conjur-env/` — the Go binary (`main.go`, `authn.go`, `authn_api_key.go`, `authn_cert.go`)
  - `manifest.yml` — buildpack-packager metadata (Go versions/SHAs per `cf_stack`)
  - `ci/`, `tests/` — CI entry points and Cucumber/shell test suites

### 2. Multi-Persona Analysis

- **Architect**:
  - This is a **supply-only** buildpack — it must never assume it runs standalone; changes should keep it stackable with a language buildpack and preserve `bin/compile` as a legacy no-op.
  - `lib/0001_retrieve-secrets.sh` uses an `__BUILDPACK_INDEX__` placeholder replaced via `sed` in `bin/supply` — flag any change that hardcodes an index path instead of using the placeholder.
  - Both auth modes (API key and `authn-cert` mTLS/SPIFFE) must stay supported end-to-end, selected via the `authn_type` credential field in `VCAP_SERVICES` — flag changes that special-case one path and silently break the other.
- **Security**:
  - Secret values must never be echoed or logged — `lib/0001_retrieve-secrets.sh` explicitly disables `xtrace` and sanitizes error messages via `export_err`/`conjur_env_err`; flag any new logging/debug statement that could leak a fetched secret value.
  - Environment variable names derived from `secrets.yml` must remain valid shell identifiers before being `export`ed.
  - No hardcoded Conjur credentials, API keys, or certs in `conjur-env/` or test fixtures.
- **Quality**:
  - Go changes in `conjur-env/` need corresponding tests (`main_test.go`, `authn_cert_test.go` style) and should pass `./ci/test_conjur-env` (`go vet` + `go test -coverprofile` + lint).
  - Shell changes to `lib/0001_retrieve-secrets.sh` need corresponding coverage in `./tests/retrieve-secrets/start`.
  - Behavior-affecting changes should be checked against `./ci/test_integration` (Cucumber, tag `not @integration`) expectations, even if a full run requires the local CF+Conjur dev stack.

### 3. Standard Compliance

- When bumping the Go version for `conjur-env`, all four locations must move together: `conjur-env/Dockerfile`, `conjur-env/go.mod`, `manifest.yml`, and `lib/install_go.sh` — flag a version change in only one of these.
- `CONJUR_BUILDPACK_BYPASS_SERVICE_CHECK=true` is CI-only (skips the `VCAP_SERVICES` binding check) — flag if it leaks into a non-test code path.
- `manifest.yml`'s `include_files` list drives what actually ships in the packaged buildpack zip — flag new files that should ship but aren't added there.
- `bin/compile` is a deprecated no-op kept only for legacy Heroku/CF compatibility — don't add new logic to it.

### 4. Report Generation

Generate `PR_REVIEW_REPORT.md` **fully filled in** — never leave a `<placeholder>` in the final file. Derive each field automatically:

- **One-line summary**: one imperative-mood sentence describing what the diff does.
- **Desired Outcome**: if the branch name or commit messages reference a ticket/issue, pull context from it (e.g. `gh issue view <n>` for a `#<n>` reference); otherwise infer the motivating problem directly from what the diff changes.
- **Implemented Changes**: build this directly from the actual `git diff` and the Multi-Persona Analysis findings (step 2) — what changed, why, and how a reviewer should verify it (manual/integration steps if any). State screenshots as "N/A" unless the diff touches UI.
- **Connected Issue/Story**: scan the current branch name and the commit messages on this branch for a ticket pattern like `[A-Z]+-\d+` (e.g. `CP-26179`, `CNJR-1234`) and use the first match as the "CyberArk internal issue ID". If a GitHub issue number is referenced (`#123` in a branch name or commit message), use it for "Resolves #123" — otherwise omit that line entirely rather than leaving a placeholder.
- **Definition of Done** — tick exactly one box per subsection based on what the diff actually contains, never by default:
  - *Changelog*: tick "CHANGELOG has been updated" only if `CHANGELOG.md` (or equivalent) is part of the diff; tick "does not require a CHANGELOG update" only if the change has no user-facing effect.
  - *Test coverage*: tick "includes new tests" only if test files changed in the diff; tick "does not require tests" only if the change is non-functional (docs, config, tooling, skills).
  - *Documentation*: tick "docs were updated" only if README/docs files changed; tick "does not require updating documentation" only if nothing user-facing changed.
  - *Behavior*: tick "no behavior was changed" only if the diff is limited to non-runtime files (skills, docs, CI config, tests); otherwise tick one of the other two options and explain why.
  - *Security*: tick "no security aspects" only if the Security persona (step 2) genuinely found nothing; otherwise tick the appropriate other option and reference the finding.

If a field genuinely cannot be determined, write a short explicit note (e.g. "No linked issue found in branch/commit history") instead of leaving an angle-bracket placeholder.

The output must follow this structure when pasted into a PR description:

```markdown
<filled summary>

### Desired Outcome

<filled>

### Implemented Changes

- <filled bullets>

### Connected Issue/Story

Resolves #<n>  (omit this line if none found)

CyberArk internal issue ID: <filled or note>

### Definition of Done

#### Changelog
- [ ] The CHANGELOG has been updated, or
- [ ] This PR does not include user-facing changes and doesn't require a CHANGELOG update

#### Test coverage
- [ ] This PR includes new unit and integration tests to go with the code changes, or
- [ ] The changes in this PR do not require tests

#### Documentation
- [ ] Docs (e.g. READMEs) were updated in this PR
- [ ] A follow-up issue to update official docs has been filed here: [insert issue ID]
- [ ] This PR does not require updating any documentation

#### Behavior
- [ ] This PR changes product behavior and has been reviewed by a PO, or
- [ ] These changes are part of a larger initiative that will be reviewed later, or
- [ ] No behavior was changed with this PR

#### Security
- [ ] Security architect has reviewed the changes in this PR,
- [ ] These changes are part of a larger initiative with a separate security review, or
- [ ] There are no security aspects to these changes
```

(the checkboxes above are the template shape — in the actual generated report, exactly one box per subsection must be ticked `[x]`, the rest `[ ]`)

## Usage Instructions

Run this command before every commit that finalizes a task:

```
@ai-pr-review run analysis and generate report
```

## Port.io Metadata

- **Maturity Level**: Tier 2 (Repeatable Skill)
- **Service**: Secrets Manager Integrations
