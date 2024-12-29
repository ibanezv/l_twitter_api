package follower

import (
	"errors"

	"github.com/ibanezv/littletwitter/pkg/logger"
	"github.com/ibanezv/littletwitter/pkg/repository"
)

var ErrFollowerAlreadyExists = errors.New("follower already exists")
var ErrFollowerNotValid = errors.New("follower not valid")

type FollowerService interface {
	Set(*FollowerUsers) error
}

type Follower struct {
	repository repository.Repo
	log        logger.Interface
}

func NewFollowerService(repo repository.Repo, l logger.Interface) FollowerService {
	return &Follower{repo, l}
}

func (f *Follower) Set(follower *FollowerUsers) error {
	if follower.UserId == follower.UserIdFollowed {
		return ErrFollowerNotValid
	}

	flw, err := f.repository.GetFollower(FollowerToDB(follower))
	if err != nil {
		f.log.Error("error in repository getting follower %+v", err)
		return err
	}

	if len(flw) > 0 {
		return ErrFollowerAlreadyExists
	}

	err = f.repository.SetFollower(FollowerToDB(follower))
	if err != nil {
		f.log.Error("error saving follower: %+v", err)
		return err
	}

	return nil
}
