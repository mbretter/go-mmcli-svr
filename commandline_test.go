package main

import (
	"bytes"
	"errors"
	"flag"
	"github.com/mbretter/go-mmcli-svr/backend"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"os"
	"testing"
)

func TestCommandLine_New(t *testing.T) {
	var buff bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buff, nil))

	backendMock := backend.NewBackendMock(t)
	cmd := NewCommandLine(logger, backendMock)

	assert.Equal(t, defaultListen, cmd.Listen)
	assert.Equal(t, defaultGpsRefresh, cmd.GpsRefresh)
	assert.Empty(t, cmd.LocationGatherings)
	assert.Same(t, logger, cmd.log)
	assert.Same(t, backendMock, cmd.backend)
}

func TestCommandLine_Parse(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		settings *CommandLine
		error    string
	}{
		{
			name: "No args",
			args: []string{},
			settings: &CommandLine{
				Listen: defaultListen,
			},
		},
		{
			name: "Listen",
			args: []string{"-listen=192.168.1.22:12345"},
			settings: &CommandLine{
				Listen: "192.168.1.22:12345",
			},
		},
		{
			name: "GPS refresh",
			args: []string{"-gps-refresh=5"},
			settings: &CommandLine{
				Listen:     defaultListen,
				GpsRefresh: 5,
			},
		},
		{
			name: "Location enable agps-msa",
			args: []string{"-location-enable=agps-msa"},
			settings: &CommandLine{
				Listen:             defaultListen,
				LocationGatherings: []string{"agps-msa", "gps-nmea"}, // assistet gps needs gps-nmea
			},
		},
		{
			name: "Location enable gps-nmea",
			args: []string{"-location-enable=gps-nmea"},
			settings: &CommandLine{
				Listen:             defaultListen,
				LocationGatherings: []string{"gps-nmea"}, // assistet gps needs gps-nmea
			},
		},
		{
			name: "Location enable agps-msa",
			args: []string{"-location-enable=agps-msb"},
			settings: &CommandLine{
				Listen:             defaultListen,
				LocationGatherings: []string{"agps-msb", "gps-nmea"}, // assistet gps needs gps-nmea
			},
		},
		{
			name: "Location enable multi",
			args: []string{"-location-enable=agps-msa,gps-nmea,gps-raw,3gpp"},
			settings: &CommandLine{
				Listen:             defaultListen,
				LocationGatherings: []string{"3gpp", "agps-msa", "gps-nmea", "gps-raw"},
			},
		},
		{
			name: "Unknown location",
			args: []string{"-location-enable=xxx"},
			settings: &CommandLine{
				Listen: defaultListen,
			},
			error: "unknown location type: xxx",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			os.Args = append([]string{"cmd"}, tt.args...)
			// reset flag, otherwise redefinition errors might be thrown
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError) //flags are now reset

			cmd := &CommandLine{}
			err := cmd.Parse()

			assert.Equal(t, tt.settings, cmd)
			if len(tt.error) > 0 {
				assert.Equal(t, tt.error, err.Error())
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestCommandLine_Activate(t *testing.T) {
	tests := []struct {
		name      string
		execArgs  []string
		execError error
		logs      []string
		settings  *CommandLine
	}{
		{
			name:     "GPS refresh",
			execArgs: []string{"--location-set-gps-refresh-rate=10"},
			logs:     []string{`level=INFO msg="Set gps-refresh" refresh=10`},
			settings: &CommandLine{
				GpsRefresh: 10,
			},
		},
		{
			name:      "GPS refresh error",
			execArgs:  []string{"--location-set-gps-refresh-rate=10"},
			execError: errors.New("failed"),
			logs:      []string{`level=INFO msg="Set gps-refresh" refresh=10`, `level=ERROR msg="Failed to execute command"`},
			settings: &CommandLine{
				GpsRefresh: 10,
			},
		},
		{
			name:     "Location 3gpp enable",
			execArgs: []string{"--location-enable-3gpp"},
			logs:     []string{`level=INFO msg="Enable location gatherings" location=3gpp`, `level=INFO msg=mmcli message=success`},
			settings: &CommandLine{
				LocationGatherings: []string{"3gpp"},
			},
		},
		{
			name:     "Location enable agps",
			execArgs: []string{"--location-enable-agps-msa", "--location-enable-gps-nmea"},
			logs:     []string{`level=INFO msg="Enable location gatherings" location=agps-msa,gps-nmea`, `level=INFO msg=mmcli message=success`},
			settings: &CommandLine{
				LocationGatherings: []string{"agps-msa", "gps-nmea"},
			},
		},
		{
			name:      "Location enable agps error",
			execArgs:  []string{"--location-enable-agps-msa", "--location-enable-gps-nmea"},
			execError: errors.New("failed"),
			logs:      []string{`level=INFO msg="Enable location gatherings" location=agps-msa,gps-nmea`, `level=ERROR msg="mmcli failed" command="mmcli --location-enable-agps-msa --location-enable-gps-nmea" error=failed`},
			settings: &CommandLine{
				LocationGatherings: []string{"agps-msa", "gps-nmea"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buff bytes.Buffer
			logger := slog.New(slog.NewTextHandler(&buff, nil))

			backendMock := backend.NewBackendMock(t)

			args := make([]any, len(tt.execArgs))
			for i, arg := range tt.execArgs {
				args[i] = arg
			}
			backendMock.EXPECT().ExecModem("", args...).Return([]byte("success"), tt.execError)

			tt.settings.backend = backendMock
			tt.settings.log = logger
			tt.settings.Activate()

			for _, log := range tt.logs {
				assert.Contains(t, buff.String(), log)
			}
		})
	}
}
