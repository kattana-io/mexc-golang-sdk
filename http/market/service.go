package mexchttpmarket

import (
	"context"
	"errors"
	"fmt"
	mexchttp "github.com/kattana-io/mexc-golang-sdk/http"
	"strconv"
	"time"
)

type Service struct {
	client                    *mexchttp.Client
	syncTimeDeltaMilliSeconds int64
}

func New(ctx context.Context, client *mexchttp.Client) (*Service, error) {
	s := &Service{client: client}

	err := s.syncServerTime(ctx)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Service) syncServerTime(ctx context.Context) error {
	r, err := s.Time(ctx)
	if err != nil {
		return fmt.Errorf("get server time: %w", err)
	}

	if r.ServerTime == 0 {
		return errors.New("server time is empty")
	}

	s.syncTimeDeltaMilliSeconds = time.Now().UnixMilli() - r.ServerTime

	return nil
}

func (s *Service) getTimestamp() string {
	return strconv.FormatInt(time.Now().UnixMilli()-s.syncTimeDeltaMilliSeconds, 10)
}

// SyncServerTime синхронизирует время сервера.
