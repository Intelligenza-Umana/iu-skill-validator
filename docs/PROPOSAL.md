# Project proposal: AI-agent skill pre-release validator

## Problem and affected users
Skill authors and maintainers need to distinguish well-formed skill packages from skills that actually behave as expected in an agent host.

## Evidence and examples
The [project registry](https://github.com/Intelligenza-Umana/intelligenza-umana) identifies workflows with unverified host behavior. This is a verification gap, not evidence of broken installations. Use valid and deliberately invalid local fixtures to demonstrate the static checks.

## Existing solutions and gap
Manual reviews and repository-specific scripts are useful, but static validation and runtime evidence may be conflated or left undocumented.

## Proposed approach and alternatives
Build a small pre-release checker for package structure and local references, and document a separate manual host verification procedure. Alternatives include continuing manual review or maintaining independent scripts per repository.

## Scope and non-goals
- First skill: [Project Lifecycle Skill](https://github.com/Intelligenza-Umana/iu-project-lifecycle-skill).
- One skill format, one repository and one named host/version initially.
- Choose the host and runtime verification procedure together.
- No multi-host support, automatic behavioral guarantees or production integrations in the first prototype.

## Security, privacy and integrity
Use disposable synthetic fixtures and a controlled runtime workspace. Do not access secrets, execute untrusted fixture instructions outside that workspace, or modify unrelated host state. Report limitations and unexpected effects.

## Minimal prototype and acceptance criteria
- Check the agreed package structure and local references.
- Demonstrate reproducible tests with valid and invalid fixtures and expected results.
- Manually install and exercise the skill on one named host/version; record environment, skill revision, procedure, expected and observed behavior.
- Keep static pass/fail, runtime pass/fail and untested behavior separate.
- Document commands, evidence and limitations. A static pass must never be presented as proof of runtime correctness.

## Contributions and maintenance
@developersnufi proposed the idea and volunteered to lead initial implementation. @fabiozagaria will support design, implementation, review and testing. Keep maintenance small; discuss further checks and host profiles after the prototype.

## Status
Initial static CLI prototype implemented. No runtime verification is claimed yet.

Original discussion: https://github.com/Intelligenza-Umana/intelligenza-umana/issues/5
