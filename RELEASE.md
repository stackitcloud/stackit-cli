# Release

## Release cycle

A release should be created at least every 2 weeks. 

## Release creation

> [!IMPORTANT]
> Consider informing / syncing with the team before creating a new release.

1. Check out latest main branch on your machine
2. Create git tag: `git tag vX.X.X`
3. Push the git tag: `git push origin --tags`
4. The [release pipeline](https://github.com/stackitcloud/stackit-cli/actions/workflows/release.yaml) will build the release and publish it on GitHub
5. Ensure the release was created properly using the [releases page](https://github.com/stackitcloud/stackit-cli/releases)

