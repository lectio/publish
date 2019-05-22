package publish

import (
	"context"
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

// // satisfies TestConfigurator interface
// func (suite *PublishSuite) StopAfterTestItemsCount(ctx context.Context) uint {
// 	return 100
// }

// // satisfies TestConfigurator interface
// func (suite *PublishSuite) SimulateLinkScores(ctx context.Context) bool {
// 	return true
// }

func (suite *PublishSuite) TestDropmarkToMarkdown() {
	ctx := context.Background()

	clpr := &CommandLineProgressReporter{maxErrors: 10}

	bpc := markdown.NewBasePathConfigurator("test_001")
	linkFactory := link.NewFactory()

	publisher, err := NewMarkdownPublisher(ctx, 25, linkFactory, bpc, clpr, suite)
	suite.Nil(err, "No error should be reported")

	err = publisher.Publish(ctx, "https://shah.dropmark.com/616548.json", clpr, suite)
	suite.Nil(err, "No error should be reported")
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(PublishSuite))
}
