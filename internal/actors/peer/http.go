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
	port   int
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
			"%s:%d",
			p.host,
			p.port,
		),
		Handler: sm,
	}

	// inspired by:
	// https://github.com/openshift/origin/blob/67ef8497bbcd4f7ea8bc4e0e2daa75ba0c613f20/examples/hello-openshift/hello_openshift.go
	go func() {
		err := p.server.ListenAndServe()

		if err != nil {
			panic("ListenAndServe: " + err.Error())
		}
	}()

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
		"%s://%s:%d",
		p.scheme,
		p.host,
		p.port,
	)
}
