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

type REST struct {
	client  *http.Client
	baseUrl string
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

func (r *REST) AddRecord(rec *storage.Record) (err error) {
	jsonB, err := json.Marshal(rec)
	if err != nil {
		return
	}
	put, err := http.NewRequest("PUT", r.baseUrl+"/transfer", bytes.NewBuffer(jsonB))
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

func (r *REST) GetRecords() (rec []storage.Record, err error) {
	resp, err := r.client.Get(r.baseUrl + "/list")
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

func (r *REST) RemoveRecord(ids ...int) (err error) {
	for _, id := range ids {
		del, errL := http.NewRequest("DELETE", fmt.Sprintf("%s/remove/%d", r.baseUrl, id), nil)
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

func (r *REST) AggregateRecords() (sales []storage.Sale, err error) {
	resp, err := r.client.Get(r.baseUrl + "/dashboard")
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

func (r *REST) Init() (err error) {
	return
}

func NewRestDAO(url string) (r *REST) {
	r = &REST{
		client: &http.Client{
			Timeout: 20 * time.Second,
		},
		baseUrl: url,
	}
	return
}
