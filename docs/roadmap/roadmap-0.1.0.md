# Roadmap — v0.1.0

This roadmap defines the implementation plan for **Mipe 0.1.0**.

The objective of this release is not to deliver the complete long-term vision, but to establish a working foundation while exploring how Codex behaves in practice. Each milestone is intentionally small and answers one architectural question before introducing additional complexity.

## M1 — Bootstrap

**Goal**

Replace the current shell-based runtime with a dedicated Go bootstrap application.

**Questions**

- Can Go fully replace the current initialization scripts?
- What is the minimum responsibility of the bootstrap?
- Which responsibilities belong elsewhere?

**Success Criteria**

- Bootstrap initializes development environment
- Codex launches successfully

## M2 — Session Detection

**Goal**

Determine how Mipe identifies the lifecycle of a development session.

**Questions**

- What defines the beginning of a session?
- What defines the end of a session?
- Which source should be considered authoritative?
- Can sessions be detected without explicit registration?

**Success Criteria**

- Active sessions are detected
- Completed sessions are detected

## M3 — Transcript

**Goal**

Understand how Codex stores conversations and expose them through the backend.

**Questions**

- Where are transcripts stored?
- Can they be streamed while a session is running?
- How should they be persisted?

**Success Criteria**

- Session transcripts are captured
- Transcript is available through the backend

## M4 — Usage

**Goal**

Capture session usage metrics.

**Questions**

- Where are token statistics available?
- How frequently are they updated?
- How should they be associated with sessions?

**Success Criteria**

- Input and output tokens are captured
- Usage is associated with sessions

## M5 — Knowledge Extraction

**Goal**

Generate useful information from completed sessions.

**Questions**

- What prompt produces useful summaries?
- What constitutes a useful note?
- When should extraction occur?

**Success Criteria**

- Session summaries are generated
- Session notes are generated

## M6 — Developer Journal

**Goal**

Provide a simple interface for exploring current and previous development sessions.

**Questions**

- What information is most useful during an active session?
- What information is most valuable after a session ends?

**Success Criteria**

- Active session is displayed
- Session history is displayed
- Session details are displayed
