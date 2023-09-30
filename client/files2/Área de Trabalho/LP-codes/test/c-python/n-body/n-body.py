import numpy as np

# Parâmetros da simulação
num_bodies = 3
num_steps = 100
dt = 0.1

# Massas e posições iniciais
masses = np.array([10, 5, 3])
positions = np.array([[0, 0], [5, 0], [0, 8]], dtype=np.float64)  # Corrigido o tipo de dados
velocities = np.zeros((num_bodies, 2), dtype=np.float64)  # Corrigido o tipo de dados

def calculate_forces(positions, masses):
    forces = np.zeros_like(positions, dtype=np.float64)
    G = 1.0  # Constante gravitacional

    for i in range(num_bodies):
        for j in range(num_bodies):
            if i != j:
                r = positions[j] - positions[i]
                r_magnitude = np.linalg.norm(r)
                force_magnitude = (G * masses[i] * masses[j]) / (r_magnitude ** 2)
                force = force_magnitude * r / r_magnitude
                forces[i] += force
    
    return forces

def update_positions_and_velocities(positions, velocities, forces, masses, dt):
    accelerations = forces / masses[:, np.newaxis]
    velocities += accelerations * dt
    positions += velocities * dt

# Simulação
for step in range(num_steps):
    print(f"Step {step + 1}:")
    for i in range(num_bodies):
        print(f"Body {i + 1} - Position: {positions[i]}, Velocity: {velocities[i]}")
    
    forces = calculate_forces(positions, masses)
    update_positions_and_velocities(positions, velocities, forces, masses, dt)

print("Simulation completed.")
