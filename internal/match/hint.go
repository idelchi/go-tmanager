package match

import "strconv"

// Hint represents a pattern used to match asset names.
// It can be a regular expression or a simple string pattern.
type Hint struct {
	Pattern string // Pattern to match against the asset's name.
	Weight  string // Weight used to adjust the score for non-mandatory hints.
	Must    bool   // Indicates if the hint is mandatory for a match.

	weightInt int `json:"-" mapstructure:"-" yaml:"-"`
}

// SetWeight converts the weight string to an integer.
func (h *Hint) SetWeight() error {
	val, err := strconv.Atoi(h.Weight)
	if err != nil {
		return err
	}

	h.weightInt = val

	return nil
}

// GetWeight returns the weight of the hint.
func (h Hint) GetWeight() int {
	return h.weightInt
}
