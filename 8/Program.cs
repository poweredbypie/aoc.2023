using System.Collections.Generic;
using System;
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

            this.Name = groups[1].Value;
            this.Left = groups[2].Value;
            this.Right = groups[3].Value;
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
        Move GetMove()
        {
            return _moves[_currMove];
        }
        // Returns the number of moves required to fully follow the map
        public int Follow()
        {
            const string start = "AAA";
            const string end = "ZZZ";

            for (var current = start; current != end; NextMove())
            {
                var node = _nodes[current];
                var move = GetMove();
                if (move == Move.Left)
                {
                    current = node.Left;
                }
                else
                {
                    current = node.Right;
                }
            }

            return _moveCount;
        }
    }

    static void Main(string[] args)
    {
        var lines = File.ReadAllLines("input");
        var map = new Map(lines);
        Console.WriteLine($"Moves needed to follow map is {map.Follow()}");
    }
}
