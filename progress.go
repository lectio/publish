package publish

import (
	"context"
	"fmt"
	"gopkg.in/cheggaaa/pb.v1"
	"io"
	"sync"
)

// ReaderProgressReporter is sent to this package's methods if activity progress reporting is expected for an io.Reader
type ReaderProgressReporter interface {
	StartReportableReaderActivityInBytes(ctx context.Context, summary string, exepectedBytes int64, inputReader io.Reader) io.Reader
	CompleteReportableActivityProgress(ctx context.Context, summary string)
}

// BoundedProgressReporter is one observation method for live reporting of long-running processes where the upper bound is known
type BoundedProgressReporter interface {
	StartReportableActivity(ctx context.Context, summary string, expectedItems int)
	IncrementReportableActivityProgress(ctx context.Context)
	IncrementReportableActivityProgressBy(ctx context.Context, incrementBy int)
	CompleteReportableActivityProgress(ctx context.Context, summary string)
}

type ExceptionReporter interface {
	ReportError(context.Context, error) bool
	MaxErrorsReached(context.Context) bool
	ReportWarning(ctx context.Context, code, message string) bool
}

type SilentProgressReporter struct{}

func (pr SilentProgressReporter) StartReportableActivity(ctx context.Context, summary string, expectedItems int) {
}

func (pr SilentProgressReporter) StartReportableReaderActivityInBytes(ctx context.Context, summary string, exepectedBytes int64, inputReader io.Reader) io.Reader {
	return inputReader
}

func (pr SilentProgressReporter) IncrementReportableActivityProgress(ctx context.Context) {
}

func (pr SilentProgressReporter) IncrementReportableActivityProgressBy(ctx context.Context, incrementBy int) {
}

func (pr SilentProgressReporter) CompleteReportableActivityProgress(ctx context.Context, summary string) {
}

func (pr SilentProgressReporter) ReportError(context.Context, error) bool {
	return true
}

func (pr SilentProgressReporter) MaxErrorsReached(context.Context) bool {
	return false
}

func (pr SilentProgressReporter) ReportWarning(ctx context.Context, code, message string) bool {
	return true
}

type SummaryProgressReporter struct {
	prefix         string
	mu             sync.RWMutex
	errorsReported uint
	maxErrors      uint
}

func (pr *SummaryProgressReporter) StartReportableActivity(ctx context.Context, summary string, expectedItems int) {
	fmt.Printf("%s%s\n", pr.prefix, summary)
}

func (pr *SummaryProgressReporter) StartReportableReaderActivityInBytes(ctx context.Context, summary string, exepectedBytes int64, inputReader io.Reader) io.Reader {
	fmt.Printf("%s%s\n", pr.prefix, summary)
	return inputReader
}

func (pr *SummaryProgressReporter) IncrementReportableActivityProgress(ctx context.Context) {
}

func (pr *SummaryProgressReporter) IncrementReportableActivityProgressBy(ctx context.Context, incrementBy int) {
}

func (pr *SummaryProgressReporter) CompleteReportableActivityProgress(ctx context.Context, summary string) {
	fmt.Printf("%s%s\n", pr.prefix, summary)
}

func (pr *SummaryProgressReporter) ReportError(ctx context.Context, err error) bool {
	pr.mu.Lock()
	pr.errorsReported++
	pr.mu.Unlock()
	fmt.Printf("%s%v\n", pr.prefix, err.Error())
	return !pr.MaxErrorsReached(ctx)
}

func (pr *SummaryProgressReporter) MaxErrorsReached(context.Context) bool {
	pr.mu.RLock()
	defer pr.mu.RUnlock()
	return pr.maxErrors > 0 && pr.errorsReported > pr.maxErrors
}

func (pr *SummaryProgressReporter) ReportWarning(ctx context.Context, code, message string) bool {
	fmt.Printf("%s%s %s\n", pr.prefix, code, message)
	return true
}

type CommandLineProgressReporter struct {
	prefix    string
	mu        sync.RWMutex
	bar       *pb.ProgressBar
	errList   []error
	warnList  []struct{ code, message string }
	maxErrors int
}

func (pr *CommandLineProgressReporter) StartReportableActivity(ctx context.Context, summary string, expectedItems int) {
	fmt.Printf("%s%s\n", pr.prefix, summary)
	pr.bar = pb.StartNew(expectedItems)
	pr.bar.ShowCounters = true
}

func (pr *CommandLineProgressReporter) StartReportableReaderActivityInBytes(ctx context.Context, summary string, exepectedBytes int64, inputReader io.Reader) io.Reader {
	pr.bar = pb.New(int(exepectedBytes)).SetUnits(pb.U_BYTES)
	pr.bar.Start()
	return pr.bar.NewProxyReader(inputReader)
}

func (pr *CommandLineProgressReporter) IncrementReportableActivityProgress(ctx context.Context) {
	pr.bar.Increment()
}

func (pr *CommandLineProgressReporter) IncrementReportableActivityProgressBy(ctx context.Context, incrementBy int) {
	pr.bar.Add(incrementBy)
}

func (pr *CommandLineProgressReporter) CompleteReportableActivityProgress(ctx context.Context, summary string) {
	pr.bar.FinishPrint(fmt.Sprintf("%s%s\n", pr.prefix, summary))

	if len(pr.errList) > 0 {
		for _, err := range pr.errList {
			fmt.Printf("%s%v\n", pr.prefix, err.Error())
		}
	}

	if len(pr.warnList) > 0 {
		for _, warning := range pr.warnList {
			fmt.Printf("%s%s %s\n", pr.prefix, warning.code, warning.message)
		}
	}
}

func (pr *CommandLineProgressReporter) ReportError(ctx context.Context, err error) bool {
	pr.mu.Lock()
	pr.errList = append(pr.errList, err)
	pr.mu.Unlock()
	return !pr.MaxErrorsReached(ctx)
}

func (pr *CommandLineProgressReporter) MaxErrorsReached(context.Context) bool {
	pr.mu.RLock()
	defer pr.mu.RUnlock()
	return pr.maxErrors > 0 && len(pr.errList) > pr.maxErrors
}

func (pr *CommandLineProgressReporter) ReportWarning(ctx context.Context, code, message string) bool {
	pr.mu.Lock()
	pr.warnList = append(pr.warnList, struct{ code, message string }{code: code, message: message})
	pr.mu.Unlock()
	return true
}
