package goselenium

import (
	"errors"
	"testing"
)

/*
	DismissAlert() Tests
*/

func Test_AlertDismissAlert_InvalidSessionIdResultsInError(t *testing.T) {
	api := &testableAPIService{
		jsonToReturn:  "",
		errorToReturn: nil,
	}

	d := setUpDriver(setUpDefaultCaps(), api)

	_, err := d.DismissAlert()
	if err == nil || !IsSessionIDError(err) {
		t.Errorf(sessionIDErrorText)
	}
}

func Test_AlertDismissAlert_CommunicationErrorIsReturnedCorrectly(t *testing.T) {
	api := &testableAPIService{
		jsonToReturn:  "",
		errorToReturn: errors.New("An error :<"),
	}

	d := setUpDriver(setUpDefaultCaps(), api)
	d.sessionID = "12345"

	_, err := d.DismissAlert()
	if err == nil || !IsCommunicationError(err) {
		t.Errorf(apiCommunicationErrorText)
	}
}

func Test_AlertDismissAlert_UnmarshallingErrorIsReturnedCorrectly(t *testing.T) {
	api := &testableAPIService{
		jsonToReturn:  "Invalid JSON!",
		errorToReturn: nil,
	}

	d := setUpDriver(setUpDefaultCaps(), api)
	d.sessionID = "12345"

	_, err := d.DismissAlert()
	if err == nil || !IsUnmarshallingError(err) {
		t.Errorf(unmarshallingErrorText)
	}
}

func Test_AlertDismissAlert_CorrectResponseIsReturned(t *testing.T) {
	api := &testableAPIService{
		jsonToReturn: `{
			"state": "success",
			"value": "8"
		}`,
		errorToReturn: nil,
	}

	d := setUpDriver(setUpDefaultCaps(), api)
	d.sessionID = "12345"

	resp, err := d.DismissAlert()
	if err != nil || resp.State != "success" {
		t.Errorf(correctResponseErrorText)
	}
}

/*
	AcceptAlert() Tests
*/

func Test_AlertAcceptAlert_InvalidSessionIdResultsInError(t *testing.T) {
	api := &testableAPIService{
		jsonToReturn:  "",
		errorToReturn: nil,
	}

	d := setUpDriver(setUpDefaultCaps(), api)

	_, err := d.AcceptAlert()
	if err == nil || !IsSessionIDError(err) {
		t.Errorf(sessionIDErrorText)
	}
}

func Test_AlertAcceptAlert_CommunicationErrorIsReturnedCorrectly(t *testing.T) {
	api := &testableAPIService{
		jsonToReturn:  "",
		errorToReturn: errors.New("An error :<"),
	}

	d := setUpDriver(setUpDefaultCaps(), api)
	d.sessionID = "12345"

	_, err := d.AcceptAlert()
	if err == nil || !IsCommunicationError(err) {
		t.Errorf(apiCommunicationErrorText)
	}
}

func Test_AlertAcceptAlert_UnmarshallingErrorIsReturnedCorrectly(t *testing.T) {
	api := &testableAPIService{
		jsonToReturn:  "Invalid JSON!",
		errorToReturn: nil,
	}

	d := setUpDriver(setUpDefaultCaps(), api)
	d.sessionID = "12345"

	_, err := d.AcceptAlert()
	if err == nil || !IsUnmarshallingError(err) {
		t.Errorf(unmarshallingErrorText)
	}
}

func Test_AlertAcceptAlert_CorrectResponseIsReturned(t *testing.T) {
	api := &testableAPIService{
		jsonToReturn: `{
			"state": "success",
			"value": "8"
		}`,
		errorToReturn: nil,
	}

	d := setUpDriver(setUpDefaultCaps(), api)
	d.sessionID = "12345"

	resp, err := d.AcceptAlert()
	if err != nil || resp.State != "success" {
		t.Errorf(correctResponseErrorText)
	}
}