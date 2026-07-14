# Concept

## Motivation

This project started as a small helper inside another repository where I was experimenting with using OpenAI Codex as part of my everyday development workflow.

At first, I just wanted Codex to be easy to use. I wanted a development environment that already had the tools I needed, a curated view of the repository instead of unrestricted access, and a place to keep project-specific instructions, tasks, and configuration. As I kept using it, I added more pieces that made the workflow smoother, like initialization logic, reusable tasks, hooks, and a few guardrails.

The turning point came when I started another project and realized I was about to copy almost everything over. The runtime stayed the same, while only the project instructions, tasks, and workspace changed. That made it clear this had grown beyond being just another folder in a repository.

This project is simply giving that shared foundation its own home. Instead of rebuilding the same setup for every project, I want one place where I can improve it over time and reuse it wherever I need it.

## Core Idea

The core idea is to separate the infrastructure required to run Codex from the parts that belong to an individual project.

The reusable pieces, such as the container image, runtime initialization, permission handling, persistent state, and common tasks and hooks, should live in one place instead of being copied into every repository.

Each project should only describe its own context by defining the workspace, providing project-specific instructions, and extending the shared behavior where needed.

As I use Codex Devkit across more projects, I expect the shared foundation to grow naturally. New functionality should exist because multiple projects benefit from it, not because it might be useful someday.

## Design Goals

The goal is to provide a solid foundation that can be reused across projects without getting in the way. I'm keeping the scope small and focused on making Codex easier to use.

A few ideas should guide the project:

- Prefer sensible defaults over endless configuration
- Keep the runtime and the consuming project responsible for different things
- Don't add abstractions until multiple projects actually need them
- Choose simple solutions over clever ones
- Build from real experience, not anticipated use cases

Ultimately, Codex Devkit should reduce the effort required to bring Codex into a new project while remaining small enough to understand and maintain.

## Open Questions

There are several design decisions that I intentionally want to leave open until I gain more experience using it.

Some of the questions I'm exploring are:

- Where should the boundary between the runtime and the consuming project be?
- Which tasks and hooks should be shared by default?
- How should project-specific extensions integrate with the shared runtime?
- What belongs in the runtime, and what belongs in the project's `.codex` directory?

These are questions I'll revisit as the project evolves.
