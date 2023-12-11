class Coord {
  row: number;
  col: number;
  constructor(row: number, col: number) {
    this.row = row;
    this.col = col;
  }
  manhattan(other: Coord) {
    return Math.abs(other.row - this.row) + Math.abs(other.col - this.col);
  }
}

class GalaxyMap {
  lines: string[];
  constructor(str: string) {
    this.lines = str.split("\n").filter((line) => line != "");
    console.log(this.lines);
  }
  expand() {
    // Expand all empty rows
    for (let row = 0; row < this.lines.length; row += 1) {
      console.log(`Row ${row} / ${this.lines.length}`);
      if (this.lines[row].includes("#")) {
        continue;
      }
      // Duplicate the line if no galaxies exist in it
      this.lines.splice(row, 0, this.lines[row]);
      // Double up the row so we can skip over the newly created one
      row += 1;
    }
    // Expand all empty columns
    const colIncludes = (col: number, val: string) => {
      for (let row = 0; row < this.lines.length; row += 1) {
        if (this.lines[row][col] == val) {
          return true;
        }
      }
      return false;
    };
    const colInsert = (col: number, val: string) => {
      for (let row = 0; row < this.lines.length; row += 1) {
        this.lines[row] = this.lines[row].slice(0, col) + val +
          this.lines[row].slice(col);
      }
    };
    for (let col = 0; col < this.lines[0].length; col += 1) {
      console.log(`Col ${col} / ${this.lines[0].length}`);
      if (colIncludes(col, "#")) {
        continue;
      }
      colInsert(col, ".");
      // Double up the column so we can skip over the newly created one
      col += 1;
    }
  }
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
  toString(): string {
    return this.lines.join("\n");
  }
}

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
  console.log(all, all.length);
  let sum = 0;
  for (const pair of all) {
    sum += pair[0].manhattan(pair[1]);
  }
  console.log(`Sum of all manhattan distances is ${sum}`);
}

main();
