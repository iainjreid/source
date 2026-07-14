# Source

Source is an experimental Git server that stores Git objects in a database
rather than on a filesystem.

## Introduction

Traditionally, Git hosting platforms rely on the local filesystem as the
durable store for repository data. As a deployment grows, concerns such as
replication, backups, failover, and storage management can only be
addressed by introducing increasingly sophisticated storage infrastructure
around that filesystem.

Source explores a different approach, delegating object storage entirely
to a database. By using a database as the system of record, Source can
inherit the durability, replication, backup, and recovery mechanisms that
databases already provide, while avoiding the operational complexity of
distributed filesystems.

The goal is not to change Git itself, but to reconsider how Git servers
are built. Git remains the protocol, the object model, and the user
experience, but the underlying storage architecture, however, is designed
around software that has spent decades addressing the complexities of
availability and data integrity at scale.

Source is intentionally narrow in scope. While it includes a lightweight
web interface for browsing repositories, it does not aim to compete with
full-featured software forges. Features such as issue tracking, CI/CD,
project management, and social coding are deliberately outside the core
mission of the project.

Instead, the focus is on providing a performant, scalable Git server that
can be deployed as a single binary and serve as a foundation upon which
other tools and applications can be built.

At present, PostgreSQL is the only supported database backend. Support for
additional databases will be introduced in the future, but the immediate
goal is to deliver a stable beta release and validate the architectural
ideas that motivated the project.

## Getting Started

Source requires a PostgreSQL database to store repository data. While
setting up a PostgreSQL database is beyond the scope of this document, the
following command should provide a suitable instance for experimentation
and development[^1].

```sh
docker run -e POSTGRES_HOST_AUTH_METHOD=trust -p 5433:5432 -d postgres
```

### Building Source

Building Source requires Go 1.25+, Node.js 24.16+, Git, and Make (optional but
recommended).

#### With Make

When using Make, all build steps are handled automatically, including installing
the frontend dependencies, building the frontend assets, and compiling the
standalone Source binary.

```sh
# Clone the repository
git clone https://github.com/iainjreid/source.git && cd source

# Build Source
make
```

#### Without Make

If Make is not available on your system, the following build steps can be run
manually. Certain build options are omitted for brevity but the resulting binary
will still operate in the same way.

```sh
# Clone the repository
git clone https://github.com/iainjreid/source.git && cd source

# Install the frontend dependencies
npm ci --ignore-scripts

# Build the frontend assets
npm run build

# Compile the binary
go build -buildvcs=true	-tags=standalone ./cmd/src
```

### Initialising the Database

Before Source can serve repositories, the database schema must be created
and at least one repository must have been imported.

The command below will create the required schemas and clone Source's own
source code into the database (you can change this by providing a URL to
a repository of your own choice).

```sh
./src manage --db "postgresql://postgres@localhost" --setup \
             --clone https://github.com/iainjreid/source.git
```

### Starting the Server

Once initialised, you can go ahead and start your Source instance.

```sh
./src start --db "postgresql://postgres@localhost"
```

The web interface can be found at <http://localhost:8080>, but the port
can be specified to suit your needs along with a handful of other options
documented under the `./src start --help` command.

### Enabling SSH Access

Source includes a built-in SSH server. It does not require OpenSSH to be
installed, but it does require a valid private key to operate.

To generate a private key for the server, run the following and be sure to
leave the passphrase empty[^2].

```sh
ssh-keygen -m PEM -t rsa -b 4096 -f ./key # ... or any other path
```

This key can be passed to Source using the `--ssh-key` flag along with the
other required runtime flags.

```sh
./src start --db "postgresql://postgres@localhost" \
            --ssh-id-path ./key
```

## Troubleshooting

### SSH agent socket unavailable

```txt
SSH agent requested but SSH_AUTH_SOCK not-specified
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
Sourcehut, although it lacks functionality beyond simply browsing
a repository.

### Future development

* Issue management via Git notes, with the ability to discuss and provide
  feedback within the repository itself.

* In repository support for patch-by-patch reviews, a highly successful
  approach to peer reviewing changes used often within Git mailing lists.

* Support partial checkouts without requiring changes to how Git stores
  objects internally.

### Current limitations

* Synchronising objects from upstream remotes over SSH will require
  additional work to support knownhosts first.

* Git LFS is not yet supported. Given that LFS is an extension to Git, and
  not core behaviour, implementing this is not currently a high-priority
  feature.

  That said, there's nothing stopping Source from supporting LFS. Many
  databases offer a range of options to support large file storage,
  including PostgreSQL.

## License

Source is licensed under the Apache License. See [LICENSE](./LICENSE) for
details.

[^1]: Setting `POSTGRES_HOST_AUTH_METHOD=trust` disables password
    authentication and should only be used on personal or well secured
    machines. Official information about trust authentication can be
    [found here](https://www.postgresql.org/docs/18/auth-trust.html).

[^2]: Supporting password protected private keys would require the user to
    provide plain text passwords in order for the program to decrypt them,
    rendering any illusion of increased security redundant.
