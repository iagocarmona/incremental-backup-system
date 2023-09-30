#include <iostream>

using namespace std;

class Shape {
    public:
        virtual void draw1() = 0;
};

class Rectangle: public Shape {
    public:
        void draw1() { cout << "Retângulo 1\n"; }
        void draw2() { cout << "Retângulo 2\n"; }
};

class Square: public Rectangle {
    public:
        void draw1() { cout << "Quadrado 1\n"; }
        void draw2() { cout << "Quadrado 2\n"; }
};

int main() {
    Shape* s = new Rectangle;
    s->draw1(); // vinculação dinâmica
    
    s = new Square;
    s->draw1(); // vinculação dinâmica
    
    Rectangle *r = new Rectangle;
    r->draw2(); // vinculação estática
    
    r = new Square;
    r->draw2(); // vinculação estática
}
\