package dawg

import "testing"

var epsilon = 0.001

func TestFuzz(t *testing.T) {
	tests := []struct {
		query   string
		choices []string
	}{
		{"api", []string{"api", "my-api", "a-phrase-that-includes-all-letters-in-order", "angular-tooltips"}},
		{"map", []string{"my-api", "api"}},
		{"open", []string{"Open", "my-openthingy"}},
	}

	for _, tt := range tests {
		for i := 0; i < len(tt.choices)-1; i++ {
			choice0, choice1 := tt.choices[i], tt.choices[i+1]
			score0 := score(choice0, tt.query)
			score1 := score(choice1, tt.query)
			if score0 < score1 {
				t.Errorf("expected score for choice %s (%f) to be greater than score for choice %s (%f) for query %s", choice0, score0, choice1, score1, tt.query)
			}
		}
	}
}
