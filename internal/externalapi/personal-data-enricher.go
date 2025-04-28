package externalapi

import "context"

type EnrichPersonalData interface {
	// GetPersonAge returns person age by name
	GetPersonAge(ctx context.Context, name string) (int, error)
	// GetPersonGender returns person gender by name
	GetPersonGender(ctx context.Context, name string) (string, error)
	// GetPersonNationality returns person nationality by name
	GetPersonNationality(ctx context.Context, name string) (string, error)
}
