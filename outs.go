package poker

import "fmt"

// CalculateOuts returns the cards from the remaining deck that would
// improve the hand type of the given hole cards + board combination.
// holeCards must be exactly 2 cards. board must be 3, 4, or 5 cards.
// Returns the list of out cards and their count.
func CalculateOuts(holeCards []Card, board []Card) ([]Card, error) {
	if len(holeCards) != 2 {
		return nil, fmt.Errorf("holeCards must be exactly 2 cards, got %d", len(holeCards))
	}
	if len(board) < 3 || len(board) > 5 {
		return nil, fmt.Errorf("board must be 3, 4, or 5 cards, got %d", len(board))
	}

	if len(board) == 5 {
		return nil, nil
	}

	baseCards := make([]Card, len(holeCards)+len(board))
	copy(baseCards, holeCards)
	copy(baseCards[len(holeCards):], board)

	currentType := evaluateHandType(baseCards)

	deck := NewDeck()
	for _, c := range holeCards {
		deck.RemoveCard(c)
	}
	for _, c := range board {
		deck.RemoveCard(c)
	}

	var outs []Card
	candidateCards := make([]Card, len(baseCards)+1)

	for _, card := range deck.Cards {
		copy(candidateCards, baseCards)
		candidateCards[len(baseCards)] = card
		candidateType := evaluateHandType(candidateCards)
		if candidateType > currentType {
			outs = append(outs, card)
		}
	}

	return outs, nil
}

// evaluateHandType returns the HandType for the given cards.
// Uses NewBestMadeHand for 7-card hands (O(1) lookup table),
// and falls back to evaluate() for other card counts.
// evaluate() modifies slices in place, so a copy is made to protect the caller.
func evaluateHandType(cards []Card) HandType {
	if len(cards) == 7 {
		return NewBestMadeHand(cards).Type()
	}
	cardsCopy := make([]Card, len(cards))
	copy(cardsCopy, cards)
	handType, _, _ := evaluate(cardsCopy)
	return handType
}
