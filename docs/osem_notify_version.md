## osem_notify version

Get build and version information

### Synopsis

osem_notify version returns its build and version information

```
osem_notify version [flags]
```

### Options

```
  -h, --help   help for version
```

### Options inherited from parent commands

```
  -a, --api string         openSenseMap API to query against (default "https://api.opensensemap.org")
  -c, --config string      path to config file (default $HOME/.osem_notify.yml)
  -d, --debug              enable verbose logging
  -l, --logformat string   log format, can be plain or json (default "plain")
  -n, --notify             if set, will send out notifications.
                           Otherwise results are printed to stdout only.
                           You might want to run 'osem_notify debug notifications' first to verify everything works.
                           
```

### SEE ALSO

* [osem_notify](osem_notify.md)	 - Root command displaying help

###### Auto generated by spf13/cobra on 26-Jun-2018