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
