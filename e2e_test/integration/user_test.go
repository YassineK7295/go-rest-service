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

func Test_UserGet_UserExistsWithNoGroup(t *testing.T) {
	// create user
	randStr := util.RandStringBytes(32)
	payload := `{"first_name":"` + randStr + `", "last_name":"` + randStr + `", "userid":"` + randStr + `", "groups":null}`
	statusCode, err := h.SendPostRequest(e.URL, "/users", payload)

	assert.Nil(t, err)
	assert.Equal(t, 201, statusCode)

	// retrieve user
	var user model.RestUser
	r, err := http.Get(fmt.Sprintf("%s/users/%s", e.URL, randStr))
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Fatal(err)
		return
	}
	defer r.Body.Close()

	assert.Nil(t, err)
	assert.Equal(t, 200, r.StatusCode)

	// check fields
	assert.Equal(t, randStr, user.FirstName)
	assert.Equal(t, randStr, user.LastName)
	assert.Equal(t, randStr, user.UserId)
	assert.Equal(t, &[]string{}, user.Groups)
}

func Test_UserGet_UserExistsWithGroup(t *testing.T) {
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

	// retrieve user
	var user model.RestUser
	r, err := http.Get(fmt.Sprintf("%s/users/%s", e.URL, randStr))
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Fatal(err)
		return
	}
	defer r.Body.Close()

	assert.Nil(t, err)
	assert.Equal(t, 200, r.StatusCode)

	// check fields
	assert.Equal(t, randStr, user.FirstName)
	assert.Equal(t, randStr, user.LastName)
	assert.Equal(t, randStr, user.UserId)
	assert.Equal(t, &[]string{groupName}, user.Groups)
}

func Test_UserGet_UserDoesNotExists(t *testing.T) {
	userId := util.RandStringBytes(32)

	r, err := http.Get(fmt.Sprintf("%s/users/%s", e.URL, userId))

	assert.Nil(t, err)
	assert.Equal(t, 404, r.StatusCode)
}

func Test_UserPost_WithoutGroup(t *testing.T) {
	// create user
	randStr := util.RandStringBytes(32)
	payload := `{"first_name":"` + randStr + `", "last_name":"` + randStr + `", "userid":"` + randStr + `", "groups":null}`
	statusCode, err := h.SendPostRequest(e.URL, "/users", payload)

	assert.Nil(t, err)
	assert.Equal(t, 201, statusCode)

	// retrieve user
	var user model.RestUser
	r, err := http.Get(fmt.Sprintf("%s/users/%s", e.URL, randStr))
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Fatal(err)
		return
	}
	defer r.Body.Close()

	assert.Nil(t, err)
	assert.Equal(t, 200, r.StatusCode)

	// check fields
	assert.Equal(t, randStr, user.FirstName)
	assert.Equal(t, randStr, user.LastName)
	assert.Equal(t, randStr, user.UserId)
	assert.Equal(t, &[]string{}, user.Groups)
}

func Test_UserPost_WithGroup(t *testing.T) {
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

	// retrieve user
	var user model.RestUser
	r, err := http.Get(fmt.Sprintf("%s/users/%s", e.URL, randStr))
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Fatal(err)
		return
	}
	defer r.Body.Close()

	assert.Nil(t, err)
	assert.Equal(t, 200, r.StatusCode)

	// check fields
	assert.Equal(t, randStr, user.FirstName)
	assert.Equal(t, randStr, user.LastName)
	assert.Equal(t, randStr, user.UserId)
	assert.Equal(t, &[]string{groupName}, user.Groups)
}

func Test_UserPost_InvalidPayload(t *testing.T) {
	payload := `{"first_name":"", "last_name":"", "userid":"", "groups":[""]}`
	statusCode, err := h.SendPostRequest(e.URL, "/users", payload)

	assert.Nil(t, err)
	assert.Equal(t, 400, statusCode)
}

func Test_UserPost_UserExists(t *testing.T) {
	randStr := util.RandStringBytes(32)
	payload := `{"first_name":"` + randStr + `", "last_name":"` + randStr + `", "userid":"` + randStr + `", "groups":null}`
	statusCode, err := h.SendPostRequest(e.URL, "/users", payload)

	assert.Nil(t, err)
	assert.Equal(t, 201, statusCode)

	payload = `{"first_name":"` + randStr + `", "last_name":"` + randStr + `", "userid":"` + randStr + `", "groups":null}`
	statusCode, err = h.SendPostRequest(e.URL, "/users", payload)

	assert.Nil(t, err)
	assert.Equal(t, 400, statusCode)
}

func Test_UserPost_GroupDoesNotExists(t *testing.T) {
	groupName := util.RandStringBytes(32)
	randStr := util.RandStringBytes(32)

	payload := `{"first_name":"` + randStr + `", "last_name":"` + randStr + `", "userid":"` + randStr + `", "groups":["` + groupName + `"]}`
	statusCode, err := h.SendPostRequest(e.URL, "/users", payload)

	assert.Nil(t, err)
	assert.Equal(t, 201, statusCode)

	// retrieve user - group omitted because it does not exist
	var user model.RestUser
	r, err := http.Get(fmt.Sprintf("%s/users/%s", e.URL, randStr))
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Fatal(err)
		return
	}
	defer r.Body.Close()

	assert.Nil(t, err)
	assert.Equal(t, 200, r.StatusCode)

	// check fields
	assert.Equal(t, randStr, user.FirstName)
	assert.Equal(t, randStr, user.LastName)
	assert.Equal(t, randStr, user.UserId)
	assert.Equal(t, &[]string{}, user.Groups)
}

