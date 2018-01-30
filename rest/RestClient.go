package rest

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/jdevelop/go-coin-ticker/storage"
)

//REST defines the REST service to start on specific URI.
type REST struct {
	client  *http.Client
	baseURL string
}

func errorFromResponse(resp *http.Response) (err error) {
	if resp.StatusCode != 200 {
		str, errLoc := ioutil.ReadAll(resp.Body)
		if errLoc != nil {
			return errLoc
		}
		err = errors.New(string(str))
	}
	return
}

//AddRecord sends the PUT request to the REST service.
func (r *REST) AddRecord(rec *storage.Record) (err error) {
	jsonB, err := json.Marshal(rec)
	if err != nil {
		return
	}
	put, err := http.NewRequest("PUT", r.baseURL+"/transfer", bytes.NewBuffer(jsonB))
	put.Header.Set("Content-Type", "application/json")
	if err != nil {
		return
	}
	resp, err := r.client.Do(put)
	if err != nil {
		return
	}
	err = errorFromResponse(resp)
	return
}

//GetRecords requests the records from the remote REST service.
func (r *REST) GetRecords() (rec []storage.Record, err error) {
	resp, err := r.client.Get(r.baseURL + "/list")
	if err != nil {
		return
	}
	err = errorFromResponse(resp)
	if err == nil {
		bs, errL := ioutil.ReadAll(resp.Body)
		if errL != nil {
			return nil, errL
		}
		err = json.Unmarshal(bs, &rec)
		if err != nil {
			return
		}
	}
	return
}

//RemoveRecord sends DELETE method to remove the records.
func (r *REST) RemoveRecord(ids ...int) (err error) {
	for _, id := range ids {
		del, errL := http.NewRequest("DELETE", fmt.Sprintf("%s/remove/%d", r.baseURL, id), nil)
		if errL != nil {
			return errL
		}
		resp, errL := r.client.Do(del)
		if err != nil {
			return errL
		}
		err = errorFromResponse(resp)
		if err != nil {
			break
		}
	}
	return
}

//AggregateRecords returns the collection of sales grouped by account.
func (r *REST) AggregateRecords() (sales []storage.Sale, err error) {
	resp, err := r.client.Get(r.baseURL + "/dashboard")
	if err != nil {
		return
	}
	err = errorFromResponse(resp)
	var dashboard Dashboard
	if err == nil {
		bs, errL := ioutil.ReadAll(resp.Body)
		if errL != nil {
			return nil, errL
		}
		err = json.Unmarshal(bs, &dashboard)
		if err != nil {
			return
		}
		for _, pr := range dashboard.Symbols {
			sales = append(sales, storage.Sale{
				Account: pr.Symbol,
				Amount:  pr.Value,
			})
		}
	}
	return
}

//Init initializes the REST service as required by the backend. Does nothing by default.
func (r *REST) Init() (err error) {
	return
}

//NewRestDAO creates the REST interface on given URL.
func NewRestDAO(url string) (r *REST) {
	r = &REST{
		client: &http.Client{
			Timeout: 20 * time.Second,
		},
		baseURL: url,
	}
	return
}
