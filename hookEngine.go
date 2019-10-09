package captainhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

//TODO: Document
type HookEngine struct {
	log         *log.Logger
	endpointSvc EndpointService
}

//TODO: Document
func NewHookEngine(log *log.Logger, ec *EndpointService) *HookEngine {
	return &HookEngine{log, *ec}
}

//TODO: Document
func (h *HookEngine) Hook(w http.ResponseWriter, r *http.Request) {

	var name string
	var ok bool

	// Extract variables from request.
	vars := mux.Vars(r)

	if name, ok = vars["id"]; !ok {
		w.WriteHeader(http.StatusNotFound)
		h.log.Println("no id provided.")
		return
	}
	h.log.Println("processing webhook:", name)

	var endpoint *Endpoint

	endpoint, err := h.endpointSvc.Endpoint(name)

	if err != nil {
		h.log.Println("error getting endpoint", name)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	b, err := ioutil.ReadAll(r.Body)

	if err != nil {
		h.log.Println("unable to get body from request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body := fmt.Sprintf("%s", b)
	decodedBody, err := url.QueryUnescape(body)

	if err != nil {
		h.log.Fatal("unable to URL decode body")
	}

	rules, err := endpoint.GetRules()

	if err != nil {
		h.log.Println("unable to enumerate rules, endpoint", endpoint.Name)
		w.WriteHeader(http.StatusInternalServerError)
	}

	dataBag := make(map[string]interface{})
	err = json.Unmarshal([]byte(decodedBody), &dataBag)

	if err != nil {
		h.log.Println("unable to unmarshal json")
	}

	var request bytes.Buffer

	for _, r := range rules {

		AssignFunction(&r)

		err := r.Function(&request, dataBag, &r)

		if err != nil {
			h.log.Println(r, "failed to execute template.", err)
			continue
		}
		h.log.Println("rendered template: ", request.String())
		h.log.Println("forwarding to", r.Destination)

		_, err = http.Post(r.Destination, "application/json", &request)

		if err != nil {
			h.log.Println("post request to", r.Destination, "failed.")
		}
		request.Reset()
	}

	w.WriteHeader(http.StatusNoContent)
}
