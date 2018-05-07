#include <stdio.h>

int main() {
    int age;
    scanf("%d\n", &age);
    printf("wow, you are %d years old!\n", age);

    {
        int x = 3;
        printf("%d\n", x);
        x = 4;
        int* y = &x;
        printf("Hello world %d, %d!\n", x, *y);
    }

    return 0;
}