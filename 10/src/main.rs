use petgraph::{
    graph::{DiGraph, NodeIndex},
    visit::{depth_first_search, DfsEvent::*},
};
use std::{collections::HashMap, fs::read_to_string, io::Result};

#[derive(Clone, Debug, PartialEq, Eq, Hash)]
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
    lines: Vec<Vec<char>>,
    rows: usize,
    cols: usize,
}

#[derive(PartialEq, Eq)]
enum Dir {
    Up,
    Left,
    Right,
    Down,
}

impl Dir {
    fn flip(&self) -> Dir {
        match self {
            Self::Up => Self::Down,
            Self::Left => Self::Right,
            Self::Right => Self::Left,
            Self::Down => Self::Up,
        }
    }
    fn list(c: char) -> Vec<Dir> {
        match c {
            '|' => vec![Self::Up, Self::Down],
            '-' => vec![Self::Left, Self::Right],
            'L' => vec![Self::Up, Self::Right],
            'J' => vec![Self::Up, Self::Left],
            '7' => vec![Self::Left, Self::Down],
            'F' => vec![Self::Right, Self::Down],
            '.' => vec![],
            // Starting position we will allow ALL neighbors before solving
            'S' => vec![Self::Up, Self::Left, Self::Right, Self::Down],
            _ => panic!("Current character does not match any pipe value"),
        }
    }
}

type DirFn = fn(&Map, &Pair) -> Option<Pair>;

impl Map {
    fn new(file: &str) -> Option<Self> {
        let lines = read_to_string(file)
            .ok()?
            .lines()
            .map(|line| line.chars().collect::<Vec<_>>())
            .collect::<Vec<_>>();
        let rows = lines.len();
        let cols = lines[0].len();
        Some(Self { lines, rows, cols })
    }
    fn at(&self, pair: &Pair) -> char {
        self.lines[pair.row][pair.col]
    }
    fn iter(&self) -> MapIter<'_> {
        MapIter::new(self)
    }
    fn start(&self) -> Pair {
        self.iter().find(|node| self.at(&node) == 'S').unwrap()
    }
    fn up(&self, pair: &Pair) -> Option<Pair> {
        (pair.row > 0).then(|| Pair::new(pair.row - 1, pair.col))
    }
    fn down(&self, pair: &Pair) -> Option<Pair> {
        (pair.row + 1 < self.rows).then(|| Pair::new(pair.row + 1, pair.col))
    }
    fn left(&self, pair: &Pair) -> Option<Pair> {
        (pair.col > 0).then(|| Pair::new(pair.row, pair.col - 1))
    }
    fn right(&self, pair: &Pair) -> Option<Pair> {
        (pair.col + 1 < self.cols).then(|| Pair::new(pair.row, pair.col + 1))
    }
    fn check_dir(&self, pair: &Pair, dir: Dir) -> Option<Pair> {
        let func = Map::func_for_dir(&dir);
        func(self, pair).and_then(|other| {
            let flipped = dir.flip();
            Dir::list(self.at(&other))
                .iter()
                .any(|dir| *dir == flipped)
                .then(|| other)
        })
    }
    fn check(&self, pair: &Pair, dirs: Vec<Dir>) -> Vec<Pair> {
        dirs.into_iter()
            .filter_map(|dir| self.check_dir(pair, dir))
            .collect()
    }
    fn func_for_dir(dir: &Dir) -> DirFn {
        match dir {
            Dir::Up => Map::up,
            Dir::Left => Map::left,
            Dir::Right => Map::right,
            Dir::Down => Map::down,
        }
    }
    fn adjs(&self, pair: &Pair) -> Vec<Pair> {
        let char = self.at(pair);
        let list = Dir::list(char);
        self.check(pair, list)
    }
}

struct MapIter<'a> {
    map: &'a Map,
    curr: Pair,
    nexted: bool,
}

impl<'a> MapIter<'a> {
    fn new(map: &'a Map) -> Self {
        MapIter {
            map,
            curr: Pair::new(0, 0),
            nexted: false,
        }
    }
}

impl<'a> Iterator for MapIter<'a> {
    type Item = Pair;
    fn next(&mut self) -> Option<Self::Item> {
        let mut row = self.curr.row;
        let mut col = self.curr.col + 1;
        if !self.nexted {
            self.nexted = true;
            col = self.curr.col;
        }
        if col >= self.map.cols {
            row += 1;
            if row >= self.map.rows {
                return None;
            } else {
                col = 0;
            }
        }
        self.curr = Pair::new(row, col);
        Some(self.curr.clone())
    }
}

fn main() -> Result<()> {
    let map = Map::new("input").unwrap();
    // Starting node (S)
    let start = map.start();

    let mut graph = DiGraph::<Pair, ()>::new();
    let mut pair_to_idx = HashMap::<Pair, NodeIndex>::new();
    let mut idx_to_pair = HashMap::<NodeIndex, Pair>::new();
    // Add all of the nodes in the map into a graph
    // Also put them into a HashMap so I can look it up later
    for node in map.iter() {
        let index = graph.add_node(node.clone());
        pair_to_idx.insert(node.clone(), index.clone());
        idx_to_pair.insert(index.clone(), node.clone());
    }
    // Add all the edges in the map
    for node in map.iter() {
        let index = pair_to_idx
            .get(&node)
            .expect("Couldn't lookup node that should exist already");
        let adjs = map.adjs(&node);
        for adj in adjs {
            let adj_index = pair_to_idx
                .get(&adj)
                .expect("Couldn't lookup node that should exist already");
            graph.update_edge(*index, *adj_index, ());
        }
    }
    let start_index = pair_to_idx.get(&start).unwrap();
    let mut parents = vec![NodeIndex::end(); graph.node_count()];
    let end: std::result::Result<(), ()> =
        depth_first_search(&graph, Some(*start_index), |event| match event {
            TreeEdge(u, v) => {
                let u_coord = idx_to_pair.get(&u).unwrap();
                let u_val = map.at(u_coord);
                let v_coord = idx_to_pair.get(&v).unwrap();
                let v_val = map.at(v_coord);
                println!("{:?} ({}) <- {:?} ({})", u_coord, u_val, v_coord, v_val);
                parents[v.index()] = u;
                Ok(())
            }
            BackEdge(u, v) if v == *start_index => {
                let u_coord = idx_to_pair.get(&u).unwrap();
                let u_val = map.at(u_coord);
                let v_coord = idx_to_pair.get(&v).unwrap();
                let v_val = map.at(v_coord);
                println!(
                    "Hit cycle: {:?} ({}) <- {:?} ({})",
                    u_coord, u_val, v_coord, v_val
                );
                Ok(())
            }
            _ => Ok(()),
        });
    // let val = map.at(idx_to_pair.get(&end).unwrap());
    //println!("Got end node value of {:?} ({})", end, val);
    // Find the cycles in the graph
    Ok(())
}
