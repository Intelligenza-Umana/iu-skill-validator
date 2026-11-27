# IU Skill Validator

An open-source pre-release validation tool for AI-agent skills, separating static package checks from runtime verification.

**Status:** Initial design. No released implementation or verified host compatibility yet.

## Goal
Help skill authors check package structure and local references, then record separate evidence of actual agent-host behavior.

## First milestone
Start with [Project Lifecycle Skill](https://github.com/Intelligenza-Umana/iu-project-lifecycle-skill), one skill format and one named host/version. Test valid and invalid fixtures, then document a manual runtime verification.

- [Project proposal](docs/PROPOSAL.md)
- [Contributing](CONTRIBUTING.md)
- [Original discussion](https://github.com/Intelligenza-Umana/intelligenza-umana/issues/5)

The implementation language, architecture and first host are open for joint design.

## License
MIT. Part of [Intelligenza Umana](https://github.com/Intelligenza-Umana).
