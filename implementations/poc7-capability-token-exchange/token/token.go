package token

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"

	"promisegrid.dev/wire-lab/implementations/poc7-capability-token-exchange/protocol"
)

const (
	TransferBearer          = "bearer"
	TransferNonTransferable = "non-transferable"
	OutcomeKept             = "kept"
	OutcomeBroken           = "broken"
	OutcomeRefused          = "refused"
)

// Token is a signed promise by the issuer, not a global permission object.
// Intent: Model capability-like access as an issuer's own promise so redemption
// outcomes can update local trust without inventing a central authority.
// Source: DI-tugih
type Token struct {
	ID           string `json:"id"`
	Issuer       string `json:"issuer"`
	OriginalPeer string `json:"original_peer"`
	ResourceKind string `json:"resource_kind"`
	ResourceID   string `json:"resource_id"`
	TransferRule string `json:"transfer_rule"`
	PublicKeyHex string `json:"public_key_hex"`
	SignatureHex string `json:"signature_hex"`
}

// IsZero reports whether the wire message carries no promise token.
func (token Token) IsZero() bool {
	return token.ID == "" && token.Issuer == "" && token.SignatureHex == ""
}

// Event records what one observer saw. It is evidence, not global truth.
// Intent: Keep promise-accounting local to the observing agent while still
// making demo outcomes auditable by humans and later tests. Source: DI-tugih
type Event struct {
	Observer string `json:"observer"`
	Event    string `json:"event"`
	Outcome  string `json:"outcome"`
	TokenID  string `json:"token_id,omitempty"`
	Detail   string `json:"detail"`
}

// Issuer owns the local status of tokens it promised.
// Intent: Make revocation and redemption issuer-local so other agents observe
// outcomes and adjust trust instead of consulting a shared status ledger.
// Source: DI-tugih
type Issuer struct {
	Name     string
	tokens   map[string]Token
	revoked  map[string]bool
	redeemed map[string]bool
	events   []Event
}

// Wallet holds tokens and local trust scores for one observing agent.
// Intent: Keep exchange rates and trust changes private to the wallet owner
// rather than making them objective network state. Source: DI-tugih
type Wallet struct {
	Owner  string
	tokens map[string]Token
	events []Event
	trust  map[string]int
}

func NewIssuer(name string) *Issuer {
	return &Issuer{
		Name:     name,
		tokens:   make(map[string]Token),
		revoked:  make(map[string]bool),
		redeemed: make(map[string]bool),
	}
}

func (issuer *Issuer) Issue(id string, originalPeer string, resourceKind string, resourceID string, transferRule string) (Token, error) {
	if transferRule != TransferBearer && transferRule != TransferNonTransferable {
		return Token{}, fmt.Errorf("unknown transfer rule %q", transferRule)
	}
	token := Token{
		ID:           id,
		Issuer:       issuer.Name,
		OriginalPeer: originalPeer,
		ResourceKind: resourceKind,
		ResourceID:   resourceID,
		TransferRule: transferRule,
	}
	signature, signErr := SignToken(token)
	if signErr != nil {
		return Token{}, signErr
	}
	token.PublicKeyHex = hex.EncodeToString(deterministicPublicKey(token.Issuer))
	token.SignatureHex = signature
	issuer.tokens[id] = token
	issuer.events = append(issuer.events, Event{Observer: issuer.Name, Event: "token_issued", Outcome: OutcomeKept, TokenID: id, Detail: issuer.Name + " promised " + resourceKind + ":" + resourceID})
	return token, nil
}

func (issuer *Issuer) Revoke(id string, reason string) error {
	if _, ok := issuer.tokens[id]; !ok {
		return fmt.Errorf("issuer has no local revocation promise state for token %s", id)
	}
	issuer.revoked[id] = true
	issuer.events = append(issuer.events, Event{Observer: issuer.Name, Event: "token_revoked", Outcome: OutcomeKept, TokenID: id, Detail: reason})
	return nil
}

func (issuer *Issuer) Redeem(holder string, presented Token) Event {
	event := Event{Observer: issuer.Name, Event: "token_redeemed", TokenID: presented.ID}
	stored, ok := issuer.tokens[presented.ID]
	switch {
	case !ok:
		event.Outcome = OutcomeRefused
		event.Detail = "unknown token"
	case stored.SignatureHex != presented.SignatureHex || stored.PublicKeyHex != presented.PublicKeyHex || VerifyToken(presented) != nil:
		event.Outcome = OutcomeRefused
		event.Detail = "token proof did not match issuer promise"
	case issuer.revoked[presented.ID]:
		event.Outcome = OutcomeBroken
		event.Detail = "issuer revoked this token before redemption"
	case issuer.redeemed[presented.ID]:
		event.Outcome = OutcomeRefused
		event.Detail = "token was already redeemed"
	case stored.TransferRule == TransferNonTransferable && stored.OriginalPeer != holder:
		event.Outcome = OutcomeRefused
		event.Detail = "non-transferable token presented by " + holder + " instead of " + stored.OriginalPeer
	default:
		issuer.redeemed[presented.ID] = true
		event.Outcome = OutcomeKept
		event.Detail = issuer.Name + " kept access promise for " + holder
	}
	issuer.events = append(issuer.events, event)
	return event
}

