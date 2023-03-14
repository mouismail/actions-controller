## GitHub App permissions

### Repository permissions

- **Actions** Workflows, workflow runs and artifacts. `read & write`
- **Administration** Repository creation, deletion, settings, teams, and collaborators. `read & write`
- **Contents** Repository contents, commits, branches, downloads, releases, and merges. `read`
- **Issues** Issues and related comments, assignees, labels, and milestones. `read & write`
- **Metadata** (mandatory) Search repositories, list collaborators, and access repository metadata. `read`
- **Workflows** Update GitHub Action workflow files. `read & write`

### Organization permissions

- **Administration** Manage access to an organization. `read & write`
- **Members** Organization members and teams. `read`
- **Self-hosted runners** View and manage Actions self-hosted runners available to an organization. `read & write`

### User permissions

- **Email addresses** Manage a user's email addresses. `read`

## GitHub App events

### Subscribe to events

- **Meta** When this App is deleted and the associated hook is removed.
- **Workflow dispatch** A manual workflow run is requested.
- **Workflow job** Workflow job queued, requested or completed on a repository.
- **Workflow run** Workflow run requested or completed on a repository.

## Additional information

### Where can this GitHub App be installed?

- **Any account** Allow this GitHub App to be installed by any user or organization.
