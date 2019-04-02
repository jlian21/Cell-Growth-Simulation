package main

import (
	"bufio"
	"fmt"
	"image"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type Cell struct {
	celltype    int
	x, y        float64
	density     float64
	signalLevel int
	edges []*VEdge
}

type GameBoard struct {
	cells []Cell
	width float64
	zone  []Zone
	maze  []Rectangle
}

type Rectangle struct {
	x, y          float64
	width, height float64
}

type Zone struct {
	shape    string
	strength float64
	centrex  float64
	centrey  float64
	radius   float64
}

type TwoClusterBoard struct {
	cells       [][]Cell
	width       float64
	totalsignal int
}

// intersection of edges when creating Voronoi diagrams
type VPoint struct {
  	x float64
  	y float64
	exist bool
}

// site presenting a cell for creating Voronoi diagram
type VSite struct {
	x float64
	y float64
	cell *Cell
}

// edge for creating Voronoi diagram
type VEdge struct {
	start *VPoint
  	end *VPoint
  	direction *VPoint
  	left *Cell
  	right *Cell
  	neighbour *VEdge
  	f float64
  	g float64
}

// parabola for creating Voronoi diagram
type VParabola struct {
  	isLeaf bool
  	site *VSite
  	edge *VEdge
  	cEvent *VEvent
  	parent *VParabola
  	left *VParabola
  	right *VParabola
}

// put into queue, presenting site event or circle event
type VEvent struct {
  	site *VSite
	cell *Cell
  	arch *VParabola
	d float64
	se bool
}




func main() {
	rand.Seed(time.Now().UnixNano())

	if os.Args[1] == "OneCluster" {

		strategy := os.Args[2]

		inputs := ReadFromFile("OneClusterInputs.txt")

		if len(inputs) != 10 {
			panic("Error: wrong number of input arguments!")
			os.Exit(1)
		}

		// converting initialcells from inputs; int
		initialcells, err1 := strconv.Atoi(inputs[0])
		if err1 != nil {
			fmt.Println("Error: cannot convert initialcells!")
			os.Exit(1)
		}

		// converting numGens from inputs; int
		numGens, err2 := strconv.Atoi(inputs[1])
		if err2 != nil {
			fmt.Println("Error: cannot convert numGens!")
			os.Exit(1)
		}

		// converting searchRadius from inputs; float64
		searchRadius, err3 := strconv.ParseFloat(inputs[2], 64)
		if err3 != nil {
			fmt.Println("Error: cannot convert searchRadius!")
			os.Exit(1)
		}

		// converting birthRadius from inputs; float64
		birthRadius, err4 := strconv.ParseFloat(inputs[3], 64)
		if err4 != nil {
			fmt.Println("Error: cannot convert birthRadius!")
			os.Exit(1)
		}

		// converting deathRadius from inputs; float64
		deathRadius, err5 := strconv.ParseFloat(inputs[4], 64)
		if err5 != nil {
			fmt.Println("Error: cannot convert deathRadius!")
			os.Exit(1)
		}

		// converting birthrate from inputs; float64
		birthrate, err6 := strconv.ParseFloat(inputs[5], 64)
		if err6 != nil {
			fmt.Println("Error: cannot convert birthrate!")
			os.Exit(1)
		}

		// converting deathrate from inputs; float64
		deathrate, err7 := strconv.ParseFloat(inputs[6], 64)
		if err7 != nil {
			fmt.Println("Error: cannot convert deathrate!")
			os.Exit(1)
		}

		// converting width from inputs; float64
		width, err8 := strconv.ParseFloat(inputs[7], 64)
		if err8 != nil {
			fmt.Println("Error: cannot convert width!")
			os.Exit(1)
		}

		// converting numZones from inputs; int
		numZones, err9 := strconv.Atoi(inputs[8])
		if err9 != nil {
			fmt.Println("Error: cannot convert numZones!")
			os.Exit(1)
		}

		// converting addmaze from inputs; int
		addmaze, err10 := strconv.Atoi(inputs[9])
		if err10 != nil {
			fmt.Println("Error: cannot convert addmaze!")
			os.Exit(1)
		}

		//initial num of cells, initial birth radius
		initialboard := InitializeBoard(initialcells, birthRadius, width)
		initialboard = initialboard.AddZone(numZones)
		if addmaze == 1 {
			initialboard = initialboard.MakeMaze()
		}

		start := time.Now()
		//For each update, save the board
		boardList := UpdateBoard(initialboard, numGens, searchRadius, birthRadius, deathRadius, birthrate, deathrate, strategy)
		elapsed := time.Since(start)
		fmt.Println("Finish generating boardList", elapsed)

		//generate image for each board
		var imagelists []image.Image
		for i := range boardList {
			if i%3 == 0 {
				imagelists = append(imagelists, DrawBoard(boardList[i]).img)
			}
		}
		fmt.Println("Generating image...")

		//make gif file
		Process(imagelists, "OneCluster")

	} else if os.Args[1] == "AutoGenerate" {
		/*
			it will automatically generate an input file
		*/
		strategy := os.Args[2]
		AutoGenerator()

		inputs := ReadFromFile("input.txt")

		if len(inputs) != 10 {
			panic("Error: wrong number of input arguments!")
			os.Exit(1)
		}

		// converting initialcells from inputs; int
		initialcells, err1 := strconv.Atoi(inputs[0])
		if err1 != nil {
			fmt.Println("Error: cannot convert initialcells!")
			os.Exit(1)
		}

		// converting numGens from inputs; int
		numGens, err2 := strconv.Atoi(inputs[1])
		if err2 != nil {
			fmt.Println("Error: cannot convert numGens!")
			os.Exit(1)
		}

		// converting searchRadius from inputs; float64
		searchRadius, err3 := strconv.ParseFloat(inputs[2], 64)
		if err3 != nil {
			fmt.Println("Error: cannot convert searchRadius!")
			os.Exit(1)
		}

		// converting birthRadius from inputs; float64
		birthRadius, err4 := strconv.ParseFloat(inputs[3], 64)
		if err4 != nil {
			fmt.Println("Error: cannot convert birthRadius!")
			os.Exit(1)
		}

		// converting deathRadius from inputs; float64
		deathRadius, err5 := strconv.ParseFloat(inputs[4], 64)
		if err5 != nil {
			fmt.Println("Error: cannot convert deathRadius!")
			os.Exit(1)
		}

		// converting birthrate from inputs; float64
		birthrate, err6 := strconv.ParseFloat(inputs[5], 64)
		if err6 != nil {
			fmt.Println("Error: cannot convert birthrate!")
			os.Exit(1)
		}

		// converting deathrate from inputs; float64
		deathrate, err7 := strconv.ParseFloat(inputs[6], 64)
		if err7 != nil {
			fmt.Println("Error: cannot convert deathrate!")
			os.Exit(1)
		}

		// converting width from inputs; float64
		width, err8 := strconv.ParseFloat(inputs[7], 64)
		if err8 != nil {
			fmt.Println("Error: cannot convert width!")
			os.Exit(1)
		}

		// converting numZones from inputs; int
		numZones, err9 := strconv.Atoi(inputs[8])
		if err9 != nil {
			fmt.Println("Error: cannot convert numZones!")
			os.Exit(1)
		}

		// converting addmaze from inputs; int
		addmaze, err10 := strconv.Atoi(inputs[9])
		if err10 != nil {
			fmt.Println("Error: cannot convert addmaze!")
			os.Exit(1)
		}

		//initial num of cells, initial birth radius
		initialboard := InitializeBoard(initialcells, birthRadius, width)
		initialboard = initialboard.AddZone(numZones)
		if addmaze == 1 {
			initialboard = initialboard.MakeMaze()
		}
		fmt.Println(initialboard)
		start := time.Now()
		//For each update, save the board
		boardList := UpdateBoard(initialboard, numGens, searchRadius, birthRadius, deathRadius, birthrate, deathrate, strategy)
		elapsed := time.Since(start)
		fmt.Println("Finish generating boardList", elapsed)

		//generate image for each board
		var imagelists []image.Image
		for i := range boardList {
			if i%3 == 0 {
				imagelists = append(imagelists, DrawBoard(boardList[i]).img)
			}
		}
		fmt.Println("Generating image...")

		//make gif file
		Process(imagelists, "AutoGenerate")

	} else if os.Args[1] == "TwoCluster" {

		inputs := ReadFromFile("TwoClusterInputs.txt")

		// checking if there are correct number of inputs
		if len(inputs) != 8 {
			panic("Error: wrong number of input arguments!")
			os.Exit(1)
		}

		// converting initialcells from inputs; int
		initialcells, err1 := strconv.Atoi(inputs[0])
		if err1 != nil {
			fmt.Println("Error: cannot convert initialcells!")
			os.Exit(1)
		}
		// converting numGens from inputs; int
		numGens, err2 := strconv.Atoi(inputs[1])
		if err2 != nil {
			fmt.Println("Error: cannot convert numGens!")
			os.Exit(1)
		}

		// converting searchRadius from inputs; float64
		searchRadius, err3 := strconv.ParseFloat(inputs[2], 64)
		if err3 != nil {
			fmt.Println("Error: cannot convert searchRadius!")
			os.Exit(1)
		}

		// converting birthRadius from inputs; float64
		birthRadius, err4 := strconv.ParseFloat(inputs[3], 64)
		if err4 != nil {
			fmt.Println("Error: cannot convert birthRadius!")
			os.Exit(1)
		}

		// converting deathRadius from inputs; float64
		deathRadius, err5 := strconv.ParseFloat(inputs[4], 64)
		if err5 != nil {
			fmt.Println("Error: cannot convert deathRadius!")
			os.Exit(1)
		}

		// converting birthrate from inputs; float64
		birthrate, err6 := strconv.ParseFloat(inputs[5], 64)
		if err6 != nil {
			fmt.Println("Error: cannot convert birthrate!")
			os.Exit(1)
		}

		// converting deathrate from inputs; float64
		deathrate, err7 := strconv.ParseFloat(inputs[6], 64)
		if err7 != nil {
			fmt.Println("Error: cannot convert deathrate!")
			os.Exit(1)
		}

		// converting width from inputs; float64
		width, err8 := strconv.ParseFloat(inputs[7], 64)
		if err8 != nil {
			fmt.Println("Error: cannot convert width!")
			os.Exit(1)
		}

		start := time.Now()

		initialboard := InitializeTwoClusterBoard(initialcells, birthRadius, width)

		boardList := initialboard.UpdateBoard(numGens, searchRadius, birthRadius, deathRadius, birthrate, deathrate)

		elapsed := time.Since(start)

		fmt.Println("Finish generating boardList", elapsed)

		var imagelists []image.Image

		for i := range boardList {
			if i%2 == 0 {
				imagelists = append(imagelists, boardList[i].DrawBoard().img)
			}
		}

		fmt.Println("Generating image...")

		Process(imagelists, "TwoClusterSS")
	}
}

/*
	ReadFromFile reads the input txt file, and returns a slice of string.
*/

func ReadFromFile(filename string) []string {
	//open file
	input, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error: cannot read the input")
		os.Exit(1)
	}

	//scan each lines
	lines := make([]string, 0)
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if scanner.Err() != nil {
		fmt.Println("Error: cannot complete file reading ")
		os.Exit(1)
	}

	//close file
	input.Close()

	inputs := make([]string, 0)

	//divide each line by " " and read the parameters
	for i := 0; i < len(lines); i++ {
		var s []string
		s = strings.Split(lines[i], " ")
		inputs = append(inputs, s[1])
	}
	return inputs
}

