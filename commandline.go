package main

import (
	"flag"
	"fmt"
	"github.com/mbretter/go-mmcli-svr/backend"
	"log/slog"
	"slices"
	"strings"
)

const (
	defaultListen         = "127.0.0.1:8743"
	defaultLocationEnable = ""
	defaultGpsRefresh     = 0
)

type CommandLine struct {
	Listen             string
	LocationGatherings []string
	GpsRefresh         int

	log     *slog.Logger
	backend backend.Backend
}

func NewCommandLine(logger *slog.Logger, backend backend.Backend) *CommandLine {
	return &CommandLine{
		log:     logger,
		backend: backend,
		Listen:  defaultListen,
	}
}

func (c *CommandLine) Parse() error {
	var listen = flag.String("listen", defaultListen, "listen: <ip:port|:port>")
	var enableLocation = flag.String("location-enable", defaultLocationEnable, "enable location gathering: <all|gps-nmea|gps-raw|3gpp|agps‐msa|agps‐msb>")
	var gpsRefresh = flag.Int("gps-refresh", defaultGpsRefresh, "gps refresh rate in seconds")
	flag.Parse()

	if *listen != "" {
		c.Listen = *listen
	}

	if *gpsRefresh > 0 {
		c.GpsRefresh = *gpsRefresh
	}

	if *enableLocation != "" {
		for _, typ := range strings.Split(*enableLocation, ",") {
			if slices.Contains(c.LocationGatherings, typ) {
				continue
			}
			switch typ {
			case "agps-msa":
				c.LocationGatherings = append(c.LocationGatherings, "agps-msa", "gps-nmea")
			case "agps-msb":
				c.LocationGatherings = append(c.LocationGatherings, "agps-msb", "gps-nmea")
			case "gps-nmea":
				c.LocationGatherings = append(c.LocationGatherings, "gps-nmea")
			case "gps-raw":
				c.LocationGatherings = append(c.LocationGatherings, "gps-raw")
			case "3gpp":
				c.LocationGatherings = append(c.LocationGatherings, "3gpp")
			default:
				return fmt.Errorf("unknown location type: %s", typ)
			}
		}
	}

	slices.Sort(c.LocationGatherings)
	c.LocationGatherings = slices.Compact(c.LocationGatherings)

	return nil
}

func (c *CommandLine) Activate() {
	if c.GpsRefresh > 0 {
		c.log.Info("Set gps-refresh", "refresh", c.GpsRefresh)
		_, err := c.backend.ExecModem("", fmt.Sprintf("--location-set-gps-refresh-rate=%d", c.GpsRefresh))
		if err != nil {
			c.log.Error("Failed to execute command", "error", err)
		}
	}

	if len(c.LocationGatherings) > 0 {
		c.log.Info("Enable location gatherings", "location", strings.Join(c.LocationGatherings, ","))
		var args []string
		for _, typ := range c.LocationGatherings {
			args = append(args, "--location-enable-"+typ)
		}
		c.log.Debug("mmcli", "command", "mmcli "+strings.Join(args, " "))
		out, err := c.backend.ExecModem("", args...)
		if err == nil {
			c.log.Info("mmcli", "message", strings.Trim(string(out), "\n "))
		} else {
			c.log.Error("mmcli failed", "command", "mmcli "+strings.Join(args, " "), "error", err)
		}

	}
}
