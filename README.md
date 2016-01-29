# Blazon

Internal tool to let track maintainers create duplicate issues in all relevant language track repositories.

## Installation

If you have Go:

    go get -u github.com/exercism/blazon

If not, grab the [most recent release](https://github.com/kytrinyx/blazon/releases/latest).

## Configuration

Create a personal api token, **giving it the scope "public_repo"**.
Save it to an environment variable named `BLAZON_GITHUB_API_TOKEN`.