//AutoGenerator randomly generate each parameters and output to a txt file
func AutoGenerator() {
	//randomly generate each parameters
	initialcells := rand.Intn(10) + 1
	numGens := -15*initialcells + 225
	searchRadius := rand.Float64()*15 + 10.0
	birthRadius := rand.Float64()*5 + 5.0
	deathRadius := rand.Float64() + 1.0
	birthrate := rand.Float64() * 0.5
	deathrate := rand.Float64() * 0.5
	width := 500.0
	numZones := rand.Intn(10) + 1
	addmaze := rand.Intn(2)

	//create a file to write in
	outfile, err := os.Create("input.txt")
	if err != nil {
		fmt.Println("error saving file")
		os.Exit(1)
	}

	//write each parameters
	fmt.Fprintln(outfile, "initialcells:", initialcells)
	fmt.Fprintln(outfile, "numGens:", numGens)
	fmt.Fprintln(outfile, "searchRadius:", searchRadius)
	fmt.Fprintln(outfile, "birthRadius:", birthRadius)
	fmt.Fprintln(outfile, "deathRadius:", deathRadius)
	fmt.Fprintln(outfile, "birthrate:", birthrate)
	fmt.Fprintln(outfile, "deathrate:", deathrate)
	fmt.Fprintln(outfile, "width:", width)
	fmt.Fprintln(outfile, "numZones:", numZones)
	fmt.Fprintln(outfile, "addmaze:", addmaze)
}

//AddZone takes in numZones and generate numZones zone randomly
func (board GameBoard) AddZone(numZones int) GameBoard {
	//create slice of circle zones
	board.zone = make([]Zone, numZones)
	for i := 0; i < len(board.zone); i++ {
		board.zone[i] = CircleZone(board.width)
	}

	return board
}

//CircleZone takes in board width and random generate a zone in board
func CircleZone(width float64) Zone {
	var circle Zone
	circle.shape = "circle"
	//random dicide inhibit or improve strength
	circle.strength = rand.Float64()*2 - 1
	circle.centrex = rand.Float64() * width
	circle.centrey = rand.Float64() * width
	circle.radius = rand.Float64() * 50

	return circle
}

//MakeMaze define rectangles in board and build up a maze
func (board GameBoard) MakeMaze() GameBoard {
	board.maze = make([]Rectangle, 6)

	board.maze[0].x = 0
	board.maze[0].y = 0
	board.maze[0].width = board.width / 6
	board.maze[0].height = board.width

	board.maze[1].x = 0
	board.maze[1].y = 5 * board.width / 6
	board.maze[1].width = board.width
	board.maze[1].height = board.width / 6

	board.maze[2].x = 5 * board.width / 6
	board.maze[2].y = 0
	board.maze[2].width = board.width / 6
	board.maze[2].height = board.width

	board.maze[3].x = 2 * board.width / 6
	board.maze[3].y = 0
	board.maze[3].width = 4 * board.width / 6
	board.maze[3].height = board.width / 6

	board.maze[4].x = 2 * board.width / 6
	board.maze[4].y = 0
	board.maze[4].width = board.width / 6
	board.maze[4].height = 4 * board.width / 6

	board.maze[5].x = 2 * board.width / 6
	board.maze[5].y = 3 * board.width / 6
	board.maze[5].width = 2 * board.width / 6
	board.maze[5].height = board.width / 6

	return board
}

// GenerateCell takes width as input, and returns a randomly-generated cell
// within a circle of radius width/2
func (board GameBoard) GenerateCell(centerX, centerY, radius float64) Cell {
	//make random angel and radius
	a := rand.Float64() * 2.0 * math.Pi
	r := (2.0*rand.Float64() - 1.0) * radius

	//make a new cell and difine its location
	var newCell Cell
	newCell.x = r*math.Cos(a) + centerX
	newCell.y = r*math.Sin(a) + centerY

	//check bondary
	if newCell.x < 0.0 {
		newCell.x = 0.0
	}
	if newCell.x > board.width {
		newCell.x = board.width
	}
	if newCell.y < 0.0 {
		newCell.y = 0.0
	}
	if newCell.y > board.width {
		newCell.y = board.width
	}
	return newCell
}

/*
	CalculateDistance takes c1 and c2 as input, and returns the distance between the two cells
	using Euclidean method.
*/

func CalculateDistance(c1, c2 Cell) float64 {
	c1x := c1.x
	c1y := c1.y
	c2x := c2.x
	c2y := c2.y
	return math.Sqrt((c1x-c2x)*(c1x-c2x) + (c1y-c2y)*(c1y-c2y))
}

// InitializeBoard takes numCells and birthRadius as inputs, and returns GameBoard
// with randomly-generated cells within the radius of birthRadius
func InitializeBoard(numCells int, birthRadius, width float64) GameBoard {
	//make an initial board
	var board GameBoard
	board.cells = make([]Cell, numCells)
	board.width = width

	//randomly generate numCells cells in board
	for i := 0; i <= numCells-1; i++ {
		board.cells[i] = board.GenerateCell(board.width/2, board.width/2, birthRadius)
	}
	return board
}

// UpdateBoard takes an initialBoard, numGens, searchRadius, birthRadius,
// deathRadius, birthrate and deathrate as inputs, and returns an updated GameBoard
func UpdateBoard(initialBoard GameBoard, numGens int, searchRadius, birthRadius, deathRadius, birthrate, deathrate float64, strategy string) []GameBoard {
	boards := make([]GameBoard, 0)
	boards = append(boards, initialBoard)
	for i := 1; i <= numGens; i++ {
		currentboard := boards[i-1].UpdateOneBoard(searchRadius, birthRadius, deathRadius, birthrate, deathrate, strategy)
		boards = append(boards, currentboard)
	}
	return boards
}

