#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

// More type safety
#define realloc(type, buf, len) (type*)realloc(buf, len * sizeof(type))
#define malloc(type, len) (type*)malloc(len * sizeof(type))

// -- Util functions --
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
        buf = realloc(char, buf, len);
        buf[len - 1] = c;
    }
    // Add null terminator
    buf = realloc(char, buf, len + 1);
    buf[len] = 0;
    return buf;
}

// Clone a string, caller owns new string
char* str_clone(const char* str) {
    size_t len = strlen(str);
    char* buf = malloc(char, len + 1);
    strncpy(buf, str, len + 1);
    return buf;
}

typedef struct {
    char* label;
    int focal;
} Lens;

// -- Box struct / methods --
typedef struct {
    Lens* lenses;
    size_t length;
} Box;

Box* box_new(void) {
    Box* box = malloc(Box, 1);
    box->length = 0;
    box->lenses = malloc(Lens, 0);
    return box;
}

Lens* box_at(Box* box, size_t index) {
    if (index >= box->length) {
        printf("Index passed (%zu) is past box bounds (%zu elements)\n", index, box->length);
    }
    return &box->lenses[index];
}

void box_deinit(Box* box) {
    for (size_t i = 0; i < box->length; i += 1) {
        // Free all the labels we own
        free(box_at(box, i)->label);
    }
}

// Returns -1 if the index isn't found (max int value for size_t)
size_t box_find(Box* box, char* label) {
    for (size_t i = 0; i < box->length; i += 1) {
        if (!strcmp(box_at(box, i)->label, label)) {
            return i;
        }
    }
    return (size_t)-1;
}

void box_set(Box* box, char* label, int focal) {
    size_t pos = box_find(box, label);
    Lens* lens = NULL;
    if (pos == (size_t)-1) {
        // Add a new lens struct at the end
        box->length += 1;
        box->lenses = realloc(Lens, box->lenses, box->length);
        lens = box_at(box, box->length - 1);
    }
    else {
        lens = box_at(box, pos);
        // Free the label of the last lens first
        free(lens->label);
    }
    // Overwrite values
    lens->label = label;
    lens->focal = focal;
}

void box_remove(Box* box, char* label) {
    size_t pos = box_find(box, label);
    if (pos == (size_t)-1) {
        // Didn't find the label so we can just stop
        return;
    }

    // Free the lens label at the position (we'll lose the pointer after shifting)
    free(box_at(box, pos)->label);

    // Shift all the boxes over, overwriting the one at [pos]
    for (size_t i = pos; i < box->length - 1; i += 1) {
        *box_at(box, i) = *box_at(box, i + 1);
    }
    // There's some garbage at the end of the array but whatever
    // We'll either overwrite it later and it won't be read because we're mutating length
    box->length -= 1;
}

size_t box_power(Box* box) {
    size_t power = 0;
    for (size_t i = 0; i < box->length; i += 1) {
        power += (i + 1) * box_at(box, i)->focal;
    }
    return power;
}

enum {
    HASH_MOD = 256,
};

// -- HashMap struct / methods --
typedef struct {
    Box** boxes;
} HashMap;

size_t hash(char* string) {
    size_t value = 0;

    size_t len = strlen(string);
    for (size_t i = 0; i < len; i += 1) {
        value += string[i];
        value *= 17;
        value %= HASH_MOD;
    }
    return value;
}

Box* hash_map_box_for(HashMap* map, char* label) {
    return map->boxes[hash(label)];
}

void hash_map_remove(HashMap* map, char* label) {
    printf("Removing label %s\n", label);
    Box* box = hash_map_box_for(map, label);
    box_remove(box, label);
}

void hash_map_set(HashMap* map, char* label, int focal) {
    printf("Adding label %s with focal %d\n", label, focal);
    Box* box = hash_map_box_for(map, label);
    box_set(box, label, focal);
}

void hash_map_op(HashMap* map, const char* string) {
    char* op = strpbrk(string, "-=");
    if (op == NULL) {
        printf("Couldn't find operator symbol in input string %s\n", string);
        exit(1);
    }

    // Create a label string
    size_t len = (size_t)(op - string) + 1;
    char* label = malloc(char, len);
    strncpy(label, string, len - 1);
    // Add null terminator
    label[len - 1] = 0;

    if (*op == '-') {
        hash_map_remove(map, label);
        // Since this isn't being added it's owned by us and we need to clean it up
        free(label);
    }
    else {
        int focal = 0;
        if (sscanf(op + 1, "%d", &focal) != 1) {
            printf("Couldn't find focal value in string %s\n", op + 1);
            exit(1);
        }
        // The label now belongs to the new Lens struct added
        hash_map_set(map, label, focal);
    }
}

HashMap* hash_map_new(void) {
    HashMap* map = malloc(HashMap, 1);
    map->boxes = malloc(Box*, HASH_MOD);
    for (size_t i = 0; i < HASH_MOD; i += 1) {
        map->boxes[i] = box_new();
    }
    return map;
}

void hash_map_deinit(HashMap* map) {
    for (size_t i = 0; i < HASH_MOD; i += 1) {
        // Free each box
        box_deinit(map->boxes[i]);
        free(map->boxes[i]);
    }
    // Free box array
    free(map->boxes);
}

HashMap* hash_map_from(const char* line) {
    HashMap* map = hash_map_new();
    // Duplicate the string because strtok mutates
    char* buf = str_clone(line);

    for (char* op = strtok(buf, ","); op; op = strtok(NULL, ",")) {
        hash_map_op(map, op);
    }
    free(buf);
    return map;
}

size_t hash_map_power(HashMap* map) {
    size_t power = 0;
    for (size_t i = 0; i < HASH_MOD; i += 1) {
        power += box_power(map->boxes[i]) * (i + 1);
    }
    return power;
}

// Part A
size_t hash_line(const char* line) {
    size_t sum = 0;

    // Duplicate the string because strtok mutates
    char* buf = str_clone(line);

    for (char* split = strtok(buf, ","); split; split = strtok(NULL, ",")) {
        sum += hash(split);
    }

    // Cleanup and exit
    free(buf);
    return sum;
}

int main(void) {
    FILE* file = fopen("input", "r");
    if (file == NULL) {
        printf("Couldn't open input file\n");
        exit(1);
    }
    char* line = get_line(file);
    // Part A
    printf("Sum of hash of all strings in line is %zu\n", hash_line(line));

    // Part B
#if 1
    HashMap* map = hash_map_from(line);
#else
    HashMap* map = hash_map_from("rn=1,cm-,qp=3,cm=2,qp-,pc=4,ot=9,ab=5,pc-,pc=6,ot=7");
#endif
    printf("Power of resulting hash map is %zu\n", hash_map_power(map));

    // Cleanups
    hash_map_deinit(map);
    free(map);
    free(line);
    fclose(file);
}
