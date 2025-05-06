package forms

import (
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/test", nil)
	form := New(r.PostForm)

	isValid := form.Valid()
	if !isValid {
		t.Error("Got invalid when Expected form to be valid")
	}
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/test", nil)
	form := New(r.PostForm)

	form.Required("a", "b", "c")
	isValid := form.Valid()
	if isValid {
		t.Error("Form shows valid when required fields are not present")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "b")
	postedData.Add("c", "c")

	r = httptest.NewRequest("POST", "/test", nil)

	r.PostForm = postedData
	form = New(r.PostForm)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("Form shows invalid when required fields are present")
	}
}

func TestForm_Has(t *testing.T) {

	postedData := url.Values{}

	form := New(postedData)
	has := form.Has("a")

	if has {
		t.Error("Form shows true when field is not present")
	}


	postedData.Add("a", "a")
	form = New(postedData)

	has = form.Has("a")
	if !has {
		t.Error("Form shows false when field is present")
	}
}

func TestForm_MinLength(t *testing.T) {

	postedData := url.Values{}
	form := New(postedData)
	form.MinLength("x", 10)

	if form.Valid() {
		t.Error("Form shows minlength when field is not present")
	}

	isError := form.Errors.Get("x")
	if isError == "" {
		t.Error("Form shows no error when field is not present")
	}

	postedData.Add("some_filed", "some_value")
	form = New(postedData)

	form.MinLength("some_filed", 100)

	if form.Valid() {
		t.Error("Form shows minlength of 100 when data is shorter")
	}

	postedData = url.Values{}
	postedData.Add("another_filed", "abc123")
	form = New(postedData)
	form.MinLength("another_filed", 1)

	if !form.Valid() {
		t.Error("Form shows minlength of 1 is not met when it is")
	}

	isError = form.Errors.Get("another_filed")
	if isError != "" {
		t.Error("Form should not have error but got an error")
	}
}

func TestForm_IsEmail(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)

	form.IsEmail("x")
	if form.Valid() {
		t.Error("Form shows valid when field is not present")
	}

	form = New(postedData)
	postedData.Add("email", "test@test.com")
	form.IsEmail("email")

	if !form.Valid() {
		t.Error("Form shows invalid when email is valid")
	}

	postedData = url.Values{}
	form = New(postedData)
	postedData.Add("email", "test@test")
	form.IsEmail("email")

	if form.Valid() {
		t.Error("Form shows valid when email is in-valid")
	}
}