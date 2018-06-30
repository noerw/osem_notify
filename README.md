# osem_notify 🔆🌡📡📈  ⚠ 📲

Cross platform command line application to run health checks against sensor stations registered on [openSenseMap.org](https://opensensemap.org).

This tool lets you automatically check if senseBoxes are still runnning correctly,
and when that's not the case, notifies you.
Currently, email notifications are implemented, but other transports can be added easily.
Implemented health checks are [described below](#available-healthchecks), and new ones can be added just as easily (given some knowledge of programming).

The tool has multiple modes of operation:

- `osem_notify check boxes`: run one-off checks on boxes
- `osem_notify watch boxes`: check boxes continuously.

Run `osem_notify help` or check the manual in the [docs/](docs/osem_notify.md) directory for more details.

## get it
Download a build from the [releases page](https://github.com/noerw/osem_notify/releases/).
You can run the application by running `./osem_notify*` in a terminal in your downloads directory.

On unix platforms you may need to make it executable, and can add it to your `$PATH` for convenience, so it is always callable via `osem_notify`:
```sh
chmod +x osem_notify*
sudo mv osem_notify* /usr/bin/osem_notify
```

## configure it
The tool works out of the box for basic functionality, but must be configured to set up notifications.
Configuration can be done via a YAML file located at `~/.osem_notify.yml` or through environment variables.
Run `osem_notify help config` for details and an example configuration.

### available healthchecks

`type`               | description
---------------------|------------
`measurement_age`    | Alert when a sensor has not submitted measurements within a given duration.
`measurement_faulty` | Alert when a sensor's last reading was a presumably faulty value (e.g. broken / disconnected sensor).
`measurement_min`    | Alert when a sensor's last measurement is lower than a given value.
`measurement_max`    | Alert when a sensor's last measurement is higher than a given value.

### available notification transports
`transport` | `options`
------------|------------
`email`     | `recipients`: list of email addresses

Want more? [add it](#contribute)!

## build it
Want to use `osem_notify` on a platform where no builds are provided?

Assuming you have golang installed, run
```sh
go get -v -d ./
go build main.go
```

For cross-compilation, check [this guide](https://dave.cheney.net/2015/08/22/cross-compilation-with-go-1-5) out.

## contribute
Contributions are welcome!
Check out the following locations for plugging in new functionality:

- new notification transports: [core/notifiers.go](core/notifiers.go)
- new health checks: [core/healthcheck*.go](core/healthchecks.go)
- new commands: [cmd/](cmd/)

Before committing and submitting a pull request, please run `go fmt ./ cmd/ core/`.

## license
GPL-3.0 Norwin Roosen
