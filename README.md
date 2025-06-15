# GoChess - A Chess Engine in Go

A chess engine and move generation library written in pure Go with no external dependencies.

## Features

- **Complete Chess Move Generation**
  - All standard piece movements (Pawn, Knight, Bishop, Rook, Queen, King)
  - Special moves: castling, en passant, pawn promotion
  - Legal move validation with check detection
  - Attack square generation

- **Board Representation**
  - Efficient 64-square array representation
  - FEN (Forsyth-Edwards Notation) parsing support
  - Algebraic notation utilities (e.g., "e2" → square 12)

- **Move Validation**
  - Perft testing framework for move generation verification
  - Comprehensive test suite with standard chess positions

## Installation

```bash
git clone https://github.com/yourusername/gochess.git
cd gochess
go build
```

## Usage

### As a Command Line Tool

Run the basic perft test:

```bash
go run main.go
```

Run the comprehensive test suite:

```bash
go run test_perft.go
```

### As a Library

```go
import (
    "github.com/yourusername/gochess/board"
    "github.com/yourusername/gochess/engine"
)

// Create a new board from FEN
b, err := board.FromFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
if err != nil {
    log.Fatal(err)
}

// Generate all legal moves
moves := b.GenerateLegalMoves()

// Make a move
if len(moves) > 0 {
    newBoard := b.Copy()
    newBoard.MakeMove(moves[0])
}

// Run perft test
nodes := engine.Perft(b, 4)
fmt.Printf("Perft(4): %d nodes\n", nodes)
```

## Project Structure

```
gochess/
├── board/              # Core chess logic
│   ├── board.go       # Board representation and move generation
│   ├── move.go        # Move structure and utilities
│   ├── fen.go         # FEN parser
│   └── piece.go       # Piece-related constants and utilities
├── engine/            # Chess engine components
│   ├── perft.go       # Perft implementation
│   ├── evaluate.go    # Position evaluation (placeholder)
│   └── search.go      # Search algorithm (placeholder)
├── uci/               # UCI protocol implementation
│   └── uci.go         # UCI interface (placeholder)
├── main.go            # Example usage
└── test_perft.go      # Perft test suite
```

## API Reference

### Board Package

#### Types

- `Board` - Chess board representation
- `Move` - Move representation with from/to squares and flags
- `Color` - White (0) or Black (1)
- `PieceType` - King, Queen, Rook, Bishop, Knight, or Pawn

#### Key Functions

```go
// Create a board from FEN string
func FromFEN(fen string) (*Board, error)

// Generate all legal moves for the current position
func (b *Board) GenerateLegalMoves() []Move

// Apply a move to the board (modifies board in-place)
func (b *Board) MakeMove(m Move)

// Create a deep copy of the board
func (b *Board) Copy() *Board

// Check if the current side is in check
func (b *Board) IsInCheck() bool

// Get the piece at a given square
func (b *Board) PieceAt(sq int) (PieceType, Color)
```

### Engine Package

```go
// Run perft test to count leaf nodes at given depth
func Perft(b *board.Board, depth int) int
```

## Development Status

### Completed ✅

- Full legal move generation
- FEN parsing
- Basic board operations
- Perft testing framework

### In Progress 🚧

- Debugging edge cases in special moves (castling, en passant)
- Some perft tests failing on complex positions

### TODO 📋

- [ ] UCI protocol implementation
- [ ] Position evaluation function
- [ ] Alpha-beta search algorithm
- [ ] Opening book support
- [ ] Endgame tablebase support
- [ ] Time management
- [ ] Transposition tables

## Testing

The project uses perft (performance test) to validate move generation correctness:

```bash
# Run all perft tests
go run test_perft.go
```

Current test positions include:
- Starting position
- Kiwipete position (complex middle game)
- Various endgame positions
- Positions testing specific rules (en passant, castling, promotions)

## Contributing

Contributions are welcome! Areas that need work:

1. **Bug Fixes**: Help debug failing perft tests
2. **Engine Development**: Implement evaluation and search algorithms
3. **UCI Protocol**: Complete the UCI interface
4. **Performance**: Optimize move generation and board operations
5. **Documentation**: Improve code comments and examples

## License

[Add your license here]

## Acknowledgments

- Perft test positions from the chess programming wiki
- Inspired by various open-source chess engines