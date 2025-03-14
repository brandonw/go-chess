package main

import (
	"fmt"
	"github.com/rivo/tview"
	"github.com/gdamore/tcell/v2"
	"os"
	"log"
)

const (
	sixbysixPawn   =
`      
      

  ()  
  )(  
 /__\ `
	sixbysixRook   =
`      
      
  II  
  )(  
  )(  
 /__\ `
	sixbysixKnight =
`      
      
 _,,  
"-=\\~
  )(  
 /__\ `
	sixbysixBishop =
`      
      
  ()  
  )(  
  )(  
 /__\ `
	sixbysixQueen  =
`  .   
  ww  
  ()  
  )(  
  )(  
 /__\ `
	sixbysixKing   =
`  +   
  ww  
  ()  
  )(  
  )(  
 /__\ `

 	sevenbyninePawn =
`       
       
       
       
   _   
  ( )  
  ) (  
 (   ) 
[_____]`
 	sevenbynineRook =
`       
       
       
  |||  
  | |  
  ) (  
  ) (  
 (   ) 
[_____]`
 	sevenbynineKnight =
`        
        
        
   _/)  
  /. |= 
 /_) |= 
(/ ) (= 
  (   ) 
 [_____]`
 	sevenbynineBishop =
`       
       
       
   ^   
  / \  
  ) (  
  ) (  
 (   ) 
[_____]`
 	sevenbynineQueen =
`   o   
  ^^^  
  ( )  
  ) (  
 (   ) 
  ) (  
  ) (  
 (   ) 
[_____]`
 	sevenbynineKing =
`   +   
  ^^^  
  ( )  
  ) (  
 (   ) 
  ) (  
  ) (  
 (   ) 
[_____]`

	elevenbytenPawn = 
`           
           
           
           
    __     
   (  )    
    ><     
   |  |    
  /    \   
 |______|  `
	elevenbytenRook =
`           
           
 _   _   _ 
| |_| |_| |
 \       / 
  |     |  
  |     |  
  |     |  
 /       \ 
|_________|`
	elevenbytenKnight =
`           
           
  |\__     
 /   o\__  
|    ___=' 
|    \     
 \    \    
  >    \   
 /      \  
|________| `
	elevenbytenBishop =
`           
     o     
   /\^/\   
  |  /  )  
  | /  /   
   Y  /    
   |  |    
   |  |    
  /    \   
 |______|  `
	elevenbytenQueen =
`           
|\ ,''. /| 
| '''''' | 
 \      /  
  |    |   
  |    |   
  |    |   
  |    |   
 /      \  
|________| `
	elevenbytenKing =
`  =||=     
|\ ,''. /| 
| '''''' | 
 \      /  
  |    |   
  |    |   
  |    |   
  |    |   
 /      \  
|________| `
)

type PieceSet struct {
	minX int
	minY int
	pieces []string
}

func generatePieceSets() []*PieceSet {
	elevenbytenPieces := PieceSet{
		minX: 11,
		minY: 10,
		pieces: []string{
			elevenbytenPawn,
			elevenbytenRook,
			elevenbytenKnight,
			elevenbytenBishop,
			elevenbytenQueen,
			elevenbytenKing,
		},
	}
	sevenbyninePieces := PieceSet{
		minX: 7,
		minY: 9,
		pieces: []string{
			sevenbyninePawn,
			sevenbynineRook,
			sevenbynineKnight,
			sevenbynineBishop,
			sevenbynineQueen,
			sevenbynineKing,
		},
	}
	sixbysixPieces := PieceSet{
		minX: 6,
		minY: 6,
		pieces: []string{
			sixbysixPawn,
			sixbysixRook,
			sixbysixKnight,
			sixbysixBishop,
			sixbysixQueen,
			sixbysixKing,
		},
	}

	return []*PieceSet{
		&elevenbytenPieces,
		&sevenbyninePieces,
		&sixbysixPieces,
	}
}

func Start(game Game) {
	f, err := os.OpenFile("./log.txt", os.O_WRONLY | os.O_APPEND | os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	logger := log.New(f, "tviewui", log.LstdFlags)
	pieceSets := generatePieceSets()

	greenBbgColor := tcell.NewHexColor(0xB2D2A4)
	beigeBgColor := tcell.NewHexColor(0xC4AE78)
	whitePieceColor := tcell.NewHexColor(0xFFFFFF)
	blackPieceColor := tcell.NewHexColor(0x000000)

	app := tview.NewApplication()
	app.EnableMouse(false)
	app.EnablePaste(false)

	outer := tview.NewFlex()
	outer.SetBorder(true);

	board := tview.NewGrid()
	board.SetSize(8, 8, 0, 0)
	board.SetBorder(false)
	board.SetBorders(false)
	board.SetGap(0, 0)

	squares := make([][]*tview.TextView, 8)
	for y := range 8 {
		squares[y] = make([]*tview.TextView, 8)
		for x := range 8 {
			square := tview.NewTextView()
			board.AddItem(square, y, x, 1, 1, 0, 0, false)
			squares[y][x] = square

			square.SetBorder(false)
			square.SetBorderPadding(0, 0, 0, 0)
			bgColor := greenBbgColor
			if y % 2 == x % 2 {
				bgColor = beigeBgColor
			}
			square.SetBackgroundColor(bgColor)
		}
	}

	status := tview.NewGrid()
	status.SetColumns(0)
	status.SetRows(0, 0, 0, 0)
	status.SetBorders(true)

	outer.AddItem(board, 0, 1, false)
	outer.AddItem(status, 0, 1, false)

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
		logger.Printf("AfterDrawFunc: detected screen size change from %vx%v to %vx%v", width, height, nWidth, nHeight)
		width, height = nWidth, nHeight

		// Assume that if the grid has cells of inconsistent sizes due to the screen size not dividing evenly, the first
		// cell of the grid will be the smallest of the entire grid.
		firstSquare := squares[0][0]
		_, _, sw, sh := firstSquare.GetRect()
		var chosenPieceSet *PieceSet
		for _, pieceSet := range pieceSets {
			if sw >= pieceSet.minX && sh >= pieceSet.minY {
				chosenPieceSet = pieceSet
				break;
			}
		}

		if (chosenPieceSet == nil) {
			logger.Panic("No piece set fit your screen size")
		}

		// Now that we have updated sizes, queue a redraw
		logger.Print("Updating layout based on screen size change")
		outer.SetTitle(fmt.Sprintf(
			" Chess [tviewui[] screen:%vx%v square:%vx%v piece:%vx%v ",
			width, height, sw, sh, chosenPieceSet.minX, chosenPieceSet.minY,
		))

		for y := range 8 {
			for x := range 8 {
				p, o := game.GetCoord(CartesianCoord{x, y})
				if o {
					color := blackPieceColor
					if p.color == White {
						color = whitePieceColor
					}
					// tview y=0 corresponds to the top row of the grid, not the bottom. We must transform between tview
					// coord to game coord where y=0 is the bottom row.
					square := squares[7-y][x]
					square.SetTextColor(color)
					square.SetText(chosenPieceSet.pieces[p.pieceType])
					topPadding := (sh - chosenPieceSet.minY) / 2
					leftPadding := (sw - chosenPieceSet.minX) / 2
					square.SetBorderPadding(topPadding, 0, leftPadding, 0)
				}
			}
		}

		go func() {
			app.QueueUpdateDraw(func() {
				logger.Print("AfterDrawFunc: initiating Draw() with changed ui configuration")
			})
		}()
	})
	
	app.SetRoot(outer, true)
	if err := app.Run(); err != nil {
		logger.Panic(err)
	}
}
