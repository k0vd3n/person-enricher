package externalapi

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newPersonalDataEnricherWithClient(c HTTPClient) EnrichPersonalData {
	return &personalDataEnricher{client: c}
}

type mockClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

func makeResponse(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Body:       io.NopCloser(bytes.NewBufferString(body)),
		Header:     make(http.Header),
	}
}

func TestGetPersonAge(t *testing.T) {
	tests := []struct {
		name    string
		client  HTTPClient
		wantAge int
		wantErr string
	}{
		{
			name: "success",
			client: &mockClient{DoFunc: func(req *http.Request) (*http.Response, error) {
				return makeResponse(200, `{"count":3800,"name":"Anna","age":28}`), nil
			}},
			wantAge: 28,
		},
		{
			name: "non-200 status",
			client: &mockClient{DoFunc: func(req *http.Request) (*http.Response, error) {
				return makeResponse(503, `service down`), nil
			}},
			wantErr: "agify returned non-200 status code",
		},
		{
			name: "decode error",
			client: &mockClient{DoFunc: func(req *http.Request) (*http.Response, error) {
				return makeResponse(200, `not-json`), nil
			}},
			wantErr: "could not decode agify response",
		},
		{
			name: "empty age field",
			client: &mockClient{DoFunc: func(req *http.Request) (*http.Response, error) {
				return makeResponse(200, `{"count":0,"name":"X","age":null}`), nil
			}},
			wantErr: "agify returned empty age field",
		},
		{
			name: "client error",
			client: &mockClient{DoFunc: func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("network fail")
			}},
			wantErr: "could not send agify request",
		},
		
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			e := newPersonalDataEnricherWithClient(tc.client)
			age, err := e.GetPersonAge(context.Background(), "Anna")

			if tc.wantErr != "" {
				assert.ErrorContains(t, err, tc.wantErr)
				assert.Zero(t, age)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.wantAge, age)
			}
		})
	}
}

func TestGetPersonGender(t *testing.T) {
	tests := []struct {
		name      string
		client    HTTPClient
		wantGender string
		wantErr   string
	}{
		{
			name: "success",
			client: &mockClient{DoFunc: func(req *http.Request) (*http.Response, error) {
				return makeResponse(200, `{"count":1000,"name":"Ivan","gender":"male","probability":1.00}`), nil
			}},
			wantGender: "male",
		},
		{
			name: "non-200 status",
			client: &mockClient{DoFunc: func(req *http.Request) (*http.Response, error) {
				return makeResponse(404, `not found`), nil
			}},
			wantErr: "genderize returned non-200 status code",
		},
		{
			name: "decode error",
			client: &mockClient{DoFunc: func(req *http.Request) (*http.Response, error) {
				return makeResponse(200, `invalid-json`), nil
			}},
			wantErr: "could not decode genderize response",
		},
		{
			name: "empty gender field",
			client: &mockClient{DoFunc: func(req *http.Request) (*http.Response, error) {
				return makeResponse(200, `{"count":0,"name":"X","gender":null,"probability":0}`), nil
			}},
			wantErr: "genderize returned empty gender field",
		},
		{
			name: "client error",
			client: &mockClient{DoFunc: func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("timeout")
			}},
			wantErr: "could not send genderize request",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			e := newPersonalDataEnricherWithClient(tc.client)
			g, err := e.GetPersonGender(context.Background(), "Ivan")

			if tc.wantErr != "" {
				assert.ErrorContains(t, err, tc.wantErr)
				assert.Empty(t, g)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.wantGender, g)
			}
		})
	}
}

func TestGetPersonNationality(t *testing.T) {
	tests := []struct {
		name     string
		client   HTTPClient
		wantNat  string
		wantErr  string
	}{
		{
			name: "success",
			client: &mockClient{DoFunc: func(req *http.Request) (*http.Response, error) {
				return makeResponse(200, `{"count":150,"name":"Ivan","country":[{"country_id":"RU","probability":1.0}]}`), nil
			}},
			wantNat: "RU",
		},
		{
			name: "non-200 status",
			client: &mockClient{DoFunc: func(req *http.Request) (*http.Response, error) {
				return makeResponse(500, `error`), nil
			}},
			wantErr: "nationalize returned non-200 status code",
		},
		{
			name: "decode error",
			client: &mockClient{DoFunc: func(req *http.Request) (*http.Response, error) {
				return makeResponse(200, `{bad json`), nil
			}},
			wantErr: "could not decode nationalize response",
		},
		{
			name: "empty country field",
			client: &mockClient{DoFunc: func(req *http.Request) (*http.Response, error) {
				return makeResponse(200, `{"count":0,"name":"X","country":[]}`), nil
			}},
			wantErr: "nationalize returned empty country field",
		},
		{
			name: "client error",
			client: &mockClient{DoFunc: func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("conn refused")
			}},
			wantErr: "could not send nationalize request",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			e := newPersonalDataEnricherWithClient(tc.client)
			n, err := e.GetPersonNationality(context.Background(), "Ivan")

			if tc.wantErr != "" {
				assert.ErrorContains(t, err, tc.wantErr)
				assert.Empty(t, n)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.wantNat, n)
			}
		})
	}
}
