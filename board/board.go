package board

type Color int

const (
	White Color = iota
	Black
)

type Piece int

const (
	Empty Piece = iota
	Pawn
	Knight
	Bishop
	Rook
	Queen
	King
)

type Square struct {
	Piece Piece
	Color Color
}

type Board struct {
	Squares         [64]Square
	SideToMove      Color
	CastlingRights  string // e.g., "KQkq"
	EnPassantTarget int    // -1 if none
	HalfmoveClock   int
	FullmoveNumber  int
}

func (b *Board) GeneratePseudoLegalMoves() []Move {
	var moves []Move
	for i, square := range b.Squares {
		if square.Piece == Empty || square.Color != b.SideToMove {
			continue
		}

		switch square.Piece {
		case Pawn:
			moves = append(moves, b.generatePawnMoves(i)...)
		case Knight:
			moves = append(moves, b.generateKnightMoves(i)...)
		case Bishop:
			moves = append(moves, b.generateBishopMoves(i)...)
		case Rook:
			moves = append(moves, b.generateRookMoves(i)...)
		case Queen:
			moves = append(moves, b.generateQueenMoves(i)...)
		case King:
			moves = append(moves, b.generateKingMoves(i)...)
		}
	}
	return moves
}

func (b *Board) generatePawnMoves(from int) []Move {
	moves := []Move{}
	dir := 8
	startRank := 1
	promotionRank := 6

	if b.SideToMove == Black {
		dir = -8
		startRank = 6
		promotionRank = 1
	}

	to := from + dir
	if b.Squares[to].Piece == Empty {
		if to/8 == promotionRank {
			for _, p := range []Piece{Queen, Rook, Bishop, Knight} {
				moves = append(moves, Move{From: from, To: to, Promote: p})
			}
		} else {
			moves = append(moves, Move{From: from, To: to})
		}

		// Double move
		if from/8 == startRank &&
			b.Squares[from+dir].Piece == Empty &&
			b.Squares[from+2*dir].Piece == Empty {
			moves = append(moves, Move{From: from, To: from + 2*dir})
		}
	}

	// Captures
	for _, side := range []int{-1, 1} {
		capSq := from + dir + side
		if capSq < 0 || capSq >= 64 {
			continue
		}
		fromFile := from % 8
		toFile := (from + side) % 8
		if abs(fromFile-toFile) != 1 {
			continue
		}
		target := b.Squares[capSq]
		if target.Piece != Empty && target.Color != b.SideToMove {
			if capSq/8 == promotionRank {
				for _, p := range []Piece{Queen, Rook, Bishop, Knight} {
					moves = append(moves, Move{From: from, To: capSq, Promote: p, Capture: true})
				}
			} else {
				moves = append(moves, Move{From: from, To: capSq, Capture: true})
			}
		}
		// En passant
		if capSq == b.EnPassantTarget {
			behind := capSq - dir
			if behind >= 0 && behind < 64 {
				captured := b.Squares[behind]
				if captured.Piece == Pawn && captured.Color != b.SideToMove {
					moves = append(moves, Move{From: from, To: capSq, EnPassant: true})
				}
			}
		}
	}

	return moves
}

func (b *Board) generateKnightMoves(from int) []Move {
	offsets := []int{-17, -15, -10, -6, 6, 10, 15, 17}
	moves := []Move{}
	fromRank, fromFile := from/8, from%8

	for _, off := range offsets {
		to := from + off
		if to < 0 || to >= 64 {
			continue
		}
		toRank, toFile := to/8, to%8
		if abs(fromRank-toRank) > 2 || abs(fromFile-toFile) > 2 {
			continue
		}

		dest := b.Squares[to]
		if dest.Piece == Empty || dest.Color != b.SideToMove {
			moves = append(moves, Move{From: from, To: to, Capture: dest.Piece != Empty})
		}
	}

	return moves
}

func (b *Board) generateBishopMoves(from int) []Move {
	return b.generateSlidingMoves(from, []int{-9, -7, 7, 9})
}

func (b *Board) generateRookMoves(from int) []Move {
	return b.generateSlidingMoves(from, []int{-8, -1, 1, 8})
}

func (b *Board) generateQueenMoves(from int) []Move {
	return b.generateSlidingMoves(from, []int{-9, -8, -7, -1, 1, 7, 8, 9})
}

