package notify

import (
	"time"

	"github.com/RedHatInsights/host-metering/hostinfo"
	"github.com/prometheus/prometheus/prompb"
)

type NotifyStatus int

const (
	StatusOK NotifyStatus = iota
	StatusRecoverableError
	StatusNonRecoverableError
)

type NotifyResult struct {
	Status NotifyStatus
	err    error
}

func (e *NotifyResult) Err() error {
	return e.err
}

func (e *NotifyResult) Error() string {
	if e.err == nil {
		return ""
	}
	return e.err.Error()
}

func OKResult() *NotifyResult {
	return &NotifyResult{StatusOK, nil}
}

func RecoverableErrorResult(err error) *NotifyResult {
	return &NotifyResult{StatusRecoverableError, err}
}

func NonRecoverableErrorResult(err error) *NotifyResult {
	return &NotifyResult{StatusNonRecoverableError, err}
}

type Notifier interface {
	Notify(samples []prompb.Sample, hostinfo *hostinfo.HostInfo) *NotifyResult

	// HostChanged tells notifier that related information on host has changed
	HostChanged()
}

func FilterSamplesByAge(samples []prompb.Sample, maxAge time.Duration) []prompb.Sample {
	treshold := time.Now().UnixMilli() - int64(maxAge.Milliseconds())
	for idx, sample := range samples {
		if sample.Timestamp >= treshold {
			return samples[idx:]
		}
	}
	return []prompb.Sample{}
}
