package publish

import (
	"context"
	"fmt"
	"gopkg.in/cheggaaa/pb.v1"
	"io"
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

type SummaryProgressReporter struct{ prefix string }

func (pr SummaryProgressReporter) StartReportableActivity(ctx context.Context, summary string, expectedItems int) {
	fmt.Printf("[%s] %s\n", pr.prefix, summary)
}

func (pr SummaryProgressReporter) StartReportableReaderActivityInBytes(ctx context.Context, summary string, exepectedBytes int64, inputReader io.Reader) io.Reader {
	fmt.Printf("[%s] %s\n", pr.prefix, summary)
	return inputReader
}

func (pr SummaryProgressReporter) IncrementReportableActivityProgress(ctx context.Context) {
}

func (pr SummaryProgressReporter) IncrementReportableActivityProgressBy(ctx context.Context, incrementBy int) {
}

func (pr SummaryProgressReporter) CompleteReportableActivityProgress(ctx context.Context, summary string) {
	fmt.Printf("[%s] %s\n", pr.prefix, summary)
}

type CommandLineProgressReporter struct {
	prefix string
	bar    *pb.ProgressBar
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
}
