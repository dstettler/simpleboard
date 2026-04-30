package handler

import (
	"github.com/google/uuid"
	"testing"
	"simpleboard/internal/repository"
	"simpleboard/internal/auth"
)

// --- colorFromID ---

func TestColorFromID(t *testing.T) {

	var u1 uint;
	var u2 uint;
	var gu1 uuid.UUID;
	var gu2 uuid.UUID;

	u1 = 1
	u2 = 2

	gu1, _ = uuid.Parse("92f9f208-91e0-4f54-a324-fd0cac7e6007")
	gu2, _ = uuid.Parse("82f9f208-91e0-4f54-a324-fd0cac7e6007")

	g1 := repository.Game {
		WhitePlayerID: u1,
		BlackPlayerID: u2,
	}

	g2 := repository.Game {
		WhiteGuestID:  gu1.String(),
		BlackGuestID:  gu2.String(),
	}

	g3 := repository.Game {
		WhitePlayerID: u1,
		BlackGuestID:  gu2.String(),
	}

	g4 := repository.Game {}

	a11 := auth.Claims {
		UserID: &u1,
	}
	a12 := auth.Claims {
		UserID: &u2,
	}
	a21 := auth.Claims {
		GuestID: &gu1,
	}
	a22 := auth.Claims {
		GuestID: &gu2,
	}

	cases := []struct {
		g repository.Game
		c auth.Claims
		wants string
		wantb bool
	}{
		{g1, a11, "w", true},
		{g1, a12, "b", true},
		{g1, a21, "", false},
		{g1, a22, "", false},
		{g2, a11, "", false},
		{g2, a12, "", false},
		{g2, a21, "w", true},
		{g2, a22, "b", true},
		{g3, a11, "w", true},
		{g3, a12, "", false},
		{g3, a21, "", false},
		{g3, a22, "b", true},
		{g4, a11, "", false},
		{g4, a12, "", false},
		{g4, a21, "", false},
		{g4, a22, "", false},
	}
	for _, tc := range cases {
		gots, gotb := colorFromID(&tc.g, &tc.c)
		if gots != tc.wants {
			t.Errorf("got %s, want %s", gots, tc.wants)
		}
		if gotb != tc.wantb {
			t.Errorf("got %t, want %t", gotb, tc.wantb)
		}
	}
}
