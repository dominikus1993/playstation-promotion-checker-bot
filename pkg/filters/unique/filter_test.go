package unique

import (
	"testing"

	"github.com/dominikus1993/playstation-promotion-checker-bot/pkg/data"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestUnique(t *testing.T) {
	filter := NewUniqeFilter()
	games := []data.PlaystationStoreGame{data.NewPlaystationStoreGame("test", "test", 10, 20), data.NewPlaystationStoreGame("test", "test", 10, 10), data.NewPlaystationStoreGame("test", "test", 10, 30)}
	stream := lo.SliceToChannel(10, games)

	subject := filter.Filter(t.Context(), stream)

	slice := lo.ChannelToSlice(subject)

	assert.Len(t, slice, 1)
}
