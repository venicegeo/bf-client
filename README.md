# bf-client

This is the Beachfront command line app and associated libraries (in Go)

# Usage

```
NAME:
   beachfront - access the Beachfront services

USAGE:
   main [global options] command [command options] [arguments...]

VERSION:
   0.0.0

COMMANDS:
     catalog, cat      access catalog (imagery feed) services
     job               access job services
     coastline, coast  access coastline data
     algorithm, alg    access the algorithm services
     help, h           Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

# Examples

* `beachfront` catalog --info landsat



# TO DO

* Add progress meter for downloads
* Fix the base-path-prefix problem
* Support submit job
* Support delete job
* Support parameters on search images: cloud cover, bbox, dates
* Support feeds other than Planet (when BF does)
* move all Planet into a Planet class
* gets from Planet prob need pagination support
