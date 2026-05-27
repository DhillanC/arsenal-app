# Examples

Use these examples as a quality bar for issues and pull requests created from this template.

## User Story

**Title:** `[Story]: Track work on a project board`

**As a:** project contributor

**I need:** a clear project board workflow

**So that:** I can understand what is ready, in progress, in review, and done

**Details and assumptions:**

- The team uses GitHub Issues and pull requests.
- The board has Backlog, Ready, In Progress, In Review, and Done columns.
- Pull requests are linked to issues.

**Acceptance criteria:**

```gherkin
Given a contributor opens the project board
When they look at an issue card
Then they can understand the current status of the work
```

## Bug Report

**Title:** `[Bug]: Issue has no acceptance criteria`

**Summary:** A newly created issue can be submitted without enough detail to review or implement.

**Steps to reproduce:**

1. Open a new issue.
2. Leave acceptance criteria blank.
3. Submit the issue.

**Expected behavior:** The template requires acceptance criteria before submission.

**Actual behavior:** The issue can be created without completion checks.

## Technical Task

**Title:** `[Task]: Add Markdown validation workflow`

**Task:** Add a pull request check that validates Markdown files.

**Why it matters:** Documentation is part of the template, so broken Markdown should be caught before merge.

**Done when:**

- [ ] Markdown validation runs on pull requests.
- [ ] The workflow checks all Markdown files.
- [ ] The README mentions the check.

## Pull Request

### Summary

Adds a project board guide with recommended columns, ready and done definitions, and working rules.

### Linked Issue

Closes #12

### Checklist

- [x] The change is focused and easy to review.
- [x] Acceptance criteria are met.
- [x] Documentation was updated.
- [x] Manual checks were completed.

### Notes for Reviewers

Review the column names and done criteria to confirm they match the team's workflow.
