import numpy as np
import matplotlib.pyplot as plt

def mandelbrot(c, max_iter):
    z = c
    for i in range(max_iter):
        if abs(z) > 2.0:
            return i
        z = z * z + c
    return max_iter

def generate_mandelbrot(width, height, x_min, x_max, y_min, y_max, max_iter):
    image = np.zeros((height, width), dtype=np.uint8)
    for x in range(width):
        for y in range(height):
            real = x_min + (x / width) * (x_max - x_min)
            imag = y_min + (y / height) * (y_max - y_min)
            c = complex(real, imag)
            color = mandelbrot(c, max_iter)
            image[y, x] = np.uint8(color)
    return image

width = 800
height = 600
x_min, x_max = -2.5, 1.5
y_min, y_max = -1.5, 1.5
max_iter = 1000

mandelbrot_image = generate_mandelbrot(width, height, x_min, x_max, y_min, y_max, max_iter)

plt.imshow(mandelbrot_image, cmap='inferno', extent=(x_min, x_max, y_min, y_max))
plt.colorbar(label='Iterations')
plt.title('Mandelbrot Set')
plt.xlabel('Real')
plt.ylabel('Imaginary')
plt.show()
