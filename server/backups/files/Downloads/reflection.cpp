#include <iostream>

using namespace std;

class MyClass {};

int main() {
	cout << typeid( int ).name() << endl;
	int x;
	cout << typeid( x ).name() << endl;
	cout << typeid( 2 + 2.8 ).name() << endl;
    string frase = "Olá mundo!";
    cout << typeid( frase ).name() << endl;
    MyClass mc;
    cout << typeid( mc ).name() << endl;
}