#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

// getline() isn't cross platform so
char* get_line(FILE* input) {
    char* buf = NULL;
    size_t len = 0;
    while (true) {
        char c = fgetc(input);
        if (c == '\n' || c == EOF) {
            break;
        }
        len += 1;
        buf = realloc(buf, len);
        buf[len - 1] = c;
    }
    // Add null terminator
    buf = realloc(buf, len + 1);
    buf[len] = 0;
    return buf;
}

size_t hash(char* string) {
    size_t value = 0;

    size_t len = strlen(string);
    for (size_t i = 0; i < len; i += 1) {
        value += string[i];
        value *= 17;
        value %= 256;
    }
    return value;
}

// This mutates line because strtok edits the string given
size_t hash_line(char* line) {
    size_t sum = 0;

    for (char* split = strtok(line, ","); split; split = strtok(NULL, ",")) {
        sum += hash(split);
    }
    return sum;
}

int main(void) {
    FILE* file = fopen("input", "r");
    if (file == NULL) {
        exit(1);
    }
    char* line = get_line(file);
    printf("Sum of hash of all strings in line is %zu\n", hash_line(line));
    free(line);
    fclose(file);
}
