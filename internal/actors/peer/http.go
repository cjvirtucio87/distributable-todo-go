package actors

import (
	"bytes"
	"cjvirtucio87/distributed-todo-go/internal/dto"
	"cjvirtucio87/distributed-todo-go/internal/rlogging"
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

func (p *httpPeer) AddEntries(e dto.EntryInfo) (bool, error) {
	if jsonStr, err := json.Marshal(e); err != nil {
		return false, err
	} else if res, err := http.Post(
		fmt.Sprintf(
			"%s/log/addEntries",
			p.Url(),
		),
		ContentApplicationJson,
		bytes.NewBuffer(jsonStr),
	); err != nil {
		return false, err
	} else {
		return (res.StatusCode == http.StatusOK), nil
	}
}

func (p *httpPeer) Entry(id int) (dto.Entry, error) {
	var entry dto.Entry
	e := map[string]int{
		entryIdKey: id,
	}

	if jsonStr, err := json.Marshal(e); err != nil {
		return entry, err
	} else {
		if res, err := http.Post(
			fmt.Sprintf(
				"%s/log/entry",
				p.Url(),
			),
			ContentApplicationJson,
			bytes.NewBuffer(jsonStr),
		); err != nil {
			return entry, err
		} else {
			defer res.Body.Close()

			return entry, json.NewDecoder(res.Body).Decode(&entry)
		}
	}
}

func (p *httpPeer) Init() error {
	p.basicPeer.Init()

	sm := http.NewServeMux()

	sm.HandleFunc(
		"/followers/count",
		func(rw http.ResponseWriter, req *http.Request) {
			if result, err := p.basicPeer.PeerCount(); err != nil {
				p.respondWithFailure(
					rw,
					err.Error(),
					http.StatusBadRequest,
				)
			} else {
				fmt.Fprintf(
					rw,
					"%d",
					result,
				)
			}
		},
	)

	sm.HandleFunc(
		"/log/count",
		func(rw http.ResponseWriter, req *http.Request) {
			if result, err := p.basicPeer.LogCount(); err != nil {
				p.respondWithFailure(
					rw,
					err.Error(),
					http.StatusBadRequest,
				)
			} else {
				fmt.Fprintf(
					rw,
					"%d",
					result,
				)
			}
		},
	)

	sm.HandleFunc(
		"/log/entry",
		func(rw http.ResponseWriter, req *http.Request) {
			defer req.Body.Close()

			var entryMap map[string]int

			if err := json.NewDecoder(req.Body).Decode(&entryMap); err != nil {
				p.respondWithFailure(
					rw,
					err.Error(),
					http.StatusBadRequest,
				)
			}

			entryId := entryMap[entryIdKey]

			if entry, err := p.basicPeer.Entry(entryId); err != nil {
				msg := fmt.Sprintf(
					"unable to retrieve entry with id %d\n",
					entryId,
				)

				p.respondWithFailure(
					rw,
					msg,
					http.StatusBadRequest,
				)

				p.basicPeer.rlogger.Errorf(msg)
			} else {
				rw.Header().Set(
					HeaderContentType,
					ContentApplicationJson,
				)

				rw.WriteHeader(http.StatusOK)

				if err := json.NewEncoder(rw).Encode(entry); err != nil {
					msg := err.Error()

					p.respondWithFailure(
						rw,
						msg,
						http.StatusInternalServerError,
					)

					p.basicPeer.rlogger.Errorf(msg)
				}
			}
		},
	)

	sm.HandleFunc(
		"/log/addEntries",
		func(rw http.ResponseWriter, req *http.Request) {
			defer req.Body.Close()

			decoder := json.NewDecoder(req.Body)

			var e dto.EntryInfo

			if err := decoder.Decode(&e); err != nil {
				p.respondWithFailure(
					rw,
					err.Error(),
					http.StatusBadRequest,
				)
			}

			if result, err := p.basicPeer.AddEntries(e); err != nil {
				p.respondWithFailure(
					rw,
					fmt.Sprintf(
						"error adding entries: %s\n",
						err,
					),
					http.StatusBadRequest,
				)
			} else if result {
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

			if err := decoder.Decode(&m); err != nil {
				p.respondWithFailure(
					rw,
					err.Error(),
					http.StatusBadRequest,
				)
			}

			if result, err := p.basicPeer.Send(m); err != nil {
				p.respondWithFailure(
					rw,
					fmt.Sprintf(
						"error attempting to send message, \n%s",
						err.Error(),
					),
					http.StatusBadRequest,
				)
			} else if result {
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
		ErrorLog: log.New(
			rlogging.NewWriterLogger(p.basicPeer.rlogger),
			"",
			0,
		),
		Handler: sm,
	}

	return p.server.ListenAndServe()
}

func (p *httpPeer) LogCount() (int, error) {
	if res, err := http.Get(
		fmt.Sprintf(
			"%s/log/count",
			p.Url(),
		),
	); err != nil {
		return 0, err
	} else {
		defer res.Body.Close()

		if body, err := ioutil.ReadAll(res.Body); err != nil {
			return 0, err
		} else if result, err := strconv.Atoi(string(body)); err != nil {
			return 0, err
		} else {
			return result, nil
		}
	}
}

func (p *httpPeer) PeerCount() (int, error) {
	if res, err := http.Get(
		fmt.Sprintf(
			"%s/followers/count",
			p.Url(),
		),
	); err != nil {
		return 0, err
	} else {
		defer res.Body.Close()

		if body, err := ioutil.ReadAll(res.Body); err != nil {
			return 0, err
		} else if result, err := strconv.Atoi(string(body)); err != nil {
			return 0, err
		} else {
			return result, nil
		}
	}
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

	if err := json.NewEncoder(rw).Encode(errPayload); err != nil {
		p.basicPeer.rlogger.Errorf(err.Error())
	}
}

func (p *httpPeer) Send(m dto.Message) (bool, error) {
	if jsonStr, err := json.Marshal(m); err != nil {
		return false, err
	} else if res, err := http.Post(
		fmt.Sprintf(
			"%s/log/send",
			p.Url(),
		),
		ContentApplicationJson,
		bytes.NewBuffer(jsonStr),
	); err != nil {
		return false, err
	} else {
		return (res.StatusCode == http.StatusOK), nil
	}
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
