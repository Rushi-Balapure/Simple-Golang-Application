# 🎮 Memory Match - Neon Edition

A beautiful, cyberpunk-themed memory card matching game built with Go!

![Go Version](https://img.shields.io/badge/Go-1.18+-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/License-MIT-ff2d95)

## ✨ Features

- 🃏 **Classic Memory Game** - Match pairs of cards to win
- 🎨 **Stunning Neon Aesthetics** - Cyberpunk-inspired design with glowing effects
- 📊 **Live Leaderboard** - Compete for the top spot
- 🎯 **3 Difficulty Levels** - Easy (6 pairs), Medium (8 pairs), Hard (10 pairs)
- ⚡ **Fast & Responsive** - Built with Go's powerful HTTP server
- 📱 **Mobile Friendly** - Play on any device

## 🚀 Quick Start

### Prerequisites

- Go 1.18 or higher

### Run the Game

```bash
go run main.go
```

Then open your browser and navigate to:

```
http://localhost:8080
```

## 🎯 How to Play

1. Enter your name (optional)
2. Select difficulty level
3. Click **START GAME**
4. Click on cards to flip them
5. Remember the positions and match pairs
6. Complete all matches with minimum moves to top the leaderboard!

## 🛠️ Tech Stack

- **Backend**: Go (net/http)
- **Frontend**: Vanilla HTML, CSS, JavaScript
- **Fonts**: Orbitron, Rajdhani (Google Fonts)
- **No external dependencies!**

## 📁 Project Structure

```
.
├── main.go          # Go server with embedded HTML game
├── go.mod           # Go module file
├── LICENSE          # MIT License
└── README.md        # This file
```

## 🎨 Design Highlights

- Animated grid background with 3D perspective
- Floating neon glow orbs
- Smooth card flip animations
- Pulsing neon typography
- Glassmorphism UI elements
- Responsive design for all screen sizes

## 🏆 Leaderboard

The game tracks:
- Number of moves (lower is better)
- Time taken to complete
- Top 10 players are displayed

## 📜 License

MIT License - feel free to use, modify, and distribute!

---

Made with 💜 and Go
