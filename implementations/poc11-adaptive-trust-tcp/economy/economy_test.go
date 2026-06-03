package economy

import "testing"

func TestEvaluatorAcceptsPositiveReciprocalValue(t *testing.T) {
	decision := Evaluator{}.Decide(Offer{PromisedValue: 2, ReciprocalValue: 4, OpportunityCost: 3, Trust: 1, Budget: 2, Capacity: 1})
	if !decision.PromiseWorthMaking {
		t.Fatalf("positive reciprocal value should be worth promising: %#v", decision)
	}
}

func TestEvaluatorRejectsOpportunityCost(t *testing.T) {
	decision := Evaluator{}.Decide(Offer{PromisedValue: 1, ReciprocalValue: 1, OpportunityCost: 5, Trust: 0, Budget: 2, Capacity: 1})
	if decision.PromiseWorthMaking {
		t.Fatalf("high opportunity cost should reject promise: %#v", decision)
	}
}

func TestEvaluatorRejectsLowTrust(t *testing.T) {
	decision := Evaluator{}.Decide(Offer{PromisedValue: 10, ReciprocalValue: 10, OpportunityCost: 1, Trust: -4, Budget: 2, Capacity: 1})
	if decision.PromiseWorthMaking {
		t.Fatalf("very low trust should reject promise: %#v", decision)
	}
}
