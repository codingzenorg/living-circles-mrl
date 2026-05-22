# Issue / PR Linking Guidance (MRL-Managed Repositories)

## Purpose

Define the local resolution shape when a tracked issue, campaign task, or finding is answered by repository changes.

## Core Rule

When repository changes are required:

- open a pull request in this repository
- link the originating issue, campaign issue, or finding issue in the PR body
- link the PR back from the originating issue, campaign record, or finding record
- keep the issue or record as the tracking artifact and the PR as the implementation artifact

## Scope

This guidance applies to repositories managed under the MRL contract.

It does not require a PR for outcomes that are resolved entirely by documentation, discussion, or decision without repository changes.

## Notes

- Prefer explicit bidirectional linkage over relying on implied closure alone.
- If the repository uses GitHub closing keywords, they are additive, not a substitute for the explicit link back to the PR.
