# Source

A horizontally scalable Git server.

## Introduction

This project aims to fully decouple the underlying durable storage beneath
Git repositories from the Git protocol itself, to reduce both the
technical and economical costs of horizontally scaling a Git server.

## Getting started

A connection to a Postgres instance is required to run Source, so please
ensure you have a working database before continuing. Instructions on how
to setup a Postgres server is beyond the scope of this project, but for
brevity, the following should enable most test cases[^1].

```sh
docker run -e POSTGRES_HOST_AUTH_METHOD=trust -p 5433:5432 -d postgres
```

In the current absence of an install script, the project also requires
a working Go environment supporting version 1.25 or higher.

Clone the project, install the dependencies, and build the program.

```sh
# Clone
git clone https://github.com/iainjreid/source.git && cd source

# Install
go mod download

# Build
go build -tags "standalone"
```

With the project built, and a Postgres instance standing by you should be
able to start `source` without any further steps other than passing the
connection URI as a runtime argument.

```sh
source --db-uri "postgresql://postgres@localhost"
```

On start-up, the program will clone its own sourcecode into the database,
this will change in future updates when the web interface supports
importing.

## Troubleshooting

### SSH agent socket unavailable

```txt
SSH agent requested but SSH_AUTH_SOCK not-specified"
```

If you see an error like this, it's highly likely that either your SSH
agent it not running, or the socket is not available in your terminal
session.

In either scenario, the below snippet will solve the issue, but it might
not persist between terminal sessions.

```sh
eval $(ssh-agent)
```

## Development status

Currently, Source supports all of the required behaviour needed to
collaborate successfully via Git. Cloning, pushing, pulling, fetching, are
all operational with no noticeable performance bottlenecks at this stage.

The web interface is limited, but is on par performance-wise with
Sourcehut, although lacks functionality beyond simply browsing
a repository.

## Future development

* Issue management via Git notes, with the ability to discuss and provide
  feedback within the repository itself.

* In repository support for patch-by-patch reviews, a highly successful
  approach to peer reviewing changes used often within Git mailing lists.

* Support partial checkouts without requiring changes to how Git stores
  objects internally.

## Current limitations

* Synchronising objects from upstream remotes over SSH will require
  additional work to support knownhosts first.

* Git LFS is not yet supported. Given that LFS is an extension to Git, and
  not core behaviour, implementing this is not currently a high-priority
  feature.

  That said, there's nothing stopping Source from supporting LFS. Many
  databases offer a range of options to support large file storage,
  including PostgreSQL.

[^1]: Setting `POSTGRES_HOST_AUTH_METHOD=trust` disables password
    authentication and should only be used on personal or well secured
    machines. Official information about trust authentication can be
    [found here](https://www.postgresql.org/docs/18/auth-trust.html).

