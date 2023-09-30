#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>
#include <math.h>

#define WIDTH 800
#define HEIGHT 600
#define MAX_ITER 1000

typedef struct {
    uint8_t r, g, b;
} Color;

Color colormap[MAX_ITER];

void init_colormap() {
    for (int i = 0; i < MAX_ITER; i++) {
        colormap[i].r = i % 256;
        colormap[i].g = (i * 7) % 256;
        colormap[i].b = (i * 13) % 256;
    }
}

void generate_mandelbrot(uint8_t *image) {
    double x_scale = (3.5 / WIDTH);
    double y_scale = (2.0 / HEIGHT);
    double x_offset = -2.5;
    double y_offset = -1.0;
    
    for (int y = 0; y < HEIGHT; y++) {
        for (int x = 0; x < WIDTH; x++) {
            double real = x * x_scale + x_offset;
            double imag = y * y_scale + y_offset;
            double z_real = real;
            double z_imag = imag;
            
            int iter = 0;
            while (iter < MAX_ITER) {
                double z_real_temp = z_real * z_real - z_imag * z_imag + real;
                double z_imag_temp = 2 * z_real * z_imag + imag;
                z_real = z_real_temp;
                z_imag = z_imag_temp;
                
                if ((z_real * z_real + z_imag * z_imag) > 4.0) {
                    break;
                }
                iter++;
            }
            
            Color color = colormap[iter];
            image[(y * WIDTH + x) * 3] = color.b;
            image[(y * WIDTH + x) * 3 + 1] = color.g;
            image[(y * WIDTH + x) * 3 + 2] = color.r;
        }
    }
}

void save_image(uint8_t *image, const char *filename) {
    FILE *file = fopen(filename, "wb");
    if (!file) {
        perror("Error opening file");
        return;
    }
    
    int header_size = 54;
    int data_size = WIDTH * HEIGHT * 3;
    int file_size = header_size + data_size;
    
    uint8_t header[54] = {
        'B', 'M',        // Signature
        file_size, 0, 0, 0, // File size in bytes
        0, 0, 0, 0,       // Reserved
        54, 0, 0, 0,      // Offset to start of pixel data
        40, 0, 0, 0,      // Header size
        WIDTH, 0, 0, 0,   // Image width
        HEIGHT, 0, 0, 0,  // Image height
        1, 0,             // Number of color planes
        24, 0,            // Bits per pixel (3 bytes)
        0, 0, 0, 0,       // Compression method
        data_size, 0, 0, 0, // Image size
        0, 0, 0, 0,       // Horizontal resolution (pixels per meter)
        0, 0, 0, 0,       // Vertical resolution (pixels per meter)
        0, 0, 0, 0,       // Number of colors in palette
        0, 0, 0, 0        // Important colors
    };
    
    fwrite(header, sizeof(uint8_t), 54, file);
    fwrite(image, sizeof(uint8_t), data_size, file);
    fclose(file);
}

int main() {
    uint8_t *image = (uint8_t *)malloc(WIDTH * HEIGHT * 3 * sizeof(uint8_t));
    
    if (!image) {
        perror("Error allocating memory");
        return 1;
    }
    
    init_colormap();
    generate_mandelbrot(image);
    save_image(image, "mandelbrot.bmp");
    
    free(image);
    return 0;
}
