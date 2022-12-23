# sink

Welcome to the sink repository! This is not intended to be public and is primarily a learning project (learning Go + CLI + GitHub APIs).

## Setup

1. Install Go
2. [Optional] Install https://taskfile.dev

You can run the application with `go run . [command] [options]`. You may need to run `go mod tidy` first to download the required packages.

You can `go build .` the application and place it in your `$PATH`. The `task build` command will do this, placing in the `~/.local/bin` directory (and assuming that's already in your path).

## FAQs

### Why?

This is primarily a learning project. We should have a more formal discussion if we want to take a dependency on this tool for managing `dbt-labs` repositories.

It's intended to perform GitHub tasks that are tedious to do manually but common for dbt Labs maintainers -- currently, copying an issue across multiple repos (e.g. "Support Python 3.11" in `dbt-core` and each adapter repo) and syncing label definitions (color and description) across all repositories owned by the Core team. These can be performed with:

```bash
sink issue 6147 \
    --org dbt-labs \
    --source-repo dbt-core \
    --target-repos dbt-bigquery,dbt-snowflake,dbt-redshift,dbt-spark
```

and

```bash
sink labels \
    --org dbt-labs \
    --source-repo dbt-core \
    --team Core \
    --exclude-team-repos core-team,"schemas.getdbt.com" \
    --extra-repos dbt-server
```

**Important**: `--dry-run=true` by default to avoid accidents. Commands will print out what they would have done, and you can subsequently run with `--dry-run=false` to perform the actions.

Note that an existing object is required -- this tool is not intended for creating those objects, but only syncing them across multiple repositories.

While this was made to work in a single organization, it could be easily refactored to allow operations across organizations (at least for public repositories), so we could also copying to `dbt-databricks` and other repos outside of our org.

You can use a `config.yaml` (or other extensions) to set defaults for flags:

```yaml
org: dbt-labs
team: Core
source-repo: dbt-core
target-repos:
    - dbt-server
```

### How?

- Go
    - `cobra` and `viper` for CLI
    - `go-gh` for GitHub operations

### Why not a `gh` CLI extension?

I think that'd be cool -- something like `gh issue sync/copy [options]`.

As I'm writing this I'm learning there's already `gh label clone <source-repository>`, rip. There's no `gh issue clone` or operations on multiple repos.

### Is this maintained?

Sure!

### Is future work planned?

Kinda!

