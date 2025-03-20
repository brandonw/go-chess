package main

import (
	"fmt"
	"github.com/rivo/tview"
	"github.com/gdamore/tcell/v2"
	"os"
	"log"
)

type State struct {
	app *tview.Application
	pieceSet *PieceSet
	squareWidth int
	squareHeight int
	game *Game
	squares [][]*tview.TextView
	currentPlayer *tview.TextView
	currentPlayerStatus *tview.TextView
	greenBbgColor tcell.Color
	beigeBgColor tcell.Color
	whitePieceColor tcell.Color
	blackPieceColor tcell.Color
	squareDefaultGreenStyle tcell.Style
	squareDefaultBeigeStyle tcell.Style
	squareHighlightStyle tcell.Style
	squareValidMoveStyle tcell.Style
	logger *log.Logger
}

func UpdateBoardUi (state *State) {
	for y := range 8 {
		for x := range 8 {
			square := state.squares[7-y][x]
			bgColor := state.beigeBgColor
			if y % 2 == x % 2 {
			bgColor = state.greenBbgColor
			}
			square.SetBackgroundColor(bgColor)

			p, occupied := state.game.GetCoord(CartesianCoord{x, y})
			if occupied {
				color := state.blackPieceColor
				if p.color == White {
					color = state.whitePieceColor
				}
				// tview y=0 corresponds to the top row of the grid, not the bottom. We must transform between tview
				// coord to game coord where y=0 is the bottom row.
				square.SetTextColor(color)
				square.SetText(state.pieceSet.pieces[p.pieceType])
				topPadding := (state.squareHeight - state.pieceSet.minY) / 2
				leftPadding := (state.squareWidth - state.pieceSet.minX) / 2
				square.SetBorderPadding(topPadding, 0, leftPadding, 0)
			} else {
				square.SetText("")
			}
		}
	}

	state.currentPlayer.SetText(state.game.currentPlayer.String())
	state.currentPlayerStatus.SetText(state.game.currentPlayerStatus)
}

// Checks that the positions entered are valid and that they are owned by the current player. textToCheck will contain
// 1-4 runes, where 1-2 is the piece to move, and 3-4 is the destination. If the 2nd rune does not correspond to a piece
// owned by the current player, it will be rejected. If the 3rd or 4th rune does not corrrespond to a valid move the
// selected piece can make, it will be rejected.
func MoveChecker (textToCheck string, lastChar rune, state *State) bool {
	var p Piece

	// Quick checks for generally valid input first
	if len(textToCheck) % 2 == 1 && (lastChar < 'a' || lastChar > 'h') {
		return false
	}
	if len(textToCheck) % 2 == 0 && (lastChar < '1' || lastChar > '8') {
		return false
	}
	if len(textToCheck) > 4 {
		return false
	}

	if len(textToCheck) >= 2 {
		c := Coord(textToCheck[0:2])
		if !c.IsValid() {
			state.logger.Panicf("bug in move checking validation: %+v", c)
		}
		pos := c.AsCartesianCoord()
		// Check if the position has a piece owned by current player
		var occupied bool
		p, occupied = state.game.GetCoord(pos)
		if !occupied || p.color != state.game.currentPlayer {
			return false
		}
	}

	if len(textToCheck) >= 3 {
		validMoves := state.game.GetValidMovesForPiece(p)

		if len(textToCheck) == 3 {
			for _, v := range validMoves {
				if v.dest.X == int(textToCheck[2] - 'a') {
					return true
				}
			}
			return false
		}

		c := Coord(textToCheck[2:4])
		if !c.IsValid() {
			state.logger.Panicf("bug in move checking validation: %+v", c)
		}
		pos := c.AsCartesianCoord()
		for _, v := range validMoves {
			if v.dest.X == pos.X && v.dest.Y == pos.Y {
				return true
			}
		}
		return false
	}

	return true
}