// UpdateOneBoard takes a currentBoard, searchRadius, birthRadius, deathRadius,
// birthrate, deathrate as inputs, and returns an updated board
func (currentBoard GameBoard) UpdateOneBoard(searchRadius, birthRadius, deathRadius, birthrate, deathrate float64, strategy string) GameBoard {
	//calculate cells density in currentboard
	if strategy == "CountDensity" {
	currentBoard.cells = CountDensity(currentBoard.cells, searchRadius)
	} else if strategy == "Voronoi" {
	currentBoard = Voronoi(&currentBoard)
	}	else {
		panic("wrong input")
	}

	//make a new board
	var newboard1 GameBoard
	newboard1 = CopyBoard(currentBoard)

	//sort cells in newboard
	newboard1.cells = Sorting(newboard1.cells)

	//decide in currentboard cells, before which key cell need to born a new cell, and after which key cell need to die
	birthkey := int(float64(len(currentBoard.cells)) * (1 - birthrate))
	deathkey := int(float64(len(currentBoard.cells)) * deathrate)

	//cell born and move and death
	newboard1.cells = newboard1.Born(birthRadius, searchRadius, birthkey)
	newboard1.cells = newboard1.Move(birthkey, deathkey, searchRadius)
	newboard1.cells = newboard1.Death(deathRadius)

	return newboard1
}

//CopyBoard make a new board and copy all the properties of the input board
func CopyBoard(board GameBoard) GameBoard {
	//make a new board
	var newboard GameBoard

	//copy cells information from currentboard to newboard1 in a new store space
	newboard.cells = make([]Cell, 0)
	for i := range board.cells {
		newboard.cells = append(newboard.cells, board.cells[i])
	}

	//copy zones
	newboard.zone = make([]Zone, 0)
	for i := range board.zone {
		newboard.zone = append(newboard.zone, board.zone[i])
	}

	//copy maze
	newboard.maze = make([]Rectangle, 0)
	for i := range board.maze {
		newboard.maze = append(newboard.maze, board.maze[i])
	}

	//copy properties
	newboard.width = board.width

	return newboard
}

//if distance between two cells is smaller than deathRadius, the cell with greater density died
func (board GameBoard) Death(deathRadius float64) []Cell {
	board.cells = Sorting(board.cells)
	for i := 0; i <= len(board.cells)-1; i++ {
		for j := i; j <= len(board.cells)-1; j++ {
			if i != j && CalculateDistance(board.cells[i], board.cells[j]) < deathRadius {
				board.cells = append(board.cells[:i], board.cells[i+1:]...)
			}
		}
	}
	return board.cells
}

//SurvivalRate takes in each cells' lifespan and calculate its correspoding survial rate
func SurvivalRate(generation int) float64 {
	survival := math.Pow(math.E, 1.1*(float64(generation)/3.2)-9.8)
	return survival
}

//CountDensity takes a list of Cell and searchRadius as input, and returns a slice of Cell with density.
func CountDensity(cells []Cell, searchRadius float64) []Cell {
	for i := range cells {
		cells[i].density = 0
		for j := range cells {
			if CalculateDistance(cells[i], cells[j]) < searchRadius {
				cells[i].density += 1
			}
		}
		cells[i].density -= 1 //exclude itself
	}
	return cells
}

//Sorting takes a slice of Cell and create a order from most dense to least dense.
func Sorting(cells []Cell) []Cell {
	//make min heap
	//i = parent index; a = child index
	for j := len(cells)/2 - 1; j >= 0; j-- { //range from largest parent node to root
		//heapify down
		i := j
		for 2*i+1 <= len(cells)-1 { //while cells[i] still have child, heapify down
			a := 2*i + 1               //left child
			if 2*i+2 <= len(cells)-1 { //if right child exists
				if cells[2*i+2].density < cells[2*i+1].density {
					a = 2*i + 2
				}
			}
			if cells[i].density > cells[a].density {
				cells[i], cells[a] = cells[a], cells[i] //if parent < child, heapify down
				i = a
			} else {
				break
			}
		}
	}

	//Sorting
	length := len(cells) - 1 //the last node's index
	for length > 0 {         //while there's still nodes in heaps
		cells[0], cells[length] = cells[length], cells[0] //swap root with the last node
		length--
		i := 0
		for 2*i+1 <= length { //while cells[i] still have child, heapify down
			a := 2*i + 1         //left child
			if 2*i+2 <= length { //if right child exists
				if cells[2*i+2].density < cells[2*i+1].density {
					a = 2*i + 2
				}
			}
			if cells[i].density > cells[a].density {
				cells[i], cells[a] = cells[a], cells[i] //if parent < child, heapify down
				i = a
			} else {
				break //if no need to swap, break loop
			}
		}
	}
	return cells
}

// Born takes a slice of Cell, birthRadius, and birthkey as inputs,
// and returns an updated slice of Cell
func (board GameBoard) Born(birthRadius, searchRadius float64, birthkey int) []Cell {
	length := len(board.cells)
	//for every cell in birthrate, born one cell
	for i := birthkey; i <= length-1; i++ {
		//randomly generate 10 cells and choose the one with least density
		a := make([]Cell, 10)
		min := len(board.cells)
		choice := 0
		for j := 0; j <= 9; j++ {
			a[j] = board.GenerateCell(board.cells[i].x, board.cells[i].y, birthRadius)
			for k := range board.cells {
				if CalculateDistance(a[j], board.cells[k]) < searchRadius {
					a[j].density += 1
				}
			}
			if a[j].density < float64(min) {
				choice = j
			}
		}

		//decide total survival rate for the new cell
		survival := 1.0

		//check maze edge
		for z := 0; z < len(board.maze); z++ {
			xdiff := a[choice].x - board.maze[z].x
			ydiff := a[choice].y - board.maze[z].y
			if xdiff > 0 && xdiff < board.maze[z].width && ydiff > 0 && ydiff < board.maze[z].height {
				survival = 0
				break
			}
		}

		//check zone
		if survival != 0 {
			for z := 0; z < len(board.zone); z++ {
				xdiff := a[choice].x - board.zone[z].centrex
				ydiff := a[choice].y - board.zone[z].centrey
				if math.Sqrt((xdiff)*(xdiff)+(ydiff)*(ydiff)) < board.zone[z].radius { //if cell generate in zone
					survival += board.zone[z].strength
				}
			}
		}

		//limit the change to [0, 2]
		if survival < 0 {
			survival = 0
		} else if survival > 2 {
			survival = 2
		}

		//decide if to keep the new cell or pormote cell birth based on survival rate
		if survival < 1 {
			if rand.Float64() <= survival { //keep the cell with survival possibility
				board.cells = append(board.cells, a[choice])
			}
		} else if survival >= 1 {
			board.cells = append(board.cells, a[choice])
			if rand.Float64() < survival-1 { //born another cell with survival possibility
				bonus := board.GenerateCell(board.cells[i].x, board.cells[i].y, birthRadius)
				board.cells = append(board.cells, bonus)
			}
		}
	}
	return board.cells
}

//Move takes repel and attract points index and calculate each cell's movemnet in forcefield
func (board GameBoard) Move(birthkey, deathkey int, searchRadius float64) []Cell {
	newcells := make([]Cell, 0)
	for i := 0; i <= len(board.cells)-1; i++ {
		newcells = append(newcells, board.cells[i])
	}

	for i := range board.cells {
		xmove := 0.0
		ymove := 0.0
		//cells from birthkey to the end are defined as attract points
		for j := birthkey; j <= (len(board.cells)+birthkey)/2; j++ {
			//one cell is only affected by attarct points in its search radius
			if i != j && CalculateDistance(board.cells[i], board.cells[j]) < searchRadius {
				distance := CalculateDistance(board.cells[i], board.cells[j])
				attractForce := 1.0 / (distance * distance)
				if attractForce > 1.0 {
					attractForce = 1.0
				}
				xratio := (board.cells[j].x - board.cells[i].x) / distance
				yratio := (board.cells[j].y - board.cells[i].y) / distance
				xmove += xratio * attractForce
				ymove += yratio * attractForce
			}
		}

		//cells before deathkey are difined as repel points
		for k := 0; k < deathkey; k++ {
			//one cell is only affected by repel points in its search radius
			if i != k && CalculateDistance(board.cells[i], board.cells[k]) < searchRadius {
				distance := CalculateDistance(board.cells[i], board.cells[k])
				repelForce := -1.0 / (distance * distance)
				if repelForce < -1.0 {
					repelForce = -1.0
				}
				xratio := (board.cells[k].x - board.cells[i].x) / distance
				yratio := (board.cells[k].y - board.cells[i].y) / distance
				xmove += xratio * repelForce
				ymove += yratio * repelForce
			}
		}

		//limit cell in board
		newcells[i].x += xmove
		newcells[i].y += ymove
		if newcells[i].x < 0.0 {
			newcells[i].x = 0.0
		}
		if newcells[i].x > board.width {
			newcells[i].x = board.width
		}
		if newcells[i].y < 0.0 {
			newcells[i].y = 0.0
		}
		if newcells[i].y > board.width {
			newcells[i].y = board.width
		}
	}

	return newcells
}

