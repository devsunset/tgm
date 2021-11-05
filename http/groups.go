package http

import (
	"net/http"
	"sort"
	"tgm/groups"
)

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

var groupDeleteHandler = withAdmin(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	return http.StatusOK, nil
})

var groupPostHandler = withAdmin(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	return http.StatusOK, nil
})
