package actors

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type httpPeer struct {
	basicPeer
	server *http.Server
	host   string
	port   string
	scheme string
}

func (p *httpPeer) Init() error {
	sm := http.NewServeMux()

	sm.HandleFunc(
		"/followers/count",
		func(rw http.ResponseWriter, req *http.Request) {
			fmt.Fprintf(
				rw,
				"%d",
				p.basicPeer.PeerCount(),
			)
		},
	)

	p.server = &http.Server{
		Addr: fmt.Sprintf(
			"%s:%s",
			p.host,
			p.port,
		),
		Handler: sm,
	}

	go p.server.ListenAndServe()

	return nil
}

func (p *httpPeer) PeerCount() int {
	res, err := http.Get(
		fmt.Sprintf(
			"%s/followers/count",
			p.Url(),
		),
	)

	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Fatal(err)
	}

	result, err := strconv.Atoi(string(body))

	if err != nil {
		log.Fatal(err)
	}

	return result
}

func (p *httpPeer) Url() string {
	return fmt.Sprintf(
		"%s://%s:%s",
		p.scheme,
		p.host,
		p.port,
	)
}
