## GitHub App permissions


300k repos

any action event (dispatch, run, job) => repo (trigger)

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

## Prerequisites

- Go version 1.16 or later
- Docker version 20.10 or later
- Git version 2.30 or later

## Usage

### Run locally

The app will listen on port 3000 by default. You can change it by setting the `HTTP_PORT` environment variable.

To test the app, open a browser and visit http://localhost:3000.

```bash
$ make dev
```

### Run in Docker

To build the app binary with Docker and `ldflags`, use the docker build command with a tag name for your image and some arguments for ldflags:

```bash
$ make build
```

This will create a Docker image named `actions-controller:latest` with some arguments for ldflags.

To run the app in a container, use the docker run command with the tag name you used to build the image:

```bash
$ make start
```

The app will listen on port 3000 inside the container and map it to port 3000 on your host machine.

To test the app, open a browser and visit http://localhost:3000.

To see the app version information, visit http://localhost:3000/version.


