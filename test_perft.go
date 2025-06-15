package main

import (
    "fmt"
    "github.com/addie/gochess/board"
    "github.com/addie/gochess/engine"
)

type PerftTest struct {
    fen      string
    name     string
    expected []int // expected results for depths 1-6
}

func main() {
    tests := []PerftTest{
        {
            name: "Starting position",
            fen:  "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
            expected: []int{20, 400, 8902, 197281, 4865609, 119060324},
        },
        {
            name: "Kiwipete",
            fen:  "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq -",
            expected: []int{48, 2039, 97862, 4085603, 193690690},
        },
        {
            name: "Position 3",
            fen:  "8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - -",
            expected: []int{14, 191, 2812, 43238, 674624},
        },
        {
            name: "Position 4",
            fen:  "r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1",
            expected: []int{6, 264, 9467, 422333, 15833292},
        },
    }

    for _, test := range tests {
        fmt.Printf("\n%s\n", test.name)
        fmt.Printf("FEN: %s\n", test.fen)
        
        b, err := board.LoadFEN(test.fen)
        if err != nil {
            fmt.Printf("Error loading FEN: %v\n", err)
            continue
        }

        maxDepth := 4 // Limit to depth 4 for speed
        for depth := 1; depth <= maxDepth && depth <= len(test.expected); depth++ {
            if test.expected[depth-1] == 0 {
                break
            }
            
            result := engine.Perft(b, depth)
            expected := test.expected[depth-1]
            
            if result == expected {
                fmt.Printf("Depth %d: %d ✓\n", depth, result)
            } else {
                fmt.Printf("Depth %d: %d (expected %d) ✗\n", depth, result, expected)
            }
        }
    }
}