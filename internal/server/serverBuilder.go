package server

import (
	"fmt"
	"net/http"
)

type ServerBuilder interface {
	Init() ServerBuilder
	SetAddress(address string) ServerBuilder
	AddFileServer(basePath string) ServerBuilder
	AddGet(route string, handler func(w http.ResponseWriter, r *http.Request)) ServerBuilder
	AddPost(route string, handler func(w http.ResponseWriter, r *http.Request)) ServerBuilder
	AddPut(route string, handler func(w http.ResponseWriter, r *http.Request)) ServerBuilder
	AddDelete(route string, handler func(w http.ResponseWriter, r *http.Request)) ServerBuilder
	Build() (*http.Server, error)
}

type serverBuilder struct {
	mux           *http.ServeMux
	address       string
	isInitialized bool
	error         error
}

func CreateServerBuilder() ServerBuilder {
	return &serverBuilder{}
}

func (sb *serverBuilder) Init() ServerBuilder {
	if sb.error != nil {
		return sb
	}

	if sb.isInitialized {
		sb.error = fmt.Errorf("builder is already initialized")
		return sb
	}

	sb.isInitialized = true
	sb.mux = http.NewServeMux()

	return sb
}

func (sb *serverBuilder) SetAddress(address string) ServerBuilder {
	if sb.error != nil {
		return sb
	}

	if !sb.isInitialized {
		sb.error = fmt.Errorf("builder is not initialized")
		return sb
	}

	sb.address = address

	return sb
}

func (sb *serverBuilder) AddFileServer(basePath string) ServerBuilder {
	if sb.error != nil {
		return sb
	}

	if !sb.isInitialized {
		sb.error = fmt.Errorf("builder is not initialized")
		return sb
	}

	if basePath == "" {
		sb.error = fmt.Errorf("base path cannot be empty")
		return sb
	}

	sb.mux.Handle("/", http.FileServer(http.Dir(basePath)))

	return sb
}

func (sb *serverBuilder) AddGet(route string, handler func(w http.ResponseWriter, r *http.Request)) ServerBuilder {
	if sb.error != nil {
		return sb
	}

	if !sb.isInitialized {
		sb.error = fmt.Errorf("builder is not initialized")
		return sb
	}

	sb.mux.HandleFunc(route, handler)

	return sb
}

func (sb *serverBuilder) AddPost(route string, handler func(w http.ResponseWriter, r *http.Request)) ServerBuilder {
	if sb.error != nil {
		return sb
	}

	if !sb.isInitialized {
		sb.error = fmt.Errorf("builder is not initialized")
		return sb
	}

	sb.mux.HandleFunc(route, handler)

	return sb
}

func (sb *serverBuilder) AddPut(route string, handler func(w http.ResponseWriter, r *http.Request)) ServerBuilder {
	if sb.error != nil {
		return sb
	}

	if !sb.isInitialized {
		sb.error = fmt.Errorf("builder is not initialized")
		return sb
	}

	sb.mux.HandleFunc(route, handler)

	return sb
}

func (sb *serverBuilder) AddDelete(route string, handler func(w http.ResponseWriter, r *http.Request)) ServerBuilder {
	if sb.error != nil {
		return sb
	}

	if !sb.isInitialized {
		sb.error = fmt.Errorf("builder is not initialized")
		return sb
	}

	sb.mux.HandleFunc(route, handler)

	return sb
}

func (sb *serverBuilder) Build() (*http.Server, error) {
	if sb.error != nil {
		return nil, sb.error
	}

	if !sb.isInitialized {
		sb.error = fmt.Errorf("builder is not initialized")
		return nil, sb.error
	}

	server := http.Server{}
	server.Handler = sb.mux
	server.Addr = sb.address

	return &server, nil
}
