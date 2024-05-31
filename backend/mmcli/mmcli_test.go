package mmcli

import (
	"errors"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type CmdMock struct {
	mock.Mock
}

func TestMmcli_Exec(t *testing.T) {
	tests := []struct {
		name          string
		args          []string
		expected      []byte
		expectedErr   error
		execArgs      []any
		execReturn    []byte
		execReturnErr error
	}{
		{
			name:       "Success",
			args:       []string{"arg1", "arg2"},
			expected:   []byte("Success"),
			execArgs:   []any{"-J", "arg1", "arg2"},
			execReturn: []byte("Success"),
		},
		{
			name:          "Error",
			args:          []string{"arg1", "arg2"},
			expectedErr:   errors.New("failed"),
			execArgs:      []any{"-J", "arg1", "arg2"},
			execReturn:    []byte(""),
			execReturnErr: errors.New("failed"),
		},
		{
			name:          "ExitError",
			args:          []string{"arg1", "arg2"},
			expectedErr:   errors.New("exec failed"),
			execArgs:      []any{"-J", "arg1", "arg2"},
			execReturn:    []byte(""),
			execReturnErr: &exec.ExitError{Stderr: []byte("exec failed\n")},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockExec := NewExecCommandMock(t)
			mockOutput := NewExecCommandOutputMock(t)
			mockExec.EXPECT().Execute("mmcli", tt.execArgs...).Return(mockOutput)
			mockOutput.EXPECT().Output().Return(tt.execReturn, tt.execReturnErr)

			cli := Provide()
			cli.exec = mockExec.Execute

			got, err := cli.Exec(tt.args...)

			if tt.expectedErr == nil {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, got)
			} else {
				assert.EqualError(t, err, tt.expectedErr.Error())
			}
		})
	}
}

func TestMmcli_ExecModem(t *testing.T) {
	tests := []struct {
		name          string
		modem         string
		args          []string
		expected      []byte
		expectedErr   error
		execArgs      []any
		execReturn    []byte
		execReturnErr error
	}{
		{
			name:       "Success with modem",
			modem:      "0",
			args:       []string{"arg1", "arg2"},
			expected:   []byte("Success"),
			execArgs:   []any{"-J", "-m", "0", "arg1", "arg2"},
			execReturn: []byte("Success"),
		},
		{
			name:       "Success without modem",
			modem:      "",
			args:       []string{"arg1", "arg2"},
			expected:   []byte("Success"),
			execArgs:   []any{"-J", "-m", "/org/freedesktop/ModemManager1/Modem/0", "arg1", "arg2"},
			execReturn: []byte("Success"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockExec := NewExecCommandMock(t)
			if tt.modem == "" {
				mockOutput := NewExecCommandOutputMock(t)
				mockExec.EXPECT().Execute("bash", "-c", "mmcli -L -J | jq -r '.\"modem-list\"[0]'").Return(mockOutput)
				mockOutput.EXPECT().Output().Return([]byte("/org/freedesktop/ModemManager1/Modem/0"), tt.execReturnErr)
			}

			mockOutput := NewExecCommandOutputMock(t)
			mockExec.EXPECT().Execute("mmcli", tt.execArgs...).Return(mockOutput)
			mockOutput.EXPECT().Output().Return(tt.execReturn, tt.execReturnErr)

			cli := Provide()
			cli.exec = mockExec.Execute

			got, err := cli.ExecModem(tt.modem, tt.args...)

			if tt.expectedErr == nil {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, got)
			} else {
				assert.EqualError(t, err, tt.expectedErr.Error())
			}
		})
	}
}
