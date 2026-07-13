# Deployment and miscellany

This folder includes deployment related and other "batteries included"
knick-knacks that will hopefully encourage those of a curious nature to evaluate
Source for their needs.

## Running Source with Docker Compose

With [Docker Engine] and [Docker Compose] installed, the following instruction
will start a complete Source installation that's ideal for personal or
evaluation use.

```sh
# Copy the example configuration to make your own changes
cp contrib/docker/compose.yml .

# Start the service
docker compose up
```

[Docker Engine]: https://docs.docker.com/engine
[Docker Compose]: https://docs.docker.com/compose

### Enabling SSH access

To enable SSH access, either `SOURCE_SSH_ID_PATH` or `SOURCE_SSH_ID` must be
supplied. Whichever value you choose to supply, a private SSH key will be needed
to secure the identity of the server.

Do not use a private key used by any other service. Before continuing, create
a new key to current standards and store that key somewhere safe.

```sh
ssh-keygen -m PEM -t rsa -b 4096 -f ./key # ... or any other path
```

#### With `SOURCE_SSH_ID_PATH` (recommended)

Secrets such as an SSH private key can easily be mounted securely using this
approach. Simply define a new secret at the end of the `compose.yml` file that
references the newly create SSH private key.

```yaml
secrets:
  ssh_id:
    file: ./key # ... the path to the private key
```

With the secret defined, expose it to the `source` service and provide it with
a path to the file that contains it.

```yaml
services:
    source:
        ...
        environment:
            SOURCE_SSH_ID_PATH: /run/secrets/ssh_id
        ...
        secrets:
            - ssh_id
```

#### With `SOURCE_SSH_ID` (not recommended)

Although not recommended, it is possible to paste the contents of the private
key into your `compose.yml` file.

```yaml
services:
    source:
        ...
        environment:
            SOURCE_SSH_ID: |
                -----BEGIN RSA PRIVATE KEY-----
                ...
                -----END RSA PRIVATE KEY-----
```

### Reloading changes

If you have made changes to your `compose.yml` file and wish to see those
reflected in your Source instance, then recreate the project.

```sh
docker compose up --force-recreate -d
```

Alternatively, if you have pulled upstream changes or made changes yourself and
wish to rebuild Source entirely, then you will also need to pass the `--build`
flag.

```sh
docker compose up --force-recreate --build -d
```

None of these actions will impact the data stored in persistent volumes.

## Running Source with `systemd`

The files in the `contrib/systemd` directory provide a minimal example of how
Source can be managed using `systemd`.

These examples are intended as a starting point and, while suitable for personal
use or for evaluation, should not be considered canonical configurations.

System administrators and package maintainers are encouraged to adapt them to
their platform's conventions and operational requirements rather than copying
them verbatim.

The following assumes you have built and installed Source using the instructions
provided in the root `README.md` file.

### Setup

Because Source accepts traffic from external actors, it should be run under
a dedicated service account, and follow the principle of least privilege to
limit the impact of a potential service compromise.

The command below creates a system user without a home directory and with
interactive logins disabled:

```sh
sudo useradd --system --no-create-home --shell /usr/sbin/nologin source
```

With the service account created, the unit file and associated configuration
file can be installed.

```sh
sudo install -Dm644 contrib/systemd/source.service /etc/systemd/system/source.service
sudo install -Dm640 contrib/systemd/source.env /etc/source/source.env
```

Runtime configuration is supplied via environment variables defined in
`source.env`. Edit `/etc/source/source.env` to suit your database, environment,
and desired behaviour.

### Running

With the service configured, reload the service manager and enable Source:

```sh
sudo systemctl daemon-reload
sudo systemctl enable --now source
```

Check the service status:

```sh
systemctl status source
journalctl -u source -f
```

[^1]: Supporting password protected private keys would require the user to
    provide plain text passwords in order for the program to decrypt them,
    rendering any illusion of increased security redundant.
