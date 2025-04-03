package common

import "github.com/pulumi/pulumi-go-provider/infer"

type Logging string

const (
	LogStdout          Logging = "stdout"
	LogStderr          Logging = "stderr"
	LogStdoutAndStderr Logging = "stdoutAndStderr"
	NoLogging          Logging = "none"
)

func (Logging) Values() []infer.EnumValue[Logging] {
	return []infer.EnumValue[Logging]{
		{Name: string(LogStdout), Value: LogStdout, Description: "Capture stdout in logs but not stderr"},
		{Name: string(LogStderr), Value: LogStderr, Description: "Capture stderr in logs but not stdout"},
		{Name: string(LogStdoutAndStderr), Value: LogStdoutAndStderr, Description: "Capture stdout and stderr in logs"},
		{Name: string(NoLogging), Value: NoLogging, Description: "Capture no logs"},
	}
}

func (l *Logging) ShouldLogStdout() bool {
	return l == nil || *l == LogStdout || *l == LogStdoutAndStderr
}
func (l *Logging) ShouldLogStderr() bool {
	return l == nil || *l == LogStderr || *l == LogStdoutAndStderr
}
