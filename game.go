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
// Returns (x, y) movement relative to self
func (c Color) Forward() (int, int) {
	switch c {
	case White:
		return 0, 1
	case Black:
		return 0, -1
	default:
		return 0, 0
	}
}
func (c Color) Backward() (int, int) {
	switch c {
	case White:
		return 0, -1
	case Black:
		return 0, 1
	default:
		return 0, 0
	}
}
func (c Color) Left() (int, int) {
	switch c {
	case White:
		return -1, 0
	case Black:
		return 1, 0
	default:
		return 0, 0
	}
}
func (c Color) Right() (int, int) {
	switch c {
	case White:
		return 1, 0
	case Black:
		return -1, 0
	default:
		return 0, 0
	}
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

type SpecialMove int
const (
	None SpecialMove = iota
	EnPassant
	Castling
)
var specialMoveName = map[SpecialMove]string{
	None: "n/a",
    EnPassant: "en passant",
    Castling: "castling",
}
func (sm SpecialMove) String() string {
	return specialMoveName[sm]
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
func (cc CartesianCoord) AsCoord() Coord {
	return Coord(fmt.Sprintf("%v%v", string(rune('a' + cc.X)), string(rune('1' + cc.Y))))
}
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
	specialMove SpecialMove
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
	game.validMoves = computeValidMoves(game.currentPlayer, game.board, game.moves, true)
	return &game
}

func getCheckThreats(color Color, b Board, gameMoves []Move) []Piece {
	threats := make([]Piece, 0)
	kingcc := BitCoord(b.players[color].pieces[King]).AsCartesianCoord()
	for c := range Colors {
		c := Color(c)
		if color == c {
			continue
		}
		moves := computeValidMoves(c, b, gameMoves, false)
		for p, moves := range moves {
			for _, m := range moves {
				if m.dest == kingcc {
					threats = append(threats, p)
				}
			}
		}
	}
	return threats
}

// Mutates game state to match the chosen move to execute. Returns a human readable representation of the move that was
// executed, and whether it was executed or not.
func (g *Game) ExecuteValidMove(move ValidMove) (string, bool) {
	moves, found := g.validMoves[move.piece]
	if !found {
		return "", false
	}
	found = false
	for _, m := range moves {
		if m.dest == move.dest {
			found = true
			break
		}
	}
	if !found {
		return "", false
	}

	g.moves = append(g.moves, Move{piece: move.piece, dest: move.dest})
	notes := ""
	if origPiece, found := GetCoord(move.dest, g.board); found {
		notes = fmt.Sprintf(" [cap %v]", origPiece.pieceType)
	} else if move.specialMove == EnPassant {
		notes = fmt.Sprintf(" [%v cap %v]", move.specialMove, Pawn)
	} else if move.specialMove != None {
		notes = fmt.Sprintf(" [%v]", move.specialMove)
	}
	moveText := fmt.Sprintf("%v %v %v to %v%v", move.piece.color, move.piece.pieceType, move.piece.cc.AsCoord(), move.dest.AsCoord(), notes)
	g.board = move.newBoard
	g.currentPlayer = Color(int(g.currentPlayer + 1) % len(Colors))
	threats := getCheckThreats(g.currentPlayer, g.board, g.moves)
	g.validMoves = computeValidMoves(g.currentPlayer, g.board, g.moves, true)

	inCheck, noMoves := false, true
	if len(threats) > 0 {
		inCheck = true
	}
	for _, pieceMoves := range g.validMoves {
		if len(pieceMoves) > 0 {
			noMoves = false
			break
		}
	}
	if inCheck && noMoves {
		g.currentPlayerStatus = "CHECKMATE"
	} else if inCheck {
		g.currentPlayerStatus = "CHECK"
	} else if noMoves {
		g.currentPlayerStatus = "DRAW"
	} else {
		g.currentPlayerStatus = ""
	}
	return moveText, true
}

// Takes in a Coord and returns a (Piece, bool). The Coord arg points to a position on the board. The Piece return value
// describes the piece located at the position specified (or a zero-valued Piece if there is no piece there). The bool
// return value describes whether a piece is located at the position or not
func GetCoord(cc CartesianCoord, b Board) (Piece, bool) {
	if !cc.IsValid() {
		return Piece{}, false
	}
	bit := uint64(cc.AsBitCoord())
	for i := range [2]Color{White, Black} {
		color := Color(i)
		player := b.players[color]
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

func computeValidMoves(color Color, board Board, gameMoves []Move, removeIfIntoCheck bool) map[Piece][]ValidMove {
	player := board.players[color]
	moves := make(map[Piece][]ValidMove, 0)

	// TODO: Instead of ranging over each piece on the board, it is more efficient to instead get the number of trailing
	// zeroes of each non-zero bit table and use that to successively determine the position of each piece/bit.
	for y := range 8 {
		for x := range 8 {
			cc := CartesianCoord{x, y}
			bit := uint64(cc.AsBitCoord())

			switch  {
			case bit & player.pieces[Pawn] != 0:
				p := Piece{color, Pawn, cc}
				moves[p] = computeValidMovesForPawn(p, board, gameMoves)
			case bit & player.pieces[Rook] != 0:
				p := Piece{color, Rook, cc}
				moves[p] = computeValidMovesForRook(p, board)
			case bit & player.pieces[Knight] != 0:
				p := Piece{color, Knight, cc}
				moves[p] = computeValidMovesForKnight(p, board)
			case bit & player.pieces[Bishop] != 0:
				p := Piece{color, Bishop, cc}
				moves[p] = computeValidMovesForBishop(p, board)
			case bit & player.pieces[Queen] != 0:
				p := Piece{color, Queen, cc}
				moves[p] = computeValidMovesForQueen(p, board)
			case bit & player.pieces[King] != 0:
				p := Piece{color, King, cc}
				moves[p] = computeValidMovesForKing(p, board, gameMoves)
			}
		}
	}
	if removeIfIntoCheck {
		for piece, pieceMoves := range moves {
			n := 0
			for _, move := range pieceMoves {
				if len(getCheckThreats(color, move.newBoard, gameMoves)) == 0 {
					pieceMoves[n] = move
					n++
				}
			}
			moves[piece] = pieceMoves[:n]
		}
	}
	return moves
}

func hasPieceMoved(piece Piece, gameMoves []Move) bool {
	return len(movesOfPiece(piece, gameMoves)) > 0
}

// Returns the Moves of the given piece taken in this game in reverse order.
func movesOfPiece(piece Piece, gameMoves []Move) []Move {
	moves := make([]Move, 0)
	for _, m := range slices.Backward(gameMoves) {
		if m.dest == piece.cc {
			moves = append(moves, m)
			piece = m.piece
		}
	}
	return moves
}

func (g *Game) GetValidMovesForPiece(p Piece) []ValidMove {
	return g.validMoves[p]
}

func checkDirection (p Piece, b Board, x int, y int, onlyOne bool, requiresCapture bool, requiresMove bool) (moves []ValidMove) {
	pos := p.cc.AsBitCoord()
	moves = make([]ValidMove, 0)
	var next BitCoord = pos.To(x, y)
	for {
		if next == 0 {
			// out of board, end this dir
			break
		}
		destPiece, found := GetCoord(next.AsCartesianCoord(), b)
		if !found  {
			if !requiresCapture {
				// no piece at target position; add move and continue
				m := ValidMove{
					piece: p,
					dest: next.AsCartesianCoord(),
					newBoard: b,
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
					newBoard: b,
				}
				m.newBoard.players[p.color].pieces[p.pieceType] = m.newBoard.players[p.color].pieces[p.pieceType] &^ uint64(pos) | uint64(next)
				m.newBoard.players[destPiece.color].pieces[destPiece.pieceType] = m.newBoard.players[destPiece.color].pieces[destPiece.pieceType] &^ uint64(next)
				moves = append(moves, m)
			}
			break
		}
		next = next.To(x, y)
	}
	return
}

func checkEnPassant(p Piece, b Board, gameMoves []Move) (moves []ValidMove) {
	moves = make([]ValidMove, 0)
	pieceMoves := movesOfPiece(p, gameMoves)
	if len(pieceMoves) != 3 {
		// has not moved 3 times exactly
		return
	}
	if slices.ContainsFunc(pieceMoves, func (m Move) bool {
		delta := m.dest.Y - m.piece.cc.Y
		if delta < 0 {
			delta = 0 - delta
		}
		return delta > 1
	}) {
		// has moved 2 squares at once
		return
	}

	lastMove := gameMoves[len(gameMoves)-1]
	yDelta := lastMove.dest.Y - lastMove.piece.cc.Y
	if yDelta < 0 {
		yDelta = 0 - yDelta
	}
	if lastMove.piece.pieceType != Pawn || yDelta != 2 {
		// the last moved piece was either not a pawn, or did not advance two squares at once
		return
	}

	pos := p.cc.AsBitCoord()
	lastMovePos := lastMove.dest.AsBitCoord()
	left, right := pos.To(p.color.Left()), pos.To(p.color.Right())
	if lastMovePos == left || lastMovePos == right {
		dest := lastMove.dest.AsBitCoord().To(p.color.Forward())
		m := ValidMove{
			piece: p,
			dest: dest.AsCartesianCoord(),
			newBoard: b,
			specialMove: EnPassant,
		}
		m.newBoard.players[p.color].pieces[p.pieceType] = m.newBoard.players[p.color].pieces[p.pieceType] &^ uint64(pos) | uint64(dest)
		m.newBoard.players[lastMove.piece.color].pieces[lastMove.piece.pieceType] = m.newBoard.players[lastMove.piece.color].pieces[lastMove.piece.pieceType] &^ uint64(lastMove.dest.AsBitCoord())
		moves = append(moves, m)
	}
	return
}

func checkCastle(p Piece, b Board, gameMoves []Move, dirX int, dirY int) (moves []ValidMove) {
	moves = make([]ValidMove, 0)
	pieceMoves := movesOfPiece(p, gameMoves)
	if len(pieceMoves) != 0 {
		return
	}

	pos := p.cc.AsBitCoord()
	var next BitCoord = pos.To(dirX, dirY)
	for {
		if next == 0 {
			break
		}
		pairPiece, found := GetCoord(next.AsCartesianCoord(), b)
		if found {
			if pairPiece.pieceType == Rook && len(movesOfPiece(pairPiece, gameMoves)) == 0 {
				pairPiecePos := pairPiece.cc.AsBitCoord()
				pairPieceDest := p.cc.AsBitCoord().To(dirX, dirY)
				dest := p.cc.AsBitCoord().To(dirX*2, dirY*2)
				m := ValidMove{
					piece: p,
					dest: dest.AsCartesianCoord(),
					newBoard: b,
					specialMove: Castling,
				}
				m.newBoard.players[p.color].pieces[p.pieceType] = m.newBoard.players[p.color].pieces[p.pieceType] &^ uint64(pos) | uint64(dest)
				m.newBoard.players[p.color].pieces[pairPiece.pieceType] = m.newBoard.players[p.color].pieces[pairPiece.pieceType] &^ uint64(pairPiecePos) | uint64(pairPieceDest)
				moves = append(moves, m)
			}
			break
		}
		next = next.To(dirX, dirY)
	}

	return moves
}

func computeValidMovesForPawn(p Piece, b Board, gameMoves []Move) []ValidMove {
	// TODO: Pawn promotion. Needs to change the ValidMove struct to allow for specifying the piece type to promote to,
	// and updating the UI logic to allow for it.
	moves := make([]ValidMove, 0)

	forwardX, forwardY := p.color.Forward() // 0,1
	leftX, leftY := p.color.Left() // -1,0
	rightX, rightY := p.color.Right() // 1,0
	moves = append(moves, checkDirection(p, b, forwardX, forwardY, true, false, true)...)
	if !hasPieceMoved(p, gameMoves) {
		moves = append(moves, checkDirection(p, b, forwardX*2, forwardY*2, true, false, true)...)
	}
	moves = append(moves, checkDirection(p, b, leftX+forwardX, leftY+forwardY, true, true, false)...)
	moves = append(moves, checkDirection(p, b, rightX+forwardX, rightY+forwardY, true, true, false)...)
	moves = append(moves, checkEnPassant(p, b, gameMoves)...)
	return moves
}

func computeValidMovesForRook(p Piece, b Board) []ValidMove {
	moves := make([]ValidMove, 0)
	moves = append(moves, checkDirection(p, b, 0, 1, false, false, false)...)
	moves = append(moves, checkDirection(p, b, 0, -1, false, false, false)...)
	moves = append(moves, checkDirection(p, b, 1, 0, false, false, false)...)
	moves = append(moves, checkDirection(p, b, -1, 0, false, false, false)...)
	return moves
}

func computeValidMovesForKnight(p Piece, b Board) []ValidMove {
	moves := make([]ValidMove, 0)
	moves = append(moves, checkDirection(p, b, 1, -2, true, false, false)...)
	moves = append(moves, checkDirection(p, b, 1, 2, true, false, false)...)
	moves = append(moves, checkDirection(p, b, 2, -1, true, false, false)...)
	moves = append(moves, checkDirection(p, b, 2, 1, true, false, false)...)
	moves = append(moves, checkDirection(p, b, -1, -2, true, false, false)...)
	moves = append(moves, checkDirection(p, b, -1, 2, true, false, false)...)
	moves = append(moves, checkDirection(p, b, -2, -1, true, false, false)...)
	moves = append(moves, checkDirection(p, b, -2, 1, true, false, false)...)
	return moves
}

func computeValidMovesForBishop(p Piece, b Board) []ValidMove {
	moves := make([]ValidMove, 0)
	moves = append(moves, checkDirection(p, b, 1, 1, false, false, false)...)
	moves = append(moves, checkDirection(p, b, 1, -1, false, false, false)...)
	moves = append(moves, checkDirection(p, b, -1, 1, false, false, false)...)
	moves = append(moves, checkDirection(p, b, -1, -1, false, false, false)...)
	return moves
}

func computeValidMovesForQueen(p Piece, b Board) []ValidMove {
	moves := make([]ValidMove, 0)
	moves = append(moves, checkDirection(p, b, 0, 1, false, false, false)...)
	moves = append(moves, checkDirection(p, b, 0, -1, false, false, false)...)
	moves = append(moves, checkDirection(p, b, 1, 0, false, false, false)...)
	moves = append(moves, checkDirection(p, b, -1, 0, false, false, false)...)
	moves = append(moves, checkDirection(p, b, 1, 1, false, false, false)...)
	moves = append(moves, checkDirection(p, b, 1, -1, false, false, false)...)
	moves = append(moves, checkDirection(p, b, -1, 1, false, false, false)...)
	moves = append(moves, checkDirection(p, b, -1, -1, false, false, false)...)
	return moves
}

func computeValidMovesForKing(p Piece, b Board, gameMoves []Move) []ValidMove {
	moves := make([]ValidMove, 0)
	moves = append(moves, checkDirection(p, b, 1, 0, true, false, false)...)
	moves = append(moves, checkDirection(p, b, -1, 0, true, false, false)...)
	moves = append(moves, checkDirection(p, b, 0, 1, true, false, false)...)
	moves = append(moves, checkDirection(p, b, 0, -1, true, false, false)...)
	leftX, leftY := p.color.Left()
	rightX, rightY := p.color.Right()
	moves = append(moves, checkCastle(p, b, gameMoves, leftX, leftY)...)
	moves = append(moves, checkCastle(p, b, gameMoves, rightX, rightY)...)
	return moves
}
