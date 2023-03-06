# Contributing

## Updating Provider Documentation

The documentation is generated using [tfplugindocs](https://github.com/hashicorp/terraform-plugin-docs).

`docs/` should never be edited by hand. Instead, update the documentation via updating `Description` fields within the `schema` blocks of the provider's resources and data sources. And if needed, update the templates in `templates/`. Finally, you can run the following command to re-generate the documentation:

```bash
make docs
```

## Submitting a PR

1. Fork the repo
2. Create a PR from your fork against the `master` branch
3. Add labels to your PR (see: [Labels](.github/release-drafter.yaml))

### PR Requirements

1. Ensure that all tests pass (via commenting `/test` if you are an admin or a contributor with write access on this repo, otherwise wait for a maintainer to submit the comment. The comment will be ignored if you are not an admin or a contributor with write access on this repo. See: https://codefresh.io/docs/docs/pipelines/triggers/git-triggers/#support-for-building-pull-requests-from-forks)
2. Ensure that `make docs` has been run and the changes have been committed.
3. Ensure that `make fmt` has been run and the changes have been committed.
