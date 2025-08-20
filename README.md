# Simfight-Tactics — TFT DPS Simulator

## Overview
Simfight-Tactics is a web-first damage simulator for Teamfight Tactics (TFT). It estimates per-unit and team DPS under different comps, items, traits, stats, augments, power-ups, on-hit effects, spell cadence, targeting behavior, buffs/debuffs, range and seeded RNG. The goals are reproducibility (fixed seeds and deterministic runs), performance (parallel Monte-Carlo with a worker pool), and transparency (structured event logs), wrapped in a responsive, server-rendered UI.

## Install
### Requirements
- Go ≥ 1.25
- Git

### Get the code
```bash
git clone https://github.com/0xm0-v1/simfight-tactics.git
cd simfight-tactics
```

## Documentation

### 📦 Package Documentation

#### Units Package
The core domain model for TFT champions and their statistics.

📚 **[Units Package Guide](./internal/models/units/README.md)**

## Development

### Quick Start
```bash
# Run tests
make test

# Run linter
make lint

# Start server
make run
```

### Project Structure
```
simfight-tactics/
├── cmd/
│   ├── sft/          # CLI tool
│   └── sftd/         # HTTP server
├── internal/
│   ├── config/       # Game configuration (patches, roles)
│   │   └── set15/    # TFT Set 15 data
│   └── models/
│       └── units/    # Champion models [📚 Documentation](./internal/models/units/README.md)
└── docs/             # Additional documentation
```