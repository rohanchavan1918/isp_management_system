package controllers

import (
	"net/http"
	"net/url"
	"testing"
)

func TestSignup(t *testing.T) {
	// testtable := struct{
	// 	urlValues url.Values,

	// }{

	// }
	data := url.Values{
		"FirstName":   {"Test"},
		"LastName":    {"Test"},
		"Email":       {"encrypted5@gmail.com"},
		"MobileNo":    {"+919879879878"},
		"DateOfBirth": {"02/04/1998"},
		"Gender":      {"M"},
		"Password":    {"Tester@1234"},
	}

	res, err := http.PostForm("http://127.0.0.1:8080/api/v1/signup", data)
	if err != nil {
		t.Fatal("TEST FAILED ", err)
	}
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("User already exists ")
	}
}
