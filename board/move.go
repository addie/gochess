package board

type Move struct {
    From     int
    To       int
    Promote  Piece // For pawn promotions; Empty if not a promotion
    Capture  bool
    EnPassant bool
    Castling bool
}

func SquareToCoord(sq int) string {
    file := sq % 8
    rank := sq / 8
    return string(rune('a'+file)) + string(rune('1'+rank))
}