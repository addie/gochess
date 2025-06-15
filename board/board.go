package board

import "strings"

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
	
	// Add castling moves if king can castle
	if b.SideToMove == White {
		if b.Squares[4].Piece == King && b.Squares[4].Color == White && !b.IsInCheck(White) {
			moves = append(moves, b.generateCastlingMoves(4)...)
		}
	} else {
		if b.Squares[60].Piece == King && b.Squares[60].Color == Black && !b.IsInCheck(Black) {
			moves = append(moves, b.generateCastlingMoves(60)...)
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
	if to >= 0 && to < 64 && b.Squares[to].Piece == Empty {
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

func (b *Board) generateCastlingMoves(from int) []Move {
	moves := []Move{}
	
	if b.SideToMove == White {
		// White castles from e1
		if from != 4 { // e1
			return moves
		}
		
		// Kingside castling
		if strings.Contains(b.CastlingRights, "K") {
			// Check if squares between king and rook are empty
			if b.Squares[5].Piece == Empty && b.Squares[6].Piece == Empty {
				// Check if rook is on h1
				if b.Squares[7].Piece == Rook && b.Squares[7].Color == White {
					// Check if king doesn't move through check
					if !b.IsSquareAttacked(5, Black) && !b.IsSquareAttacked(6, Black) {
						moves = append(moves, Move{From: 4, To: 6, Castling: true})
					}
				}
			}
		}
		
		// Queenside castling
		if strings.Contains(b.CastlingRights, "Q") {
			// Check if squares between king and rook are empty
			if b.Squares[1].Piece == Empty && b.Squares[2].Piece == Empty && b.Squares[3].Piece == Empty {
				// Check if rook is on a1
				if b.Squares[0].Piece == Rook && b.Squares[0].Color == White {
					// Check if king doesn't move through check
					if !b.IsSquareAttacked(2, Black) && !b.IsSquareAttacked(3, Black) {
						moves = append(moves, Move{From: 4, To: 2, Castling: true})
					}
				}
			}
		}
	} else {
		// Black castles from e8
		if from != 60 { // e8
			return moves
		}
		
		// Kingside castling
		if strings.Contains(b.CastlingRights, "k") {
			// Check if squares between king and rook are empty
			if b.Squares[61].Piece == Empty && b.Squares[62].Piece == Empty {
				// Check if rook is on h8
				if b.Squares[63].Piece == Rook && b.Squares[63].Color == Black {
					// Check if king doesn't move through check
					if !b.IsSquareAttacked(61, White) && !b.IsSquareAttacked(62, White) {
						moves = append(moves, Move{From: 60, To: 62, Castling: true})
					}
				}
			}
		}
		
		// Queenside castling
		if strings.Contains(b.CastlingRights, "q") {
			// Check if squares between king and rook are empty
			if b.Squares[57].Piece == Empty && b.Squares[58].Piece == Empty && b.Squares[59].Piece == Empty {
				// Check if rook is on a8
				if b.Squares[56].Piece == Rook && b.Squares[56].Color == Black {
					// Check if king doesn't move through check
					if !b.IsSquareAttacked(58, White) && !b.IsSquareAttacked(59, White) {
						moves = append(moves, Move{From: 60, To: 58, Castling: true})
					}
				}
			}
		}
	}
	
	return moves
}

func (b *Board) IsSquareAttacked(sq int, byColor Color) bool {
	return len(b.attacksToSquare(sq, byColor)) > 0
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

		var attacks bool
		switch piece.Piece {
		case Pawn:
			// Pawns attack diagonally, not where they move
			attacks = b.pawnAttacks(i, sq)
		case Knight:
			moves := b.generateKnightMoves(i)
			for _, m := range moves {
				if m.To == sq {
					attacks = true
					break
				}
			}
		case Bishop:
			moves := b.generateBishopMoves(i)
			for _, m := range moves {
				if m.To == sq {
					attacks = true
					break
				}
			}
		case Rook:
			moves := b.generateRookMoves(i)
			for _, m := range moves {
				if m.To == sq {
					attacks = true
					break
				}
			}
		case Queen:
			moves := b.generateQueenMoves(i)
			for _, m := range moves {
				if m.To == sq {
					attacks = true
					break
				}
			}
		case King:
			moves := b.generateKingMoves(i)
			for _, m := range moves {
				if m.To == sq {
					attacks = true
					break
				}
			}
		}

		if attacks {
			attackers = append(attackers, i)
		}
	}

	return attackers
}

func (b *Board) pawnAttacks(from int, target int) bool {
	dir := 8
	if b.Squares[from].Color == Black {
		dir = -8
	}
	
	// Check both diagonal attack squares
	for _, side := range []int{-1, 1} {
		attackSq := from + dir + side
		if attackSq == target {
			// Verify it's a valid diagonal (not wrapping around board edge)
			fromFile := from % 8
			targetFile := target % 8
			if abs(fromFile-targetFile) == 1 && abs(from/8-target/8) == 1 {
				return true
			}
		}
	}
	
	return false
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

    // Castling
    if m.Castling {
        if moving.Color == White {
            if m.To == 6 { // Kingside
                b.Squares[7] = Square{Piece: Empty}
                b.Squares[5] = Square{Piece: Rook, Color: White}
            } else if m.To == 2 { // Queenside
                b.Squares[0] = Square{Piece: Empty}
                b.Squares[3] = Square{Piece: Rook, Color: White}
            }
        } else {
            if m.To == 62 { // Kingside
                b.Squares[63] = Square{Piece: Empty}
                b.Squares[61] = Square{Piece: Rook, Color: Black}
            } else if m.To == 58 { // Queenside
                b.Squares[56] = Square{Piece: Empty}
                b.Squares[59] = Square{Piece: Rook, Color: Black}
            }
        }
    }

    // Update en passant target
    if moving.Piece == Pawn && abs(m.To-m.From) == 16 {
        b.EnPassantTarget = (m.From + m.To) / 2
    } else {
        b.EnPassantTarget = -1
    }

    // Update castling rights
    newRights := ""
    for _, c := range b.CastlingRights {
        keep := true
        switch c {
        case 'K':
            // White kingside - lost if king or h1 rook moves
            if m.From == 4 || m.From == 7 || m.To == 7 {
                keep = false
            }
        case 'Q':
            // White queenside - lost if king or a1 rook moves
            if m.From == 4 || m.From == 0 || m.To == 0 {
                keep = false
            }
        case 'k':
            // Black kingside - lost if king or h8 rook moves
            if m.From == 60 || m.From == 63 || m.To == 63 {
                keep = false
            }
        case 'q':
            // Black queenside - lost if king or a8 rook moves
            if m.From == 60 || m.From == 56 || m.To == 56 {
                keep = false
            }
        }
        if keep {
            newRights += string(c)
        }
    }
    if newRights == "" {
        b.CastlingRights = "-"
    } else {
        b.CastlingRights = newRights
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
		prevRank := fromRank
		prevFile := fromFile
		
		for to >= 0 && to < 64 {
			toRank := to / 8
			toFile := to % 8

			// Check if we've wrapped around the board
			// For horizontal moves (delta = -1 or 1), rank should stay the same
			// For vertical moves (delta = -8 or 8), file should stay the same
			// For diagonal moves, the rank and file differences should be consistent
			
			rankDiff := toRank - prevRank
			fileDiff := toFile - prevFile
			
			// Check if the move is consistent with the delta
			if delta == 1 || delta == -1 { // Horizontal
				if rankDiff != 0 { // Wrapped around edge
					break
				}
			} else if delta == 8 || delta == -8 { // Vertical
				if fileDiff != 0 { // Should stay on same file
					break
				}
			} else { // Diagonal
				// For diagonals, abs(rankDiff) should equal abs(fileDiff) and both should be 1
				if abs(rankDiff) != 1 || abs(fileDiff) != 1 {
					break
				}
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

			prevRank = toRank
			prevFile = toFile
			to += delta
		}
	}

	return moves
}
