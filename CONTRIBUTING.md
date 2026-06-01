# Contributing

Use this guide when opening issues, changing code, or preparing pull requests for Arsenal App.

## Workflow

1. Create or pick up a GitHub issue before starting work.
2. Use the matching issue template for stories, bugs, features, epics, research spikes, or technical tasks.
3. Keep each branch focused on one outcome.
4. Link pull requests to issues with `Closes #issue-number`.
5. Document tests, manual checks, screenshots, or follow-up risks in the pull request.

## Local Checks

Run the checks that match the change:

```bash
go test ./...
npx markdownlint-cli2@latest --config .markdownlint.json "**/*.md"
ruby -e 'require "yaml"; ARGV.each { |file| YAML.load_file(file) }' .github/ISSUE_TEMPLATE/*.yml
```

For UI or template changes, also run the app locally and verify the changed view:

```bash
make dev
```

## Documentation

- Update `README.md` when setup or usage changes.
- Update `docs/TASKS.md` or the linked issue when scope changes.
- Follow `docs/MERMAID.md` for Mermaid diagrams.
- Keep project-board conventions aligned with `docs/PROJECT_BOARD.md`.

## Pull Requests

Before requesting review:

- The change is small enough to review.
- Acceptance criteria are met.
- Tests or manual checks are listed.
- Documentation is updated if needed.
- Security or privacy impact has been considered.