//DrawBoard takes in a board and draw a image
func DrawBoard(board GameBoard) Canvas {
	c := CreateNewCanvas(int(board.width), int(board.width))

	c.SetFillColor(MakeColor(0, 0, 0))
	c.ClearRect(0, 0, int(board.width), int(board.width))
	c.Fill()

	for i := range board.cells {
		c.SetFillColor(MakeColor(255, 255, 255))
		c.Circle(board.cells[i].x, board.cells[i].y, 1)
		c.Fill()
	}
	return c
}

/*
	GenerateCell takes centerX, centerY, radius, cellType as inputs, and returns a randomly generated
	cell with the designated cellType and within the radius around the point (centerX, centerY).
*/

func (board TwoClusterBoard) GenerateCell(centerX, centerY, radius float64, cellType string) Cell {
	a := rand.Float64() * 2.0 * math.Pi
	r := (2.0*rand.Float64() - 1.0) * radius

	var newCell Cell
	newCell.x = r*math.Cos(a) + centerX
	newCell.y = r*math.Sin(a) + centerY

	/*
		Set boundary for source and sink
		sources are allowed to grow on the left side of the board
		sinks are only allowed to grow on the right side of the board
	*/

	if cellType == "source" {
		newCell.celltype = 1
		if newCell.x < 0.0 {
			newCell.x = 0.0
		} else if newCell.x > board.width/2 {
			newCell.x = board.width / 2
		} else if newCell.y < 0.0 {
			newCell.y = 0.0
		} else if newCell.y > board.width {
			newCell.y = board.width
		}
	} else if cellType == "sink" {
		newCell.celltype = 2
		if newCell.x < board.width/2 {
			newCell.x = board.width / 2
		} else if newCell.x > board.width {
			newCell.x = board.width
		} else if newCell.y < 0.0 {
			newCell.y = 0.0
		} else if newCell.y > board.width {
			newCell.y = board.width
		}
	}
	return newCell
}

/*
InitializeBoard takes numCells, birthRadius, and the width of TwoClusterBoard as inputs,
and returns an initialized TwoClusterBoard with a total amount of numCells on the initial
TwoClusterBoard, and source cells are on the left center of the TwoClusterBoard within
the birthRadius, and sink cells are on the right centre of the TwoClusterBoard within the
birthRadius.
*/

func InitializeTwoClusterBoard(numCells int, birthRadius, width float64) TwoClusterBoard {
	/*
		Make an initial TwoClusterBoard
	*/
	var board TwoClusterBoard
	board.cells = make([][]Cell, 2)
	board.width = width

	/*
		Randomly generate numCells cells in board
	*/
	for i := 0; i <= numCells-1; i++ {
		cellType := rand.Intn(2)
		if cellType == 0 { // generate source
			source := board.GenerateCell(board.width/4, board.width/2, birthRadius, "source")
			board.cells[0] = append(board.cells[0], source)
		} else if cellType == 1 { // generate sink
			sink := board.GenerateCell(3*board.width/4, board.width/2, birthRadius, "sink")
			board.cells[1] = append(board.cells[1], sink)
		}
	}

	/*
		Update the totalsignal, which is calculated based on the number
		of source cells on the board
	*/

	board.totalsignal = len(board.cells[0])

	return board
}

/*
	UpdateBoard takes an initialBoard, numGens, searchRadius, birthRadius, deathRadius,
	birthrate and deathrate as inputs, and returns a collection of TwoClusterBoard during
	the update
*/

func (initialBoard TwoClusterBoard) UpdateBoard(numGens int, searchRadius, birthRadius, deathRadius, birthrate, deathrate float64) []TwoClusterBoard {
	boards := make([]TwoClusterBoard, 0)
	boards = append(boards, initialBoard)
	for i := 1; i <= numGens; i++ {
		currentboard := boards[i-1].UpdateOneBoard(searchRadius, birthRadius, deathRadius, birthrate, deathrate)
		boards = append(boards, currentboard)
	}
	return boards
}

/*
	UpdateOneBoard takes a currentBoard, searchRadius, birthRadius, deathRadius,
	birthrate, deathrate as inputs, and returns an updated board
*/

func (currentBoard TwoClusterBoard) UpdateOneBoard(searchRadius, birthRadius, deathRadius, birthrate, deathrate float64) TwoClusterBoard {

	/*
		Calculate cells density in currentboard
	*/

	for i := range currentBoard.cells {
		currentBoard.cells[i] = currentBoard.CountDensity(i, searchRadius)
	}

	/*
		Make a new board
	*/

	var newboard1 TwoClusterBoard
	newboard1 = currentBoard.CopyBoard()

	/*
		Sort cells in newboard by density
	*/

	for i := range newboard1.cells {
		newboard1.cells[i] = newboard1.SortingDensity(i)
	}

	/*
		In currentboard cells, and caluclate the birthkey and deathkey
		according to birthrate and deathrate. In the sorted []Cell (from most dense
		to least dense), cells after the birthkey need to perform born and move
		function, and cells before the deathkey cell need to perform death function.
	*/

	for i := range newboard1.cells {
		birthkey := int(float64(len(currentBoard.cells[i])) * (1 - birthrate))
		deathkey := int(float64(len(currentBoard.cells[i])) * deathrate)
		newboard1.cells[i] = newboard1.Born(birthRadius, searchRadius, birthkey, i)
		newboard1.cells[i] = newboard1.Move(birthkey, deathkey, i, searchRadius)
		if len(newboard1.cells[i]) < 5 {
			break
		} else if len(newboard1.cells[i]) >= 5 {
			newboard1.cells[i] = newboard1.Death(deathRadius, i)
		}
	}

	/*
		Sort the source cells by their X positions, which can be used to estimate if the boundary
		of two clusters are close enough for the transmission of the signal by calcultating the
		distance between the most left cell in the source cluster and the most right cell in the
		sink cluster.
	*/

	newboard1.cells[0] = newboard1.SortingXPosition(0)

	lastnode := len(newboard1.cells[0]) - 1

	newboard1.cells[1] = newboard1.SortingXPosition(1)

	if CalculateDistance(newboard1.cells[1][0], newboard1.cells[0][lastnode]) < 100.0 {

		/*
			The totalsignal availiable is calculated by adding the totalsignal with the cell born
			at each update. Then, 9/10 of them would be used in the signal transmission, and then
			1/10 would remained for the next update.
		*/

		intitalSourceNumber := len(currentBoard.cells[0])
		finalSourceNumber := len(newboard1.cells[0])
		newboard1.totalsignal = finalSourceNumber - intitalSourceNumber + newboard1.totalsignal
		availiableSignal := 9 * newboard1.totalsignal / 10
		newboard1.totalsignal = newboard1.totalsignal / 10

		/*
			The maximum amount of signal a sink can receive from the source cluster is 3. According
			to the sorted X position, the sink cell at the most left will be added with signal. If
			the sink cell reaches a signal level of 3, we proceed to the second most left sink and so
			on.
		*/

		for i := range newboard1.cells[1] {
			if newboard1.cells[1][i].signalLevel == 3 {
				continue
			} else if newboard1.cells[1][i].signalLevel < 3 {
				if newboard1.cells[1][i].signalLevel+availiableSignal >= 3 {
					signaladded := rand.Intn(3-newboard1.cells[1][i].signalLevel) + 1
					newboard1.cells[1][i].signalLevel = newboard1.cells[1][i].signalLevel + signaladded
					availiableSignal = availiableSignal - signaladded
				} else if availiableSignal+newboard1.cells[1][i].signalLevel < 3 {
					newboard1.cells[1][i].signalLevel += availiableSignal
					availiableSignal = 0
				}
			} else if availiableSignal == 0 {
				break
			}
		}

	}

	return newboard1
}

/*
	CopyBoard takes a TwoClusterBoard and returns a copied TwoClusterBoard
*/

func (board TwoClusterBoard) CopyBoard() TwoClusterBoard {

	/*
		Make a new board
	*/
	var newboard TwoClusterBoard

	/*
		Copy cells information from currentboard to newboard1 in a new store space
	*/

	newboard.cells = make([][]Cell, 2)
	for i := range board.cells {
		cells := make([]Cell, 0)
		for j := range board.cells[i] {
			cells = append(cells, board.cells[i][j])
		}
		newboard.cells[i] = append(newboard.cells[i], cells...)
	}

	/*
		Copy properties of board
	*/

	newboard.width = board.width
	newboard.totalsignal = board.totalsignal

	return newboard
}

/*
	Death takes a TwoClusterBoard, deathRadius, and k as inputs, and returns
	a slice of Cell after death representing board.cells[k]
*/

