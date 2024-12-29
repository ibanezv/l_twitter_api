package timeline

import (
	"testing"

	"github.com/ibanezv/littletwitter/internal/mocks"
	"github.com/ibanezv/littletwitter/pkg/logger"
	"github.com/ibanezv/littletwitter/pkg/repository"
	"github.com/ibanezv/littletwitter/settings"
	"github.com/stretchr/testify/assert"
)

func TestGettingTimeline(t *testing.T) {
	// given
	var tests = []struct {
		name          string
		userId        int64
		dbMock        *mocks.MockDBStorage
		dbCache       *mocks.MockCache
		errorExpected bool
		fromCache     int
	}{
		{
			name:          "getting timeline from cache success",
			userId:        1,
			dbMock:        mocks.CreateMockDB(),
			dbCache:       mocks.CreateMockCache(),
			errorExpected: false,
			fromCache:     1,
		},
		{
			name:          "getting timeline from db success",
			userId:        3,
			dbMock:        mocks.CreateMockDB(),
			dbCache:       mocks.CreateMockCache(),
			errorExpected: false,
			fromCache:     0,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			l := logger.New("test")
			appSettings, _ := settings.NewAppSettings()
			repo := repository.NewRepository(test.dbMock, test.dbCache, l)
			timelineServ := NewTimeLine(repo, l, appSettings)

			// when
			tl, err := timelineServ.Get(test.userId)

			// then
			assert.Equal(t, err != nil, test.errorExpected)
			assert.Equal(t, tl.UserId, test.userId)
			assert.Equal(t, test.fromCache, test.dbCache.CountGetTimeline)

		})
	}

}
