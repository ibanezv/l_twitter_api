package tweet

import (
	"testing"

	"github.com/ibanezv/littletwitter/internal/mocks"
	"github.com/ibanezv/littletwitter/pkg/logger"
	"github.com/ibanezv/littletwitter/pkg/repository"
	"github.com/ibanezv/littletwitter/settings"
	"github.com/stretchr/testify/assert"
)

func TestSaveTest(t *testing.T) {
	// given
	var tests = []struct {
		name          string
		newTweet      Tweet
		mockCache     *mocks.MockCache
		mockDB        *mocks.MockDBStorage
		expectedError bool
	}{
		{
			name:          "save sucessful",
			newTweet:      Tweet{UserId: 1, Text: "testing"},
			mockCache:     mocks.CreateMockCache(),
			mockDB:        mocks.CreateMockDB(),
			expectedError: false,
		},
		{
			name:          "failed because of DBError",
			newTweet:      Tweet{UserId: 2, Text: "testing"},
			mockCache:     mocks.CreateMockCache(),
			mockDB:        mocks.CreateMockDB(),
			expectedError: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			l := logger.New("test")
			repo := repository.NewRepository(test.mockDB, test.mockCache, l)
			appSettings, _ := settings.NewAppSettings()
			tweetServ := NewTwitter(repo, l, appSettings)

			// when
			tt, err := tweetServ.Add(&test.newTweet)

			// then
			assert.Equal(t, err != nil, test.expectedError)
			if !test.expectedError {
				assert.Equal(t, test.newTweet.UserId, tt.UserId)
				assert.Equal(t, test.newTweet.Text, tt.Text)
			}
		})
	}

}