func (board TwoClusterBoard) Death(deathRadius float64, k int) []Cell {

	/*
		If the distance between two cells is smaller than deathRadius,
		and the cell with greater density dies
	*/

	board.cells[k] = board.SortingDensity(k)
	for i := 0; i <= len(board.cells[k])-1; i++ {
		for j := i; j <= len(board.cells[k])-1; j++ {
			if i != j && CalculateDistance(board.cells[k][i], board.cells[k][j]) < deathRadius {
				board.cells[k] = append(board.cells[k][:i], board.cells[k][i+1:]...)
			}
		}
	}
	return board.cells[k]
}

/*
	CountDensity takes a list of Cell, searchRadius, and k as input, and returns a slice
	of Cell, board.cells[k], with updated density in each Cell.
*/

func (board TwoClusterBoard) CountDensity(k int, searchRadius float64) []Cell {
	cells := board.cells[k]
	for i := range cells {
		cells[i].density = 0
		for j := range cells {
			if i != j && CalculateDistance(cells[i], cells[j]) < searchRadius {
				cells[i].density += 1
			}
		}
	}

	/*
		Compute the growth score using density, ccording to the formula:
		score = e^density/((1+e^density)^2), derivative of the logistic growth
		curve
	*/

	for i := range cells {
		d := float64(cells[i].density) - 500.0
		growthScore := math.Exp(d) / ((1 + math.Exp(d)) * (1 + math.Exp(d)))
		cells[i].density = growthScore
	}

	return cells
}

/*
	SortingDensity takes board and i as inputs, and returns a sorted []Cell from highest
	density to lowest density using the heap structure
*/

func (board TwoClusterBoard) SortingDensity(i int) []Cell {
	cells := board.cells[i]

	/*
		Make a max heap
		i = parent node
		a = child index
	*/

	for j := len(cells)/2 - 1; j >= 0; j-- { //range from largest parent node to root
		/*
			Heapify down
		*/
		i := j
		for 2*i+1 <= len(cells)-1 { // while cells[i] still have child, heapify down
			a := 2*i + 1               // left child
			if 2*i+2 <= len(cells)-1 { // if right child exists
				if cells[2*i+2].density < cells[2*i+1].density {
					a = 2*i + 2
				}
			}
			if cells[i].density > cells[a].density {
				cells[i], cells[a] = cells[a], cells[i] // if parent < child, heapify down
				i = a
			} else {
				break
			}
		}
	}

	/*
		Sorting
	*/

	length := len(cells) - 1 // the last node's index
	for length > 0 {         // while there's still nodes in heaps
		cells[0], cells[length] = cells[length], cells[0] // swap root with the last node
		length--
		i := 0
		for 2*i+1 <= length { // while cells[i] still have child, heapify down
			a := 2*i + 1         // left child
			if 2*i+2 <= length { // if right child exists
				if cells[2*i+2].density < cells[2*i+1].density {
					a = 2*i + 2
				}
			}
			if cells[i].density > cells[a].density {
				cells[i], cells[a] = cells[a], cells[i] // if parent < child, swap
				i = a
			} else {
				break // if no node to swap, break loop
			}
		}
	}
	return cells
}

/*
	Born takes a board, birthRadius, searchRadius, birthkey, and i as inputs,
	and returns an updated slice of Cell of board.cells[l] with the newborn
	Cell.
*/

func (board TwoClusterBoard) Born(birthRadius, searchRadius float64, birthkey int, l int) []Cell {
	length := len(board.cells[l])
	var cellType string
	if l == 0 {
		cellType = "source"
	} else if l == 1 {
		cellType = "sink"
	}

	/*
		Cells that are after the birthkey in the sorted density slice will have
		a new Cell born within the birthRadius of the current Cell.
	*/

	for i := birthkey; i <= length-1; i++ {
		/*
			Randomly generate 10 cells and choose the one with least density
		*/
		a := make([]Cell, 10)
		min := float64(len(board.cells[l]))
		choice := 0
		for j := 0; j <= 9; j++ {
			a[j] = board.GenerateCell(board.cells[l][i].x, board.cells[l][i].y, birthRadius, cellType)
			for k := range board.cells[l] {
				if CalculateDistance(a[j], board.cells[l][k]) < searchRadius {
					a[j].density++
				}
			}
			if a[j].density < min {
				choice = j
			}
		}
		board.cells[l] = append(board.cells[l], a[choice])
	}
	return board.cells[l]
}

/*
	Move takes a board, birthkey, deathkey, m, and searchRadius as inputs,
	and returns a slice of Cell of board.cells[m] after performing the movement
	among cells
*/

func (board TwoClusterBoard) Move(birthkey, deathkey, m int, searchRadius float64) []Cell {
	newcells := make([]Cell, 0)
	for i := 0; i <= len(board.cells[m])-1; i++ {
		newcells = append(newcells, board.cells[m][i])
	}
	for i := range board.cells[m] {
		xmove := 0.0
		ymove := 0.0

		/*
			For the cells that is not newborn but after the birthkey, if two cells are within the
			searchRadius, cells attract each other
		*/

		for j := birthkey; j <= (len(board.cells[m])+birthkey)/2; j++ {
			if i != j && CalculateDistance(board.cells[m][i], board.cells[m][j]) < searchRadius {
				distance := CalculateDistance(board.cells[m][i], board.cells[m][j])
				attractForce := 1.0 / (distance * distance)
				if attractForce > 1.0 {
					attractForce = 1.0
				}
				xratio := (board.cells[m][j].x - board.cells[m][i].x) / distance
				yratio := (board.cells[m][j].y - board.cells[m][i].y) / distance
				xmove += xratio * attractForce
				ymove += yratio * attractForce
			}
		}

		/*
			For the cells that is before the deathkey, if two cells are within the
			searchRadius, cells repel each other
		*/

		for k := 0; k < deathkey; k++ {
			if i != k && CalculateDistance(board.cells[m][i], board.cells[m][k]) < searchRadius {
				distance := CalculateDistance(board.cells[m][i], board.cells[m][k])
				repelForce := -1.0 / (distance * distance)
				if repelForce < -1.0 {
					repelForce = -1.0
				}
				xratio := (board.cells[m][k].x - board.cells[m][i].x) / distance
				yratio := (board.cells[m][k].y - board.cells[m][i].y) / distance
				xmove += xratio * repelForce
				ymove += yratio * repelForce
			}
		}

		/*
			Calculate the net force
			Also, set the boundary for cells after the move function
		*/

		newcells[i].x += xmove
		newcells[i].y += ymove
		if newcells[i].x < 0.0 {
			newcells[i].x = 0.0
		}
		if newcells[i].x > board.width {
			newcells[i].x = board.width
		}
		if newcells[i].y < 0.0 {
			newcells[i].y = 0.0
		}
		if newcells[i].y > board.width {
			newcells[i].y = board.width
		}
	}

	return newcells
}

/*
	SortingXPosition takes a board and i as inputs, and returns a sorted slice of Cells
	after ranking the cells in board.cells[i] according to their x value with the order
	from lowest x position to highest x position
*/

func (board TwoClusterBoard) SortingXPosition(i int) []Cell {
	cells := board.cells[i]
	/*
		Make a min heap
		i = parent node
		a = child index
	*/

	for j := len(cells)/2 - 1; j >= 0; j-- { // range from largest parent node to root
		// Heapify down
		i := j
		for 2*i+1 <= len(cells)-1 { // while cells[i] still have child, heapify down
			a := 2*i + 1               // left child
			if 2*i+2 <= len(cells)-1 { // if right child exists
				if cells[2*i+2].x > cells[2*i+1].x {
					a = 2*i + 2
				}
			}
			if cells[i].x < cells[a].x {
				cells[i], cells[a] = cells[a], cells[i] // if parent < child, swap
				i = a
			} else {
				break
			}
		}
	}

	/*
		Sorting
	*/

	length := len(cells) - 1 // the last node's index
	for length > 0 {         // while there's still nodes in heaps
		cells[0], cells[length] = cells[length], cells[0] // swap root with the last node
		length--
		i := 0
		for 2*i+1 <= length { // while cells[i] still have child, heapify down
			a := 2*i + 1         // left child
			if 2*i+2 <= length { // if right child exists
				if cells[2*i+2].x > cells[2*i+1].x {
					a = 2*i + 2
				}
			}
			if cells[i].x < cells[a].x {
				cells[i], cells[a] = cells[a], cells[i] // if parent < child, swap
				i = a
			} else {
				break //if no node to swap, break loop
			}
		}
	}
	return cells
}

/*
	DrawBoard takes a board as input, and returns a Canvas with the drawing board
*/

