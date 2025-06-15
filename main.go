package main

import (
    "fmt"
    "github.com/addie/gochess/board"
    "github.com/addie/gochess/engine"
)

func main() {
    fen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
    b, _ := board.LoadFEN(fen)

    // First check perft(1) and perft(2) from starting position
    fmt.Printf("Perft(1) from starting position: %d (expected 20)\n", engine.Perft(b, 1))
    fmt.Printf("Perft(2) from starting position: %d (expected 400)\n", engine.Perft(b, 2))
    fmt.Printf("Perft(3) from starting position: %d (expected 8902)\n", engine.Perft(b, 3))
    fmt.Println()

    // Then show perft(2) for each move (which gives perft(3) breakdown)
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