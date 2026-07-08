# Deployment and miscellany

This folder includes deployment related and other "batteries included"
knick-knacks that will hopefully encourage those of a curious nature to evaluate
Source for their needs.

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
