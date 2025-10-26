package obsidian

import "testing"

func TestMomentToGoLayout(t *testing.T) {
	tests := map[string]string{
		"YYYY-MM-DD":     "2006-01-02",
		"YYYYMMDDHHmm":   "200601021504",
		"[prefix]-YYYY":  "prefix-2006",
		"YYYY-MM-DDTHH":  "2006-01-02T15",
		"YYYY-MM-DD HH":  "2006-01-02 15",
		"YYYY-MM-DD hhA": "2006-01-02 03PM",
	}

	for input, expected := range tests {
		actual, ok := momentToGoLayout(input)
		if !ok {
			t.Fatalf("expected conversion to succeed for %q", input)
		}
		if actual != expected {
			t.Fatalf("expected %q -> %q, got %q", input, expected, actual)
		}
	}
}

func TestMomentToGoLayoutUnbalancedLiteral(t *testing.T) {
	if _, ok := momentToGoLayout("[test"); ok {
		t.Fatal("expected failure for unbalanced literal")
	}
}
