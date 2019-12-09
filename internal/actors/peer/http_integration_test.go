// +build integration

package actors

import (
	"cjvirtucio87/distributed-todo-go/internal/rlogging"
	"reflect"
	"testing"
	"time"
)

func TestIntegrationLogCountHttp(t *testing.T) {
	scheme := "http"
	host := "127.0.0.1"
	leaderPort := 8080

	rlogger := rlogging.NewZapLogger()

	leader := NewHttpPeer(
		scheme,
		host,
		leaderPort,
		0,
		WithLogger(rlogger),
	)

	followers := []Peer{}

	for i := 1; i < 3; i++ {
		followers = append(
			followers,
			NewHttpPeer(
				scheme,
				host,
				leaderPort+i,
				i,
				WithLogger(rlogger),
			),
		)
	}

	channels := []chan error{}

	for _, follower := range followers {
		followerChannel := make(chan error)

		channels = append(channels, followerChannel)

		go func() {
			followerChannel <- follower.Init()
		}()

		leader.AddPeer(follower)
	}

	leaderChannel := make(chan error)

	go func() {
		leaderChannel <- leader.Init()
	}()

	timeoutChannel := make(chan error)

	go func() {
		time.Sleep(10 * time.Second)

		timeoutChannel <- nil
	}()

	channels = append(channels, leaderChannel, timeoutChannel)

	// inspired by: https://stackoverflow.com/a/19992525/5665947
	selectCases := make([]reflect.SelectCase, len(channels))

	for i, ch := range channels {
		selectCases[i] = reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(ch),
		}
	}

	_, value, _ := reflect.Select(selectCases)

	if valueInterface := value.Interface(); valueInterface != nil {
		t.Fatalf(valueInterface.(error).Error())
	}

	expectedLogCount := 0

	for _, follower := range followers {
		if actualLogCount, err := follower.LogCount(); err != nil {
			t.Fatalf("failed to retrieve log count due to error, %s\n", err.Error())
		} else if expectedLogCount != actualLogCount {
			t.Fatalf("expected %d, was %d", expectedLogCount, actualLogCount)
		}
	}
}
