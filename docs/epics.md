# EPICs

This document tracks the major product EPICs for `notes.md`.

## Implementation Priority

The recommended implementation order is:

1. Markdown reading experience
2. Vault management
3. Workspace navigation and note management
4. Notes browsing and discovery
5. CLI quick actions
6. Writing workflow integration
7. Linking and knowledge navigation
8. Settings and keybindings
9. Daily notes and capture flows
10. Git sync

## Epic 1: Markdown Reading Experience

Make reading notes in the terminal feel intentional and useful from the start.

- render Markdown cleanly in the terminal
- support headings, lists, code blocks, links, and blockquotes
- improve layout for long notes
- make the reading UI feel strong enough to justify a notes-first product
- add theming support later if the renderer is stable

## Epic 2: Vault Management

Allow the user to connect the app to a Markdown vault and make the vault the primary source of truth.

- choose a vault folder path
- detect and suggest likely vault folders
- stop hardcoding the pages directory
- load notes dynamically from the configured vault
- persist the selected vault in app settings

## Epic 3: Workspace Navigation and Note Management

Make the app usable as a real workspace once a vault is connected.

- add a sidebar for notes, folders, or core navigation sections
- support tabs for open notes
- allow fast switching between open notes
- add basic note management such as create, rename, move, and delete
- consider pinned notes or pinned tabs later

## Epic 4: Notes Browsing and Discovery

Make it fast to find, open, and move between notes.

- list notes from the vault
- add fuzzy search
- support recent notes
- support note metadata such as tags if available
- add quick navigation between core views

## Epic 5: CLI Quick Actions

Allow users to search, inspect, and open notes without launching the full TUI.

- add CLI commands for quick search
- add CLI commands for quick open
- return useful results directly in the terminal
- make CLI actions fast enough for frequent shell usage
- reuse vault and note discovery logic from the main app
- consider quick capture commands later if the core commands are strong

## Epic 6: Writing Workflow Integration

Support writing without turning the project into a full editor.

- open notes in the user's preferred editor
- default to `$EDITOR` when available
- fall back to `nano` if no editor is configured
- create new notes from inside the app
- support quick note capture

## Epic 7: Linking and Knowledge Navigation

Support the workflows that make note systems more useful than plain files.

- detect wikilinks and regular Markdown links
- navigate from one note to another
- show backlinks
- support unresolved links gracefully

## Epic 8: Settings and Keybindings

Create a stable configuration surface for the app.

- settings screen or settings file
- configure vault path
- configure preferred editor
- configure sync behavior
- define default keybindings
- add custom keybinding support later

## Epic 9: Daily Notes and Capture Flows

Make the app useful for recurring note-taking habits.

- daily note creation
- open today's note directly
- quick capture into inbox or scratch note
- templates for common note types

## Epic 10: Git Sync

Add opt-in sync workflows for users who store their vault in Git.

- detect whether a vault is a git repository
- show sync status
- allow manual pull/push from the app
- handle conflicts conservatively
- consider automated sync only after manual sync is reliable
