package postgresdb

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/Masterminds/squirrel"

	"github.com/ibanezv/littletwitter/config"
	"github.com/ibanezv/littletwitter/pkg/db"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	defaultMaxPoolSize   = 1
	defaultConnAttempts  = 10
	defaultConnTimeout   = time.Second
	follwersActiveStatus = "active"
)

type PostgresConn struct {
	conn Postgres
}

type Postgres struct {
	maxPoolSize  int
	connAttempts int
	connTimeout  time.Duration

	Builder squirrel.StatementBuilderType
	Pool    *pgxpool.Pool
}

func NewDbEngine(cfg *config.Config) (*PostgresConn, error) {
	pg := Postgres{
		maxPoolSize:  defaultMaxPoolSize,
		connAttempts: defaultConnAttempts,
		connTimeout:  defaultConnTimeout,
	}
	MaxPoolSize(cfg.PG.PoolMax)

	pg.Builder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	poolConfig, err := pgxpool.ParseConfig(cfg.PG.URL)
	if err != nil {
		return nil, fmt.Errorf("postgres - NewPostgres - pgxpool.ParseConfig: %w", err)
	}

	poolConfig.MaxConns = int32(pg.maxPoolSize)
	poolConfig.ConnConfig.Database = cfg.PG.Name
	for pg.connAttempts > 0 {
		pg.Pool, err = pgxpool.ConnectConfig(context.Background(), poolConfig)
		if err == nil {
			break
		}

		log.Printf("Postgres is trying to connect, attempts left: %d", pg.connAttempts)

		time.Sleep(pg.connTimeout)

		pg.connAttempts--
	}

	if err != nil {
		return nil, fmt.Errorf("postgres - NewPostgres - connAttempts == 0: %w", err)
	}

	return &PostgresConn{pg}, nil
}

func (p *PostgresConn) GetTweet(id int64) (*db.DbTweet, error) {
	sql, _, err := p.conn.Builder.
		Select("id, user_id, text, datetime").
		From("tweets").
		Where("id = " + strconv.FormatInt(id, 10)).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("Timeline - GetTimeline - r.Builder: %w", err)
	}

	rows, err := p.conn.Pool.Query(context.Background(), sql)
	if err != nil {
		return nil, fmt.Errorf("Timeline - GetTweet - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	tweet := db.DbTweet{}

	err = rows.Scan(&tweet.Id, &tweet.Text, &tweet.UserId, &tweet.DateTime)
	if err != nil {
		return nil, fmt.Errorf("Timeline - GetTweet - rows.Scan: %w", err)
	}

	return &tweet, nil
}

func (p *PostgresConn) SetTweet(tweet *db.DbTweet) (*db.DbTweet, error) {
	newTweet := db.DbTweet{}
	sql, args, err := p.conn.Builder.
		Insert("tweets").
		Columns("user_id, text, datetime").
		Values(tweet.UserId, tweet.Text, time.Now()).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("PostresDB - insert error: %w", err)
	}

	err = p.conn.Pool.QueryRow(context.Background(), sql+" RETURNING id,user_id,text,datetime", args...).Scan(&newTweet.Id, &newTweet.UserId, &newTweet.Text, &newTweet.DateTime)
	if err != nil {
		return nil, fmt.Errorf("PostresDB - r.Pool.Exec: %w", err)
	}
	return &newTweet, nil
}

