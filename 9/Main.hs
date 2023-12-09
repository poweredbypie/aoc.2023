import Data.Maybe (mapMaybe)
import Data.Tuple (swap)
import Text.Read (readMaybe)

lineToNums :: String -> [Int]
lineToNums = mapMaybe readMaybe . words

predict :: [Int] -> Int
predict nums =
  let pairs = zip nums (tail nums)
      diffs = map (uncurry (-) . swap) pairs
   in if all (== 0) diffs
        then last nums
        else predict diffs + last nums

predictAll :: [String] -> Int
predictAll = sum . map (predict . lineToNums)

retrodict :: [Int] -> Int
retrodict nums =
  let pairs = zip nums (tail nums)
      diffs = map (uncurry (-) . swap) pairs
   in if all (== 0) diffs
        then head nums
        else head nums - retrodict diffs

retrodictAll :: [String] -> Int
retrodictAll = sum . map (retrodict . lineToNums)

main =
  do
    file <- readFile "input"
    putStr "Sum of predicted values is "
    print $ predictAll $ lines file
    putStr "Sum of retrodicted values is "
    print $ retrodictAll $ lines file