func (board TwoClusterBoard) DrawBoard() Canvas {
	c := CreateNewCanvas(int(board.width), int(board.width))

	c.SetFillColor(MakeColor(0, 0, 0))
	c.ClearRect(0, 0, int(board.width), int(board.width))
	c.Fill()

	/*
		The source cell shows a blue color, and the sink cell shows different levels of
		yellowness according to its signal level
	*/

	for i := range board.cells {
		for j := range board.cells[i] {
			if board.cells[i][j].celltype == 1 { // source
				c.SetFillColor(MakeColor(38, 226, 220))
			} else if board.cells[i][j].celltype == 2 && board.cells[i][j].signalLevel == 0 { // sink
				c.SetFillColor(MakeColor(105, 105, 0))
			} else if board.cells[i][j].celltype == 2 && board.cells[i][j].signalLevel == 1 {
				c.SetFillColor(MakeColor(155, 155, 0))
			} else if board.cells[i][j].celltype == 2 && board.cells[i][j].signalLevel == 2 {
				c.SetFillColor(MakeColor(205, 205, 0))
			} else if board.cells[i][j].celltype == 2 && board.cells[i][j].signalLevel == 3 {
				c.SetFillColor(MakeColor(255, 255, 0))
			}
			c.Circle(board.cells[i][j].x, board.cells[i][j].y, 1)
			c.Fill()
		}
	}
	return c
}

/*
Implement Fortune's algorithm to create Voronoi diagram. it takes a board as input and generate a boards
of cells with correspoding edges.
*/
func Voronoi(board *GameBoard) GameBoard {
	edges := GetEdges(board)
	board = EdgestoCell(board, edges)
	if len(board.cells) == 1 {
		(*board).cells[0].density = board.width * board.width
	} else {
		for i := range board.cells {
			(*board).cells[i].density = 1 / (GetArea(board.cells[i], board.width) + 0.1)
		}
	}
	return *board
}
/*
func Voronoi(board *GameBoard) GameBoard {
  	edges := GetEdges(board)
  	board = EdgestoCell(board, edges)
  	for i := range board.cells {
    	(*board).cells[i].density = GetArea(board.cells[i])
  	}
  	return *board
}
*/

/*
takes a slice of generated edges and add them to corresponding cells on the input GameBoard.
*/
func EdgestoCell(board *GameBoard, edges []*VEdge) *GameBoard {
	var newBoard GameBoard
	newBoard.width = board.width
	newBoard.zone = board.zone
	maps := make(map[*Cell]int)
	newCells := make([]Cell, 0)

	for i := range edges {
		if _, v := maps[edges[i].left]; v == false {
			maps[edges[i].left] = 1
			edges[i].left.edges = append(edges[i].left.edges, edges[i])
		} else if v == true {
			check := false
			for j := range edges[i].left.edges {
				if edges[i].left != nil && edges[i].start != nil && edges[i].end != nil {
					if edges[i].left.edges[j].start != nil && edges[i].left.edges[j].end != nil {
						if edges[i].left.edges[j].start.x == edges[i].start.x && edges[i].left.edges[j].start.y == edges[i].start.y && edges[i].left.edges[j].end.x == edges[i].end.x && edges[i].left.edges[j].end.y == edges[i].end.y {
							check = true
							break
						}
					}

				}

			}
			if !check {
				maps[edges[i].left] += 1
				edges[i].left.edges = append(edges[i].left.edges, edges[i])
			}
		}

		if _, v := maps[edges[i].right]; v == false {
			maps[edges[i].right] = 1
			edges[i].right.edges = append(edges[i].right.edges, edges[i])
		} else if v == true {
			check := false
			for j := range edges[i].right.edges {
				if edges[i].right != nil && edges[i].start != nil && edges[i].end != nil {
					if edges[i].right.edges[j].start != nil && edges[i].right.edges[j].end != nil {
						if edges[i].right.edges[j].start.x == edges[i].start.x && edges[i].right.edges[j].start.y == edges[i].start.y && edges[i].right.edges[j].end.x == edges[i].end.x && edges[i].right.edges[j].end.y == edges[i].end.y {
							check = true
							break
						}
					}

				}
			}
			if !check {
				maps[edges[i].right] += 1
				edges[i].right.edges = append(edges[i].right.edges, edges[i])
			}
		}
	}

		for k, _ := range maps {
			newCells = append(newCells, *k)
		}

		newBoard.cells = newCells

		return &newBoard
}


/*
main part of generating Voronoi diagram
*/
func GetEdges(board *GameBoard) []*VEdge {
  edges := make([]*VEdge, 0)
  width := board.width
  var tree *VParabola  // start a new tree of parabolas
  //var Y float64
  queue := make([]VEvent, 0)

  for i := range board.cells {
    e := InitiationSEvent(&board.cells[i])
    queue = Push(queue, e) //heapify up
  }


  for len(queue) != 0 {
  	var e VEvent
    queue, e = Pop(queue) // get and delete the tree and heapify down
    if e.se == true {
      queue, edges = InsertParabola(&tree, e, queue, edges, width) //
    } else {
      queue, edges = RemoveParabola(tree, e, queue, edges, width)
    }
  }

  for i := range edges {
	  if edges[i].start.x != 0 && edges[i].start.x != width && edges[i].start.y == width {
		  var boundary VEdge
		  boundary.start = InitiationVPoint(0, width)
		  boundary.end = InitiationVPoint(width, width)
		  boundary.left = edges[i].left
		  boundary.right = edges[i].right
		  edges = append(edges, &boundary)
	  }
  }

  edges = FinishEdge(tree, width, edges)

  for i := range edges {
    if edges[i].neighbour != nil {
      edges[i].start = edges[i].neighbour.end
    }
  }

  return edges
}


/*
add a parabola when the sweeping ling scan a cell and create a site event
*/
func InsertParabola(tree **VParabola, e VEvent, queue []VEvent, edges []*VEdge, width float64) ([]VEvent, []*VEdge) {
	Y := e.site.y
	site := e.site

	// current tree has no node
	if *tree == nil {
		*tree = InitiationParabola1(site)
		return queue, edges
	}

	// case where two site co-exist in one site event
	if (*tree).isLeaf == true && Y == (*tree).site.y {
		newParabolaLeft := InitiationParabola1((*tree).site)
		newParabolaRight := InitiationParabola1(site)
		(*tree).left = newParabolaLeft
		newParabolaLeft.parent = *tree
		(*tree).right = newParabolaRight
		newParabolaRight.parent = *tree

		startX := ((*tree).site.x + site.x)/2
		startY := width
		start := InitiationVPoint(startX, startY)

		(*tree).isLeaf = false
		if (*tree).site.x > site.x {
			(*tree).edge = InitiationVEdge(start, site, (*tree).site)
		} else {
			(*tree).edge = InitiationVEdge(start, (*tree).site, site)
		}
		edges = append(edges, (*tree).edge)
		return queue, edges
	}

	// most common case
	p := GetCorrespondingPar((*tree), site.x, Y)

	if Y == p.site.y {
		newParabolaLeft := InitiationParabola1(p.site)
		newParabolaRight := InitiationParabola1(site)
		p.left = newParabolaLeft
		newParabolaLeft.parent = p
		p.right = newParabolaRight
		newParabolaRight.parent = p

		startX := (p.site.x + site.x)/2
		startY := width
		start := InitiationVPoint(startX, startY)

		p.isLeaf = false
		if p.site.x > site.x {
			p.edge = InitiationVEdge(start, site, p.site)
		} else {
			p.edge = InitiationVEdge(start, p.site, site)
		}
		edges = append(edges, p.edge)
		return queue, edges
	}
	p.isLeaf = false

	/*
	if p.cEvent != nil { // is CEvent
	queue = Delete(queue, *p.cEvent)
	p.cEvent = nil
	}
	*/

	// create "two" edges
	startX := site.x
	startY := GetYByX(p.site, site, Y)
	start := InitiationVPoint(startX, startY)
	el := InitiationVEdge(start, p.site, site)
	er := InitiationVEdge(start, site, p.site)
	el.neighbour = er

	// modify tree
	p0 := InitiationParabola1(p.site)
	p1 := InitiationParabola1(site)
	p2 := InitiationParabola1(p.site)

	p.left = InitiationParabola2()
	p.left.parent = p

	p.left.left = p0
	p0.parent = p.left
	p.left.right = p1
	p1.parent = p.left
	p.right = p2
	p2.parent = p

	p.edge = er
	p.left.edge = el
	edges = append(edges, el)

	queue = CheckCircle(p0, e.site.y, queue)
	queue = CheckCircle(p2, e.site.y, queue)

	return queue, edges
}

