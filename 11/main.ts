// Run with `deno run --allow-read main.ts`

class Coord {
  row: number;
  col: number;
  constructor(row: number, col: number) {
    this.row = row;
    this.col = col;
  }
  // For console.log
  toString() {
    return `(${this.row}, ${this.col})`;
  }
}

class GalaxyMap {
  lines: string[];
  emptyRows: number[] = [];
  emptyCols: number[] = [];
  weight = 1;
  constructor(str: string) {
    this.lines = str.split("\n").filter((line) => line != "");
  }
  // Find all empty rows and columns, and "expand" the galaxy
  // This really just finds all the rows and columns so it can be used with GalaxyMap.manhattan
  expand() {
    // Find all empty rows
    for (let row = 0; row < this.lines.length; row += 1) {
      if (this.lines[row].includes("#")) {
        continue;
      }
      this.emptyRows.push(row);
    }
    const colIncludes = (col: number, val: string) => {
      for (let row = 0; row < this.lines.length; row += 1) {
        if (this.lines[row][col] == val) {
          return true;
        }
      }
      return false;
    };
    // Find all empty columns
    for (let col = 0; col < this.lines[0].length; col += 1) {
      if (colIncludes(col, "#")) {
        continue;
      }
      this.emptyCols.push(col);
    }
    console.log(`Empty rows: ${this.emptyRows}`);
    console.log(`Empty cols: ${this.emptyCols}`);
  }
  // Set the weight of an empty column's expansion. Affects the manhattan distance calculator
  setWeight(weight: number) {
    this.weight = weight;
  }
  // Find all galaxies (#) in the map
  galaxies(): Coord[] {
    const list = [];
    for (let row = 0; row < this.lines.length; row += 1) {
      for (let col = 0; col < this.lines[row].length; col += 1) {
        if (this.lines[row][col] == "#") {
          list.push(new Coord(row, col));
        }
      }
    }
    return list;
  }
  // Find the weighted manhattan distance between two points in the galaxy.
  // Takes into account any empty rows and their expansion weights
  manhattan(one: Coord, two: Coord): number {
    const between = (val: number, low: number, hi: number): boolean => {
      if (low > hi) {
        return between(val, hi, low);
      }
      return val >= low && val <= hi;
    };
    const weight = this.weight - 1;
    const rowAdd =
      this.emptyRows.filter((row) => between(row, one.row, two.row)).length;
    const colAdd =
      this.emptyCols.filter((col) => between(col, one.col, two.col)).length;
    const rows = Math.abs(one.row - two.row);
    const cols = Math.abs(one.col - two.col);
    return rows + cols + (rowAdd + colAdd) * weight;
  }
  // For console.log
  toString(): string {
    return this.lines.join("\n");
  }
}

// Thanks to https://stackoverflow.com/a/20765091
// (list.length) C 2, returns all pair combos in the input list
// Does not contain any duplicates (e.g. if [one, two] is in the list, [two, one] isn't)
function combos<T>(list: T[]): [T, T][] {
  const out: [T, T][] = [];
  for (let i = 0; i < list.length; i += 1) {
    for (let j = i + 1; j < list.length; j += 1) {
      out.push([list[i], list[j]]);
    }
  }
  return out;
}

async function main() {
  const utf8 = new TextDecoder("utf-8");
  const contents = utf8.decode(await Deno.readFile("input"));
  const map = new GalaxyMap(contents);
  map.expand();
  const all = combos(map.galaxies());

  // Part A
  map.setWeight(2);
  const sumTwo = all.reduce(
    (val, pair) => val + map.manhattan(pair[0], pair[1]),
    0,
  );
  console.log(
    `Sum of all manhattan distances with empty expansion weight 2 is ${sumTwo}`,
  );

  // Part B
  map.setWeight(1_000_000);
  const sumMill = all.reduce(
    (val, pair) => val + BigInt(map.manhattan(pair[0], pair[1])),
    BigInt(0),
  );
  console.log(
    `Sum of all manhattan distances with empty expansion weight 1,000,000 is ${sumMill}`,
  );
}

main();
