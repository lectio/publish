package publish

import (
	"context"
	"github.com/lectio/dropmark"
	"github.com/lectio/link"
	"github.com/lectio/markdown"
	"testing"

	"github.com/stretchr/testify/suite"
)

type PublishSuite struct {
	suite.Suite
}

func (suite *PublishSuite) SetupSuite() {
}

func (suite *PublishSuite) TearDownSuite() {
}

// satisfies TestConfigurator interface
func (suite *PublishSuite) StopAfterTestItemsCount(ctx context.Context) uint {
	return 10
}

// satisfies TestConfigurator interface
func (suite *PublishSuite) SimulateLinkScores(ctx context.Context) bool {
	return true
}

func (suite *PublishSuite) TestDropmarkToMarkdown() {
	ctx := context.Background()

	var rpr dropmark.ReaderProgressReporter = &CommandLineProgressReporter{prefix: "[TESTING] "}
	var bpr dropmark.BoundedProgressReporter = &CommandLineProgressReporter{prefix: "[TESTING] "}

	bpc := markdown.NewBasePathConfigurator("test_001")
	linkFactory := link.NewFactory()

	publisher, err := NewMarkdownPublisher(ctx, true, linkFactory, bpc, rpr, bpr, suite)
	suite.Nil(err, "No error should be reported")

	err = publisher.Publish(ctx, "https://shah.dropmark.com/616548.json", rpr, bpr, suite)
	suite.Nil(err, "No error should be reported")
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(PublishSuite))
}
