use std::collections::HashMap;
use std::fs::read_to_string;

#[derive(Clone, Debug, Eq, Hash, Ord, PartialEq, PartialOrd)]
enum Card {
    J,
    Number(u32),
    T,
    Q,
    K,
    A,
}

impl Card {
    fn from_char(char: char) -> Option<Card> {
        // Too lazy to make this more concise idc
        match char {
            'J' => Some(Card::J),
            '2' => Some(Card::Number(2)),
            '3' => Some(Card::Number(3)),
            '4' => Some(Card::Number(4)),
            '5' => Some(Card::Number(5)),
            '6' => Some(Card::Number(6)),
            '7' => Some(Card::Number(7)),
            '8' => Some(Card::Number(8)),
            '9' => Some(Card::Number(9)),
            'T' => Some(Card::T),
            'Q' => Some(Card::Q),
            'K' => Some(Card::K),
            'A' => Some(Card::A),
            _ => None,
        }
    }
}

#[derive(Debug, Eq, Ord, PartialEq, PartialOrd)]
enum HandKind {
    HighCard,
    OnePair,
    TwoPair,
    ThreeOfAKind,
    FullHouse,
    FourOfAKind,
    FiveOfAKind,
}

#[derive(Debug, Eq, Ord, PartialEq, PartialOrd)]
struct Hand {
    kind: HandKind,
    cards: Vec<Card>,
    bid: u32,
}

impl Hand {
    fn hand_kind(cards: &Vec<Card>) -> HandKind {
        assert!(cards.len() == 5);

        let (jokers, others): (Vec<_>, Vec<_>) = cards.iter().partition(|card| **card == Card::J);

        let mut map = HashMap::<Card, u32>::new();
        for card in others {
            *map.entry(card.clone()).or_insert(0) += 1;
        }

        if map.len() == 0 {
            // This means we have 5 jokers
            return HandKind::FiveOfAKind;
        }

        let mut max_same = *map.iter().max_by_key(|pair| pair.1).unwrap().1;
        // We let the jokers become part of the max, whatever it is
        max_same += jokers.len() as u32;
        match map.len() {
            1 => HandKind::FiveOfAKind,
            // 2 cards can only have (1, 4) or (2, 3) combos with 5 cards
            // (1, 4) or (4, 1) case
            2 if max_same == 4 => HandKind::FourOfAKind,
            // (2, 3) or (3, 2) case
            2 => HandKind::FullHouse,
            // 3 cards can be (1, 1, 3) or (1, 2, 2)
            // (1, 1, 3) or (1, 3, 1) or (3, 1, 1) case
            3 if max_same == 3 => HandKind::ThreeOfAKind,
            // (1, 2, 2) or (2, 1, 2) or (2, 2, 1) case
            3 => HandKind::TwoPair,
            // 4 cards has only the (1, 1, 1, 2) case
            4 => HandKind::OnePair,
            // 5 cards has only the (1, 1, 1, 1, 1) case
            5 => HandKind::HighCard,
            _ => panic!("Map length was out of range"),
        }
    }
    fn new(string: String) -> Option<Hand> {
        let split = string.split_once(" ")?;
        let cards = split
            .0
            .chars()
            .filter_map(Card::from_char)
            .collect::<Vec<_>>();

        let bid = split.1.parse::<u32>().ok()?;
        let kind = Hand::hand_kind(&cards);

        Some(Hand { kind, cards, bid })
    }
}

fn main() -> std::io::Result<()> {
    let mut hands = read_to_string("input")?
        .lines()
        .map(String::from)
        .filter_map(Hand::new)
        .collect::<Vec<_>>();

    hands.sort();
    let sum = hands
        .iter()
        .enumerate()
        .fold(0, |sum, pair| sum + ((pair.0 as u32 + 1) * pair.1.bid));

    println!("Sum of bids by rank is {}", sum);

    Ok(())
}
