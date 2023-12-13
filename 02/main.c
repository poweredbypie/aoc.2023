#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

bool cubes_possible(int count, char* color) {
#define check_color(str, lit, max)                \
    do {                                          \
        if (!memcmp(str, lit, sizeof(lit) - 1)) { \
            return count <= max;                  \
        }                                         \
    } while (false);

    check_color(color, "red", 12);
    check_color(color, "green", 13);
    check_color(color, "blue", 14);

#undef check_color

    printf("Unexpected color %s, exiting\n", color);
    exit(1);
}

bool round_possible(const char* round, const char* end) {
    for (const char* cubes = round - 1; cubes < end && cubes != NULL; cubes = strchr(cubes, ',')) {
        // Skip past the comma (the start edge case is fixed by round - 1)
        cubes += 1;
        int count;
        // Hardcoded because the longest color should be 5 chars long
        char color[41];
        color[40] = 0;
        if (sscanf(cubes, "%d %40[^,;]", &count, color) != 2) {
            printf("Couldn't find cube count or color for %s, exiting\n", cubes);
            exit(1);
        }
        // Skip space
        if (!cubes_possible(count, color)) {
            return false;
        }
    }
    return true;
}

// Part A
int game_possible(const char* line) {
    int id;
    if (sscanf(line, "Game %d", &id) != 1) {
        printf("Couldn't parse ID of game, exiting\n");
        exit(1);
    }

    for (const char* round = strchr(line, ':'); round != NULL; round = strchr(round, ';')) {
        // Skip past the character we found since we don't want it
        round += 1;
        const char* end = strchr(round, ';');
        if (end == NULL) {
            // Some large value so the end ptr for round possible isn't bad
            end = (const char*)-1;
        }

        if (!round_possible(round, end)) {
            // 0 means the round isn't possible and doesn't affect the sum
            return 0;
        }
    }

    return id;
}

int max_of(int one, int two) {
    if (one > two) {
        return one;
    }
    else {
        return two;
    }
}

int round_max(const char* round, const char* end, const char* for_color) {
    int max = 0;
    for (const char* cubes = round - 1; cubes < end && cubes != NULL; cubes = strchr(cubes, ',')) {
        // Skip past the comma (the start edge case is fixed by round - 1)
        cubes += 1;
        int count;
        // Hardcoded because the longest color should be 5 chars long
        char color[41];
        color[40] = 0;
        if (sscanf(cubes, "%d %40[^,;]", &count, color) != 2) {
            printf("Couldn't find cube count or color for %s, exiting\n", cubes);
            exit(1);
        }
        if (!memcmp(color, for_color, strlen(for_color) - 1)) {
            max = max_of(count, max);
        }
    }
    return max;
}

// Part B
// Kid named writing reusable code
int game_power(const char* line) {
    int max_red = 0;
    int max_green = 0;
    int max_blue = 0;
    for (const char* round = strchr(line, ':'); round != NULL; round = strchr(round, ';')) {
        // Skip past the character we found since we don't want it
        round += 1;
        const char* end = strchr(round, ';');
        if (end == NULL) {
            // Some large value so the end ptr for round possible isn't bad
            end = (const char*)-1;
        }
        max_red = max_of(max_red, round_max(round, end, "red"));
        max_green = max_of(max_green, round_max(round, end, "green"));
        max_blue = max_of(max_blue, round_max(round, end, "blue"));
    }

    return max_red * max_green * max_blue;
}

int main(void) {
    FILE* file = fopen("input", "r");
    if (file == NULL) {
        return 1;
    }
    // Again, hardcoded because longest line is definitely shorter than 400 chars
    char line[401];
    // Null terminator
    line[400] = 0;
    int possible = 0;
    int power = 0;
    while (fscanf(file, " %400[^\n]", line) != EOF) {
        possible += game_possible(line);
        power += game_power(line);
    };
    printf("Sum of possible game IDs is %d\n", possible);
    printf("Sum of power of each game is %d\n", power);
    fclose(file);
}
