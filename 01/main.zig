const std = @import("std");
const fs = std.fs;
const mem = std.mem;
const out = std.io.getStdOut().writer();
const Alloc = std.heap.GeneralPurposeAllocator(.{});

fn partOne(line: []const u8) ?u64 {
    const values = "123456789";
    const first = mem.indexOfAny(u8, line, values) orelse return null;
    const last = mem.lastIndexOfAny(u8, line, values) orelse return null;
    return (line[first] - '0') * 10 + (line[last] - '0');
}

fn sliceToNum(slice: []const u8) ?u64 {
    if (slice[0] >= '1' and slice[0] <= '9') {
        return slice[0] - '0';
    } else if (mem.startsWith(u8, slice, "one")) {
        return 1;
    } else if (mem.startsWith(u8, slice, "two")) {
        return 2;
    } else if (mem.startsWith(u8, slice, "three")) {
        return 3;
    } else if (mem.startsWith(u8, slice, "four")) {
        return 4;
    } else if (mem.startsWith(u8, slice, "five")) {
        return 5;
    } else if (mem.startsWith(u8, slice, "six")) {
        return 6;
    } else if (mem.startsWith(u8, slice, "seven")) {
        return 7;
    } else if (mem.startsWith(u8, slice, "eight")) {
        return 8;
    } else if (mem.startsWith(u8, slice, "nine")) {
        return 9;
    }
    return null;
}

fn findFirstNum(line: []const u8) ?u64 {
    for (0..line.len) |i| {
        if (sliceToNum(line[i..])) |num| {
            return num;
        }
    }
    return null;
}

fn findLastNum(line: []const u8) ?u64 {
    // This is so janky
    var i: isize = @intCast(line.len - 1);
    while (i >= 0) : (i -= 1) {
        if (sliceToNum(line[@intCast(i)..])) |num| {
            return num;
        }
    }
    return null;
}

fn partTwo(line: []const u8) ?u64 {
    const tens = findFirstNum(line) orelse return null;
    const ones = findLastNum(line) orelse return null;
    return tens * 10 + ones;
}

pub fn main() !void {
    var gpa = Alloc{};
    defer _ = gpa.deinit();
    var alloc = gpa.allocator();
    const input = try fs.cwd().readFileAlloc(alloc, "input", 2 << 20);
    defer alloc.free(input);
    var lines = mem.splitScalar(u8, input, '\n');
    var sumOne: u64 = 0;
    var sumTwo: u64 = 0;
    while (lines.next()) |line| {
        if (line.len == 0) {
            continue;
        }
        sumOne += partOne(line) orelse return error.FailedToParse;
        sumTwo += partTwo(line) orelse return error.FailedToParse;
    }
    try out.print("Sum of all numeric only values is {}\n", .{sumOne});
    try out.print("Sum of all numeric / alphabetic values is {}\n", .{sumTwo});
}
