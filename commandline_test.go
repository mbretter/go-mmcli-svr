package main

import (
	"bytes"
	"errors"
	"github.com/mbretter/go-mmcli-svr/backend"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"testing"
)

func TestCommandLine_Activate(t *testing.T) {
	tests := []struct {
		name      string
		execArgs  []string
		execError error
		logs      []string
		settings  *CommandLine
	}{
		{
			name:     "GPS refresh success",
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
			name:     "Location 3gpp enable success",
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
