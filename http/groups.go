package http

import (
	"encoding/json"
	"net/http"
	"sort"
	"tgm/errors"
	"tgm/groups"
)

type requestData struct {
	What  string   `json:"what"`  // Answer to: what data type?
	Which []string `json:"which"` // Answer to: which fields?
	Data  string   `json:"data"`  // Answer to: which fields?
}

func getParameter(_ http.ResponseWriter, r *http.Request) (*requestData, error) {
	if r.Body == nil {
		return nil, errors.ErrEmptyRequest
	}

	req := &requestData{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		return nil, err
	}

	if req.What != "group" {
		return nil, errors.ErrInvalidDataType
	}

	return req, nil
}

var groupsGetHandler = withAdmin(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	groups, err := groups.Gets()
	if err != nil {
		return http.StatusInternalServerError, err
	}
	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Gid < groups[j].Gid
	})
	return renderJSON(w, r, groups)
})

var groupPostHandler = withAdmin(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	req, err := getParameter(w, r)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	if req.Data == "" {
		return http.StatusInternalServerError, err
	}

	err = groups.Save(req.Data)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, err
})

var groupDeleteHandler = withAdmin(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	req, err := getParameter(w, r)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	if req.Data == "" {
		return http.StatusInternalServerError, err
	}

	err = groups.Delete(req.Data)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
})
