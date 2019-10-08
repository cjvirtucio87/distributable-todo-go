package actors

import (
	"bytes"
	"encoding/json"
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

func (p *httpPeer) AddEntries(e EntryInfo) bool {
	jsonStr, err := json.Marshal(e)

	res, err := http.Post(
		fmt.Sprintf(
			"%s/log/addEntries",
			p.Url(),
		),
		"application/json",
		bytes.NewBuffer(jsonStr),
	)

	if err != nil {
		log.Fatal(err)
	}

	return res.StatusCode == http.StatusOK
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

	sm.HandleFunc(
		"/log/count",
		func(rw http.ResponseWriter, req *http.Request) {
			fmt.Fprintf(
				rw,
				"%d",
				p.basicPeer.LogCount(),
			)
		},
	)

	sm.HandleFunc(
		"/log/addEntries",
		func(rw http.ResponseWriter, req *http.Request) {
			defer req.Body.Close()

			decoder := json.NewDecoder(req.Body)

			var e EntryInfo

			err := decoder.Decode(&e)

			if err != nil {
				rw.Header().Set(
					"Content-Type",
					"application/json",
				)

				rw.WriteHeader(http.StatusBadRequest)

				errPayload := map[string]string{
					"error": err.Error(),
				}

				json.NewEncoder(rw).Encode(errPayload)
			}

			result := p.basicPeer.AddEntries(e)

			if result {
				rw.WriteHeader(http.StatusOK)
			} else {
				rw.Header().Set(
					"Content-Type",
					"application/json",
				)

				rw.WriteHeader(http.StatusBadRequest)

				errPayload := map[string]string{
					"error": "failed to send message",
				}

				json.NewEncoder(rw).Encode(errPayload)
			}
		},
	)

	sm.HandleFunc(
		"/log/send",
		func(rw http.ResponseWriter, req *http.Request) {
			defer req.Body.Close()

			decoder := json.NewDecoder(req.Body)

			var m Message

			err := decoder.Decode(&m)

			if err != nil {
				rw.Header().Set(
					"Content-Type",
					"application/json",
				)

				rw.WriteHeader(http.StatusBadRequest)

				errPayload := map[string]string{
					"error": err.Error(),
				}

				json.NewEncoder(rw).Encode(errPayload)
			}

			result := p.basicPeer.Send(m)

			if result {
				rw.WriteHeader(http.StatusOK)
			} else {
				rw.Header().Set(
					"Content-Type",
					"application/json",
				)

				rw.WriteHeader(http.StatusBadRequest)

				errPayload := map[string]string{
					"error": "failed to send message",
				}

				json.NewEncoder(rw).Encode(errPayload)
			}
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

func (p *httpPeer) LogCount() int {
	res, err := http.Get(
		fmt.Sprintf(
			"%s/log/count",
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

func (p *httpPeer) Send(m Message) bool {
	jsonStr, err := json.Marshal(m)

	res, err := http.Post(
		fmt.Sprintf(
			"%s/log/send",
			p.Url(),
		),
		"application/json",
		bytes.NewBuffer(jsonStr),
	)

	if err != nil {
		log.Fatal(err)
	}

	return res.StatusCode == http.StatusOK
}

func (p *httpPeer) Url() string {
	return fmt.Sprintf(
		"%s://%s:%d",
		p.scheme,
		p.host,
		p.port,
	)
}
