# GitHub Organization Contributors Action

This action fetch contributors of a specified GitHub organization, and output to a file in `YAML`, `TOML` or `JSON` format.

## Inputs

| Parameter |  Type  | Default | Required | Description                                                                              |
| --------- | :----: | :-----: | :------: | ---------------------------------------------------------------------------------------- |
| `org`     | string |    -    |    Y     | The organization name.                                                                   |
| `output`  | string |    -    |    Y     | The output filename, i.e. `contributors.toml`, `contributors.yaml`, `contributors.json`. |

## Example

```yaml
on:
  workflow_dispatch:
  schedule:
    - cron: '0 */6 * * *'

jobs:
  contributors:
    permissions:
      contents: write
    runs-on: ubuntu-latest
    name: List contributors
    steps:
      - name: Checkout
        uses: actions/checkout@v3

      - name: Get contributors
        uses: razonyang/github-org-contributors-action@main
        id: contributors
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        with:
          org: 'hugomods'
          output: 'data/contributors.toml'

      - run: cat data/contributors.toml

      - uses: stefanzweifel/git-auto-commit-action@v4
        with:
          commit_message: "chore: update data/contributors.toml"
          file_pattern: 'data/contributors.toml'
```
