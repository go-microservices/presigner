package server

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/minodisk/presigner/options"
	"github.com/minodisk/presigner/publisher"
)

func Serve(o options.Options) (err error) {
	http.Handle("/", Index{o})
	err = http.ListenAndServe(fmt.Sprintf(":%d", o.Port), nil)
	return
}

type Index struct {
	Options options.Options
}

func (i Index) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	resp, err := i.process(r.Method, r.Body)
	if err != nil {
		if coder, ok := err.(Coder); ok {
			resp = NewErrorResp(coder.Code(), err)
		} else {
			resp = NewErrorResp(500, err)
		}
	}

	header := w.Header()
	// header.Set("Access-Control-Allow-Origin", "*")
	// header.Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	// header.Set("Access-Control-Allow-Headers", "Origin, Content-Type")

	if resp.Body == nil {
		w.WriteHeader(resp.Code())
		return
	}

	header.Set("Content-Type", "application/json")

	b, err := json.Marshal(resp.Body)
	if err != nil {
		log.Printf("fail to marshal JSON: %+v", err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(resp.Code())

	if _, err = w.Write(b); err != nil {
		log.Printf("fail to write response body: %+v", err)
		return
	}

	log.Printf("response: %v", string(b))
}

func (i Index) process(method string, body io.ReadCloser) (*Resp, error) {
	switch method {
	default:
		return nil, NewMethodNotAllowed(method)
	case "GET", "OPTIONS":
		return NewResp(http.StatusOK, nil), nil
	case "POST":
		b, err := ioutil.ReadAll(body)
		if err != nil {
			return nil, NewBadRequest(err)
		}
		var pub publisher.Publisher
		if err = json.Unmarshal(b, &pub); err != nil {
			return nil, NewBadRequest(err)
		}
		res, err := pub.Publish(i.Options)
		if err != nil {
			return nil, NewBadRequest(err)
		}
		return NewResp(http.StatusOK, res), nil
	}
}