func (p *PostgresConn) GetTweetsByUser(userId int64, top uint64) ([]db.DbTweet, error) {
	sql, _, err := p.conn.Builder.
		Select("id, user_id, text, datetime").
		From("tweets").
		Where("user_id = " + strconv.FormatInt(userId, 10)).
		OrderBy("id DESC").
		Limit(top).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("Timeline - GetTimeline - r.Builder: %w", err)
	}

	rows, err := p.conn.Pool.Query(context.Background(), sql)
	if err != nil {
		return nil, fmt.Errorf("Timeline - GetTimeline - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	list := make([]db.DbTweet, 0)

	for rows.Next() {
		tweet := db.DbTweet{}

		err = rows.Scan(&tweet.Id, &tweet.UserId, &tweet.Text, &tweet.DateTime)
		if err != nil {
			return nil, fmt.Errorf("Timeline - GetTimeline - rows.Scan: %w", err)
		}

		list = append(list, tweet)
	}

	return list, nil
}

func (p *PostgresConn) SetFollower(f *db.DbFollower) error {
	sql, args, err := p.conn.Builder.
		Insert("followers").
		Columns("user_id, user_id_followed, status, datetime").
		Values(f.UserId, f.UserFollowedId, follwersActiveStatus, time.Now()).
		ToSql()

	if err != nil {
		return fmt.Errorf("PostresDB - insert error: %w", err)
	}

	_, err = p.conn.Pool.Exec(context.Background(), sql, args...)
	var pgErr *pgconn.PgError
	if err != nil {
		if errors.As(err, &pgErr) {
			return fmt.Errorf("PostresDB - constraintError: %s", pgErr.ConstraintName)
		}
		return fmt.Errorf("PostresDB - r.Pool.Exec: %w", err)
	}

	return nil
}

func (p *PostgresConn) GetFollower(follower *db.DbFollower) ([]db.DbFollower, error) {
	sql, _, err := p.conn.Builder.
		Select("user_id,user_id_followed, status, datetime").
		From("followers").
		Where("user_id= " + strconv.FormatInt(follower.UserId, 10) + " AND user_id_followed=" + strconv.FormatInt(follower.UserFollowedId, 10) + " AND status='active'").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("Followers - GetFollower - r.Builder: %w", err)
	}

	rows, err := p.conn.Pool.Query(context.Background(), sql)
	if err != nil {
		return nil, fmt.Errorf("Followers - GetFollower - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	list := make([]db.DbFollower, 0)
	for rows.Next() {
		var userId int64
		var userIdFollowed int64
		var dateInit time.Time
		var status string
		err = rows.Scan(&userId, &userIdFollowed, &status, &dateInit)
		if err != nil {
			return nil, fmt.Errorf("Followers - GetFollowers - rows.Scan: %w", err)
		}

		list = append(list, db.DbFollower{UserId: userId, UserFollowedId: userIdFollowed, DateInit: dateInit})
	}
	return list, nil
}

func (p *PostgresConn) GetFollowers(userIdFollowed int64) ([]int64, error) {
	sql, _, err := p.conn.Builder.
		Select("user_id").
		From("followers").
		Where("user_id_followed = " + strconv.FormatInt(userIdFollowed, 10)).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("Followers - GetFollowers - r.Builder: %w", err)
	}

	rows, err := p.conn.Pool.Query(context.Background(), sql)
	if err != nil {
		return nil, fmt.Errorf("Followers - GetFollowers - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	list := make([]int64, 0)

	for rows.Next() {
		var userID int64
		err = rows.Scan(&userID)
		if err != nil {
			return nil, fmt.Errorf("Followers - GetFollowers - rows.Scan: %w", err)
		}

		list = append(list, userID)
	}
	return list, nil
}

func (p *PostgresConn) GetFollowing(userIdFollower int64) ([]int64, error) {
	sql, _, err := p.conn.Builder.
		Select("user_id_followed").
		From("followers").
		Where("user_id = " + strconv.FormatInt(userIdFollower, 10) + " AND status='active'").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("Followers - GetFollowing - r.Builder: %w", err)
	}

	rows, err := p.conn.Pool.Query(context.Background(), sql)
	if err != nil {
		return nil, fmt.Errorf("Followers - GetFollowing - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	list := make([]int64, 0)

	for rows.Next() {
		var userID int64
		err = rows.Scan(&userID)
		if err != nil {
			return nil, fmt.Errorf("Followers - GetFollowing - rows.Scan: %w", err)
		}

		list = append(list, userID)
	}
	return list, nil
}

func (p *PostgresConn) Close() {
	if p.conn.Pool != nil {
		p.conn.Pool.Close()
	}
}