func (b *Board) generateKingMoves(from int) []Move {
	deltas := []int{-9, -8, -7, -1, 1, 7, 8, 9}
	moves := []Move{}
	fromRank, fromFile := from/8, from%8

	for _, d := range deltas {
		to := from + d
		if to < 0 || to >= 64 {
			continue
		}
		toRank, toFile := to/8, to%8
		if abs(toRank-fromRank) > 1 || abs(toFile-fromFile) > 1 {
			continue
		}
		dest := b.Squares[to]
		if dest.Piece == Empty || dest.Color != b.SideToMove {
			moves = append(moves, Move{From: from, To: to, Capture: dest.Piece != Empty})
		}
	}

	return moves
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (b *Board) IsInCheck(color Color) bool {
	kingIndex := -1
	for i, sq := range b.Squares {
		if sq.Piece == King && sq.Color == color {
			kingIndex = i
			break
		}
	}
	if kingIndex == -1 {
		return true // This shouldn't happen in a real game, but treat as invalid
	}

	// Check for attackers
	oppColor := White
	if color == White {
		oppColor = Black
	}

	attackers := b.attacksToSquare(kingIndex, oppColor)
	return len(attackers) > 0
}

func (b *Board) attacksToSquare(sq int, byColor Color) []int {
	attackers := []int{}

	for i, piece := range b.Squares {
		if piece.Color != byColor || piece.Piece == Empty {
			continue
		}

		var moves []Move
		switch piece.Piece {
		case Pawn:
			moves = b.generatePawnMoves(i)
		case Knight:
			moves = b.generateKnightMoves(i)
		case Bishop:
			moves = b.generateBishopMoves(i)
		case Rook:
			moves = b.generateRookMoves(i)
		case Queen:
			moves = b.generateQueenMoves(i)
		case King:
			moves = b.generateKingMoves(i)
		}

		for _, m := range moves {
			if m.To == sq {
				attackers = append(attackers, i)
			}
		}
	}

	return attackers
}

func (b *Board) Copy() *Board {
	newBoard := *b
	newBoard.Squares = b.Squares
	return &newBoard
}

func (b *Board) ApplyMove(m Move) {
    moving := b.Squares[m.From]

    // Basic move
    b.Squares[m.To] = moving
    b.Squares[m.From] = Square{Piece: Empty, Color: White} // or Black – doesn't matter

    // Promotion
    if m.Promote != Empty {
        b.Squares[m.To] = Square{Piece: m.Promote, Color: moving.Color}
    }

    // En Passant
    if m.EnPassant {
        var capSq int
        if moving.Color == White {
            capSq = m.To - 8
        } else {
            capSq = m.To + 8
        }
        b.Squares[capSq] = Square{Piece: Empty}
    }

    // Castling (not yet implemented)
    if m.Castling {
        // TODO: move rook
    }

    // Update en passant target
    if moving.Piece == Pawn && abs(m.To-m.From) == 16 {
        b.EnPassantTarget = (m.From + m.To) / 2
    } else {
        b.EnPassantTarget = -1
    }

    // Toggle side to move
    b.SideToMove = 1 - b.SideToMove
}

func (b *Board) GenerateLegalMoves() []Move {
	legal := []Move{}
	candidates := b.GeneratePseudoLegalMoves()

	for _, m := range candidates {
		copy := b.Copy()
		copy.ApplyMove(m)
		if !copy.IsInCheck(b.SideToMove) {
			legal = append(legal, m)
		}
	}

	return legal
}

func (b *Board) generateSlidingMoves(from int, deltas []int) []Move {
	moves := []Move{}
	fromRank := from / 8
	fromFile := from % 8

	for _, delta := range deltas {
		to := from + delta
		for to >= 0 && to < 64 {
			toRank := to / 8
			toFile := to % 8

			// Break if wrapped around edge
			if abs(toRank-fromRank) > 7 || abs(toFile-fromFile) > 7 {
				break
			}

			dest := b.Squares[to]
			if dest.Piece == Empty {
				moves = append(moves, Move{From: from, To: to})
			} else {
				if dest.Color != b.SideToMove {
					moves = append(moves, Move{From: from, To: to, Capture: true})
				}
				break
			}

			to += delta
		}
	}

	return moves
}
