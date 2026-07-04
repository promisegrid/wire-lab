package economy

import (
	"strings"
	"testing"
	"time"

	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/store"
)

func TestStoragePaymentTokenRedeemsOnce(t *testing.T) {
	aliceCAS, aliceErr := store.Open(t.TempDir())
	if aliceErr != nil {
		t.Fatalf("open alice CAS: %v", aliceErr)
	}
	frankCAS, frankErr := store.Open(t.TempDir())
	if frankErr != nil {
		t.Fatalf("open frank CAS: %v", frankErr)
	}
	objectCID := store.CIDText(store.CIDForBytes([]byte("release root")))
	now := time.Date(2026, 7, 4, 12, 0, 0, 0, time.UTC)
	token := StoragePaymentToken{
		Issuer:       "alice",
		Subject:      "frank",
		Scope:        "poc18-release-retention",
		ObjectCID:    objectCID,
		Value:        5,
		Unit:         "storage_credit",
		ExpiresUnix:  now.Add(time.Hour).Unix(),
		Nonce:        "poc18-storage-payment-token-1",
		Transferable: true,
	}
	issued, issueErr := IssueStoragePaymentToken(aliceCAS, token)
	if issueErr != nil {
		t.Fatalf("issue storage payment token: %v", issueErr)
	}
	tokenCID, parseErr := store.ParseCIDText(issued.CID)
	if parseErr != nil {
		t.Fatalf("parse token CID: %v", parseErr)
	}
	if !aliceCAS.Has(tokenCID) {
		t.Fatalf("issuer CAS does not contain token CID %s", issued.CID)
	}
	expected := ExpectedStoragePayment{
		Issuer:    "alice",
		Subject:   "frank",
		Scope:     "poc18-release-retention",
		ObjectCID: objectCID,
		Value:     5,
		Unit:      "storage_credit",
	}
	ledger := NewLedger()
	redemption, redeemErr := ledger.RedeemStoragePaymentToken(frankCAS, issued.Bytes, expected, "frank", now)
	if redeemErr != nil {
		t.Fatalf("redeem storage payment token: %v", redeemErr)
	}
	if !redemption.SignatureVerified || !redemption.Redeemed || redemption.TokenCID != issued.CID {
		t.Fatalf("unexpected redemption report: %#v", redemption)
	}
	if !frankCAS.Has(tokenCID) {
		t.Fatalf("redeemer CAS does not contain token CID %s", issued.CID)
	}
	if _, replayErr := ledger.RedeemStoragePaymentToken(frankCAS, issued.Bytes, expected, "frank", now); replayErr == nil || !strings.Contains(replayErr.Error(), "already redeemed") {
		t.Fatalf("replay should be rejected, got %v", replayErr)
	}
}

func TestStoragePaymentTokenRejectsWrongSubject(t *testing.T) {
	aliceCAS, aliceErr := store.Open(t.TempDir())
	if aliceErr != nil {
		t.Fatalf("open alice CAS: %v", aliceErr)
	}
	frankCAS, frankErr := store.Open(t.TempDir())
	if frankErr != nil {
		t.Fatalf("open frank CAS: %v", frankErr)
	}
	objectCID := store.CIDText(store.CIDForBytes([]byte("release root")))
	now := time.Date(2026, 7, 4, 12, 0, 0, 0, time.UTC)
	issued, issueErr := IssueStoragePaymentToken(aliceCAS, StoragePaymentToken{
		Issuer:       "alice",
		Subject:      "mallory",
		Scope:        "poc18-release-retention",
		ObjectCID:    objectCID,
		Value:        5,
		Unit:         "storage_credit",
		ExpiresUnix:  now.Add(time.Hour).Unix(),
		Nonce:        "poc18-storage-payment-token-2",
		Transferable: true,
	})
	if issueErr != nil {
		t.Fatalf("issue storage payment token: %v", issueErr)
	}
	_, redeemErr := NewLedger().RedeemStoragePaymentToken(frankCAS, issued.Bytes, ExpectedStoragePayment{
		Issuer:    "alice",
		Subject:   "frank",
		Scope:     "poc18-release-retention",
		ObjectCID: objectCID,
		Value:     5,
		Unit:      "storage_credit",
	}, "frank", now)
	if redeemErr == nil || !strings.Contains(redeemErr.Error(), "subject=mallory") {
		t.Fatalf("wrong subject should be rejected, got %v", redeemErr)
	}
}
