import Data.Char (digitToInt, isDigit)
import Data.Maybe (mapMaybe)

parseCharA :: Char -> Maybe Int
parseCharA char
  | isDigit char = Just $ digitToInt char
  | otherwise = Nothing

parseLineA :: String -> [Int]
parseLineA = mapMaybe parseCharA

matchLineA :: String -> Int
matchLineA line =
  let parsed = parseLineA line
      tens = head parsed * 10
      ones = last parsed
   in tens + ones

partA :: [String] -> Int
partA lines = sum $ map matchLineA lines

main :: IO ()
main =
  do
    file <- readFile "input"
    putStr "Part A result is: "
    print $ partA $ lines file