// Updates UI with highlights for potential pieces, selected piece, and valid moves for selected piece.
func GridStateUpdater (text string, state *State) {
	validMoves := []ValidMove{}

	var px1, py1, px2, py2 int
	var cc1 CartesianCoord
	if len(text) == 1 {
		px1 = int(text[0]-'a')
	}
	if len(text) > 1 {
		c := Coord(text[0:2])
		if !c.IsValid() {
			state.logger.Panicf("bug in move checking validation: %+v", c)
		}
		cc1 = c.AsCartesianCoord()
		px1 = cc1.X
		py1 = 7 - cc1.Y
	}
	if len(text) >= 2 && len(text) < 4 {
		// we need valid moves if only the target piece was selected, or if the first part of the destination was
		// selected
		p, _ := state.game.GetCoord(cc1)
		validMoves = state.game.GetValidMovesForPiece(p)
	}
	if len(text) == 3 {
		px2 = int(text[2]-'a')
	}
	if len(text) > 3 {
		c := Coord(text[2:4])
		if !c.IsValid() {
			state.logger.Panicf("bug in move checking validation: %+v", c)
		}
		pos := c.AsCartesianCoord()
		px2 = pos.X
		py2 = 7 - pos.Y
	}

	var targetStyle tcell.Style
	for y := range 8 {
		for x := range 8 {
			square := state.squares[y][x]

			// target is default style based on square until we match a specific override style
			targetStyle = state.squareDefaultGreenStyle
			if y % 2 == x % 2 {
				targetStyle = state.squareDefaultBeigeStyle
			}

			if len(text) == 1 && px1 == x {
				// highlight the chosen column
				targetStyle = state.squareHighlightStyle
			} else if len(text) >= 2 && px1 == x && py1 == y {
				// highlight the chosen piece
				targetStyle = state.squareHighlightStyle
			}
			if len(text) == 4 && px2 == x && py2 == y {
				targetStyle = state.squareValidMoveStyle
			}

			square.Box.SetBorderStyle(targetStyle)
		}
	}

	if len(text) >= 2 {
		// We can now iterate over valid moves and override those specific squares
		for _, v := range validMoves {
			square := state.squares[7-v.dest.Y][v.dest.X]
			square.Box.SetBorderStyle(state.squareValidMoveStyle)
		}
	}
}

func ProcessMove(key tcell.Key, inputField *tview.InputField, state *State) {
	if key == tcell.KeyESC {
		inputField.SetText("")
		return
	}
	text := inputField.GetText()
	if key == tcell.KeyEnter {
		if len(text) != 4 {
			return
		}
		p := Coord(text[0:2])
		d := Coord(text[2:4])
		if !p.IsValid() || !d.IsValid() {
			state.logger.Panicf("unexpected coord in entered input %v", text)
		}

		pcc, dcc := p.AsCartesianCoord(), d.AsCartesianCoord()
		piece, hasPiece := state.game.GetCoord(pcc)
		if !hasPiece {
			state.logger.Panicf("unexpected entered starting coord with no piece %v", text)
		}
		
		moves, found := state.game.validMoves[piece]
		if !found {
			state.logger.Panicf("unexpected piece chosen with no valid moves %v", text)
		}
		var chosenMove ValidMove
		found = false
		for _, m := range moves {
			if m.piece == piece && m.dest == dcc {
				chosenMove = m
				found = true
				break
			}
		}

		if !found {
			state.logger.Panicf("unexpected entered dest coord with no matching move %v", text)
		}

		state.logger.Printf("Found matching move, executing state change %+v", chosenMove)
		state.game.ExecuteValidMove(chosenMove)
		inputField.SetText("")
		UpdateBoardUi(state)
	}

}

