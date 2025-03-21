package main

import (
	"fmt"
	"testing"
	"slices"
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

				b := uint64(tt.c.AsCartesianCoord().AsBitCoord())
				if b != tt.wantBit {
					t.Errorf("AsBitCoord got %b, want %b", b, tt.wantBit)
				}
			}
		})
	}
}

func TestBitCoord(t *testing.T) {
	var tests = []struct{
		bc BitCoord
		wantIsValid bool
		wantCartesianCoord CartesianCoord
	}{
		{BitCoord(0b1_00000000_00000000_00000000_00000000_00000000_00000000_00000000), true, CartesianCoord{0, 7}},
		{BitCoord(0b10000000_00000000_00000000_00000000_00000000_00000000_00000000_00000000), true, CartesianCoord{7, 7}},
		{BitCoord(0b10000000), true, CartesianCoord{7, 0}},
		{BitCoord(0b1), true, CartesianCoord{0, 0}},
		{BitCoord(0b10000_00000000_00000000_00000000), true, CartesianCoord{4, 3}},
		{BitCoord(0b10000_00001000_00000000_00100000), false, CartesianCoord{}},
		{BitCoord(0b10000_00000000_00000000_00100000), false, CartesianCoord{}},
		{BitCoord(0b0), false, CartesianCoord{}},
	}

	for _, tt := range tests {
		t.Run(tt.bc.String(), func(t *testing.T) {
			isValid := tt.bc.IsValid()
			if isValid != tt.wantIsValid {
				t.Errorf("IsValid got %v, want %v", isValid, tt.wantIsValid)
			}

			if isValid {
				cc := tt.bc.AsCartesianCoord()
				if cc != tt.wantCartesianCoord {
					t.Errorf("AsCartesianCoord got %v, want %v", cc, tt.wantCartesianCoord)
				}
			}
		})
	}
}

