package main

import (
	"fmt"
	"math/bits"
	"slices"
	"strconv"
)

type Color int
const (
	White Color = iota
	Black
)
var Colors = []Color {
	White,
	Black,
}
var colorName = map[Color]string{
    White: "White",
    Black: "Black",
}
func (c Color) String() string {
    return colorName[c]
}

type PieceType int
const (
	Pawn PieceType = iota
	Rook
	Knight
	Bishop
	Queen
	King
)
var pieceTypeName = map[PieceType]string{
    Pawn: "pawn",
	Rook: "rook",
	Knight: "knight",
	Bishop: "bishop",
	Queen: "queen",
	King: "king",
}
func (pt PieceType) String() string {
    return pieceTypeName[pt]
}

type Piece struct {
	color Color
	pieceType PieceType
	cc CartesianCoord
}
func (p Piece) String() string {
    return fmt.Sprintf("%v %v @ %v", p.color, p.pieceType, p.cc)
}

type Coord string
func (c Coord) IsValid() bool {
	if len(c) != 2 {
		return false
	}
	if c[0] < 'a' || c[0] > 'h' {
		return false
	}
	d, err := strconv.Atoi(string(c[1]))
	if err != nil {
		return false
	}
	if d < 1 || d > 8 {
		return false
	}
	return true
}
// Converts the Coord to a CartesianCoord. Assumes the Coord is valid.
func (c Coord) AsCartesianCoord() CartesianCoord {
	x := int(c[0] - 'a')
	y, err := strconv.Atoi(string(c[1]))
	if err != nil {
		panic(fmt.Sprintf("%v is not a valid Coord", c))
	}
	// x is already offset from 'a' as 0
	// y is parsed as 1-indexed so we transform to 0-indexed
	return CartesianCoord{x, y-1}
}

type CartesianCoord struct {
	X int
	Y int
}
func (cc CartesianCoord) String() string {
    return fmt.Sprintf("(%v,%v)", cc.X, cc.Y)
}
func (c CartesianCoord) IsValid() (isValid bool) {
	if c.X >= 0 && c.Y >= 0 && c.X <= 7 && c.Y <= 7 { 
		isValid = true
	}
	return
}
// Converts the CartesianCoord to a BitCoord. Assumes CartesianCoord is valid.
func (cc CartesianCoord) AsBitCoord() BitCoord {
	return BitCoord(0b1 << (cc.X+(8*cc.Y)))
}

type BitCoord uint64
func (bc BitCoord) String() string {
	return fmt.Sprintf("%b", bc)
}
func (bc BitCoord) IsValid() bool {
	return bc != 0 && (bc & (bc-1)) == 0
}
// Converts the BitCoord to a CartesianCoord. Assumes BitCoord is valid.
func (bc BitCoord) AsCartesianCoord() CartesianCoord {
	trailingZeroes := bits.TrailingZeros64(uint64(bc))
	return CartesianCoord{trailingZeroes-(trailingZeroes/8*8), trailingZeroes/8}
}
func (bc BitCoord) To(x, y int) BitCoord {
	var left, down bool
	if x < 0 {
		left = true
		x = 0 - x
	}
	if y < 0 {
		down = true
		y = 0 - y
	}
	if x > 7 || y > 7 {
		panic("Cannot transform to a position more than 7 spaces on a given axis")
	}

	result := bc

	if x != 0 {
		cc := bc.AsCartesianCoord()
		if left {
			result = result >> x
		} else {
			result = result << x
		}
		// Moving on the x axis is only allowed within the same row.
		newCc := result.AsCartesianCoord()
		if cc.Y != newCc.Y {
			return 0b0
		}
	}
	if y != 0 {
		if down {
			result = result >> (8*y)
		} else {
			result = result << (8*y)
		}
	}
	
	return result
}

type Move struct {
	piece Piece
	dest CartesianCoord
}

type ValidMove struct {
	piece Piece
	dest CartesianCoord
	newBoard Board
}

type Game struct {
	currentPlayer Color
	currentPlayerStatus string
	validMoves map[Piece][]ValidMove
	board Board
	moves []Move
}

// Board is a struct with no pointers to ensure cloning is easy.
type Board struct {
	players [2]Player
}

type Player struct {
	pieces [6]uint64
}

func NewGame() *Game {
	game := Game{
		currentPlayer: White,
		board: Board{
			players: [...]Player{
				{
					pieces: [...]uint64{
						0b11111111 << (8*1),
						0b10000001 << (8*0),
						0b01000010 << (8*0),
						0b00100100 << (8*0),
						0b00001000 << (8*0),
						0b00010000 << (8*0),
					},
				},
				{
					pieces: [...]uint64{
						0b11111111 << (8*6),
						0b10000001 << (8*7),
						0b01000010 << (8*7),
						0b00100100 << (8*7),
						0b00001000 << (8*7),
						0b00010000 << (8*7),
					},
				},
			},
		},
		moves: make([]Move, 0),
	}
	game.validMoves = game.PrecomputeValidMoves()
	return &game
}

