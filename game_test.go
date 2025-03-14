package main

import (
	"fmt"
	"testing"
)

func TestCoord(t *testing.T) {
	var tests = []struct{
		c Coord
		wantIsValid bool
		wantCartesianCoord CartesianCoord
		wantBit uint64
	}{
		{"a1", true, CartesianCoord{0, 0}, 0b1},
		{"a8", true, CartesianCoord{0, 7}, 0b1_00000000_00000000_00000000_00000000_00000000_00000000_00000000},
		{"d4", true, CartesianCoord{3, 3}, 0b1000_00000000_00000000_00000000},
		{"d5", true, CartesianCoord{3, 4}, 0b1000_00000000_00000000_00000000_00000000},
		{"e4", true, CartesianCoord{4, 3}, 0b10000_00000000_00000000_00000000},
		{"e5", true, CartesianCoord{4, 4}, 0b10000_00000000_00000000_00000000_00000000},
		{"h1", true, CartesianCoord{7, 0}, 0b10000000},
		{"h4", true, CartesianCoord{7, 3}, 0b10000000_00000000_00000000_00000000},
		{"h8", true, CartesianCoord{7, 7}, 0b10000000_00000000_00000000_00000000_00000000_00000000_00000000_00000000},
		{"A8", false, CartesianCoord{}, 0},
		{"18", false, CartesianCoord{}, 0},
		{"a9", false, CartesianCoord{}, 0},
		{"Ra0", false, CartesianCoord{}, 0},
		{"a01", false, CartesianCoord{}, 0},
	}

	for _, tt := range tests {
		t.Run(string(tt.c), func(t *testing.T) {
			isValid := tt.c.IsValid()
			if isValid != tt.wantIsValid {
				t.Errorf("IsValid got %v, want %v", isValid, tt.wantIsValid)
			}

			if isValid {
				cc := tt.c.AsCartesianCoord()
				if cc != tt.wantCartesianCoord {
					t.Errorf("AsCartesianCoord got %v, want %v", cc, tt.wantCartesianCoord)
				}

				b := tt.c.AsCartesianCoord().AsBit()
				if b != tt.wantBit {
					t.Errorf("AsBit got %b, want %b", b, tt.wantBit)
				}
			}
		})
	}
}

func TestCartesianCoord(t *testing.T) {
	var tests = []struct{
		cc CartesianCoord
		wantIsValid bool
	}{
		{CartesianCoord{0, 0}, true},
		{CartesianCoord{7, 7}, true},
		{CartesianCoord{-1, 3}, false},
		{CartesianCoord{1, -3}, false},
		{CartesianCoord{8, 3}, false},
		{CartesianCoord{1, 8}, false},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt.cc), func(t *testing.T) {
			isValid := tt.cc.IsValid()
			if isValid != tt.wantIsValid {
				t.Errorf("IsValid got %v, want %v", isValid, tt.wantIsValid)
			}
		})
	}
}

func TestNewGame(t *testing.T) {
	g := NewGame()
	var tests = []struct{
		name string
		val uint64
		want uint64
	}{
		{"white pawns", g.board.white.pawns, 0b11111111_00000000},
		{"white rooks", g.board.white.rooks, 0b10000001},
		{"white knights", g.board.white.knights, 0b01000010},
		{"white bishops", g.board.white.bishops, 0b00100100},
		{"white queens", g.board.white.queens, 0b00001000},
		{"white king", g.board.white.king, 0b00010000},
		{"black pawns", g.board.black.pawns, 0b11111111_00000000_00000000_00000000_00000000_00000000_00000000},
		{"black rooks", g.board.black.rooks, 0b10000001_00000000_00000000_00000000_00000000_00000000_00000000_00000000},
		{"black knights", g.board.black.knights, 0b01000010_00000000_00000000_00000000_00000000_00000000_00000000_00000000},
		{"black bishops", g.board.black.bishops, 0b00100100_00000000_00000000_00000000_00000000_00000000_00000000_00000000},
		{"black queens", g.board.black.queens, 0b00001000_00000000_00000000_00000000_00000000_00000000_00000000_00000000},
		{"black king", g.board.black.king, 0b00010000_00000000_00000000_00000000_00000000_00000000_00000000_00000000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.val != tt.want {
				t.Errorf("got %b, want %b", tt.val, tt.want)
			}
		})
	}
}