func Test_UserPut_UserUpdated(t *testing.T) {
	// create user
	randStr := util.RandStringBytes(32)
	payload := `{"first_name":"` + randStr + `", "last_name":"` + randStr + `", "userid":"` + randStr + `", "groups":null}`
	statusCode, err := h.SendPostRequest(e.URL, "/users", payload)

	assert.Nil(t, err)
	assert.Equal(t, 201, statusCode)

	// update user
	newRandStr := util.RandStringBytes(32)
	payload = `{"first_name":"` + newRandStr + `", "last_name":"` + newRandStr + `", "userid":"` + randStr + `", "groups":null}`
	statusCode, err = h.SendPutRequest(e.URL, "/users", randStr, payload)

	assert.Nil(t, err)
	assert.Equal(t, 200, statusCode)

	// retrieve user - group omitted because it does not exist
	var user model.RestUser
	r, err := http.Get(fmt.Sprintf("%s/users/%s", e.URL, randStr))
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Fatal(err)
		return
	}
	defer r.Body.Close()

	assert.Nil(t, err)
	assert.Equal(t, 200, r.StatusCode)

	// check new fields
	assert.Equal(t, newRandStr, user.FirstName)
	assert.Equal(t, newRandStr, user.LastName)
	assert.Equal(t, randStr, user.UserId)
	assert.Equal(t, &[]string{}, user.Groups)
}

func Test_UserPut_AttemptToUpdateKey(t *testing.T) {
	// create user
	randStr := util.RandStringBytes(32)
	payload := `{"first_name":"` + randStr + `", "last_name":"` + randStr + `", "userid":"` + randStr + `", "groups":null}`
	statusCode, err := h.SendPostRequest(e.URL, "/users", payload)

	assert.Nil(t, err)
	assert.Equal(t, 201, statusCode)

	// update user
	newRandStr := util.RandStringBytes(32)
	payload = `{"first_name":"` + randStr + `", "last_name":"` + randStr + `", "userid":"` + newRandStr + `", "groups":null}`
	statusCode, err = h.SendPutRequest(e.URL, "/users", randStr, payload)

	assert.Nil(t, err)
	// userid from payload is used, new one is not found
	assert.Equal(t, 404, statusCode)
}

func Test_UserPut_InvalidPayload(t *testing.T) {
	randStr := util.RandStringBytes(32)
	payload := `{"first_name":"", "last_name":"", "userid":"", "groups":null}`
	statusCode, err := h.SendPutRequest(e.URL, "/users", randStr, payload)

	assert.Nil(t, err)
	assert.Equal(t, 400, statusCode)
}

func Test_UserPut_UserDoesNotExist(t *testing.T) {
	randStr := util.RandStringBytes(32)
	payload := `{"first_name":"asd", "last_name":"asd", "userid":"asd", "groups":null}`
	statusCode, err := h.SendPutRequest(e.URL, "/users", randStr, payload)

	assert.Nil(t, err)
	assert.Equal(t, 404, statusCode)
}

func Test_UserPut_UpdateGroup(t *testing.T) {
	// create user
	randStr := util.RandStringBytes(32)
	payload := `{"first_name":"` + randStr + `", "last_name":"` + randStr + `", "userid":"` + randStr + `", "groups":null}`
	statusCode, err := h.SendPostRequest(e.URL, "/users", payload)

	assert.Nil(t, err)
	assert.Equal(t, 201, statusCode)

	// retrieve user
	var user model.RestUser
	r, err := http.Get(fmt.Sprintf("%s/users/%s", e.URL, randStr))
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Fatal(err)
		return
	}
	defer r.Body.Close()

	assert.Nil(t, err)
	assert.Equal(t, 200, r.StatusCode)
	// no group yet
	assert.Equal(t, &[]string{}, user.Groups)
	
	// create group
	groupName := util.RandStringBytes(32)
	payload = `{"name":"` + groupName + `"}`
	statusCode, err = h.SendPostRequest(e.URL, "/groups", payload)

	assert.Nil(t, err)
	assert.Equal(t, 201, statusCode)

	// update user by adding group
	payload = `{"first_name":"` + randStr + `", "last_name":"` + randStr + `", "userid":"` + randStr + `", "groups":["` + groupName + `"]}`
	statusCode, err = h.SendPutRequest(e.URL, "/users", randStr, payload)

	assert.Nil(t, err)
	assert.Equal(t, 200, statusCode)

	// retrieve updated user
	r, err = http.Get(fmt.Sprintf("%s/users/%s", e.URL, randStr))
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Fatal(err)
		return
	}
	defer r.Body.Close()

	assert.Nil(t, err)
	assert.Equal(t, 200, r.StatusCode)
	// group updated
	assert.Equal(t, &[]string{groupName}, user.Groups)
}

func Test_UserDelete_UserDeleted(t *testing.T) {
	// create user
	randStr := util.RandStringBytes(32)
	payload := `{"first_name":"` + randStr + `", "last_name":"` + randStr + `", "userid":"` + randStr + `", "groups":null}`

	statusCode, err := h.SendPostRequest(e.URL, "/users", payload)

	assert.Nil(t, err)
	assert.Equal(t, 201, statusCode)

	// delete user
	statusCode, err = h.SendDelRequest(e.URL, "/users", randStr)

	assert.Nil(t, err)
	assert.Equal(t, 200, statusCode)

	// retrieve user
	r, err := http.Get(fmt.Sprintf("%s/users/%s", e.URL, randStr))

	assert.Nil(t, err)
	assert.Equal(t, 404, r.StatusCode)
}

func Test_UserDelete_UserDoesNotExist(t *testing.T) {
	userId := util.RandStringBytes(32)

	statusCode, err := h.SendDelRequest(e.URL, "/users", userId)

	assert.Nil(t, err)
	assert.Equal(t, 404, statusCode)
}
