package main

import (
	"reflect"
	"testing"
)

func TestPlayNormalGravityPlacesTokenAtBottom(t *testing.T) {
	g := NewGame("easy", "A", "B")

	g.Play(0)

	if g.Board[0][g.Rows-1] != 1 {
		t.Fatalf("expected token to settle at bottom row index %d, got column state %#v", g.Rows-1, g.Board[0])
	}
}

func TestPlayNormalGravityStacksTokensInOrder(t *testing.T) {
	g := NewGame("easy", "A", "B")

	g.Play(0) // Player 1
	g.Play(0) // Player 2

	if g.Board[0][g.Rows-1] != 1 {
		t.Fatalf("expected bottom cell to remain the first token, got column state %#v", g.Board[0])
	}
	if g.Board[0][g.Rows-2] != 2 {
		t.Fatalf("expected second token to rest above the first one, got column state %#v", g.Board[0])
	}
}

func TestPlayInverseGravityStacksFromTop(t *testing.T) {
	g := NewGame("easy", "A", "B")
	g.GravityInverse = true

	g.Play(0) // Player 1
	g.Play(0) // Player 2

	if g.Board[0][0] != 1 {
		t.Fatalf("expected first token to stay at the top, got column state %#v", g.Board[0])
	}
	if g.Board[0][1] != 2 {
		t.Fatalf("expected second token to attach below the first token when gravity is inverted, got column state %#v", g.Board[0])
	}
}

func TestGravityToggleReordersColumnWithoutFlippingOrder(t *testing.T) {
	g := NewGame("easy", "A", "B")

	for i := 0; i < 5; i++ {
		g.Play(0)
	}

	if !g.GravityInverse {
		t.Fatalf("expected gravity to toggle after 5 turns")
	}

	expected := []int{1, 2, 1, 2, 1, 0}
	if !reflect.DeepEqual(g.Board[0], expected) {
		t.Fatalf("expected column to be %v after gravity toggle, got %v", expected, g.Board[0])
	}
}
