package gitlab

import (
	"fmt"
	"github.com/dghubble/sling"
	"net/http"
)

type Query struct {
	// Page number (default: 1).
	Page int `url:"page,omitempty"`
	// Number of items to list per page (default: 20, max: 100).
	PerPage int `url:"per_page,omitempty"`
}

func paginate[Type interface{}](req *sling.Sling, list []Type) ([]Type, error) {
	query := &Query{
		Page: 1,
	}
	req.QueryStruct(query)
	for {
		var page = make([]Type, 0)
		var errorResponse ErrorResponse
		resp, err := req.Receive(&page, &errorResponse)
		err = handleHttpError(resp, err, errorResponse)
		if err != nil {
			return nil, err
		}
		if len(page) == 0 {
			break
		}
		list = append(list, page...)
		query.Page++
	}
	return list, nil
}

func handleHttpError(response *http.Response, err error, errorResponse ErrorResponse) error {
	if err != nil {
		return err
	}
	if response.StatusCode > 399 {
		return fmt.Errorf("received status code [%d] on url: %s\nerror: %s\nmessage: %s", response.StatusCode, response.Request.URL, errorResponse.Error, errorResponse.Message)
	}
	return nil
}
