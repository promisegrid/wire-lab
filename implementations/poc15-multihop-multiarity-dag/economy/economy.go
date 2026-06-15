package economy

// Offer is one local agent's view of a reciprocal promise opportunity.
// Intent: POC15 gives agents a minimal economics model so cooperation depends
// on local value, opportunity cost, and trust instead of a scripted global
// transaction. Source: DI-timah
type Offer struct {
	Promiser        string
	Promisee        string
	Resource        string
	PromisedValue   int
	ReciprocalValue int
	OpportunityCost int
	Trust           int
	Budget          int
	Capacity        int
}

// Decision is a local judgment about whether a promise is worth making.
type Decision struct {
	PromiseWorthMaking bool
	Score              int
	Reason             string
}

// Evaluator applies a deliberately small local scoring rule.
type Evaluator struct{}

// Decide compares expected value against opportunity cost and relationship
// risk. It does not command either agent; it only summarizes the local agent's
// own willingness to promise.
func (Evaluator) Decide(offer Offer) Decision {
	score := offer.PromisedValue + offer.ReciprocalValue + offer.Trust - offer.OpportunityCost
	if offer.Budget <= 0 {
		return Decision{PromiseWorthMaking: false, Score: score, Reason: "budget exhausted"}
	}
	if offer.Capacity <= 0 {
		return Decision{PromiseWorthMaking: false, Score: score, Reason: "capacity exhausted"}
	}
	if offer.Trust < -3 {
		return Decision{PromiseWorthMaking: false, Score: score, Reason: "relationship trust too low"}
	}
	if score <= 0 {
		return Decision{PromiseWorthMaking: false, Score: score, Reason: "opportunity cost exceeds local value"}
	}
	return Decision{PromiseWorthMaking: true, Score: score, Reason: "local value exceeds opportunity cost"}
}
