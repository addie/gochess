package engine

import (
    "github.com/addie/gochess/board"
)

func Perft(b *board.Board, depth int) int {
    if depth == 0 {
        return 1
    }

    moves := b.GenerateLegalMoves()
    nodes := 0

    for _, m := range moves {
        copy := b.Copy()
        copy.ApplyMove(m)
        nodes += Perft(copy, depth-1)
    }

    return nodes
}