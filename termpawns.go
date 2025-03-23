package main

const (
	threebythreePawn =
`   
 p 
   `
	threebythreeRook =
`   
 R 
   `
	threebythreeKnight =
`   
 N 
   `
	threebythreeBishop =
`   
 B 
   `
	threebythreeQueen =
`   
 Q 
   `
	threebythreeKing =
`   
 K 
   `


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
	threebythreePieces := PieceSet{
		minX: 3,
		minY: 3,
		pieces: []string{
			threebythreePawn,
			threebythreeRook,
			threebythreeKnight,
			threebythreeBishop,
			threebythreeQueen,
			threebythreeKing,
		},
	}

	return []*PieceSet{
		&elevenbytenPieces,
		&sevenbyninePieces,
		&sixbysixPieces,
		&threebythreePieces,
	}
}
