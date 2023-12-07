{-# LANGUAGE OverloadedStrings #-}

import Data.Maybe (mapMaybe)
import Data.Text (pack, replace, takeWhileEnd, unpack)
import Debug.Trace (traceShow)
import Text.Read (readMaybe)

data Race = Race
  { time :: Int,
    dist :: Int
  }
  deriving (Show)

recordInputs :: Race -> (Double, Double)
recordInputs race = ((t - discrim) / 2, (t + discrim) / 2)
  where
    t = fromIntegral $ time race
    d = fromIntegral $ dist race
    discrim = sqrt ((t * t) - (4 * d))

recordRange :: Race -> Int
recordRange race = floor big - ceiling small + 1
  where
    (small, big) = recordInputs race

numsInLine :: String -> [Int]
numsInLine = mapMaybe readMaybe . words

-- Part A: multiple races
getRaces :: [String] -> [Race]
getRaces lines = zipWith Race times dists
  where
    times = numsInLine $ head lines
    dists = numsInLine $ last lines

numFromLine :: String -> Int
numFromLine = read . unpack . takeWhileEnd (/= ':') . replace " " "" . pack

-- Part B: one race
getRace :: [String] -> Race
getRace lines = Race time dist
  where
    time = numFromLine $ head lines
    dist = numFromLine $ last lines

partB :: [String] -> Int
partB = recordRange . getRace

partA :: [String] -> Int
partA = product . map recordRange . getRaces

main =
  do
    file <- readFile "input"
    putStr "Answer to part A is: "
    print $ partA $ lines file
    putStr "Answer to part B is: "
    print $ partB $ lines file
