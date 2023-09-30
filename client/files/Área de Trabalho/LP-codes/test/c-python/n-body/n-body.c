#include <stdio.h>
#include <math.h>

#define NUM_BODIES 3
#define NUM_STEPS 100
#define DT 0.1
#define G 1.0

typedef struct {
    double mass;
    double position[2];
    double velocity[2];
} Body;

void calculate_forces(Body bodies[NUM_BODIES], double forces[NUM_BODIES][2]) {
    for (int i = 0; i < NUM_BODIES; i++) {
        forces[i][0] = 0.0;
        forces[i][1] = 0.0;
        
        for (int j = 0; j < NUM_BODIES; j++) {
            if (i != j) {
                double r[2];
                r[0] = bodies[j].position[0] - bodies[i].position[0];
                r[1] = bodies[j].position[1] - bodies[i].position[1];
                double r_magnitude = sqrt(r[0] * r[0] + r[1] * r[1]);
                double force_magnitude = (G * bodies[i].mass * bodies[j].mass) / (r_magnitude * r_magnitude);
                forces[i][0] += force_magnitude * r[0] / r_magnitude;
                forces[i][1] += force_magnitude * r[1] / r_magnitude;
            }
        }
    }
}

void update_positions_and_velocities(Body bodies[NUM_BODIES], double forces[NUM_BODIES][2]) {
    for (int i = 0; i < NUM_BODIES; i++) {
        bodies[i].velocity[0] += forces[i][0] / bodies[i].mass * DT;
        bodies[i].velocity[1] += forces[i][1] / bodies[i].mass * DT;
        
        bodies[i].position[0] += bodies[i].velocity[0] * DT;
        bodies[i].position[1] += bodies[i].velocity[1] * DT;
    }
}

int main() {
    Body bodies[NUM_BODIES] = {
        {10.0, {0.0, 0.0}, {0.0, 0.0}},
        {5.0, {5.0, 0.0}, {0.0, 0.0}},
        {3.0, {0.0, 8.0}, {0.0, 0.0}}
    };

    double forces[NUM_BODIES][2];

    for (int step = 0; step < NUM_STEPS; step++) {
        calculate_forces(bodies, forces);
        update_positions_and_velocities(bodies, forces);
        
        printf("Step %d:\n", step + 1);
        for (int i = 0; i < NUM_BODIES; i++) {
            printf("Body %d: Position (%lf, %lf)\n", i, bodies[i].position[0], bodies[i].position[1]);
        }
        printf("\n");
    }

    return 0;
}
