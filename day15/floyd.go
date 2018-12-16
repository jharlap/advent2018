package main

type PathCalculator struct {
	dist [][]int
	next [][]int
}

func PathCalculatorFrom(g *Game) *PathCalculator {
	p := &PathCalculator{
		dist: matrix(g.Map.width*g.Map.height, PathMaxDist),
		next: matrix(g.Map.width*g.Map.height, -1),
	}

	/*
	   procedure FloydWarshallWithPathReconstruction ()
	      for each edge (u,v)
	         dist[u][v] ← w(u,v)  // the weight of the edge (u,v)
	         next[u][v] ← v
	      for k from 1 to |V| // standard Floyd-Warshall implementation
	         for i from 1 to |V|
	            for j from 1 to |V|
	               if dist[i][j] > dist[i][k] + dist[k][j] then
	                  dist[i][j] ← dist[i][k] + dist[k][j]
	                  next[i][j] ← next[i][k]
	*/
	for _, e := range g.Edges() {
		p.dist[e.U][e.V] = 1
		p.next[e.U][e.V] = e.V
	}

	n := len(p.dist)
	for k := 0; k < n; k++ {
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				if p.dist[i][j] > p.dist[i][k]+p.dist[k][j] {
					p.dist[i][j] = p.dist[i][k] + p.dist[k][j]
					p.next[i][j] = p.next[i][k]
				}
			}
		}
	}
	return p
}

func (p *PathCalculator) Dist(u, v int) int {
	return p.dist[u][v]
}

func (p *PathCalculator) ShortestPath(u, v int) []int {
	/*
	   procedure Path(u, v)
	      if next[u][v] = null then
	          return []
	      path = [u]
	      while u ≠ v
	          u ← next[u][v]
	          path.append(u)
	      return path
	*/
	if p.next[u][v] == -1 {
		return nil
	}

	r := []int{u}
	for u != v {
		u = p.next[u][v]
		r = append(r, u)
	}
	return r
}

func matrix(n, v int) [][]int {
	m := make([][]int, n)
	for i := 0; i < n; i++ {
		m[i] = make([]int, n)
		for j := 0; j < n; j++ {
			m[i][j] = v
		}
	}
	return m
}