func TestGetCoord(t *testing.T) {
	var tests = []struct{
		c Coord
		wantPiece Piece
		wantIsOccupied bool
	}{
		{"a1", Piece{White, Rook}, true},
		{"b1", Piece{White, Knight}, true},
		{"c1", Piece{White, Bishop}, true},
		{"d1", Piece{White, Queen}, true},
		{"e1", Piece{White, King}, true},
		{"f1", Piece{White, Bishop}, true},
		{"g1", Piece{White, Knight}, true},
		{"h1", Piece{White, Rook}, true},

		{"a2", Piece{White, Pawn}, true},
		{"b2", Piece{White, Pawn}, true},
		{"c2", Piece{White, Pawn}, true},
		{"d2", Piece{White, Pawn}, true},
		{"e2", Piece{White, Pawn}, true},
		{"f2", Piece{White, Pawn}, true},
		{"g2", Piece{White, Pawn}, true},
		{"h2", Piece{White, Pawn}, true},

		{"a3", Piece{}, false},
		{"b3", Piece{}, false},
		{"c3", Piece{}, false},
		{"d3", Piece{}, false},
		{"e3", Piece{}, false},
		{"f3", Piece{}, false},
		{"g3", Piece{}, false},
		{"h3", Piece{}, false},

		{"a4", Piece{}, false},
		{"b4", Piece{}, false},
		{"c4", Piece{}, false},
		{"d4", Piece{}, false},
		{"e4", Piece{}, false},
		{"f4", Piece{}, false},
		{"g4", Piece{}, false},
		{"h4", Piece{}, false},

		{"a5", Piece{}, false},
		{"b5", Piece{}, false},
		{"c5", Piece{}, false},
		{"d5", Piece{}, false},
		{"e5", Piece{}, false},
		{"f5", Piece{}, false},
		{"g5", Piece{}, false},
		{"h5", Piece{}, false},

		{"a6", Piece{}, false},
		{"b6", Piece{}, false},
		{"c6", Piece{}, false},
		{"d6", Piece{}, false},
		{"e6", Piece{}, false},
		{"f6", Piece{}, false},
		{"g6", Piece{}, false},
		{"h6", Piece{}, false},

		{"a7", Piece{Black, Pawn}, true},
		{"b7", Piece{Black, Pawn}, true},
		{"c7", Piece{Black, Pawn}, true},
		{"d7", Piece{Black, Pawn}, true},
		{"e7", Piece{Black, Pawn}, true},
		{"f7", Piece{Black, Pawn}, true},
		{"g7", Piece{Black, Pawn}, true},
		{"h7", Piece{Black, Pawn}, true},

		{"a8", Piece{Black, Rook}, true},
		{"b8", Piece{Black, Knight}, true},
		{"c8", Piece{Black, Bishop}, true},
		{"d8", Piece{Black, Queen}, true},
		{"e8", Piece{Black, King}, true},
		{"f8", Piece{Black, Bishop}, true},
		{"g8", Piece{Black, Knight}, true},
		{"h8", Piece{Black, Rook}, true},
	}
	g := NewGame()
	for _, tt := range tests {
		testname := fmt.Sprintf("%v", tt.c)
		t.Run(testname, func(t *testing.T) {
			p, isOccupied := g.GetCoord(tt.c.AsCartesianCoord())
			if tt.wantIsOccupied != isOccupied {
				t.Fatalf("isOccupied got %v, want %v", isOccupied, tt.wantIsOccupied)
			}
			if tt.wantIsOccupied {
				if tt.wantPiece != p {
					t.Errorf("piece got %v, want %v", p, tt.wantPiece)
				}
			}
		})
	}
}
