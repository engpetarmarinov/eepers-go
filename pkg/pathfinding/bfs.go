package pathfinding

import "github.com/engpetarmarinov/eepers-go/pkg/world"

// Point represents a point in the pathfinding grid.
type Point struct {
	X, Y int
}

// BFS finds the shortest path from start to end using Breadth-First Search.
func BFS(grid [][]world.Cell, start, end Point) []Point {
	queue := [][]Point{{start}}
	visited := make(map[Point]bool)
	visited[start] = true

	for len(queue) > 0 {
		path := queue[0]
		queue = queue[1:]
		node := path[len(path)-1]

		if node.X == end.X && node.Y == end.Y {
			return path
		}

		for _, dir := range []Point{{0, 1}, {0, -1}, {1, 0}, {-1, 0}} {
			nextNode := Point{X: node.X + dir.X, Y: node.Y + dir.Y}
			if isValid(nextNode, grid) && !visited[nextNode] {
				newPath := make([]Point, len(path))
				copy(newPath, path)
				newPath = append(newPath, nextNode)
				queue = append(queue, newPath)
				visited[nextNode] = true
			}
		}
	}

	return nil
}

func isValid(p Point, grid [][]world.Cell) bool {
	return p.Y >= 0 && p.Y < len(grid) && p.X >= 0 && p.X < len(grid[0]) && grid[p.Y][p.X] == world.CellFloor
}

// ComputeDistanceMap creates a distance map from the target position.
// The map contains -1 for unreachable cells, 0 for the target, and positive integers for distance.
// stepsLimit limits how far the pathfinding will search.
// stepLengthLimit allows moving multiple cells in one direction per step.
// canStandFunc checks if a position is valid for the entity.
func ComputeDistanceMap(
	grid [][]world.Cell,
	target Point,
	targetSize Point,
	stepsLimit int,
	stepLengthLimit int,
	canStandFunc func(Point) bool,
) [][]int {
	height := len(grid)
	width := len(grid[0])

	// Initialize distance map with -1 (unreachable)
	distMap := make([][]int, height)
	for i := range distMap {
		distMap[i] = make([]int, width)
		for j := range distMap[i] {
			distMap[i][j] = -1
		}
	}

	// Queue for BFS
	queue := []Point{}

	// Mark all positions where the entity could overlap with the target as distance 0
	for dy := 0; dy < targetSize.Y; dy++ {
		for dx := 0; dx < targetSize.X; dx++ {
			pos := Point{X: target.X - dx, Y: target.Y - dy}
			if canStandFunc(pos) {
				distMap[pos.Y][pos.X] = 0
				queue = append(queue, pos)
			}
		}
	}

	// BFS with step length
	directions := []Point{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}

	for len(queue) > 0 {
		pos := queue[0]
		queue = queue[1:]

		currentDist := distMap[pos.Y][pos.X]
		if currentDist >= stepsLimit {
			continue
		}

		// Try each direction
		for _, dir := range directions {
			newPos := Point{X: pos.X + dir.X, Y: pos.Y + dir.Y}

			// Try moving up to stepLengthLimit cells in this direction
			for step := 1; step <= stepLengthLimit; step++ {
				if !canStandFunc(newPos) {
					break
				}

				// If this cell hasn't been visited yet
				if distMap[newPos.Y][newPos.X] == -1 {
					distMap[newPos.Y][newPos.X] = currentDist + 1
					queue = append(queue, newPos)
				} else {
					// Already visited, stop extending in this direction
					break
				}

				// Move to next position in the same direction
				newPos = Point{X: newPos.X + dir.X, Y: newPos.Y + dir.Y}
			}
		}
	}

	return distMap
}
