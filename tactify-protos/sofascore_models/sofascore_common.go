package sofascore_models

// Common types tactify-protos across multiple domains

type Country struct {
	Name string `json:"name,omitempty"`
}

type TeamColors struct {
	PrimaryColor   string `json:"primary,omitempty"`
	SecondaryColor string `json:"secondary,omitempty"`
	TextColor      string `json:"text,omitempty"`
}

type PlayerColor struct {
	Primary     string `json:"primary"`
	Number      string `json:"number"`
	Outline     string `json:"outline"`
	FancyNumber string `json:"fancyNumber"`
}

type Status struct {
	Code        int    `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
	Type        string `json:"type,omitempty"`
}
