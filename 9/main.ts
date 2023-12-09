// Run with `deno run --allow-red main.ts`

// Kind of annoying I can't nest this in the class
type Recur = (nums: number[]) => number;
type Action = (next: number | null, nums: number[]) => number;

class History {
  nums: number[];
  constructor(line: string) {
    this.nums = line.split(" ").map((str) => parseInt(str));
  }
  private recur(action: Action): Recur {
    // This needs to be a variable so I can recurse
    const func = (nums: number[]): number => {
      let allZero = true;
      const diffs: number[] = [];
      for (let i = 0; i < nums.length - 1; i += 1) {
        const diff = nums[i] - nums[i + 1];
        allZero = allZero && diff == 0;
        diffs.push(nums[i + 1] - nums[i]);
      }
      return action(allZero ? null : func(diffs), nums);
    };
    return func;
  }
  // Part A: predict the number after the last one
  predict(): number {
    return this.recur((next, nums) => (next ?? 0) + nums.at(-1)!)(this.nums);
  }
  // Part B: "retrodict" the number before the first one
  retrodict(): number {
    return this.recur((next, nums) => nums.at(0)! - (next ?? 0))(this.nums);
  }
}

async function main() {
  const utf8 = new TextDecoder("utf-8");
  const contents = utf8.decode(await Deno.readFile("input"));
  const histories = contents
    .split("\n")
    .filter((line) => line != "")
    .map((line) => new History(line));
  const predicted = histories.reduce(
    (sum, history) => sum + history.predict(),
    0,
  );
  const retrodicted = histories.reduce(
    (sum, history) => sum + history.retrodict(),
    0,
  );
  console.log(`Sum of all predicted values is ${predicted}`);
  console.log(`Sum of all predicted values is ${retrodicted}`);
}

main();
