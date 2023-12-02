import re

def parse(num_str: str) -> int:
    convert = {
        'one': 1,
        'two': 2,
        'three': 3,
        'four': 4,
        'five': 5,
        'six': 6,
        'seven': 7,
        'eight': 8,
        'nine': 9
    }

    try:
        return int(num_str)
    except:
        converted = convert.get(num_str)
        if converted is None:
            raise ValueError
        return converted

def parse_line(line: str) -> int:
    # Allow overlapping
    found = re.findall('(?=(one|two|three|four|five|six|seven|eight|nine|[1-9]))', line)
    num = parse(found[0]) * 10 + parse(found[-1])
    return num

def main():
    sum = 0

    with open('input') as file:
        lines = file.readlines()
        for line in lines:
            sum += parse_line(line)
        
    print(f'The sum is {sum}')

if __name__ == "__main__":
    main()