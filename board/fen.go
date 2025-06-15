package board

import (
	"fmt"
	"strings"
	"strconv"
)

func LoadFEN(fen string) (*Board, error) {
	parts := strings.Split(fen, " ")
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid FEN string: needs at least position and side to move")
	}

	b := &Board{
		EnPassantTarget: -1,
		HalfmoveClock: 0,
		FullmoveNumber: 1,
	}

	// 1. Parse board position
	rows := strings.Split(parts[0], "/")
	if len(rows) != 8 {
		return nil, fmt.Errorf("invalid board position")
	}

	for r := 0; r < 8; r++ {
        row := rows[r]
        file := 0
        for _, c := range row {
            if c >= '1' && c <= '8' {
                file += int(c - '0')
            } else {
                index := (7-r)*8 + file
                p, color := pieceFromFENChar(c)
                b.Squares[index] = Square{
					Piece: p, Color: color,
				}
                file++
            }
        }
    }

    // 2. Side to move
    if parts[1] == "w" {
        b.SideToMove = White
    } else {
        b.SideToMove = Black
    }

    // 3. Castling rights
    if len(parts) > 2 {
        b.CastlingRights = parts[2]
    } else {
        b.CastlingRights = "-"
    }

    // 4. En passant target
    if len(parts) > 3 && parts[3] != "-" {
        file := parts[3][0] - 'a'
        rank := parts[3][1] - '1'
        b.EnPassantTarget = int(rank*8 + file)
    }

    // 5. Halfmove clock
    if len(parts) > 4 {
        half, err := strconv.Atoi(parts[4])
        if err != nil {
            return nil, err
        }
        b.HalfmoveClock = half
    }

    // 6. Fullmove number
    if len(parts) > 5 {
        full, err := strconv.Atoi(parts[5])
        if err != nil {
            return nil, err
        }
        b.FullmoveNumber = full
    }

    return b, nil
}