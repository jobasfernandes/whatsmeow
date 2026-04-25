// Copyright (c) 2026 Tulir Asokan
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package whatsmeow

import (
	"errors"
	"fmt"
	"testing"
)

func TestClassifyServerError_KnownCodes(t *testing.T) {
	cases := []struct {
		code     int
		sentinel error
	}{
		{463, ErrPrivacyTokenMissing},
		{475, ErrMessageCapped},
		{479, ErrAddressingStale},
		{421, ErrServerRateLimited},
	}
	for _, tc := range cases {
		err := classifyServerError(tc.code)
		if !errors.Is(err, tc.sentinel) {
			t.Errorf("code %d: errors.Is(err, %v) = false; want true", tc.code, tc.sentinel)
		}
		if !errors.Is(err, ErrServerReturnedError) {
			t.Errorf("code %d: errors.Is(err, ErrServerReturnedError) = false; want true (backwards-compat)", tc.code)
		}
		var wae *WAServerError
		if !errors.As(err, &wae) {
			t.Errorf("code %d: errors.As(&WAServerError) = false; want true", tc.code)
			continue
		}
		if wae.Code != tc.code {
			t.Errorf("code %d: WAServerError.Code = %d; want %d", tc.code, wae.Code, tc.code)
		}
	}
}

func TestClassifyServerError_UnknownCode(t *testing.T) {
	err := classifyServerError(599)
	if !errors.Is(err, ErrServerReturnedError) {
		t.Fatalf("unknown code must still match ErrServerReturnedError")
	}
	if errors.Is(err, ErrPrivacyTokenMissing) {
		t.Fatalf("unknown code must not match ErrPrivacyTokenMissing")
	}
	var wae *WAServerError
	if !errors.As(err, &wae) || wae.Code != 599 {
		t.Fatalf("expected WAServerError{Code:599}, got %v", err)
	}
}

func TestWAServerError_WrappedByCaller(t *testing.T) {
	wrapped := fmt.Errorf("send failed: %w", classifyServerError(463))
	if !errors.Is(wrapped, ErrPrivacyTokenMissing) {
		t.Fatalf("wrapped error must still match ErrPrivacyTokenMissing through %%w chain")
	}
}