func (issuer *Issuer) Events() []Event {
	events := make([]Event, len(issuer.events))
	copy(events, issuer.events)
	return events
}

func NewWallet(owner string) *Wallet {
	return &Wallet{
		Owner:  owner,
		tokens: make(map[string]Token),
		trust:  make(map[string]int),
	}
}

func (wallet *Wallet) Add(token Token, detail string) {
	wallet.tokens[token.ID] = token
	wallet.events = append(wallet.events, Event{Observer: wallet.Owner, Event: "token_received", Outcome: OutcomeKept, TokenID: token.ID, Detail: detail})
}

// Holds reports whether this local wallet currently holds a token.
func (wallet *Wallet) Holds(id string) bool {
	_, ok := wallet.tokens[id]
	return ok
}

func (wallet *Wallet) Transfer(id string, recipient *Wallet) (Token, error) {
	token, ok := wallet.tokens[id]
	if !ok {
		return Token{}, fmt.Errorf("%s does not hold token %s", wallet.Owner, id)
	}
	if token.TransferRule != TransferBearer {
		return Token{}, fmt.Errorf("token %s is not bearer-transferable", id)
	}
	delete(wallet.tokens, id)
	recipient.Add(token, wallet.Owner+" transferred bearer token to "+recipient.Owner)
	wallet.events = append(wallet.events, Event{Observer: wallet.Owner, Event: "token_transferred", Outcome: OutcomeKept, TokenID: id, Detail: "transferred to " + recipient.Owner})
	return token, nil
}

func (wallet *Wallet) ApplyRedemption(event Event) {
	issuer := ""
	if token, ok := wallet.tokens[event.TokenID]; ok {
		issuer = token.Issuer
	}
	if issuer == "" {
		issuer = "unknown"
	}
	switch event.Outcome {
	case OutcomeKept:
		wallet.trust[issuer] += 1
	case OutcomeBroken, OutcomeRefused:
		wallet.trust[issuer] -= 1
	}
	wallet.events = append(wallet.events, Event{Observer: wallet.Owner, Event: "local_trust_updated", Outcome: OutcomeKept, TokenID: event.TokenID, Detail: fmt.Sprintf("%s trust in %s is now %d after %s", wallet.Owner, issuer, wallet.trust[issuer], event.Outcome)})
}

// Quote prices one issuer's token against another using only this wallet's
// holdings and trust history. Intent: Keep POC7 economics local and observable
// without creating a central exchange or global exchange rate. Source: DI-fibok
func (wallet *Wallet) Quote(issuer string, wantedIssuer string) ExchangeOffer {
	issuerTrust := wallet.trust[issuer]
	wantedTrust := wallet.trust[wantedIssuer]
	offeredHeld := wallet.countByIssuer(issuer)
	wantedHeld := wallet.countByIssuer(wantedIssuer)
	rate := 1
	if wantedTrust > issuerTrust {
		rate += wantedTrust - issuerTrust
	}
	if issuerTrust < 0 {
		rate += 2
	}
	if offeredHeld == 0 {
		rate++
	}
	if wantedHeld > offeredHeld {
		rate++
	}
	return ExchangeOffer{Observer: wallet.Owner, OfferedIssuer: issuer, WantedIssuer: wantedIssuer, OfferedCount: rate, WantedCount: 1}
}

func (wallet *Wallet) countByIssuer(issuer string) int {
	count := 0
	for _, heldToken := range wallet.tokens {
		if heldToken.Issuer == issuer {
			count++
		}
	}
	return count
}

func (wallet *Wallet) Trust(issuer string) int {
	return wallet.trust[issuer]
}

func (wallet *Wallet) Events() []Event {
	events := make([]Event, len(wallet.events))
	copy(events, wallet.events)
	return events
}

// ExchangeOffer is one wallet's local quote, not a market price.
// Intent: Demonstrate floating peer-local exchange without a central exchange
// or global price oracle. Source: DI-tugih
type ExchangeOffer struct {
	Observer      string `json:"observer"`
	OfferedIssuer string `json:"offered_issuer"`
	WantedIssuer  string `json:"wanted_issuer"`
	OfferedCount  int    `json:"offered_count"`
	WantedCount   int    `json:"wanted_count"`
}