func TestBitCoordTo(t *testing.T) {
	var tests = []struct{
		cc CartesianCoord
		x int
		y int
		wantBitCoord BitCoord
	}{
		{CartesianCoord{0, 7}, 2, -2, CartesianCoord{2, 5}.AsBitCoord()},
		{CartesianCoord{0, 4}, 0, 3, CartesianCoord{0, 7}.AsBitCoord()},
		{CartesianCoord{0, 6}, 0, -6, CartesianCoord{0, 0}.AsBitCoord()},
		{CartesianCoord{5, 4}, -3, 2, CartesianCoord{2, 6}.AsBitCoord()},
		{CartesianCoord{5, 4}, -5, 0, CartesianCoord{0, 4}.AsBitCoord()},
		{CartesianCoord{7, 3}, -2, 1, CartesianCoord{5, 4}.AsBitCoord()},
		{CartesianCoord{0, 7}, 0, 1, 0b0},
		{CartesianCoord{0, 4}, 0, 5, 0b0},
		{CartesianCoord{0, 4}, 4, 6, 0b0},
		{CartesianCoord{5, 4}, -6, 0, 0b0},
		{CartesianCoord{5, 4}, -6, 2, 0b0},
		{CartesianCoord{5, 4}, 3, 0, 0b0},
		{CartesianCoord{5, 4}, 3, 2, 0b0},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v To(%v,%v)", tt.cc, tt.x, tt.y), func(t *testing.T) {
			bc := tt.cc.AsBitCoord()
			result := bc.To(tt.x, tt.y)
			if result != tt.wantBitCoord {
				t.Errorf("To got %b, want %b", result, tt.wantBitCoord)
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
		{"white pawns", g.board.players[White].pieces[Pawn], 0b11111111_00000000},
		{"white rooks", g.board.players[White].pieces[Rook], 0b10000001},
		{"white knights", g.board.players[White].pieces[Knight], 0b01000010},
		{"white bishops", g.board.players[White].pieces[Bishop], 0b00100100},
		{"white queens", g.board.players[White].pieces[Queen], 0b00001000},
		{"white king", g.board.players[White].pieces[King], 0b00010000},
		{"black pawns", g.board.players[Black].pieces[Pawn], 0b11111111_00000000_00000000_00000000_00000000_00000000_00000000},
		{"black rooks", g.board.players[Black].pieces[Rook], 0b10000001_00000000_00000000_00000000_00000000_00000000_00000000_00000000},
		{"black knights", g.board.players[Black].pieces[Knight], 0b01000010_00000000_00000000_00000000_00000000_00000000_00000000_00000000},
		{"black bishops", g.board.players[Black].pieces[Bishop], 0b00100100_00000000_00000000_00000000_00000000_00000000_00000000_00000000},
		{"black queens", g.board.players[Black].pieces[Queen], 0b00001000_00000000_00000000_00000000_00000000_00000000_00000000_00000000},
		{"black king", g.board.players[Black].pieces[King], 0b00010000_00000000_00000000_00000000_00000000_00000000_00000000_00000000},
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
		{"a1", Piece{White, Rook, CartesianCoord{0,0}}, true},
		{"b1", Piece{White, Knight, CartesianCoord{1,0}}, true},
		{"c1", Piece{White, Bishop, CartesianCoord{2,0}}, true},
		{"d1", Piece{White, Queen, CartesianCoord{3,0}}, true},
		{"e1", Piece{White, King, CartesianCoord{4,0}}, true},
		{"f1", Piece{White, Bishop, CartesianCoord{5,0}}, true},
		{"g1", Piece{White, Knight, CartesianCoord{6,0}}, true},
		{"h1", Piece{White, Rook, CartesianCoord{7,0}}, true},

		{"a2", Piece{White, Pawn, CartesianCoord{0,1}}, true},
		{"b2", Piece{White, Pawn, CartesianCoord{1,1}}, true},
		{"c2", Piece{White, Pawn, CartesianCoord{2,1}}, true},
		{"d2", Piece{White, Pawn, CartesianCoord{3,1}}, true},
		{"e2", Piece{White, Pawn, CartesianCoord{4,1}}, true},
		{"f2", Piece{White, Pawn, CartesianCoord{5,1}}, true},
		{"g2", Piece{White, Pawn, CartesianCoord{6,1}}, true},
		{"h2", Piece{White, Pawn, CartesianCoord{7,1}}, true},

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

		{"a7", Piece{Black, Pawn, CartesianCoord{0,6}}, true},
		{"b7", Piece{Black, Pawn, CartesianCoord{1,6}}, true},
		{"c7", Piece{Black, Pawn, CartesianCoord{2,6}}, true},
		{"d7", Piece{Black, Pawn, CartesianCoord{3,6}}, true},
		{"e7", Piece{Black, Pawn, CartesianCoord{4,6}}, true},
		{"f7", Piece{Black, Pawn, CartesianCoord{5,6}}, true},
		{"g7", Piece{Black, Pawn, CartesianCoord{6,6}}, true},
		{"h7", Piece{Black, Pawn, CartesianCoord{7,6}}, true},

		{"a8", Piece{Black, Rook, CartesianCoord{0,7}}, true},
		{"b8", Piece{Black, Knight, CartesianCoord{1,7}}, true},
		{"c8", Piece{Black, Bishop, CartesianCoord{2,7}}, true},
		{"d8", Piece{Black, Queen, CartesianCoord{3,7}}, true},
		{"e8", Piece{Black, King, CartesianCoord{4,7}}, true},
		{"f8", Piece{Black, Bishop, CartesianCoord{5,7}}, true},
		{"g8", Piece{Black, Knight, CartesianCoord{6,7}}, true},
		{"h8", Piece{Black, Rook, CartesianCoord{7,7}}, true},
	}
	g := NewGame()
	for _, tt := range tests {
		testname := fmt.Sprintf("%v", tt.c)
		t.Run(testname, func(t *testing.T) {
			p, isOccupied := GetCoord(tt.c.AsCartesianCoord(), g.board)
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

// r n b q k b n r
// p p p p p p p p
// - - - - - - - -
// - - - - - - - -
// - - - - - - - -
// - - - - - - - -
// - - - - - - - -
// - - - - - - - -
// P P P P P P P P
// R N B Q K B N R
var startBoard Board = Board{
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
}

func TestPawnMoves (t *testing.T) {
	var tests = []struct{
		name string
		g Game
		p Piece
		wantMoves []ValidMove
	}{
		{
			"standard starting pawn",
			Game{
				currentPlayer: White,
				board: startBoard,
				moves: make([]Move, 0),
			},
			Piece{White, Pawn, CartesianCoord{0,1}},
			append(make([]ValidMove, 0),
				ValidMove{
					piece: Piece{White, Pawn, CartesianCoord{0,1}},
					dest: CartesianCoord{0,2},
					newBoard: Board{
						players: [...]Player{
							{
								pieces: [...]uint64{
									(0b11111110 << (8*1)) | (0b00000001 << (8*2)),
									startBoard.players[White].pieces[Rook],
									startBoard.players[White].pieces[Knight],
									startBoard.players[White].pieces[Bishop],
									startBoard.players[White].pieces[Queen],
									startBoard.players[White].pieces[King],
								},
							},
							{
								pieces: [...]uint64{
									startBoard.players[Black].pieces[Pawn],
									startBoard.players[Black].pieces[Rook],
									startBoard.players[Black].pieces[Knight],
									startBoard.players[Black].pieces[Bishop],
									startBoard.players[Black].pieces[Queen],
									startBoard.players[Black].pieces[King],
								},
							},
						},
					},
				},
				ValidMove{
					piece: Piece{White, Pawn, CartesianCoord{0,1}},
					dest: CartesianCoord{0,3},
					newBoard: Board{
						players: [...]Player{
							{
								pieces: [...]uint64{
									(0b11111110 << (8*1)) | (0b00000001 << (8*3)),
									startBoard.players[White].pieces[Rook],
									startBoard.players[White].pieces[Knight],
									startBoard.players[White].pieces[Bishop],
									startBoard.players[White].pieces[Queen],
									startBoard.players[White].pieces[King],
								},
							},
							{
								pieces: [...]uint64{
									startBoard.players[Black].pieces[Pawn],
									startBoard.players[Black].pieces[Rook],
									startBoard.players[Black].pieces[Knight],
									startBoard.players[Black].pieces[Bishop],
									startBoard.players[Black].pieces[Queen],
									startBoard.players[Black].pieces[King],
								},
							},
						},
					},
				},
			),
		},
		{
			"pawn that has moved",
			Game{
				currentPlayer: White,
				board: Board{
					players: [...]Player{
						{
							pieces: [...]uint64{
								(0b11111101 << (8*1)) | (0b00000010 << (8*2)),
								startBoard.players[White].pieces[Rook],
								startBoard.players[White].pieces[Knight],
								startBoard.players[White].pieces[Bishop],
								startBoard.players[White].pieces[Queen],
								startBoard.players[White].pieces[King],
							},
						},
						{
							pieces: [...]uint64{
								startBoard.players[Black].pieces[Pawn],
								startBoard.players[Black].pieces[Rook],
								startBoard.players[Black].pieces[Knight],
								startBoard.players[Black].pieces[Bishop],
								startBoard.players[Black].pieces[Queen],
								startBoard.players[Black].pieces[King],
							},
						},
					},
				},
				moves: append(make([]Move, 0),
					Move{
						piece: Piece{White, Pawn, CartesianCoord{1,1}},
						dest: CartesianCoord{1,2},
					},
				),
			},
			Piece{White, Pawn, CartesianCoord{1,2}},
			append(make([]ValidMove, 0),
				ValidMove{
					piece: Piece{White, Pawn, CartesianCoord{1,2}},
					dest: CartesianCoord{1,3},
					newBoard: Board{
						players: [...]Player{
							{
								pieces: [...]uint64{
									(0b11111101 << (8*1)) | (0b00000010 << (8*3)),
									startBoard.players[White].pieces[Rook],
									startBoard.players[White].pieces[Knight],
									startBoard.players[White].pieces[Bishop],
									startBoard.players[White].pieces[Queen],
									startBoard.players[White].pieces[King],
								},
							},
							{
								pieces: [...]uint64{
									startBoard.players[Black].pieces[Pawn],
									startBoard.players[Black].pieces[Rook],
									startBoard.players[Black].pieces[Knight],
									startBoard.players[Black].pieces[Bishop],
									startBoard.players[Black].pieces[Queen],
									startBoard.players[Black].pieces[King],
								},
							},
						},
					},
				},
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validMoves := computeValidMovesForPawn(tt.p, tt.g.board, tt.g.moves)
			if !slices.Equal(tt.wantMoves, validMoves) {
				t.Errorf("moves got %+v, want %+v", validMoves, tt.wantMoves)
			}
		})
	}
}
