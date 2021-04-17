package integration

import (
	"fmt"
	"net/http"
	"testing"
	"encoding/json"
	"log"

	"github.com/stretchr/testify/assert"
	e "github.com/yassinekhaliqui/go-rest-service/e2e_test"
	h "github.com/yassinekhaliqui/go-rest-service/pkg/http"
	"github.com/yassinekhaliqui/go-rest-service/pkg/util"
	"github.com/yassinekhaliqui/go-rest-service/internal/model"
)

func Test_GroupGet_GroupExists(t *testing.T) {
	// create group
	groupName := util.RandStringBytes(32)
	payload := `{"name":"` + groupName + `"}`
	statusCode, err := h.SendPostRequest(e.URL, "/groups", payload)

	assert.Nil(t, err)
	assert.Equal(t, 201, statusCode)

	// get group
	r, err := http.Get(fmt.Sprintf("%s/groups/%s", e.URL, groupName))

	assert.Nil(t, err)
	assert.Equal(t, 200, r.StatusCode)
}

func Test_GroupGet_GroupWithUsers(t *testing.T) {
	// create group
	groupName := util.RandStringBytes(32)
	payload := `{"name":"` + groupName + `"}`
	statusCode, err := h.SendPostRequest(e.URL, "/groups", payload)

	assert.Nil(t, err)
	assert.Equal(t, 201, statusCode)

	// create user
	randStr := util.RandStringBytes(32)
	payload = `{"first_name":"` + randStr + `", "last_name":"` + randStr + `", "userid":"` + randStr + `", "groups":["` + groupName + `"]}`
	statusCode, err = h.SendPostRequest(e.URL, "/users", payload)

	assert.Nil(t, err)
	assert.Equal(t, 201, statusCode)

	// get group
	r, err := http.Get(fmt.Sprintf("%s/groups/%s", e.URL, groupName))
	var restGroupMembers model.RestGroupMembers
	if err := json.NewDecoder(r.Body).Decode(&restGroupMembers); err != nil {
		log.Fatal(err)
		return
	}
	defer r.Body.Close()

	// check member is in the group
	assert.Nil(t, err)
	assert.Equal(t, 200, r.StatusCode)
	assert.Equal(t, &[]string{randStr}, restGroupMembers.UserIds)
}

func Test_GroupGet_GroupDoesNotExist(t *testing.T) {
	groupName := util.RandStringBytes(32)

	r, err := http.Get(fmt.Sprintf("%s/groups/%s", e.URL, groupName))

	assert.Nil(t, err)
	assert.Equal(t, 404, r.StatusCode)
}

func Test_GroupPost_GroupCreated(t *testing.T) {
	randStr := util.RandStringBytes(32)
	payload := `{"name":"` + randStr + `"}`
	statusCode, err := h.SendPostRequest(e.URL, "/groups", payload)

	assert.Nil(t, err)
	assert.Equal(t, 201, statusCode)
}

func Test_GroupPost_GroupAlreadyExists(t *testing.T) {
	randStr := util.RandStringBytes(32)
	payload := `{"name":"` + randStr + `"}`
	statusCode, err := h.SendPostRequest(e.URL, "/groups", payload)

	assert.Nil(t, err)
	assert.Equal(t, 201, statusCode)

	payload = `{"name":"` + randStr + `"}`
	statusCode, err = h.SendPostRequest(e.URL, "/groups", payload)

	assert.Nil(t, err)
	assert.Equal(t, 400, statusCode)
}

func Test_GroupPost_InvalidPayload(t *testing.T) {
	payload := `{"name":""}`
	statusCode, err := h.SendPostRequest(e.URL, "/groups", payload)

	assert.Nil(t, err)
	assert.Equal(t, 400, statusCode)
}