// SignToken signs a token with a deterministic POC-only Ed25519 key.
// Intent: Keep the executable demo reproducible while still making token
// mutation visible to redeemers and tests. Source: DI-tugih
func SignToken(token Token) (string, error) {
	token.SignatureHex = ""
	token.PublicKeyHex = ""
	token.PublicKeyHex = hex.EncodeToString(deterministicPublicKey(token.Issuer))
	bytes, marshalErr := CanonicalBytes(token)
	if marshalErr != nil {
		return "", marshalErr
	}
	privateKey := deterministicPrivateKey(token.Issuer)
	return hex.EncodeToString(ed25519.Sign(privateKey, bytes)), nil
}

// VerifyToken checks that the presented bytes still match the issuer promise.
// Intent: Treat proof failure as evidence about this presented promise token,
// not as a global gatekeeping decision. Source: DI-tugih; DI-tanat
func VerifyToken(token Token) error {
	signature, sigErr := hex.DecodeString(token.SignatureHex)
	if sigErr != nil {
		return sigErr
	}
	publicKey, publicErr := hex.DecodeString(token.PublicKeyHex)
	if publicErr != nil {
		return publicErr
	}
	expectedPublicKey := deterministicPublicKey(token.Issuer)
	if !equalBytes(publicKey, expectedPublicKey) {
		return fmt.Errorf("token public key does not match issuer")
	}
	bytes, marshalErr := CanonicalBytes(Token{
		ID:           token.ID,
		Issuer:       token.Issuer,
		OriginalPeer: token.OriginalPeer,
		ResourceKind: token.ResourceKind,
		ResourceID:   token.ResourceID,
		TransferRule: token.TransferRule,
		PublicKeyHex: token.PublicKeyHex,
	})
	if marshalErr != nil {
		return marshalErr
	}
	if !ed25519.Verify(publicKey, bytes, signature) {
		return fmt.Errorf("token signature failed")
	}
	return nil
}

// CanonicalBytes encodes the token promise as deterministic CBOR fields.
// Intent: Sign exact token bytes instead of Go struct layout or JSON text so
// mutation evidence is protocol-shaped. Source: DI-fibok
func CanonicalBytes(token Token) ([]byte, error) {
	return protocol.MarshalStringMap(map[string]string{
		"id":             token.ID,
		"issuer":         token.Issuer,
		"original_peer":  token.OriginalPeer,
		"resource_kind":  token.ResourceKind,
		"resource_id":    token.ResourceID,
		"transfer_rule":  token.TransferRule,
		"public_key_hex": token.PublicKeyHex,
	})
}

// Encode serializes the full signed token as deterministic CBOR.
func Encode(token Token) ([]byte, error) {
	return protocol.MarshalStringMap(map[string]string{
		"id":             token.ID,
		"issuer":         token.Issuer,
		"original_peer":  token.OriginalPeer,
		"resource_kind":  token.ResourceKind,
		"resource_id":    token.ResourceID,
		"transfer_rule":  token.TransferRule,
		"public_key_hex": token.PublicKeyHex,
		"signature_hex":  token.SignatureHex,
	})
}

// Decode rebuilds a signed token from deterministic CBOR fields.
func Decode(tokenBytes []byte) (Token, error) {
	fields, fieldsErr := protocol.UnmarshalStringMap(tokenBytes)
	if fieldsErr != nil {
		return Token{}, fieldsErr
	}
	return Token{
		ID:           fields["id"],
		Issuer:       fields["issuer"],
		OriginalPeer: fields["original_peer"],
		ResourceKind: fields["resource_kind"],
		ResourceID:   fields["resource_id"],
		TransferRule: fields["transfer_rule"],
		PublicKeyHex: fields["public_key_hex"],
		SignatureHex: fields["signature_hex"],
	}, nil
}

func deterministicPrivateKey(seedText string) ed25519.PrivateKey {
	seed := sha256.Sum256([]byte("poc7 token issuer: " + seedText))
	return ed25519.NewKeyFromSeed(seed[:])
}

func deterministicPublicKey(seedText string) ed25519.PublicKey {
	return deterministicPrivateKey(seedText).Public().(ed25519.PublicKey)
}

func equalBytes(left []byte, right []byte) bool {
	if len(left) != len(right) {
		return false
	}
	for index := range left {
		if left[index] != right[index] {
			return false
		}
	}
	return true
}

func SortedEvents(events []Event) []Event {
	sorted := make([]Event, len(events))
	copy(sorted, events)
	sort.SliceStable(sorted, func(left int, right int) bool {
		if sorted[left].Observer != sorted[right].Observer {
			return sorted[left].Observer < sorted[right].Observer
		}
		if sorted[left].TokenID != sorted[right].TokenID {
			return sorted[left].TokenID < sorted[right].TokenID
		}
		return sorted[left].Event < sorted[right].Event
	})
	return sorted
}
