package rhttp

import "net/http"

type basicClient struct {
	http.Client
}