/*
remove an arc on the tree and retrive the fixed edge
*/
func RemoveParabola(tree *VParabola, e VEvent, queue []VEvent, edges []*VEdge, width float64) ([]VEvent, []*VEdge) {
	// current event e is a circle event
	//Y := e.site.y

	// p1 is the disappearing parabola corresponding to the current circle event
	p1 := e.arch
	edgeLeft := GetLeftParent(p1)
	edgeRight := GetRightParent(p1)
	p0 := GetLeftChild(edgeLeft)
	p2 := GetRightChild(edgeRight)

/*
	if p0.cEvent != nil { // is CEvent
	queue = Delete(queue, *p0.cEvent)
	p0.cEvent = nil
	}
	if p2.cEvent != nil { // is CEvent
	queue = Delete(queue, *p2.cEvent)
	p2.cEvent = nil
	}
	*/

	// edge left and right + para left and right
	end := InitiationVPoint(e.site.x, e.site.y + e.d) //CheckCircle
	edgeLeft.edge.end = end
	edgeRight.edge.end = end

	// new edge
	var higher, par *VParabola
	par = p1
	for par != tree {
		par = par.parent
		if par == edgeLeft {
			higher = edgeLeft
		} else if par == edgeRight {
			higher = edgeRight
		}
	}
	higher.edge = InitiationVEdge(end, p0.site, p2.site)
	edges = append(edges, higher.edge)

	// modify tree
	var gparent *VParabola
	gparent = p1.parent.parent
	if p1.parent.left == p1 {
		if gparent.left == p1.parent {
			gparent.left = p1.parent.right
			p1.parent.right.parent = gparent
		}
		if gparent.right == p1.parent {
			gparent.right = p1.parent.right
			p1.parent.right.parent = gparent
		}
	} else if p1.parent.right == p1 {
		if gparent.left == p1.parent {
			gparent.left = p1.parent.left
			p1.parent.left.parent = gparent
		}
		if gparent.right == p1.parent {
			gparent.right = p1.parent.left
			p1.parent.left.parent = gparent
		}
	}

	queue = CheckCircle(p0, e.site.y, queue)
	queue = CheckCircle(p2, e.site.y, queue)

	return queue, edges
}

/*
decide whether the input parabola will be "circled".
*/
func CheckCircle(p *VParabola, Y float64, queue []VEvent) []VEvent {
	edgeLeft := GetLeftParent(p)
	edgeRight := GetRightParent(p)
	parLeft := GetLeftChild(edgeLeft)
	parRight := GetRightChild(edgeRight)


	if parLeft == nil || parRight == nil || parLeft.site == parRight.site {
		return queue
	} else if !GetEdgeIntersection(edgeLeft.edge, edgeRight.edge).exist {
		return queue
	}

	Intersection := GetEdgeIntersection(edgeLeft.edge, edgeRight.edge)
	dx := parLeft.site.x - Intersection.x
	dy := parLeft.site.y - Intersection.y
	d := math.Sqrt((dx*dx)+(dy*dy))

	if Intersection.y - d >= Y {
		return queue
	}

	// already confirmed circle event
	//point := InitiationVPoint(Intersection.x, Intersection.y-d)
	//e := InitiationCEvent(point)
	var e VEvent 	// this is a circle event
	var s VSite
	s.x = Intersection.x
	s.y = Intersection.y - d
	e.site = &s
	e.arch = p
	e.se = false
	e.d = d
	p.cEvent = &e
	queue = Push(queue, e)

	return queue
}

func GetCorrespondingPar(tree *VParabola, siteX, Y float64) *VParabola {
	var edgeX float64
	parabola := tree

	for parabola.isLeaf == false {
		edgeX = GetXOfEdge(parabola, Y)
		if edgeX > siteX {
			parabola = parabola.left
		} else {
			parabola = parabola.right
		}
	}

	return parabola
}

/*
fixed the last edges when the sweeping line reaches the bottom of the plane
*/
func FinishEdge(tree *VParabola, width float64, edges []*VEdge) []*VEdge {
	if tree.isLeaf == true {
		return edges
	}

	f := tree.edge.f
	g := tree.edge.g
	endX := 0.0
	endY := 0.0

	if tree.edge.direction.x > 0 {
		endX = width
		endY = f*endX + g
		if endY > width {
			endY = width
			endX = (endY - g) / f
		}
	} else if tree.edge.direction.x < 0 {
		endX = 0
		endY = f*endX + g
		if endY < 0 {
			endY = 0
			endX = (endY - g) / f
		}
	} else {
		//direction.x = 0
		endX = tree.edge.start.x
		endY = 0
	}

	if endX == 0 {
		var boundary VEdge
		boundary.start = InitiationVPoint(0, 0)
		boundary.end = InitiationVPoint(0, width)
		boundary.left = tree.edge.left
		boundary.right = tree.edge.right
		edges = append(edges, &boundary)
		edges = append(edges, &boundary)
	}
	if endX == width {
		var boundary VEdge
		boundary.start = InitiationVPoint(width, width)
		boundary.end = InitiationVPoint(width, 0)
		boundary.left = tree.edge.left
		boundary.right = tree.edge.right
		edges = append(edges, &boundary)
		edges = append(edges, &boundary)
	}
	if endY == 0 {
		var boundary VEdge
		boundary.start = InitiationVPoint(width, 0)
		boundary.end = InitiationVPoint(0, 0)
		boundary.left = tree.edge.left
		boundary.right = tree.edge.right
		edges = append(edges, &boundary)
		edges = append(edges, &boundary)
	}
	if endY == width {
		var boundary VEdge
		boundary.start = InitiationVPoint(width, width)
		boundary.end = InitiationVPoint(0, width)
		boundary.left = tree.edge.left
		boundary.right = tree.edge.right
		edges = append(edges, &boundary)
		edges = append(edges, &boundary)
	}

	end := InitiationVPoint(endX, endY)
	tree.edge.end = end

	edges = FinishEdge(tree.left, width, edges)
	edges = FinishEdge(tree.right, width, edges)

	return edges
}


func GetXOfEdge(parabola *VParabola, Y float64) float64 {
	var edgeX float64

	parLeft := GetLeftChild(parabola)
	parRight := GetRightChild(parabola)
	siteLeft := parLeft.site
	siteRight := parRight.site

	a1 := 1 /(2 * (siteLeft.y - Y))
	b1 := -2 * siteLeft.x / (2 * (siteLeft.y - Y))
	c1 := Y + (2 * (siteLeft.y - Y)) / 4 + siteLeft.x  * siteLeft.x / (2 * (siteLeft.y - Y))

	a2 := 1 / (2 * (siteRight.y - Y))
	b2 := -2 * siteRight.x / (2 * (siteRight.y - Y))
	c2 := Y + (2 * (siteRight.y - Y)) / 4 + siteRight.x  * siteRight.x / (2 * (siteRight.y - Y))

	a := a1 - a2
	b := b1 - b2
	c := c1 - c2
	x1 := (-b - math.Sqrt(b*b - 4*a*c)) / (2*a)
	x2 := (-b + math.Sqrt(b*b - 4*a*c)) / (2*a)

	if siteLeft.y < siteRight.y {
		edgeX = max(x1, x2)
	} else {
		edgeX = min(x1, x2)  //???
	}

	return edgeX
}

func GetYByX(site1, site2 *VSite, Y float64) float64 {
	x := site2.x
	y := 0.0


	a := 1 / (2 * (site1.y - Y))
	b := -2 * site1.x / (2 * (site1.y - Y))
	c := Y + (2 * (site1.y - Y)) / 4 + site1.x * site1.x / (2 * (site1.y - Y))

	y = a*x*x + b*x + c

	return y
}

func GetEdgeIntersection(e1, e2 *VEdge) *VPoint {
	if e1.direction.x == 0 {
		if e2.direction.x == 0 {
			var p VPoint
			p.exist = false
			return &p
		}
		p := InitiationVPoint(e1.start.x, e2.f * e1.start.x + e2.g)
		return p
	} else if e2.direction.x == 0 {
		p := InitiationVPoint(e2.start.x, e1.f * e2.start.x + e1.g)
		return p
	}

	f1 := e1.f
	g1 := e1.g

	f2 := e2.f
	g2 := e2.g

	x := - (g1 - g2) / (f1 - f2)
	y := f1 * x + g1

	p := InitiationVPoint(x, y)
	p.exist = false

	if (x - e1.start.x) / e1.direction.x < 0 {
		return p
	} else if (y - e1.start.y) / e1.direction.y < 0 {
		return p
	} else if (x - e2.start.x) / e2.direction.x < 0 {
		return p
	} else if (y - e2.start.x) / e2.direction.y < 0 {
		return p
	}

	p.exist = true
	return p
}


func max(a, b float64) float64 {
	if a >= b {
		return a
	}

	return b
}

func min(a, b float64) float64 {
	if a <= b {
		return a
	}

	return b
}


// take an "edge"
func GetLeftChild(parabola *VParabola) *VParabola {
	if parabola == nil {
		return nil
	}

	rp := parabola.left
	for rp.isLeaf == false {
		rp = rp.right
	}

	return rp
}

// take a parabola and reach its closet left “edge” on the tree
func GetLeftParent(parabola *VParabola) *VParabola {
	// empty parabola
	if parabola == nil {
		return nil
	}

	rp := parabola.parent
	tail := parabola

	// when parabola is the left child of its parent
	for rp.left == tail {
		// if there is no such closet left parent exists
		if rp.parent == nil {
			return nil
		}

		// else search for higher node to see if it meets the requirement that rp.right = tail
		tail = rp
		rp = rp.parent
	}

	return rp
}

