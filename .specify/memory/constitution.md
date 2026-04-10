<!--
  Sync Impact Report
  ==================
  Version change: (none — initial template) → 1.0.0
  Modified principles: N/A (initial ratification)
  Added sections:
    - Core Principles (5): Test-First Development, Multi-Layer Resilience,
      Resource API Consistency, Platform-Aware Performance, Living Documentation
    - Windows Compatibility & Platform Constraints
    - Development Workflow & Quality Gates
    - Governance
  Removed sections: None
  Templates requiring updates:
    - .specify/templates/plan-template.md ✅ no update needed (generic gates)
    - .specify/templates/spec-template.md ✅ no update needed (generic structure)
    - .specify/templates/tasks-template.md ✅ no update needed (TDD aligned)
    - .specify/templates/checklist-template.md ✅ no update needed (generic)
    - .specify/templates/agent-file-template.md ✅ no update needed (generic)
  Follow-up TODOs: None
-->

# Pulumi HyperV Provider Constitution

## Core Principles

### I. Test-First Development (NON-NEGOTIABLE)

- TDD is mandatory: tests MUST be written and failing before
  implementation code.
- Red-Green-Refactor cycle is strictly enforced for all resource CRUD
  operations.
- Resources MUST be testable without Hyper-V infrastructure by defining
  mockable interfaces for WMI services.
- Coverage targets: 80%+ for provider logic, 100% for error handling
  and fallback paths.
- Integration tests run on real Windows/Hyper-V in CI; unit tests MUST
  run cross-platform (Linux and Windows).
- Example programs in `examples/` serve as acceptance tests — every
  resource MUST have a working example.

### II. Multi-Layer Resilience

- All Hyper-V operations MUST follow the three-tier fallback chain:
  WMI High-Level methods -> WMI Direct API calls -> PowerShell Cmdlets.
- Each fallback transition MUST be logged with severity level and
  reason for fallback.
- Panic recovery MUST wrap all WMI operations — no unrecovered panics
  reach users.
- Errors MUST include: operation context, which fallback layer failed,
  and actionable user guidance.
- Fallback decisions MUST be deterministic and documented — never
  silent degradation.
- The provider MUST handle graceful degradation across Windows 10/11,
  Windows Server, and Azure Edition environments.

### III. Resource API Consistency

- All 4 resources (Machine, VhdFile, VirtualSwitch, NetworkAdapter)
  MUST follow identical patterns for naming, validation, defaults, and
  error reporting.
- Input property naming uses camelCase consistently across all
  resources with no drift (e.g., no `vmName` vs `MachineName`
  inconsistency).
- Every optional field MUST document its default value in the schema.
- Input validation MUST happen at the provider layer before any WMI or
  PowerShell call is made.
- Schema changes MUST be backward-compatible within a MINOR version;
  breaking changes require a MAJOR version bump and migration guide.

### IV. Platform-Aware Performance

- WMI connections MUST be reusable within a single Pulumi operation —
  no fresh connection per CRUD call.
- Every WMI and PowerShell operation MUST have a configurable timeout
  with sensible defaults.
- Windows edition detection (10/11 client, Server 2016-2025, Azure
  Edition) MUST happen once at provider initialization and inform all
  subsequent operation behavior.
- Performance-critical paths MUST have benchmark tests with documented
  baseline expectations.
- Connection pooling strategy MUST be documented and testable.

### V. Living Documentation

- Documentation MUST live alongside code — embedded markdown per
  resource, auto-synced to generated SDKs.
- Every resource property MUST have a description in the schema with
  no undocumented fields allowed.
- README, CLAUDE.md, and resource docs MUST be updated in the same PR
  as code changes — no separate "docs follow-up" PRs.
- API reference documentation MUST be auto-generated from the provider
  schema.
- Manual documentation is reserved for guides, troubleshooting, and
  architectural decisions only.
- Breaking changes MUST have migration guides published before release.

## Windows Compatibility & Platform Constraints

- Supported platforms: Windows 10 (build 1809+), Windows 11,
  Windows Server 2016/2019/2022/2025, Azure Edition.
- Provider MUST initialize successfully on all supported platforms with
  graceful degradation rather than hard failure.
- OS-specific behavior differences MUST be documented in a maintained
  compatibility matrix.
- Both PowerShell Core (`pwsh`/`pwsh.exe`) and Windows PowerShell
  (`powershell.exe`) MUST be supported via `util.FindPowerShellExe()`.
- All file paths MUST use platform-agnostic utilities.
- Hyper-V detection at startup MUST provide OS-specific guidance
  messages when services are unavailable.

## Development Workflow & Quality Gates

- All PRs MUST pass: `make lint`, `make test`, `make format` with zero
  warnings.
- Code review is required — no self-merges to master branch.
- SDK regeneration via `make codegen` is required whenever schema
  changes are made.
- golangci-lint configuration (`.golangci.yml`) is authoritative and
  MUST NOT be bypassed or weakened.
- Commit messages MUST follow conventional commit format.
- CHANGELOG.md MUST be updated for every user-facing change following
  Keep a Changelog format.
- Go code follows standard library first, external dependencies second,
  project imports last ordering.

## Governance

- Constitution supersedes all ad-hoc development practices.
- Amendments require a PR with review and approval.
- Use CLAUDE.md for runtime development guidance; constitution defines
  principles, CLAUDE.md defines procedures.
- Compliance with constitution principles MUST be verified during code
  review.

**Version**: 1.0.0 | **Ratified**: 2026-04-10 | **Last Amended**: 2026-04-10
