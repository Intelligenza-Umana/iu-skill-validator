# IU Skill Validator

An open-source pre-release validation tool for AI-agent skills, separating static package checks from runtime verification.

**Status:** Static prototype implemented; no release or verified host compatibility yet.

## Goal
Help skill authors check package structure and local references, then record separate evidence of actual agent-host behavior.

## First milestone
Start with [Project Lifecycle Skill](https://github.com/Intelligenza-Umana/iu-project-lifecycle-skill), one skill format and one named host/version. Test valid and invalid fixtures, then document a manual runtime verification.

- [Project proposal](docs/PROPOSAL.md)
- [Contributing](CONTRIBUTING.md)
- [Original discussion](https://github.com/Intelligenza-Umana/intelligenza-umana/issues/5)

The first host remains open for joint design.

## Static checks

Run the static validator with:

```sh
go run . check ./fixtures/valid
go run . check --json ./fixtures/valid
```

`SKILL.md` must begin with YAML frontmatter containing non-empty top-level `name` and `description` fields. The validator checks local references, rejects absolute paths and paths escaping the package through traversal or symlinks, and reports unreferenced non-hidden regular files (except `SKILL.md`) as `WARN`. Warnings do not change the validation result or exit code.

A `PASS` is evidence only of these static checks; it does not mean the skill was installed or tested in an agent host.

## Manual runtime report

Copy [the runtime report template](docs/runtime-report.template.json) to a separate file and fill it in only after manually testing the skill on a named host/version. Set `status` to `pass`, `fail`, `blocked`, or leave it as `not_tested`. Record the host environment, skill revision, procedure, expected behavior, and observed behavior. The validator does not create or verify this report.

## License
MIT. Part of [Intelligenza Umana](https://github.com/Intelligenza-Umana).