/*
take an "edge"
*/
func GetRightChild(parabola *VParabola) *VParabola {
	// empty parabola
	if parabola == nil {
		return nil
	}

	rp := parabola.right
	for rp.isLeaf == false {
		rp = rp.left
	}

	return rp
}

func GetRightParent(parabola *VParabola) *VParabola {
	// empty parabola
	if parabola == nil {
		return nil
	}

	rp := parabola.parent
	tail := parabola

	// when parabola is the right child of its parent
	for rp.right == tail {
		// if there is no such closet right parent exists
		if rp.parent == nil {
			return nil
		}

		// else search for higher node to see if it meets the requirement that rp.right = tail
		tail = rp
		rp = rp.parent
	}

	return rp
}

/*
calculate the area defined by generated edges by using Heron's formula
*/
func GetArea(c Cell, width float64) float64 {

  area := 0.0
  Vertexes := make([]*VPoint, 0)
  for i := range c.edges {
    Vertexes = append(Vertexes, c.edges[i].start)
    Vertexes = append(Vertexes, c.edges[i].end)
  }
  check := true
  for check {
    check = false
    for i := range Vertexes {
      index, num := CheckRepeat(Vertexes, i)
      if num >= 1 && !CheckCanvasVertexes(Vertexes[i], width){
        Vertexes = append(Vertexes[:index], Vertexes[index + 1:]...)
        check = true
        break
      } else {
        if CheckCanvasVertexes(Vertexes[i], width) && num == 0{
          Vertexes = append(Vertexes[:i], Vertexes[i + 1:]...)
          check = true
          break
        }
      }
    }
  }

  check = true
  for check {
    check = false
    for i := range Vertexes {
      index, num := CheckRepeat(Vertexes, i)
      if num >= 1 {
        Vertexes = append(Vertexes[:index], Vertexes[index + 1:]...)
        check = true
        break
      }
    }
  }

  for i := 1; i<len(Vertexes)-1; i++ {
    a := EdgeLength(Vertexes[0], Vertexes[i])
    b := EdgeLength(Vertexes[0], Vertexes[i+1])
    c := EdgeLength(Vertexes[i], Vertexes[i+1])
    area += Heron(a, b, c)
  }
  return area
}

func CheckRepeat(Vertexes []*VPoint, a int) (int, int) {
  num := 0
  var b int
  for i := range Vertexes {
		if Vertexes[i] != nil && Vertexes[a] != nil {
			if Vertexes[i].x == Vertexes[a].x && Vertexes[i].y == Vertexes[a].y && i != a{
	      num += 1
	      b = i
	    }
		}

  }
  return b, num
}

func CheckCanvasVertexes(p *VPoint, width float64) bool {
  flag := false
	if p == nil {
		return flag
	}
  if (p.x == 0.0 && p.y == 0.0) || (p.x == 0.0 && p.y == width) || (p.x == width && p.y == 0.0) || (p.x == width && p.y == width) {
    flag = true
  }
  return flag
}

func EdgeLength(v1, v2 *VPoint) float64 {
	if v1 == nil || v2 == nil {
		return 0.0
	}
  length := math.Sqrt((v1.x-v2.x)*(v1.x-v2.x) + (v1.y-v2.y)*(v1.y-v2.y))
  return length
}

func Heron(a, b, c float64) float64 {
  s := (a + b + c)/2
  area := math.Sqrt(s*(s-a)*(s-b)*(s-c))
  return area
}

func InitiationVSite(cell *Cell) *VSite {
	var site VSite

	site.x = cell.x
	site.y = cell.y
	site.cell = cell

	return &site
}

func InitiationParabola1(site *VSite) *VParabola {
	var parabola VParabola

	parabola.site = InitiationVSite(site.cell)
	parabola.isLeaf = true

	return &parabola
}

func InitiationParabola2() *VParabola {
	var parabola VParabola

	parabola.isLeaf = false

	return &parabola
}

func InitiationVEdge(start *VPoint, siteLeft, siteRight *VSite) *VEdge {
	var e VEdge

	e.start = InitiationVPoint(start.x, start.y)
	dx := -(siteRight.x - siteLeft.x)
	dy := siteRight.y - siteLeft.y
	e.direction = InitiationVPoint(dy, dx)
	e.left = siteLeft.cell
	e.right = siteRight.cell

  	e.f = (siteRight.x - siteLeft.x) / (siteLeft.y - siteRight.y)
	e.g = start.y - e.f * start.x

	return &e
}

func InitiationVPoint(x, y float64) *VPoint {
	var p VPoint

	p.x = x
	p.y = y
	p.exist = true

	return &p
}

func InitiationSEvent(cell *Cell) VEvent {
	var event VEvent

	event.cell = cell
	event.site = InitiationVSite(cell)
	event.se = true

	return event
}

/*
put an event into the queue 
*/
func Push(queue []VEvent, e VEvent) []VEvent {
  queue = append(queue, e)
  parent := len(queue)

  a := len(queue) - 1
  //fmt.Println(queue[a-1].site.y)
  for a != 0 {
    // get parent node
    if a % 2 != 0 {
      parent = (a - 1) / 2
    } else {
      parent = (a - 2) / 2
    }

    // fmt.Println(queue[parent].site.y)
    if queue[parent].site.y < queue[a].site.y {
      // swap
      queue[parent], queue[a] = queue[a], queue[parent]
      a = parent
    } else {
      break
    }
  }

  // sorting
  //fmt.Println(len(queue))

  length := len(queue) - 1
  c := 0
  for length > 0 {
    queue[0], queue[length] = queue[length], queue[0]
    length --

    i := 0
    for 2*i + 1 <= length {
      c = 2*i + 1
      if 2*i + 2 <= length {
        if queue[2*i + 2].site.y > queue[2*i + 1].site.y {
          c = 2*i + 2
        }
      }
      if queue[i].site.y < queue[c].site.y {
        queue[i], queue[c] = queue[c], queue[i]
        i = c
      } else {
        break
      }
    }
  }

  newQueue := make([]VEvent, 0)
  for i := len(queue) - 1; i >= 0; i-- {
    newQueue = append(newQueue, queue[i])
  }

  return newQueue
}

/*
  get the max value and delete it from the queue
*/
func Pop(queue []VEvent) ([]VEvent ,VEvent) {
  // get the element with max value and return it
  e := queue[0]

  newQueue:= Delete(queue, queue[0])

  return newQueue, e
}

/*
  search
  swap with the last element
  heapify down after swapping
*/
func Delete(queue []VEvent, e VEvent) []VEvent {
  var a int
  a = SearchE(queue, e) //an index
  parent := 0

  queue[a], queue[len(queue)-1] = queue[len(queue)-1], queue[a]
  queue = append(queue[:len(queue)-1])


  // moving up
  for a != 0 {
    // compare with the parent
    if a % 2 != 0 {
      parent = (a - 1) / 2
      } else {
      parent = (a - 2) / 2
      }

    if queue[parent].site.y < queue[a].site.y {
    // swap a with its current parent
      queue[parent], queue[a] = queue[a], queue[parent]
      a = parent
    } else {
    // up to the suitable position
      break
    }
  }


  // then moving down
  c := 0
  //fmt.Println("a is ", a)
  if a <= len(queue) - 1 { // avoid the case where index a is in range
    for 2*a + 1 <= len(queue) - 1 { // check the existence of child of current node
      c = 2*a + 1
      if 2*a + 2 <= len(queue) - 1 {
        if queue[2*a + 2].site.y > queue[2*a + 1].site.y {
          c = 2*a + 2
        }
      }

      if queue[c].site.y > queue[a].site.y {
        queue[a], queue[c] = queue[c], queue[a]
        a = c
      } else {
        // no need to move down
        break
      }
    }
  }

  // sorting
  length := len(queue) - 1
  c = 0
  for length > 0 {
    queue[0], queue[length] = queue[length], queue[0]
    length --

    i := 0
    for 2*i + 1 <= length {
      c = 2*i + 1
      if 2*i + 2 <= length {
        if queue[2*i + 2].site.y > queue[2*i + 1].site.y {
          c = 2*i + 2
        }
      }
      if queue[i].site.y < queue[c].site.y {
        queue[i], queue[c] = queue[c], queue[i]
        i = c
      } else {
        break
      }
    }
  }

  newQueue := make([]VEvent, 0)
  for i := len(queue) - 1; i >= 0; i-- {
    newQueue = append(newQueue, queue[i])
  }

  return newQueue
}


func SearchE(queue []VEvent, e VEvent) int {
  a := 0
  for i:=0; i<=len(queue)-1; i++ {
    if queue[i] == e {
      a = i
      break
    }
  }

  return a
}
