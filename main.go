package main

import (
    "fmt"
    "github.com/addie/gochess/board"
    "github.com/addie/gochess/engine"
)

func main() {
    fen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
    b, _ := board.LoadFEN(fen)

    total := 0
    for _, m := range b.GenerateLegalMoves() {
        copy := b.Copy()
        copy.ApplyMove(m)
        count := engine.Perft(copy, 2)
        fmt.Printf("%s -> %s: %d\n", board.SquareToCoord(m.From), board.SquareToCoord(m.To), count)
        total += count
    }

    fmt.Printf("Total = %d (expected 8902)\n", total)
}