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

## Intelligenza Umana community

**Human intelligence. Artificial assistance.**

[Organization](https://github.com/Intelligenza-Umana) · [Website](https://intelligenza-umana-web.vercel.app/) · [Discord](https://discord.gg/Zt85ZYEXAA) · [Onboarding](https://github.com/Intelligenza-Umana/intelligenza-umana/blob/main/community/ONBOARDING.md) · [Good first issues](https://github.com/search?q=org%3AIntelligenza-Umana+is%3Apublic+is%3Aissue+is%3Aopen+no%3Aassignee+label%3A%22good+first+issue%22&type=issues)

Contributions to public projects are welcome from everyone, regardless of experience. Organization membership and Discord are optional. Discuss scope in an issue before starting work, then submit a focused pull request with verification evidence and limitations. Use English for shared GitHub discussions.

To become an organization member, [submit a membership request](https://github.com/Intelligenza-Umana/intelligenza-umana/issues/new?template=membership_request.md). The Discord numeric user ID is optional; the issue is public, so you may provide it privately instead. A maintainer reviews requests manually and sends an official GitHub invitation if approved. Sentinel is configured to assign the Contributor role on Discord entry; the Member role is assigned manually only after identity and accepted GitHub membership are verified. Neither membership nor a Discord role automatically grants repository write access.

Follow the [code of conduct](https://github.com/Intelligenza-Umana/intelligenza-umana/blob/main/community/CODE_OF_CONDUCT.md) and [private security reporting policy](https://github.com/Intelligenza-Umana/intelligenza-umana/blob/main/SECURITY.md). Each project retains its own license, maturity and verification requirements.

## License
MIT. Part of [Intelligenza Umana](https://github.com/Intelligenza-Umana).
