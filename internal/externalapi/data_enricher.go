package externalapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
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
	log.Printf("personalDataEnricher.GetPersonAge: getting person age by name: %s", name)
	endpoint := fmt.Sprintf("https://api.agify.io/?name=%s", url.QueryEscape(name))
	log.Printf("personalDataEnricher.GetPersonAge: endpoint: %s", endpoint)

	// Create request
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	log.Printf("personalDataEnricher.GetPersonAge: request created: %v", req)

	// Send request
	log.Printf("personalDataEnricher.GetPersonAge: sending request")
	resp, err := p.client.Do(req)
	if err != nil {
		log.Printf("personalDataEnricher.GetPersonAge: could not send request: %v", err)
		return 0, fmt.Errorf("could not send agify request: %w", err)
	}
	log.Printf("personalDataEnricher.GetPersonAge: request sent: %v", resp)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("personalDataEnricher.GetPersonAge: agify returned non-200 status code %d: %s", resp.StatusCode, string(body))
		return 0, fmt.Errorf("agify returned non-200 status code %d: %s", resp.StatusCode, string(body))
	}

	log.Printf("personalDataEnricher.GetPersonAge: decoding agify response")
	var ar models.AgifyResponse
	if err := json.NewDecoder(resp.Body).Decode(&ar); err != nil {
		log.Printf("personalDataEnricher.GetPersonAge: could not decode agify response: %v", err)
		return 0, fmt.Errorf("could not decode agify response: %w", err)
	}

	if ar.Age != nil {
		log.Printf("personalDataEnricher.GetPersonAge: age: %d", *ar.Age)
		return *ar.Age, nil
	} else {
		log.Printf("personalDataEnricher.GetPersonAge: agify returned empty age field")
		return 0, fmt.Errorf("agify returned empty age field")
	}
}

// GetPersonGender returns the predicted gender of a person by name.
//
// It calls the Genderize API (https://genderize.io/) and returns the predicted gender.
// If the API request fails, the response cannot be decoded, or the gender field is empty, an error is returned.
func (p *personalDataEnricher) GetPersonGender(ctx context.Context, name string) (string, error) {
	// Create endpoint
	log.Printf("personalDataEnricher.GetPersonGender: getting person gender by name: %s", name)
	endpoint := fmt.Sprintf("https://api.genderize.io/?name=%s", url.QueryEscape(name))

	// Create request
	log.Printf("personalDataEnricher.GetPersonGender: creating request")
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	log.Printf("personalDataEnricher.GetPersonGender: request created: %v", req)

	// Send request
	log.Printf("personalDataEnricher.GetPersonGender: sending request")
	resp, err := p.client.Do(req)
	if err != nil {
		log.Printf("personalDataEnricher.GetPersonGender: could not send request: %v", err)
		return "", fmt.Errorf("could not send genderize request: %w", err)
	}
	log.Printf("personalDataEnricher.GetPersonGender: request sent: %v", resp)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("personalDataEnricher.GetPersonGender: genderize returned non-200 status code %d", resp.StatusCode)
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("genderize returned non-200 status code %d: %s", resp.StatusCode, string(body))
	}

	var gr models.GenderizeResponse
	log.Printf("personalDataEnricher.GetPersonGender: decoding genderize response")
	if err := json.NewDecoder(resp.Body).Decode(&gr); err != nil {
		log.Printf("personalDataEnricher.GetPersonGender: could not decode genderize response: %v", err)
		return "", fmt.Errorf("could not decode genderize response: %w", err)
	}

	if gr.Gender != nil {
		log.Printf("personalDataEnricher.GetPersonGender: gender: %s", *gr.Gender)
		return *gr.Gender, nil
	} else {
		log.Printf("personalDataEnricher.GetPersonGender: genderize returned empty gender field")
		return "", fmt.Errorf("genderize returned empty gender field")
	}
}


// GetPersonNationality returns the predicted nationality of a person by name.
//
// It calls the Nationalize API (https://nationalize.io/) and returns the predicted nationality.
// If the API request fails, the response cannot be decoded, or the nationality field is empty, an error is returned.
func (p *personalDataEnricher) GetPersonNationality(ctx context.Context, name string) (string, error) {
	// Create endpoint
	log.Printf("personalDataEnricher.GetPersonNationality: getting person nationality by name: %s", name)
	endpoint := fmt.Sprintf("https://api.nationalize.io/?name=%s", url.QueryEscape(name))
	log.Printf("personalDataEnricher.GetPersonNationality: endpoint: %s", endpoint)

	// Create request
	log.Printf("personalDataEnricher.GetPersonNationality: creating request")
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	log.Printf("personalDataEnricher.GetPersonNationality: request created: %v", req)

	// Send request
	log.Printf("personalDataEnricher.GetPersonNationality: sending request")
	resp, err := p.client.Do(req)
	if err != nil {
		log.Printf("personalDataEnricher.GetPersonNationality: could not send request: %v", err)
		return "", fmt.Errorf("could not send nationalize request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("personalDataEnricher.GetPersonNationality: nationalize returned non-200 status code %d", resp.StatusCode)
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("nationalize returned non-200 status code %d: %s", resp.StatusCode, string(body))
	}

	var nr models.NationalizeResponse
	if err := json.NewDecoder(resp.Body).Decode(&nr); err != nil {
		log.Printf("personalDataEnricher.GetPersonNationality: could not decode nationalize response: %v", err)
		return "", fmt.Errorf("could not decode nationalize response: %w", err)
	}

	if len(nr.Country) > 0 {
		log.Printf("personalDataEnricher.GetPersonNationality: country: %s", nr.Country[0].CountryID)
		return nr.Country[0].CountryID, nil
	} else {
		log.Printf("personalDataEnricher.GetPersonNationality: nationalize returned empty country field")
		return "", fmt.Errorf("nationalize returned empty country field")
	}
}