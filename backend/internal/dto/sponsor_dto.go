package dto

// SponsorStatsResponse: 
type SponsorStatsResponse struct {
	SponsorName      string `json:"sponsor_name"`
	LogoURL          string `json:"logo_url"`
	
	// Impact Metrics
	TotalDonors      int64 `json:"total_donors"`     
		LivesSaved       int64 `json:"lives_saved"`       
		RequestsFulfilled int64 `json:"requests_fulfilled"` 
	
	// Engagement
	ActiveRegions    []string `json:"active_regions"`  
}