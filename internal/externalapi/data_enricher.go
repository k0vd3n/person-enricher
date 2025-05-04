package externalapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"person-enricher/internal/models"
)

// EnrichPersonalData is an interface for enriching personal data
type EnrichPersonalData interface {
	// GetPersonAge returns person age by name
	GetPersonAge(ctx context.Context, name string) (int, error)
	// GetPersonGender returns person gender by name
	GetPersonGender(ctx context.Context, name string) (string, error)
	// GetPersonNationality returns person nationality by name
	GetPersonNationality(ctx context.Context, name string) (string, error)
}

// HTTPClient is an interface for making HTTP requests
type HTTPClient interface {
    Do(req *http.Request) (*http.Response, error)
}

// personalDataEnricher implements EnrichPersonalData
type personalDataEnricher struct {
	client HTTPClient
}

// NewPersonalDataEnricher creates a new instance of personalDataEnricher with a default HTTP client.
// It returns an EnrichPersonalData interface, which provides methods for enriching personal data
// such as age, gender, and nationality using external APIs.
func NewPersonalDataEnricher() EnrichPersonalData {
	return &personalDataEnricher{
		client: &http.Client{},
	}
}

// GetPersonAge returns person age by name
//
// It calls the Agify API (https://agify.io/) and returns person age.
// If the API returns an error or empty age field, it returns an error.
func (p *personalDataEnricher) GetPersonAge(ctx context.Context, name string) (int, error) {
	// Create endpoint
	endpoint := fmt.Sprintf("https://api.agify.io/?name=%s", url.QueryEscape(name))

	// Create request
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)

	// Send request
	resp, err := p.client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("could not send agify request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return 0, fmt.Errorf("agify returned non-200 status code %d: %s", resp.StatusCode, string(body))
	}

	var ar models.AgifyResponse
	if err := json.NewDecoder(resp.Body).Decode(&ar); err != nil {
		return 0, fmt.Errorf("could not decode agify response: %w", err)
	}

	if ar.Age != nil {
		return *ar.Age, nil
	} else {
		return 0, fmt.Errorf("agify returned empty age field")
	}
}

// GetPersonGender returns the predicted gender of a person by name.
//
// It calls the Genderize API (https://genderize.io/) and returns the predicted gender.
// If the API request fails, the response cannot be decoded, or the gender field is empty, an error is returned.
func (p *personalDataEnricher) GetPersonGender(ctx context.Context, name string) (string, error) {
	// Create endpoint
	endpoint := fmt.Sprintf("https://api.genderize.io/?name=%s", url.QueryEscape(name))

	// Create request
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)

	// Send request
	resp, err := p.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("could not send genderize request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("genderize returned non-200 status code %d: %s", resp.StatusCode, string(body))
	}

	var gr models.GenderizeResponse
	if err := json.NewDecoder(resp.Body).Decode(&gr); err != nil {
		return "", fmt.Errorf("could not decode genderize response: %w", err)
	}

	if gr.Gender != nil {
		return *gr.Gender, nil
	} else {
		return "", fmt.Errorf("genderize returned empty gender field")
	}
}


// GetPersonNationality returns the predicted nationality of a person by name.
//
// It calls the Nationalize API (https://nationalize.io/) and returns the predicted nationality.
// If the API request fails, the response cannot be decoded, or the nationality field is empty, an error is returned.
func (p *personalDataEnricher) GetPersonNationality(ctx context.Context, name string) (string, error) {
	// Create endpoint
	endpoint := fmt.Sprintf("https://api.nationalize.io/?name=%s", url.QueryEscape(name))

	// Create request
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)

	// Send request
	resp, err := p.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("could not send nationalize request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("nationalize returned non-200 status code %d: %s", resp.StatusCode, string(body))
	}

	var nr models.NationalizeResponse
	if err := json.NewDecoder(resp.Body).Decode(&nr); err != nil {
		return "", fmt.Errorf("could not decode nationalize response: %w", err)
	}

	if len(nr.Country) > 0 {
		return nr.Country[0].CountryID, nil
	} else {
		return "", fmt.Errorf("nationalize returned empty country field")
	}
}