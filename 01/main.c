#include <stdio.h>
#include <stdlib.h>
#include <string.h>

int get_value_a(char* line, int len) {
    int first = 0;
    int last = 0;
    for (int i = 0; i < len; i += 1) {
        if (line[i] < '1' || line[i] > '9') {
            continue;
        }

        int val = line[i] - '0';
        if (first == 0) {
            first = val;
        }
        last = val;
    }
    // Input does not contain any zeroes
    if (first == 0 || last == 0) {
        printf("Couldn't find at least one digit in line %s", line);
        exit(1);
    }
    return first * 10 + last;
}

int get_value_b(char* line, int len) {
    int first = 0;
    int last = 0;
    for (int i = 0; i < len; i += 1) {
        int val = 0;
        char* slice = line + i;
#define litcmp(str, lit, num)                     \
    do {                                          \
        if (!memcmp(str, lit, sizeof(lit) - 1)) { \
            val = num;                            \
        }                                         \
    } while (0)
        // I love programming in C guys! It's so fun!
        if (line[i] >= '1' && line[i] <= '9') {
            val = line[i] - '0';
        }
        else {
            litcmp(slice, "one", 1);
            litcmp(slice, "two", 2);
            litcmp(slice, "three", 3);
            litcmp(slice, "four", 4);
            litcmp(slice, "five", 5);
            litcmp(slice, "six", 6);
            litcmp(slice, "seven", 7);
            litcmp(slice, "eight", 8);
            litcmp(slice, "nine", 9);
        }
#undef litcmp
        if (val == 0) {
            continue;
        }

        if (first == 0) {
            first = val;
        }
        last = val;
    }
    return first * 10 + last;
}

int main(void) {
    FILE* file = fopen("input", "r");
    if (file == NULL) {
        return 1;
    }

    int sum_a = 0;
    int sum_b = 0;

    // This is just hardcoded because I know the max input length
    char line[401];
    // Null terminator
    line[400] = 0;
    while (fscanf(file, " %400[^\n]", line) != EOF) {
        int len = strlen(line);
        sum_a += get_value_a(line, len);
        sum_b += get_value_b(line, len);
    }

    printf("Sum of all parsed numbers (no letters) is %d\n", sum_a);
    printf("Sum of all parsed numbers (with letters) is %d\n", sum_b);

    fclose(file);
}
