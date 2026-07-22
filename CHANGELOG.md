# Changelog

All notable changes to Mipe will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/2.0.0/),
and Mipe adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [runtime-0.1.0] - 2026-07-21

### Added

- A bootstrap CLI that validates configuration and starts agents in a ready-to-use workspace
- Shared agent configuration installed safely before switching to local developer user
- Optional project setup via `<workspace>/.mipe/init/setup.sh`
- Mipe runtime images for dedicated agent environments:
  - Codex — Codex 0.144.6
  - Claude — Claude Code 2.1.216
  - Codex Java — Codex 0.144.6 and Temurin 21
  - Claude Java — Claude Code 2.1.216 and Temurin 21
  - Codex Web — Codex 0.144.6, Chromium, and Playwright MCP
  - Claude Web — Claude Code 2.1.216, Chromium, and Playwright MCP
