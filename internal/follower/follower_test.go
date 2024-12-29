package follower

import (
	"testing"

	"github.com/ibanezv/littletwitter/internal/mocks"
	"github.com/ibanezv/littletwitter/pkg/cache"
	"github.com/ibanezv/littletwitter/pkg/db"
	"github.com/ibanezv/littletwitter/pkg/logger"
	"github.com/ibanezv/littletwitter/pkg/repository"
	"github.com/stretchr/testify/assert"
)

func TestAddFollower(t *testing.T) {
	// given
	var tests = []struct {
		name        string
		follower    FollowerUsers
		mockCache   cache.Cache
		mockDB      db.DBStorage
		errorExpect bool
	}{
		{
			name:        "add sucess",
			follower:    FollowerUsers{UserId: 1, UserIdFollowed: 2},
			mockCache:   mocks.CreateMockCache(),
			mockDB:      mocks.CreateMockDB(),
			errorExpect: false,
		},
		{
			name:        "fail database error",
			follower:    FollowerUsers{UserId: 3, UserIdFollowed: 4},
			mockCache:   mocks.CreateMockCache(),
			mockDB:      mocks.CreateMockDB(),
			errorExpect: true,
		},
		{
			name:        "fail follower already exists error",
			follower:    FollowerUsers{UserId: 3, UserIdFollowed: 4},
			mockCache:   mocks.CreateMockCache(),
			mockDB:      mocks.CreateMockDB(),
			errorExpect: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			l := logger.New("")
			repo := repository.NewRepository(test.mockDB, test.mockCache, l)
			followerServ := NewFollowerService(repo, l)

			// when
			err := followerServ.Set(&test.follower)

			// then
			assert.Equal(t, err != nil, test.errorExpect)
		})
	}
}
