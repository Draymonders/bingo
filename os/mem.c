#include <unistd.h>
#include <stdio.h>
#include <stdlib.h>

extern void Spin(int t);

int main(int argc, char **argv) {
    int *p = (int*)(malloc(sizeof(int)));

    if (p == NULL) {
        exit(1);
//        return 1;
    }
    printf("(%d) p address pointer %p\n", getpid(), p);

    while (1) {
        Spin(1);
        *p = *p + 1;
        printf("(%d) *p %d\n", getpid(), *p);
    }

    return 0;
}