func Start(game *Game) {
	f, err := os.OpenFile("./log.txt", os.O_WRONLY | os.O_APPEND | os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	logger := log.New(f, "tviewui", log.LstdFlags)
	pieceSets := generatePieceSets()

	app := tview.NewApplication()
	app.EnableMouse(false)
	app.EnablePaste(false)

	outer := tview.NewFlex()
	outer.SetDirection(tview.FlexColumn)
	outer.SetBorder(true);

	board := tview.NewGrid()
	board.SetColumns(3, 0, 0, 0, 0, 0, 0, 0, 0, 3)
	board.SetRows(1, 0, 0, 0, 0, 0, 0, 0, 0, 1)
	board.SetBorder(false)
	board.SetBorders(false)
	board.SetGap(0, 0)

	// placeholder boxes in corners to keep look consistent
	board.AddItem(tview.NewBox(), 0, 0, 1, 1, 0, 0, false)
	board.AddItem(tview.NewBox(), 0, 9, 1, 1, 0, 0, false)
	board.AddItem(tview.NewBox(), 9, 0, 1, 1, 0, 0, false)
	board.AddItem(tview.NewBox(), 9, 9, 1, 1, 0, 0, false)

	squares := make([][]*tview.TextView, 8)
	currentPlayer := tview.NewTextView()
	currentPlayerStatus := tview.NewTextView()

	state := State{
		app: app,
		pieceSet: nil,
		game: game,
		squares: squares,
		currentPlayer: currentPlayer,
		currentPlayerStatus: currentPlayerStatus,
		greenBbgColor: tcell.NewHexColor(0x95B089),
		beigeBgColor: tcell.NewHexColor(0xB5A16E),
		whitePieceColor: tcell.NewHexColor(0xFFFFFF),
		blackPieceColor: tcell.NewHexColor(0x000000),
		squareDefaultGreenStyle: tcell.Style{}.Foreground(tview.Styles.PrimaryTextColor).Background(tcell.NewHexColor(0x95B089)),
		squareDefaultBeigeStyle: tcell.Style{}.Foreground(tview.Styles.PrimaryTextColor).Background(tcell.NewHexColor(0xB5A16E)),
		squareHighlightStyle: tcell.Style{}.Background(tcell.NewHexColor(0xFFFF00)).Foreground(tcell.NewHexColor(0xFFFF00)),
		squareValidMoveStyle: tcell.Style{}.Background(tcell.NewHexColor(0x008000)).Foreground(tcell.NewHexColor(0x008000)),
		logger: logger,
	}
	for y := range 8 {

		rowLabel := tview.NewTextView()
		rowLabel.SetTextAlign(tview.AlignCenter)
		rowLabel.SetText(string(byte('8'-y)))
		flx := tview.NewFlex()
		flx.SetDirection(tview.FlexRow)
		flx.SetBorder(false)
		flx.AddItem(tview.NewBox(), 0, 1, false)
		flx.AddItem(rowLabel, 3, 0, false)
		flx.AddItem(tview.NewBox(), 0, 1, false)
		board.AddItem(flx, y+1, 0, 1, 1, 0, 0, false)

		rowLabel = tview.NewTextView()
		rowLabel.SetTextAlign(tview.AlignCenter)
		rowLabel.SetText(string(byte('8'-y)))
		flx = tview.NewFlex()
		flx.SetDirection(tview.FlexRow)
		flx.SetBorder(false)
		flx.AddItem(tview.NewBox(), 0, 1, false)
		flx.AddItem(rowLabel, 3, 0, false)
		flx.AddItem(tview.NewBox(), 0, 1, false)
		board.AddItem(flx, y+1, 9, 1, 1, 0, 0, false)

		squares[y] = make([]*tview.TextView, 8)
		for x := range 8 {
			if y == 0 {
				colLabel := tview.NewTextView()
				colLabel.SetTextAlign(tview.AlignCenter)
				colLabel.SetText(string(byte('a'+x)))
				board.AddItem(colLabel, 0, x+1, 1, 1, 0, 0, false)

				colLabel = tview.NewTextView()
				colLabel.SetTextAlign(tview.AlignCenter)
				colLabel.SetText(string(byte('a'+x)))
				board.AddItem(colLabel, 9, x+1, 1, 1, 0, 0, false)
			}

			square := tview.NewTextView()
			board.AddItem(square, y+1, x+1, 1, 1, 0, 0, false)
			squares[y][x] = square

			square.SetBorder(true)
			bgColor := state.greenBbgColor
			if y % 2 == x % 2 {
				bgColor = state.beigeBgColor
			}
			square.SetBackgroundColor(bgColor)
		}
	}


	input := tview.NewInputField()
	input.SetBorder(true)
	input.SetFieldBackgroundColor(tcell.NewHexColor(0x000000))
	input.SetTitle("Move:")
	input.SetTitleAlign(tview.AlignLeft)
	input.SetAcceptanceFunc(func(textToCheck string, lastChar rune) bool {
		return MoveChecker(textToCheck, lastChar, &state)
	})
	input.SetChangedFunc(func(text string) {
		GridStateUpdater(text, &state)
	})
	input.SetDoneFunc(func(key tcell.Key) {
		ProcessMove(key, input, &state)
	})

	currentPlayerStatus.SetBorder(true)
	currentPlayerStatus.SetTitle("Player Status:")
	currentPlayerStatus.SetTitleAlign(tview.AlignLeft)

	currentPlayer.SetBorder(true)
	currentPlayer.SetTitle("Current Player:")
	currentPlayer.SetTitleAlign(tview.AlignLeft)

	history := tview.NewTextView()
	history.SetBorder(true)
	history.SetTitle("History:")
	history.SetTitleAlign(tview.AlignLeft)
	history.SetChangedFunc(func() {
		app.Draw()
	})

	status := tview.NewFlex()
	status.SetDirection(tview.FlexRow)
	status.AddItem(history, 0, 1, false)
	status.AddItem(currentPlayer, 3, 0, false)
	status.AddItem(currentPlayerStatus, 3, 0, false)
	status.AddItem(input, 3, 0, false)

	outer.AddItem(board, 0, 1, false)
	outer.AddItem(status, 40, 0, false)

	var width, height int
	app.SetAfterDrawFunc(func(screen tcell.Screen) {
		// In lieu of having access to the tcell.EventResize events directly, we infer it before draw.
		// If no change since the last draw, we can leave the size-dependent configuration as-is.
		// NOTE: sizes of widgets are not changed until after they have drawn to the screen. This means we need to hook
		// into the state of the ui after it finishes drawing, change the sizing, then trigger *another* Draw().
		nWidth, nHeight := screen.Size()
		if nWidth == width && nHeight == height {
			return
		}
		logger.Printf("detected screen size change from %vx%v to %vx%v", width, height, nWidth, nHeight)
		width, height = nWidth, nHeight

		go func() {
			app.QueueUpdateDraw(func() {
				logger.Printf("updating layout and pieceSet based on new size %vx%v", width, height)

				// Assume that if the grid has cells of inconsistent sizes due to the screen size not dividing evenly, the first
				// cell of the grid will be the smallest of the entire grid.
				firstSquare := squares[0][0]
				_, _, sw, sh := firstSquare.GetRect()
				// remove 2 to account for borders
				sw -= 2
				sh -= 2
				state.squareWidth = sw
				state.squareHeight = sh
				state.pieceSet = nil
				for _, pieceSet := range pieceSets {
					if state.squareWidth >= pieceSet.minX && state.squareHeight >= pieceSet.minY {
						state.pieceSet = pieceSet
						break;
					}
				}

				if (state.pieceSet == nil) {
					logger.Panicf("No piece set fit your screen size (%vx%v)", width, height)
				}

				outer.SetTitle(fmt.Sprintf(
					" Chess [tviewui[] screen:%vx%v square:%vx%v piece:%vx%v ",
					width, height, sw, sh, state.pieceSet.minX, state.pieceSet.minY,
				))

				UpdateBoardUi(&state)
				logger.Printf("completed layout update using square size %v,%v and piece size %v,%v", sw, sh, state.pieceSet.minX, state.pieceSet.minY)
			})
		}()
	})
	
	app.SetRoot(outer, true)
	app.SetFocus(input)
	if err := app.Run(); err != nil {
		logger.Panic(err)
	}
}
