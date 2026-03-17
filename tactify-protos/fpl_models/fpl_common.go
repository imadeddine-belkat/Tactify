package fpl_models

// Common types tactify-protos across multiple domains

// Chip represents a chip play
type Chip struct {
	Name      string `json:"name"`
	NumPlayed int    `json:"num_played"`
}

// Phase represents a phase of the season
type Phase struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	StartEvent int    `json:"start_event"`
	StopEvent  int    `json:"stop_event"`
}

// ElementStat represents element statistics metadata
type ElementStat struct {
	Label string `json:"label"`
	Name  string `json:"name"`
}
