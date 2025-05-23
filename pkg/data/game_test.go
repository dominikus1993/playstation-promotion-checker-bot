package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPErcentageCalculation(t *testing.T) {
	data := []struct {
		game     PlaystationStoreGame
		expected float64
	}{
		{game: NewPlaystationStoreGame("", "", 10.0, 20.0), expected: 50},
		{game: NewPlaystationStoreGame("", "", 10.0, 10.0), expected: 0},
		{game: NewPlaystationStoreGame("", "", 10.0, 100.0), expected: 90},
	}

	for _, d := range data {
		result := d.game.CalculatePromotionPercentage()
		assert.Equal(t, d.expected, result)
	}
}

func FuzzFormatPrice(f *testing.F) {
	data := []struct {
		price    float64
		expected string
	}{
		{price: 50.0, expected: "50.00"},
		{price: 0, expected: "0.00"},
		{price: 90.0, expected: "90.00"},
		{price: 50, expected: "50.00"},
		{price: 21.37, expected: "21.37"},
	}

	for _, d := range data {
		f.Add(d.price, d.expected)
	}

	f.Fuzz(func(t *testing.T, price float64, expected string) {
		result := formatPrice(price)
		assert.Equal(t, expected, result)
	})
}

func TestFormatPercentage(t *testing.T) {
	data := []struct {
		game     PlaystationStoreGame
		expected string
	}{
		{game: NewPlaystationStoreGame("", "", 10.0, 20.0), expected: "50.00"},
		{game: NewPlaystationStoreGame("", "", 10.0, 10.0), expected: "0.00"},
		{game: NewPlaystationStoreGame("", "", 10.0, 100.0), expected: "90.00"},
		{game: NewPlaystationStoreGame("", "", 134.99, 269.99), expected: "50.00"},
	}

	for _, d := range data {
		result := d.game.FormatPromotionPercentage()
		assert.Equal(t, d.expected, result)
	}
}

func TestGetLink(t *testing.T) {

	link := "/pl-pl/product/EP1004-CUSA00411_00-PREMIUMPACKOGGW1"
	game := NewPlaystationStoreGame("", link, 10.0, 20.0)
	subject, err := game.GetLink()

	assert.NoError(t, err)

	assert.Equal(t, "https://store.playstation.com/pl-pl/product/EP1004-CUSA00411_00-PREMIUMPACKOGGW1", subject.String())
}

func FuzzJoinPath(f *testing.F) {

	f.Add("", "", true)
	f.Add("/pl-pl/product/EP1004-CUSA00411_00-PREMIUMPACKOGGW1", "https://store.playstation.com/pl-pl/product/EP1004-CUSA00411_00-PREMIUMPACKOGGW1", false)
	f.Add("https://store.playstation.com/pl-pl/product/EP1004-CUSA00411_00-PREMIUMPACKOGGW1", "https://store.playstation.com/pl-pl/product/EP1004-CUSA00411_00-PREMIUMPACKOGGW1", false)
	f.Fuzz(func(t *testing.T, url string, expectedUrl string, isError bool) {
		result, err := joinPath(url)
		assert.Equal(t, expectedUrl, result)
		if !isError {
			assert.NoError(t, err)
		} else {
			assert.NotNil(t, err)
			assert.ErrorIs(t, err, errUrlIsEmptyErr)
		}
	})
}
