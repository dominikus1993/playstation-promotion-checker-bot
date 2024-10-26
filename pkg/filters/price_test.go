package filters

import (
	"context"
	"testing"

	"github.com/dominikus1993/playstation-promotion-checker-bot/pkg/data"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestPriceFilter(t *testing.T) {
	filter := NewPriceFilter(50)
	games := []data.PlaystationStoreGame{data.NewPlaystationStoreGame("test", "test", 10, 20), data.NewPlaystationStoreGame("test", "test", 10, 10), data.NewPlaystationStoreGame("test", "test", 10, 30)}
	stream := lo.SliceToChannel(10, games)

	subject := lo.ChannelToSlice(filter.Filter(context.TODO(), stream))

	assert.Len(t, subject, 2)
}
