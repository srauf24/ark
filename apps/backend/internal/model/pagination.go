package model

// Constants for default pagination limits
const (
	// DefaultAssetLimit is the default number of assets returned per page
	DefaultAssetLimit = 20
	// MaxAssetLimit is the maximum number of assets that can be requested per page
	MaxAssetLimit = 100

	// DefaultLogLimit is the default number of logs returned per page
	DefaultLogLimit = 50
	// MaxLogLimit is the maximum number of logs that can be requested per page
	MaxLogLimit = 200
)

// PaginationParams represents query parameters for pagination
type PaginationParams struct {
	Limit  int `query:"limit" validate:"omitempty,min=1,max=100"`
	Offset int `query:"offset" validate:"omitempty,min=0"`
}

// SetDefaults sets default values for pagination parameters
// If Limit is 0, it sets it to the provided defaultLimit
// If Offset is negative, it sets it to 0
func (p *PaginationParams) SetDefaults(defaultLimit int) {
	if p.Limit == 0 {
		p.Limit = defaultLimit
	}
	if p.Offset < 0 {
		p.Offset = 0
	}
}
