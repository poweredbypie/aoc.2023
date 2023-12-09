import Data.Maybe (mapMaybe)
import Data.Tuple (swap)
import Text.Read (readMaybe)

lineToNums :: String -> [Int]
lineToNums = mapMaybe readMaybe . words

deriv :: (Int -> [Int] -> Int) -> [Int] -> Int
deriv action nums =
  let pairs = zip nums (tail nums)
      diffs = map (uncurry (-) . swap) pairs
   in if all (== 0) diffs
        then action 0 nums
        else action (deriv action diffs) nums

predict :: [Int] -> Int
predict = deriv $ \val nums -> val + last nums

retrodict :: [Int] -> Int
retrodict = deriv $ \val nums -> head nums - val

main :: IO ()
main =
  do
    file <- readFile "input"
    let nums = map lineToNums (lines file)
    putStr "Sum of predicted values is "
    print $ sum $ map predict nums
    putStr "Sum of retrodicted values is "
    print $ sum $ map retrodict nums