func (g *Game) ExecuteValidMove(move ValidMove) bool {
	moves, found := g.validMoves[move.piece]
	if !found {
		return false
	}
	found = false
	for _, m := range moves {
		if m.dest == move.dest {
			found = true
			break
		}
	}
	if !found {
		return false
	}

	g.moves = append(g.moves, Move{piece: move.piece, dest: move.dest})
	g.board = move.newBoard
	g.currentPlayer = (g.currentPlayer + 1) % (Black + 1)
	g.validMoves = g.PrecomputeValidMoves()
	return true
}

// Takes in a Coord and returns a (Piece, bool). The Coord arg points to a position on the board. The Piece return value
// describes the piece located at the position specified (or a zero-valued Piece if there is no piece there). The bool
// return value describes whether a piece is located at the position or not
func (g *Game) GetCoord(cc CartesianCoord) (Piece, bool) {
	if !cc.IsValid() {
		return Piece{}, false
	}
	bit := uint64(cc.AsBitCoord())
	for i := range [2]Color{White, Black} {
		color := Color(i)
		player := g.board.players[color]
		if bit & player.pieces[Pawn] != 0 {
			return Piece{color, Pawn, cc}, true
		} else if bit & player.pieces[Rook] != 0 {
			return Piece{color, Rook, cc}, true
		} else if bit & player.pieces[Knight] != 0 {
			return Piece{color, Knight, cc}, true
		} else if bit & player.pieces[Bishop] != 0 {
			return Piece{color, Bishop, cc}, true
		} else if bit & player.pieces[Queen] != 0 {
			return Piece{color, Queen, cc}, true
		} else if bit & player.pieces[King] != 0 {
			return Piece{color, King, cc}, true
		}
	}
	return Piece{}, false
}

func (g *Game) PrecomputeValidMoves() map[Piece][]ValidMove {
	player := g.board.players[g.currentPlayer]
	moves := make(map[Piece][]ValidMove, 0)

	// TODO: Instead of ranging over each piece on the board, it is more efficient to instead get the number of trailing
	// zeroes of each non-zero bittable and use that to successively determine the position of each known piece.
	for y := range 8 {
		for x := range 8 {
			cc := CartesianCoord{x, y}
			bit := uint64(cc.AsBitCoord())

			switch  {
			case bit & player.pieces[Pawn] != 0:
				p := Piece{g.currentPlayer, Pawn, cc}
				moves[p] = g.computeValidMovesForPawn(p)
			case bit & player.pieces[Rook] != 0:
				p := Piece{g.currentPlayer, Rook, cc}
				moves[p] = g.computeValidMovesForRook(p)
			case bit & player.pieces[Knight] != 0:
				p := Piece{g.currentPlayer, Knight, cc}
				moves[p] = g.computeValidMovesForKnight(p)
			case bit & player.pieces[Bishop] != 0:
				p := Piece{g.currentPlayer, Bishop, cc}
				moves[p] = g.computeValidMovesForBishop(p)
			case bit & player.pieces[Queen] != 0:
				p := Piece{g.currentPlayer, Queen, cc}
				moves[p] = g.computeValidMovesForQueen(p)
			case bit & player.pieces[King] != 0:
				p := Piece{g.currentPlayer, King, cc}
				moves[p] = g.computeValidMovesForKing(p)
			}
		}
	}
	return moves
}

func (g *Game) hasPieceMoved(piece Piece) bool {
		hasMoved := false
		for _, move := range slices.Backward(g.moves) {
			if move.piece.color == piece.color && move.dest == piece.cc {
				hasMoved = true
				break
			}
		}
		return hasMoved
}

func (g *Game) GetValidMovesForPiece(p Piece) []ValidMove {
	return g.validMoves[p]
}

func checkDirection (g *Game, p Piece, x int, y int, onlyOne bool, requiresCapture bool, requiresMove bool) []ValidMove {
	pos := p.cc.AsBitCoord()
	moves := make([]ValidMove, 0)
	var next BitCoord = pos.To(x, y)
	for {
		if next == 0 {
			// out of board, end this dir
			break
		}
		destPiece, found := g.GetCoord(next.AsCartesianCoord())
		if !found  {
			if !requiresCapture {
				// no piece at target position; add move and continue
				m := ValidMove{
					piece: p,
					dest: next.AsCartesianCoord(),
					newBoard: g.board,
				}
				m.newBoard.players[p.color].pieces[p.pieceType] = m.newBoard.players[p.color].pieces[p.pieceType] &^ uint64(pos) | uint64(next)
				moves = append(moves, m)
			}
			if onlyOne {
				break
			}
		} else {
			if destPiece.color == p.color {
				// piece of current player at target pos, end this dir
				break
			} else if !requiresMove {
				// piece of another player at target pos, add move and change other player state as well and end this dir
				m := ValidMove{
					piece: p,
					dest: next.AsCartesianCoord(),
					newBoard: g.board,
				}
				m.newBoard.players[p.color].pieces[p.pieceType] = m.newBoard.players[p.color].pieces[p.pieceType] &^ uint64(pos) | uint64(next)
				m.newBoard.players[destPiece.color].pieces[destPiece.pieceType] = m.newBoard.players[destPiece.color].pieces[destPiece.pieceType] &^ uint64(next)
				moves = append(moves, m)
			}
			break
		}
		next = next.To(x, y)
	}
	return moves
}