func Test_GroupPut_UpdateGroup(t *testing.T) {
	// create group w/out members
	groupName := util.RandStringBytes(32)
	payload := `{"name":"` + groupName + `"}`
	statusCode, err := h.SendPostRequest(e.URL, "/groups", payload)

	assert.Nil(t, err)
	assert.Equal(t, 201, statusCode)

	// create user w/out groups
	randStr := util.RandStringBytes(32)
	payload = `{"first_name":"` + randStr + `", "last_name":"` + randStr + `", "userid":"` + randStr + `", "groups":null}`
	statusCode, err = h.SendPostRequest(e.URL, "/users", payload)

	assert.Nil(t, err)
	assert.Equal(t, 201, statusCode)

	// update group to add user
	statusCode, err = h.SendPutRequest(e.URL, "/groups", groupName, `{"userids":["` + randStr +`"]}`)

	assert.Nil(t, err)
	assert.Equal(t, 200, statusCode)

	// get group
	r, err := http.Get(fmt.Sprintf("%s/groups/%s", e.URL, groupName))
	var restGroupMembers model.RestGroupMembers
	if err := json.NewDecoder(r.Body).Decode(&restGroupMembers); err != nil {
		log.Fatal(err)
		return
	}
	defer r.Body.Close()

	// check member is in the group
	assert.Nil(t, err)
	assert.Equal(t, 200, r.StatusCode)
	assert.Equal(t, &[]string{randStr}, restGroupMembers.UserIds)
}

func Test_GroupPut_GroupDoesNotExist(t *testing.T) {
	groupName := util.RandStringBytes(32)

	statusCode, err := h.SendPutRequest(e.URL, "/groups", groupName, `{"userids":["lex"]}`)

	assert.Nil(t, err)
	assert.Equal(t, 404, statusCode)
}

func Test_GroupDel_GroupExist(t *testing.T) {
	// create group w/out members
	groupName := util.RandStringBytes(32)
	payload := `{"name":"` + groupName + `"}`
	statusCode, err := h.SendPostRequest(e.URL, "/groups", payload)

	assert.Nil(t, err)
	assert.Equal(t, 201, statusCode)

	// delete group
	statusCode, err = h.SendDelRequest(e.URL, "/groups", groupName)

	assert.Nil(t, err)
	assert.Equal(t, 200, statusCode)

	// try to get group
	r, err := http.Get(fmt.Sprintf("%s/groups/%s", e.URL, groupName))
	var restGroupMembers model.RestGroupMembers
	if err := json.NewDecoder(r.Body).Decode(&restGroupMembers); err != nil {
		log.Fatal(err)
		return
	}
	defer r.Body.Close()

	// check member is in the group
	assert.Nil(t, err)
	assert.Equal(t, 404, r.StatusCode)
}

func Test_GroupDel_GroupDoesNotExist(t *testing.T) {
	groupName := util.RandStringBytes(32)

	statusCode, err := h.SendDelRequest(e.URL, "/groups", groupName)

	assert.Nil(t, err)
	assert.Equal(t, 404, statusCode)
}

func Test_GroupDel_GroupWithUsers(t *testing.T) {
	// create group w/out members
	groupName := util.RandStringBytes(32)
	payload := `{"name":"` + groupName + `"}`
	statusCode, err := h.SendPostRequest(e.URL, "/groups", payload)

	assert.Nil(t, err)
	assert.Equal(t, 201, statusCode)

	// create user with group
	randStr := util.RandStringBytes(32)
	payload = `{"first_name":"` + randStr + `", "last_name":"` + randStr + `", "userid":"` + randStr + `", "groups":["` + groupName + `"]}`
	statusCode, err = h.SendPostRequest(e.URL, "/users", payload)

	assert.Nil(t, err)
	assert.Equal(t, 201, statusCode)

	// delete group
	statusCode, err = h.SendDelRequest(e.URL, "/groups", groupName)

	assert.Nil(t, err)
	assert.Equal(t, 200, statusCode)

	// try to get group
	r, err := http.Get(fmt.Sprintf("%s/groups/%s", e.URL, groupName))
	var restGroupMembers model.RestGroupMembers
	if err := json.NewDecoder(r.Body).Decode(&restGroupMembers); err != nil {
		log.Fatal(err)
		return
	}
	defer r.Body.Close()

	// check member is in the group
	assert.Nil(t, err)
	assert.Equal(t, 404, r.StatusCode)
}
