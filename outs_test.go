package poker

import (
	"slices"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCalculateOuts(t *testing.T) {
	tests := []struct {
		name         string
		holeCards    []Card
		board        []Card
		wantCount    int
		wantCards    []Card
		wantContains []Card
		wantErr      bool
	}{
		{
			name: "flush draw from straight - 9 outs (only flush/straight-flush improve)",
			holeCards: []Card{
				{RankFive, Hearts},
				{RankSix, Hearts},
			},
			board: []Card{
				{RankSeven, Hearts},
				{RankEight, Hearts},
				{RankNine, Clubs},
			},
			wantCount: 9,
		},
		{
			name: "open-ended straight draw from high card - includes pair outs",
			holeCards: []Card{
				{RankEight, Hearts},
				{RankNine, Clubs},
			},
			board: []Card{
				{RankSeven, Diamonds},
				{RankSix, Spades},
				{RankDeuce, Hearts},
			},
			wantCount: 23,
			wantContains: []Card{
				{RankFive, Hearts},
				{RankFive, Clubs},
				{RankFive, Diamonds},
				{RankFive, Spades},
				{RankTen, Hearts},
				{RankTen, Clubs},
				{RankTen, Diamonds},
				{RankTen, Spades},
			},
		},
		{
			name: "pocket pair - two pair and set outs",
			holeCards: []Card{
				{RankJack, Hearts},
				{RankJack, Clubs},
			},
			board: []Card{
				{RankDeuce, Diamonds},
				{RankFive, Spades},
				{RankNine, Hearts},
			},
			wantCount: 11,
			wantContains: []Card{
				{RankJack, Diamonds},
				{RankJack, Spades},
			},
		},
		{
			name: "turn flush draw from straight - 9 outs (uses NewBestMadeHand for 7 cards)",
			holeCards: []Card{
				{RankFive, Hearts},
				{RankSix, Hearts},
			},
			board: []Card{
				{RankSeven, Hearts},
				{RankEight, Hearts},
				{RankNine, Clubs},
				{RankKing, Diamonds},
			},
			wantCount: 9,
		},
		{
			name: "royal flush - no outs",
			holeCards: []Card{
				{RankAce, Spades},
				{RankKing, Spades},
			},
			board: []Card{
				{RankQueen, Spades},
				{RankJack, Spades},
				{RankTen, Spades},
			},
		},
		{
			name: "river board 5 cards - no outs",
			holeCards: []Card{
				{RankAce, Hearts},
				{RankKing, Hearts},
			},
			board: []Card{
				{RankDeuce, Clubs},
				{RankThree, Diamonds},
				{RankFour, Spades},
				{RankSeven, Hearts},
				{RankNine, Clubs},
			},
		},
		{
			name: "invalid hole cards - 1 card",
			holeCards: []Card{
				{RankAce, Hearts},
			},
			board: []Card{
				{RankDeuce, Clubs},
				{RankThree, Diamonds},
				{RankFour, Spades},
			},
			wantErr: true,
		},
		{
			name: "invalid board - 6 cards",
			holeCards: []Card{
				{RankAce, Hearts},
				{RankKing, Hearts},
			},
			board: []Card{
				{RankDeuce, Clubs},
				{RankThree, Diamonds},
				{RankFour, Spades},
				{RankFive, Hearts},
				{RankSix, Clubs},
				{RankSeven, Diamonds},
			},
			wantErr: true,
		},
		{
			name: "invalid board - 2 cards",
			holeCards: []Card{
				{RankAce, Hearts},
				{RankKing, Hearts},
			},
			board: []Card{
				{RankDeuce, Clubs},
				{RankThree, Diamonds},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CalculateOuts(tt.holeCards, tt.board)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(got) != tt.wantCount {
				t.Errorf("outs count = %d, want %d", len(got), tt.wantCount)
			}

			if tt.wantCards != nil {
				if diff := cmp.Diff(tt.wantCards, got); diff != "" {
					t.Errorf("outs mismatch (-want +got):\n%s", diff)
				}
			}

			for _, wc := range tt.wantContains {
				if !slices.Contains(got, wc) {
					t.Errorf("expected outs to contain %v", wc)
				}
			}
		})
	}
}

func BenchmarkCalculateOuts(b *testing.B) {
	holeCards := []Card{
		{RankFive, Hearts},
		{RankSix, Hearts},
	}
	board := []Card{
		{RankSeven, Hearts},
		{RankEight, Hearts},
		{RankNine, Clubs},
	}

	for b.Loop() {
		_, _ = CalculateOuts(holeCards, board)
	}
}
