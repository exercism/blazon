# Blazon

Internal tool to let track maintainers create duplicate issues in all relevant language track repositories.

## Installation

If you have Go:

    go get -u github.com/exercism/blazon

If not, grab the [most recent release](https://github.com/kytrinyx/blazon/releases/latest).

## Configuration

Create a personal api token, **giving it the scope "public_repo"**.
Save it to an environment variable named `BLAZON_GITHUB_API_TOKEN`.

## Usage

Run `blazon` to see the available flags.

## Tips

Before submitting the blazon issue, make sure that there is an issue that
summarizes the problem/solution. If it doesn't already exist, then create one
(it might fit into the https://github.com/exercism/x-common repository, or the
https://github.com/exercism/todo repository).

Then include a link to the summary issue in the body of the text for the
blazon issue.

This will ensure that the summary issue contains references to each of the
individual issues in each of the track repositories. Each of these links will
include an icon whose color indicates the status of the issue (open|closed).

This ensures that there is a handy, visible "todo" list that immediately shows
the overall status of the effort, and once all of the related links are red,
the summary issue can be closed.
