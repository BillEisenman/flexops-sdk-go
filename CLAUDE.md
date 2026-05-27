# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`github.com/BillEisenman/flexops-sdk-go` is the **official hand-crafted Go SDK** for the FlexOps Platform. It targets the FlexOps **Gateway BFF**. Module path `github.com/BillEisenman/flexops-sdk-go`; Go 1.22+. Current tag `v1.0.1`.

> **Gateway-targeted, not VSCS-targeted.** All hand-crafted SDKs in this family (.NET, Node, Python, Go, PHP, Ruby) hit Gateway. The Java SDK is the lone exception — it was auto-generated against VisionSuiteCoreServices and is archived as of 2026-03-08.

## Build & Run Commands

```bash
go build ./...                  # Build
go test ./...                   # Run tests
go vet ./...                    # Static analysis
go fmt ./...                    # Format
```

## Architecture

```text
Customer Go app  →  flexops-sdk-go (this repo)  →  Gateway BFF (gateway.flexops.io)
                                                   ↓
                                                   VSCS / Integrations / etc.
```

## Key files

| Path | Purpose |
|---|---|
| `flexops.go` | Primary client surface |
| `flexops_test.go` | Top-level tests |
| `auth.go` | Authentication flow / token handling |
| `errors.go` | Typed error surface (carries Gateway error envelope) |
| `carriers.go` + `carriers_dhl.go` / `carriers_fedex.go` / `carriers_ups.go` | Per-carrier convenience surfaces |
| `analytics.go`, `api_keys.go`, `email_templates.go` | Feature-area clients |
| `go.mod` | Module + Go version (1.22) |

## Conventions

- **Idiomatic Go**: errors are returned values typed via `errors.go` — don't panic.
- One file per resource family (carriers, api_keys, analytics, etc.) keeps the package navigable.
- Public API is the package surface — keep internal helpers unexported.

## Publish

Go modules publish on tag push to GitHub — there is no separate registry workflow. Bump in code, tag (`vX.Y.Z`), push the tag, and the proxy picks it up on the next `go get`.

## Related Repositories

| Repository | Purpose |
|---|---|
| **This repo** | `github.com/BillEisenman/flexops-sdk-go` |
| FlexOps Gateway | The HTTP API this SDK calls — `BillEisenman/FlexOpsGateway` |
| Sibling SDKs | `FlexOps.Sdk` (.NET), `@flexops/sdk` (Node), `flexops` (Python/Ruby), `flexops/sdk` (PHP) |
| FlexOps Developer Docs | Hosts the SDK page — `BillEisenman/FlexOpsDeveloperDocs` |
