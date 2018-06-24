# osem_notify 🔆🌡📡📈  ⚠ 📲

Cross platform command line application to run health checks against sensor stations registered on [openSenseMap.org](https://opensensemap.org).

This tool quickly runs various health checks on selected senseBoxes,
and can send out notifications via various protocols.
Specifically email is implemented, but other transports can be added easily.

The tool can also operate in watch mode, checking boxes at regular intervals.

Check the manual in the [doc/](doc/osem_notify.md) directory for a description of features.

## get it
Download a build from the [releases page](releases/).
You can immediately call the application by running `./osem_notify` in a terminal in your downloads directory.

On unix platforms you may add it to your path for convenience, so it is always callable via `osem_notify`:
```sh
sudo mv osem_notify /usr/bin/osem_notify
```

## configure it
Configuration is required to set up notification transports, and can set the default healthchecks.
Configuration can be done via a YAML file located at `~/.osem_notify.yml` or through environment variables.
Example configuration:

```yaml
defaultHealthchecks:
  notifications:
    transport: email
    options:
      recipients:
      - fridolina@example.com
      - ruth.less@example.com
  events:
    - type: "measurement_age"
      target: "all"    # all sensors
      threshold: "15m" # any duration
    - type: "measurement_faulty"
      target: "all"
      threshold: ""

# only needed when sending notifications via email
email:
  host: smtp.example.com
  port: 25
  user: foo
  pass: bar
  from: hildegunst@example.com
```

### possible values for `notifications`:
`transport` | `options`
------------|------------
`email`     | `recipients`: list of email addresses

Want more? [add it](#contribute)!

### possible values for `events[]`:

`type`               | description
---------------------|------------
`measurement_age`    | Alert when sensor `target` has not submitted measurements within `threshold` duration.
`measurement_faulty` | Alert when sensor `target`'s last reading was a presumably faulty value (e.g. broken / disconnected sensor).
`measurement_min`    | Alert when sensor `target`'s last measurement is lower than `threshold`.
`measurement_max`    | Alert when sensor `target`'s last measurement is higher than `threshold`.

`target` can be either a sensor ID, or `"all"` to match all sensors of the box.
`threshold` must be a string.

### configuration via environment variables
Instead of a YAML file, you may configure the tool through environment variables. Keys are the same as in the YAML, but:

- prefixed with `osem_notify_`
- path separator is not `.`, but `_`

Example: `OSEM_NOTIFY_EMAIL_PASS=supersecret osem_notify check boxes`

## build it
Want to use `osem_notify` on a platform where no builds are provided?

Assuming you have golang installed, run
```sh
go get -v -d ./
go build ./
```

For cross-compilation, check [this guide](https://dave.cheney.net/2015/08/22/cross-compilation-with-go-1-5) out.

## contribute
Contributions are welcome!
Check out the following locations for plugging in new functionality:

- new notification transports: [core/notifiers.go](core/notifiers.go)
- new health checks: [core/Box.go](core/Box.go)
- new commands: [cmd/](cmd/)

Before committing and submitting a pull request, please run `go fmt ./ cmd/ core/`.

## license
GPL-3.0 Norwin Roosen
