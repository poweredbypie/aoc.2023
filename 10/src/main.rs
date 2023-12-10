use petgraph::graph::{NodeIndex, UnGraph};
use std::{collections::HashMap, fs::read_to_string, io::Result};

#[derive(Clone, PartialEq, Eq, Hash)]
struct Pair {
    pub row: usize,
    pub col: usize,
}

impl Pair {
    fn new(row: usize, col: usize) -> Pair {
        Pair { row, col }
    }
}

struct Map {
    curr: Pair,
    lines: Vec<Vec<char>>,
    rows: usize,
    cols: usize,
}

impl Map {
    fn new(file: &str) -> Option<Map> {
        let lines = read_to_string(file)
            .ok()?
            .lines()
            .map(|line| line.chars().collect::<Vec<_>>())
            .collect::<Vec<_>>();
        let rows = lines.len();
        let cols = lines[0].len();
        Some(Map {
            curr: Pair::new(0, 0),
            lines,
            rows,
            cols,
        })
    }
    fn at(&self, pair: &Pair) -> char {
        self.lines[pair.row][pair.col]
    }
    fn up(&self, pair: &Pair) -> Option<Pair> {
        if pair.row <= 0 {
            None
        } else {
            Some(Pair::new(pair.row - 1, pair.col))
        }
    }
    fn down(&self, pair: &Pair) -> Option<Pair> {
        if pair.row + 1 >= self.rows {
            None
        } else {
            Some(Pair::new(pair.row + 1, pair.col))
        }
    }
    fn left(&self, pair: &Pair) -> Option<Pair> {
        if pair.col <= 0 {
            None
        } else {
            Some(Pair::new(pair.row, pair.col - 1))
        }
    }
    fn right(&self, pair: &Pair) -> Option<Pair> {
        if pair.col + 1 >= self.cols {
            None
        } else {
            Some(Pair::new(pair.row, pair.col + 1))
        }
    }
    fn adjs(&self, pair: &Pair) -> Vec<Pair> {
        let char = self.at(pair);
        match char {
            '|' => vec![self.up(pair), self.down(pair)],
            '-' => vec![self.left(pair), self.right(pair)],
            'L' => vec![self.up(pair), self.right(pair)],
            'J' => vec![self.up(pair), self.left(pair)],
            '7' => vec![self.left(pair), self.down(pair)],
            'F' => vec![self.right(pair), self.down(pair)],
            '.' => vec![],
            // Starting position we will allow ALL neighbors before solving
            'S' => vec![
                self.up(pair),
                self.left(pair),
                self.right(pair),
                self.down(pair),
            ],
            _ => panic!("Current character does not match any pipe value"),
        }
        .into_iter()
        .flatten()
        .collect()
    }
}

impl Iterator for &Map {
    type Item = Pair;
    fn next(&mut self) -> Option<Self::Item> {
        let mut row = self.curr.row;
        let mut col = self.curr.col + 1;
        if col >= self.cols {
            row += 1;
            if row >= self.rows {
                return None;
            } else {
                col = 0;
            }
        }
        return Some(Pair::new(row, col));
    }
}

fn main() -> Result<()> {
    let map = Map::new("input").unwrap();
    let mut graph = UnGraph::<Pair, ()>::new_undirected();
    let mut pair_to_idx = HashMap::<Pair, NodeIndex<u32>>::new();
    // Add all of the spots in the map into a graph
    // Also put them into a HashMap so I can look it up later
    for node in &map {
        pair_to_idx.insert(node.clone(), graph.add_node(node));
    }
    for node in &map {
        let index = pair_to_idx
            .get(&node)
            .expect("Couldn't lookup node that should exist already");
        let adjs = map.adjs(&node);
        for adj in adjs {
            let adj_index = pair_to_idx
                .get(&adj)
                .expect("Couldn't lookup node that should exist already");
            graph.add_edge(*index, *adj_index, ());
        }
    }
    Ok(())
}
