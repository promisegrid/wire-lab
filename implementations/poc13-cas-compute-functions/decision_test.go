package poc13

import "testing"

func TestResponseText(t *testing.T) {
	tests := []struct {
		name     string
		response ProviderResponse
		want     string
	}{
		{
			name:     "top level output text",
			response: ProviderResponse{OutputText: " I locally promise CAS storage "},
			want:     "I locally promise CAS storage",
		},
		{
			name: "nested output text",
			response: ProviderResponse{Output: []ProviderOutput{{
				Content: []ProviderContent{{Text: " I locally promise CID compute "}},
			}}},
			want: "I locally promise CID compute",
		},
		{
			name: "empty response",
			response: ProviderResponse{Output: []ProviderOutput{{
				Content: []ProviderContent{{Text: "   "}},
			}}},
			want: "",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.response.ResponseText(); got != test.want {
				t.Fatalf("ResponseText() = %q, want %q", got, test.want)
			}
		})
	}
}
