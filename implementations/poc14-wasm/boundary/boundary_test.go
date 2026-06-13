package boundary

import "testing"

func TestValidateWASMModuleAcceptsFixture(t *testing.T) {
	if err := ValidateWASMModule(MinimalWASMModule); err != nil {
		t.Fatalf("validate wasm fixture: %v", err)
	}
}

func TestValidateWASMModuleRejectsBadMagic(t *testing.T) {
	moduleBytes := append([]byte(nil), MinimalWASMModule...)
	moduleBytes[0] = 0xff
	if err := ValidateWASMModule(moduleBytes); err == nil {
		t.Fatalf("bad wasm magic should be rejected")
	}
}

func TestPromiseFieldsUsesSinglePromiseAction(t *testing.T) {
	fields := PromiseFields("victor", "peggy", PromiseAboutStdioAdapter, "Victor promises one stdio-carried envelope.")
	if fields["act"] != "promise" {
		t.Fatalf("act = %q, want promise", fields["act"])
	}
	if fields["field_promise_about"] != PromiseAboutStdioAdapter {
		t.Fatalf("field_promise_about = %q, want %q", fields["field_promise_about"], PromiseAboutStdioAdapter)
	}
}
