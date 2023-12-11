using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Text.RegularExpressions;

class Solution
{
    enum Move
    {
        Left,
        Right,
    }

    class Node
    {
        static Regex regex = new Regex(@"(.+) = \((.+), (.+)\)");

        public string Name;
        public string Left;
        public string Right;

        public static bool Parseable(string line)
        {
            return Node.regex.Match(line).Success;
        }
        public Node(string line)
        {
            var match = Node.regex.Match(line);
            if (!match.Success)
            {
                throw new InvalidDataException($"Couldn't match line against regex: {line}");
            }

            var groups = match.Groups;

            Name = groups[1].Value;
            Left = groups[2].Value;
            Right = groups[3].Value;
        }
    }

    class Map
    {
        Move[] _moves;
        int _currMove;
        int _moveCount;
        Dictionary<string, Node> _nodes;

        public Map(string[] lines)
        {
            _moves = lines[0]
                .Where(c => c == 'L' || c == 'R')
                .Select(c => c == 'L' ? Move.Left : Move.Right)
                .ToArray();

            _currMove = 0;
            // We call NextMove after the first move, so this fixes an edge case
            _moveCount = 0;

            _nodes = lines
                .Where(line => Node.Parseable(line))
                .Select(line => new Node(line))
                .ToDictionary(node => node.Name, node => node);
        }
        void NextMove()
        {
            _moveCount += 1;
            _currMove += 1;
            if (_currMove >= _moves.Count())
            {
                _currMove = 0;
            }
        }
        string FollowMove(string name)
        {
            var node = _nodes[name];
            return (_moves[_currMove] == Move.Left) ? node.Left : node.Right;
        }
        // Reset the number of moves and the instruction index
        void Reset()
        {
            _currMove = 0;
            _moveCount = 0;
        }
        // Returns the number of moves required to fully follow the map
        // Answer to Part A and part of Part B
        public int FollowLength(string start, Regex end)
        {
            Reset();
            var current = start;
            while (!end.IsMatch(current))
            {
                current = FollowMove(current);
                NextMove();
            }

            return _moveCount;
        }
        // Returns the number of moves required to follow the map
        // Part B: find coinciding steps for _all_ inputs
        // Using long because the value is too big otherwise
        public long FollowCoinciding(Regex start, Regex end)
        {
            // All the XXA nodes we need to iterate with
            var nodes = _nodes
                .Where(pair => start.IsMatch(pair.Key))
                .Select(pair => pair.Key)
                .ToArray();

            // Find the follow length for each of these nodes
            // This is probably wrong, but:
            // I thought each node had a "follow cycle" with differing offsets and lengths (like a FSM cycle)
            // However, running it with all the XXA inputs actually suggests that the values are the same
            // So I'm making that assumption now
            var lengths = nodes
                .Select(node => (long)FollowLength(node, end))
                .ToArray();

            // The LCM of all of the lengths is the answer
            return Solution.LeastCommonMultiple(lengths);
        }
    }
    static long GreatestCommonDenominator(long first, long second)
    {
        // Euclidean algorithm
        // https://en.wikipedia.org/wiki/Greatest_common_divisor#Euclidean_algorithm
        while (second != 0)
        {
            (first, second) = (second, first % second);
        }
        return first;
    }
    public static long LeastCommonMultiple(long[] nums)
    {
        // Taken from StackOverflow Lol
        // https://stackoverflow.com/questions/147515/least-common-multiple-for-3-or-more-numbers
        return nums.Aggregate((first, second) =>
        {
            return Math.Abs(first * second) / GreatestCommonDenominator(first, second);
        });
    }
    static void Main(string[] args)
    {
        var lines = File.ReadAllLines("input");
        var map = new Map(lines);
        var singular = map.FollowLength("AAA", new Regex("ZZZ"));
        Console.WriteLine($"Moves needed to follow map from AAA -> ZZZ is {singular}");
        var multiple = map.FollowCoinciding(new Regex("..A"), new Regex("..Z"));
        Console.WriteLine($"Moves needed to follow all XXA -> XXZ at the same time is {multiple}");
    }
}
