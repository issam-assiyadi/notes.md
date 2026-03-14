# notes.md

`notes.md` is a terminal-native notes application for Markdown vaults.

## Core Idea

The goal is not to build another text editor. The goal is to build a notes-first app that lives in the terminal.

Neovim already solves editing very well. Competing with Neovim as an editor would be a bad strategy and would push this project into rebuilding mature editor features with a worse user experience.

This project should instead focus on the problem that notes apps solve:

- living inside a vault, not inside a file tree
- browsing notes quickly
- searching and opening notes with minimal friction
- navigating links between notes
- rendering Markdown in a readable terminal UI
- supporting capture workflows such as daily notes and quick notes

Editing is still part of the workflow, but it is not the core differentiator. The app can open the user's preferred editor when they want to write, while `notes.md` remains responsible for note discovery, navigation, rendering, and vault-oriented workflows.

In short:

- this is a notes-first app
- this is not an editor-first app
- this is not trying to replace Neovim
- this is trying to make Markdown vaults feel native in the terminal

## Product Direction

The product should feel closer to "Obsidian for the terminal" than "a lightweight terminal editor".

That means the foundation of the app should be:

- user-selected vaults instead of hardcoded content folders
- strong read, browse, search, and navigation flows
- clean terminal Markdown rendering
- external editor integration for writing
- settings for vault path, editor preference, sync behavior, and keybindings

Git-based sync can be valuable, but it should be treated as a workflow feature, not the foundation of the product. A safe first version should prefer explicit sync actions before moving to automatic pull/push behavior.

## Roadmap

The project EPICs live in [docs/epics.md](/Users/iassiyadi/Side-projects/notes.md/docs/epics.md).
