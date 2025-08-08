# Simfight-Tactics — TFT DPS Simulator

## Overview
Simfight-Tactics is a web-first damage simulator for Teamfight Tactics (TFT). It estimates per-unit and team DPS under different comps, items, traits, stats, augments, power-ups, on-hit effects, spell cadence, targeting behavior, buffs/debuffs, range and seeded RNG. The goals are reproducibility (fixed seeds and deterministic runs), performance (parallel Monte-Carlo with a worker pool), and transparency (structured event logs), wrapped in a responsive, server-rendered UI.

## Install
### Requirements
- Go ≥ 1.22
- Git
- (Optional) Postgres if you plan to persist runs; otherwise file mode is fine.

### Get the code
```bash
git clone https://github.com/0xm0-v1/simfight-tactics.git
cd simfight-tactics
