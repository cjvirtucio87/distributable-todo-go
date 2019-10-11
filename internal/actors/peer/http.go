package actors

import (
	"bytes"
	"cjvirtucio87/distributed-todo-go/internal/dto"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const (
	entryIdKey             = "EntryId"
	HeaderContentType      = "Content-Type"
	ContentApplicationJson = "application/json"
)

type httpPeer struct {
	basicPeer
	server *http.Server
	host   string
	port   int
	scheme string
}

func (p *httpPeer) AddEntries(e dto.EntryInfo) bool {
	jsonStr, err := json.Marshal(e)

	res, err := http.Post(
		fmt.Sprintf(
			"%s/log/addEntries",
			p.Url(),
		),
		ContentApplicationJson,
		bytes.NewBuffer(jsonStr),
	)

	if err != nil {
		log.Fatal(err)
	}

	return res.StatusCode == http.StatusOK
}

func (p *httpPeer) Entry(id int) (dto.Entry, bool) {
	e := map[string]int{
		entryIdKey: id,
	}

	jsonStr, err := json.Marshal(e)

	if err != nil {
		log.Fatal(err)
	}

	res, err := http.Post(
		fmt.Sprintf(
			"%s/log/entry",
			p.Url(),
		),
		ContentApplicationJson,
		bytes.NewBuffer(jsonStr),
	)

	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	var result dto.Entry
	ok := true

	err = json.NewDecoder(res.Body).Decode(&result)

	if err != nil {
		log.Fatal(err)
		ok = false
	}

	if &result == nil {
		log.Fatal(
			fmt.Sprintf("unable to retrieve entry for id %d", id),
		)
		ok = false
	}

	return result, ok
}

func (p *httpPeer) Init() error {
	p.basicPeer.Init()

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
		"/log/entry",
		func(rw http.ResponseWriter, req *http.Request) {
			defer req.Body.Close()

			var entryMap map[string]int

			err := json.NewDecoder(req.Body).Decode(&entryMap)

			if err != nil {
				p.respondWithFailure(
					rw,
					err.Error(),
					http.StatusBadRequest,
				)
			}

			entryId := entryMap[entryIdKey]

			entry, ok := p.basicPeer.Entry(entryId)

			if !ok {
				msg := fmt.Sprintf(
					"unable to retrieve entry with id %d\n",
					entryId,
				)

				p.respondWithFailure(
					rw,
					msg,
					http.StatusBadRequest,
				)

				log.Fatal(msg)
			}

			rw.Header().Set(
				HeaderContentType,
				ContentApplicationJson,
			)

			rw.WriteHeader(http.StatusOK)

			err = json.NewEncoder(rw).Encode(entry)

			if err != nil {
				p.respondWithFailure(
					rw,
					err.Error(),
					http.StatusInternalServerError,
				)

				log.Fatal(err)
			}
		},
	)

	sm.HandleFunc(
		"/log/addEntries",
		func(rw http.ResponseWriter, req *http.Request) {
			defer req.Body.Close()

			decoder := json.NewDecoder(req.Body)

			var e dto.EntryInfo

			err := decoder.Decode(&e)

			if err != nil {
				p.respondWithFailure(
					rw,
					err.Error(),
					http.StatusBadRequest,
				)
			}

			result := p.basicPeer.AddEntries(e)

			if result {
				rw.WriteHeader(http.StatusOK)
			} else {
				entryStrings := []string{}

				for _, entry := range e.Entries {
					entryStrings = append(
						entryStrings,
						fmt.Sprintf(
							"%s",
							entry.Command,
						),
					)
				}

				p.respondWithFailure(
					rw,
					fmt.Sprintf(
						"failed to add entries: %s\n",
						strings.Join(entryStrings, ", "),
					),
					http.StatusBadRequest,
				)
			}
		},
	)

	sm.HandleFunc(
		"/log/send",
		func(rw http.ResponseWriter, req *http.Request) {
			defer req.Body.Close()

			decoder := json.NewDecoder(req.Body)

			var m dto.Message

			err := decoder.Decode(&m)

			if err != nil {
				p.respondWithFailure(
					rw,
					err.Error(),
					http.StatusBadRequest,
				)
			}

			result := p.basicPeer.Send(m)

			if result {
				rw.WriteHeader(http.StatusOK)
			} else {
				entryStrings := []string{}

				for _, entry := range m.Entries {
					entryStrings = append(
						entryStrings,
						fmt.Sprintf(
							"%s",
							entry.Command,
						),
					)
				}

				p.respondWithFailure(
					rw,
					fmt.Sprintf(
						"failed to send message with entries, %s\n",
						strings.Join(entryStrings, ", "),
					),
					http.StatusBadRequest,
				)
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

	go func() {
		log.Fatal(p.server.ListenAndServe())
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

func (p *httpPeer) respondWithFailure(rw http.ResponseWriter, msg string, status int) {
	rw.Header().Set(
		HeaderContentType,
		ContentApplicationJson,
	)

	rw.WriteHeader(status)

	errPayload := map[string]string{
		"error": msg,
	}

	err := json.NewEncoder(rw).Encode(errPayload)

	if err != nil {
		log.Fatal(err)
	}
}

func (p *httpPeer) Send(m dto.Message) bool {
	jsonStr, err := json.Marshal(m)

	res, err := http.Post(
		fmt.Sprintf(
			"%s/log/send",
			p.Url(),
		),
		ContentApplicationJson,
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

func (p *httpPeer) Shutdown(ctx context.Context) error {
	return p.server.Shutdown(ctx)
}
