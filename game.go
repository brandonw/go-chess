package main

import (
	"fmt"
	"strconv"
)

type Color int
const (
	White Color = iota
	Black
)
var colorName = map[Color]string{
    White: "white",
    Black: "black",
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
}
func (p Piece) String() string {
    return fmt.Sprintf("%v %v", p.color, p.pieceType)
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
func (cc CartesianCoord) AsBit() uint64 {
	return uint64(0b1 << (cc.X+(8*cc.Y)))
}

type Game struct {
	currentPlayer Color
	board Board
}

type Board struct {
	white Player
	black Player
}

type Player struct {
	pawns uint64
	rooks uint64
	knights uint64
	bishops uint64
	queens uint64
	king uint64

	caputuredPawns uint8
	caputuredRooks uint8
	caputuredKnights uint8
	caputuredBishops uint8
	caputuredQueens uint8
}

func NewGame() Game {
	game := Game{
		currentPlayer: White,
		board: Board{
			white: Player{
				pawns: 0b11111111 << (8*1),
				rooks: 0b10000001 << (8*0),
				knights: 0b01000010 << (8*0),
				bishops: 0b00100100 << (8*0),
				queens: 0b00001000 << (8*0),
				king: 0b00010000 << (8*0),
				caputuredPawns: 0,
				caputuredRooks: 0,
				caputuredKnights: 0,
				caputuredBishops: 0,
				caputuredQueens: 0,
			},
			black: Player{
				pawns: 0b11111111 << (8*6),
				rooks: 0b10000001 << (8*7),
				knights: 0b01000010 << (8*7),
				bishops: 0b00100100 << (8*7),
				queens: 0b00001000 << (8*7),
				king: 0b00010000 << (8*7),
				caputuredPawns: 0,
				caputuredRooks: 0,
				caputuredKnights: 0,
				caputuredBishops: 0,
				caputuredQueens: 0,
			},
		},
	}
	return game
}

// Takes in a Coord and returns a (Piece, bool). The Coord arg points to a position on the board. The Piece return value
// describes the piece located at the position specified (or a zero-valued Piece if there is no piece there). The bool
// return value describes whether a piece is located at the position or not
func (g Game) GetCoord(cc CartesianCoord) (p Piece, isOccupied bool) {
	if !cc.IsValid() {
		return
	}
	bit := cc.AsBit()
	if bit & g.board.white.pawns != 0 {
		p = Piece{White, Pawn}
		isOccupied  = true
	} else if bit & g.board.white.rooks != 0 {
		p = Piece{White, Rook}
		isOccupied  = true
	} else if bit & g.board.white.knights != 0 {
		p = Piece{White, Knight}
		isOccupied  = true
	} else if bit & g.board.white.bishops != 0 {
		p = Piece{White, Bishop}
		isOccupied  = true
	} else if bit & g.board.white.queens != 0 {
		p = Piece{White, Queen}
		isOccupied  = true
	} else if bit & g.board.white.king != 0 {
		p = Piece{White, King}
		isOccupied  = true
	} else if bit & g.board.black.pawns != 0 {
		p = Piece{Black, Pawn}
		isOccupied  = true
	} else if bit & g.board.black.rooks != 0 {
		p = Piece{Black, Rook}
		isOccupied  = true
	} else if bit & g.board.black.knights != 0 {
		p = Piece{Black, Knight}
		isOccupied  = true
	} else if bit & g.board.black.bishops != 0 {
		p = Piece{Black, Bishop}
		isOccupied  = true
	} else if bit & g.board.black.queens != 0 {
		p = Piece{Black, Queen}
		isOccupied  = true
	} else if bit & g.board.black.king != 0 {
		p = Piece{Black, King}
		isOccupied  = true
	}
	return
}