func (g *Game) computeValidMovesForPawn(p Piece) []ValidMove {
	// TODO: en passant
	// TODO: promotion
	moves := make([]ValidMove, 0)

	switch p.color {
	case White:
		moves = append(moves, checkDirection(g, p, 0, 1, true, false, true)...)
		if !g.hasPieceMoved(p) {
			moves = append(moves, checkDirection(g, p, 0, 2, true, false, true)...)
		}
		moves = append(moves, checkDirection(g, p, -1, 1, true, true, false)...)
		moves = append(moves, checkDirection(g, p, 1, 1, true, true, false)...)
	case Black:
		moves = append(moves, checkDirection(g, p, 0, -1, true, false, true)...)
		if !g.hasPieceMoved(p) {
			moves = append(moves, checkDirection(g, p, 0, -2, true, false, true)...)
		}
		moves = append(moves, checkDirection(g, p, -1, -1, true, true, false)...)
		moves = append(moves, checkDirection(g, p, 1, -1, true, true, false)...)
	} 

	return moves
}

func (g *Game) computeValidMovesForRook(p Piece) []ValidMove {
	moves := make([]ValidMove, 0)
	moves = append(moves, checkDirection(g, p, 0, 1, false, false, false)...)
	moves = append(moves, checkDirection(g, p, 0, -1, false, false, false)...)
	moves = append(moves, checkDirection(g, p, 1, 0, false, false, false)...)
	moves = append(moves, checkDirection(g, p, -1, 0, false, false, false)...)
	return moves
}

func (g *Game) computeValidMovesForKnight(p Piece) []ValidMove {
	moves := make([]ValidMove, 0)
	moves = append(moves, checkDirection(g, p, 1, -2, true, false, false)...)
	moves = append(moves, checkDirection(g, p, 1, 2, true, false, false)...)
	moves = append(moves, checkDirection(g, p, 2, -1, true, false, false)...)
	moves = append(moves, checkDirection(g, p, 2, 1, true, false, false)...)
	moves = append(moves, checkDirection(g, p, -1, -2, true, false, false)...)
	moves = append(moves, checkDirection(g, p, -1, 2, true, false, false)...)
	moves = append(moves, checkDirection(g, p, -2, -1, true, false, false)...)
	moves = append(moves, checkDirection(g, p, -2, 1, true, false, false)...)
	return moves
}

func (g *Game) computeValidMovesForBishop(p Piece) []ValidMove {
	moves := make([]ValidMove, 0)
	moves = append(moves, checkDirection(g, p, 1, 1, false, false, false)...)
	moves = append(moves, checkDirection(g, p, 1, -1, false, false, false)...)
	moves = append(moves, checkDirection(g, p, -1, 1, false, false, false)...)
	moves = append(moves, checkDirection(g, p, -1, -1, false, false, false)...)
	return moves
}

func (g *Game) computeValidMovesForQueen(p Piece) []ValidMove {
	moves := make([]ValidMove, 0)
	moves = append(moves, checkDirection(g, p, 0, 1, false, false, false)...)
	moves = append(moves, checkDirection(g, p, 0, -1, false, false, false)...)
	moves = append(moves, checkDirection(g, p, 1, 0, false, false, false)...)
	moves = append(moves, checkDirection(g, p, -1, 0, false, false, false)...)
	moves = append(moves, checkDirection(g, p, 1, 1, false, false, false)...)
	moves = append(moves, checkDirection(g, p, 1, -1, false, false, false)...)
	moves = append(moves, checkDirection(g, p, -1, 1, false, false, false)...)
	moves = append(moves, checkDirection(g, p, -1, -1, false, false, false)...)
	return moves
}

func (g *Game) computeValidMovesForKing(p Piece) []ValidMove {
	// TODO: Castling
	moves := make([]ValidMove, 0)
	moves = append(moves, checkDirection(g, p, 1, 0, true, false, false)...)
	moves = append(moves, checkDirection(g, p, -1, 0, true, false, false)...)
	moves = append(moves, checkDirection(g, p, 0, 1, true, false, false)...)
	moves = append(moves, checkDirection(g, p, 0, -1, true, false, false)...)
	return moves
}
