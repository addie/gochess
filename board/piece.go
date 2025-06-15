package board

func pieceFromFENChar(c rune) (Piece, Color) {
    var color Color
    if c >= 'A' && c <= 'Z' {
        color = White
    } else {
        color = Black
    }

    switch c {
    case 'p', 'P':
        return Pawn, color
    case 'n', 'N':
        return Knight, color
    case 'b', 'B':
        return Bishop, color
    case 'r', 'R':
        return Rook, color
    case 'q', 'Q':
        return Queen, color
    case 'k', 'K':
        return King, color
    }
    return Empty, color
}