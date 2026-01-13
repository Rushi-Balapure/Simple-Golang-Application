package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

// GameScore represents a player's score
type GameScore struct {
	PlayerName string    `json:"playerName"`
	Moves      int       `json:"moves"`
	TimeTaken  float64   `json:"timeTaken"`
	Timestamp  time.Time `json:"timestamp"`
}

var (
	leaderboard []GameScore
	mu          sync.RWMutex
)

func main() {
	fmt.Println("🎮 Memory Match Game Server")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("Starting server on http://localhost:8080")

	// Main game page
	http.HandleFunc("/", handleHome)

	// API endpoints
	http.HandleFunc("/api/leaderboard", handleLeaderboard)
	http.HandleFunc("/api/score", handleScore)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Memory Match | Neon Edition</title>
    <link href="https://fonts.googleapis.com/css2?family=Orbitron:wght@400;700;900&family=Rajdhani:wght@300;500;700&display=swap" rel="stylesheet">
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        :root {
            --neon-pink: #ff2d95;
            --neon-cyan: #00f5ff;
            --neon-purple: #b829dd;
            --neon-yellow: #f5ff00;
            --dark-bg: #0a0a0f;
            --card-bg: #12121a;
        }

        body {
            font-family: 'Rajdhani', sans-serif;
            background: var(--dark-bg);
            min-height: 100vh;
            overflow-x: hidden;
            color: #fff;
        }

        /* Animated background */
        .bg-grid {
            position: fixed;
            top: 0;
            left: 0;
            width: 100%;
            height: 100%;
            background-image: 
                linear-gradient(rgba(0, 245, 255, 0.03) 1px, transparent 1px),
                linear-gradient(90deg, rgba(0, 245, 255, 0.03) 1px, transparent 1px);
            background-size: 50px 50px;
            animation: gridMove 20s linear infinite;
            pointer-events: none;
            z-index: 0;
        }

        @keyframes gridMove {
            0% { transform: perspective(500px) rotateX(60deg) translateY(0); }
            100% { transform: perspective(500px) rotateX(60deg) translateY(50px); }
        }

        .bg-glow {
            position: fixed;
            width: 600px;
            height: 600px;
            border-radius: 50%;
            filter: blur(150px);
            opacity: 0.3;
            pointer-events: none;
            z-index: 0;
        }

        .glow-1 {
            top: -200px;
            left: -200px;
            background: var(--neon-pink);
            animation: float1 8s ease-in-out infinite;
        }

        .glow-2 {
            bottom: -200px;
            right: -200px;
            background: var(--neon-cyan);
            animation: float2 10s ease-in-out infinite;
        }

        @keyframes float1 {
            0%, 100% { transform: translate(0, 0); }
            50% { transform: translate(100px, 100px); }
        }

        @keyframes float2 {
            0%, 100% { transform: translate(0, 0); }
            50% { transform: translate(-100px, -100px); }
        }

        .container {
            position: relative;
            z-index: 1;
            max-width: 900px;
            margin: 0 auto;
            padding: 20px;
        }

        header {
            text-align: center;
            padding: 40px 0 30px;
        }

        h1 {
            font-family: 'Orbitron', sans-serif;
            font-size: 3.5rem;
            font-weight: 900;
            text-transform: uppercase;
            letter-spacing: 8px;
            background: linear-gradient(135deg, var(--neon-cyan), var(--neon-pink));
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
            background-clip: text;
            text-shadow: 0 0 80px rgba(0, 245, 255, 0.5);
            animation: titlePulse 2s ease-in-out infinite;
        }

        @keyframes titlePulse {
            0%, 100% { filter: brightness(1); }
            50% { filter: brightness(1.2); }
        }

        .subtitle {
            font-size: 1.1rem;
            color: var(--neon-purple);
            letter-spacing: 6px;
            margin-top: 10px;
            text-transform: uppercase;
        }

        .stats-bar {
            display: flex;
            justify-content: center;
            gap: 40px;
            margin: 30px 0;
            flex-wrap: wrap;
        }

        .stat {
            text-align: center;
            padding: 15px 30px;
            background: linear-gradient(135deg, rgba(18, 18, 26, 0.9), rgba(30, 30, 45, 0.9));
            border: 1px solid rgba(0, 245, 255, 0.3);
            border-radius: 10px;
            min-width: 140px;
            position: relative;
            overflow: hidden;
        }

        .stat::before {
            content: '';
            position: absolute;
            top: 0;
            left: -100%;
            width: 100%;
            height: 2px;
            background: linear-gradient(90deg, transparent, var(--neon-cyan), transparent);
            animation: scanline 3s linear infinite;
        }

        @keyframes scanline {
            0% { left: -100%; }
            100% { left: 100%; }
        }

        .stat-value {
            font-family: 'Orbitron', sans-serif;
            font-size: 2rem;
            font-weight: 700;
            color: var(--neon-cyan);
            text-shadow: 0 0 20px var(--neon-cyan);
        }

        .stat-label {
            font-size: 0.85rem;
            color: #888;
            text-transform: uppercase;
            letter-spacing: 2px;
            margin-top: 5px;
        }

        .game-board {
            display: grid;
            grid-template-columns: repeat(4, 1fr);
            gap: 15px;
            max-width: 500px;
            margin: 0 auto;
            padding: 30px;
            background: linear-gradient(135deg, rgba(18, 18, 26, 0.8), rgba(10, 10, 15, 0.9));
            border-radius: 20px;
            border: 1px solid rgba(0, 245, 255, 0.2);
            box-shadow: 
                0 0 60px rgba(0, 245, 255, 0.1),
                inset 0 0 60px rgba(0, 0, 0, 0.5);
        }

        .card {
            aspect-ratio: 1;
            perspective: 1000px;
            cursor: pointer;
        }

        .card-inner {
            position: relative;
            width: 100%;
            height: 100%;
            transition: transform 0.6s cubic-bezier(0.4, 0, 0.2, 1);
            transform-style: preserve-3d;
        }

        .card.flipped .card-inner {
            transform: rotateY(180deg);
        }

        .card-front, .card-back {
            position: absolute;
            width: 100%;
            height: 100%;
            backface-visibility: hidden;
            border-radius: 12px;
            display: flex;
            align-items: center;
            justify-content: center;
            font-size: 2.5rem;
        }

        .card-back {
            background: linear-gradient(135deg, #1a1a2e, #16213e);
            border: 2px solid rgba(0, 245, 255, 0.3);
            box-shadow: 
                0 0 20px rgba(0, 245, 255, 0.2),
                inset 0 0 20px rgba(0, 245, 255, 0.05);
        }

        .card-back::before {
            content: '?';
            font-family: 'Orbitron', sans-serif;
            font-size: 2rem;
            color: var(--neon-cyan);
            text-shadow: 0 0 15px var(--neon-cyan);
            opacity: 0.7;
        }

        .card-front {
            background: linear-gradient(135deg, #1a1a2e, #0f0f1a);
            border: 2px solid var(--neon-pink);
            transform: rotateY(180deg);
            box-shadow: 0 0 30px rgba(255, 45, 149, 0.4);
        }

        .card.matched .card-front {
            border-color: var(--neon-yellow);
            box-shadow: 0 0 30px rgba(245, 255, 0, 0.5);
            animation: matchPulse 0.5s ease-out;
        }

        @keyframes matchPulse {
            0% { transform: rotateY(180deg) scale(1); }
            50% { transform: rotateY(180deg) scale(1.1); }
            100% { transform: rotateY(180deg) scale(1); }
        }

        .card:hover:not(.flipped):not(.matched) .card-back {
            border-color: var(--neon-pink);
            box-shadow: 0 0 30px rgba(255, 45, 149, 0.4);
        }

        .btn {
            font-family: 'Orbitron', sans-serif;
            font-size: 1rem;
            font-weight: 700;
            padding: 15px 40px;
            border: none;
            border-radius: 8px;
            cursor: pointer;
            text-transform: uppercase;
            letter-spacing: 3px;
            transition: all 0.3s ease;
            position: relative;
            overflow: hidden;
        }

        .btn-primary {
            background: linear-gradient(135deg, var(--neon-pink), var(--neon-purple));
            color: white;
            box-shadow: 0 0 30px rgba(255, 45, 149, 0.4);
        }

        .btn-primary:hover {
            transform: translateY(-3px);
            box-shadow: 0 0 50px rgba(255, 45, 149, 0.6);
        }

        .btn-secondary {
            background: transparent;
            color: var(--neon-cyan);
            border: 2px solid var(--neon-cyan);
            box-shadow: 0 0 20px rgba(0, 245, 255, 0.2);
        }

        .btn-secondary:hover {
            background: rgba(0, 245, 255, 0.1);
            box-shadow: 0 0 40px rgba(0, 245, 255, 0.4);
        }

        .controls {
            display: flex;
            justify-content: center;
            gap: 20px;
            margin-top: 30px;
            flex-wrap: wrap;
        }

        /* Modal */
        .modal {
            display: none;
            position: fixed;
            top: 0;
            left: 0;
            width: 100%;
            height: 100%;
            background: rgba(0, 0, 0, 0.9);
            z-index: 100;
            align-items: center;
            justify-content: center;
            animation: fadeIn 0.3s ease;
        }

        .modal.active {
            display: flex;
        }

        @keyframes fadeIn {
            from { opacity: 0; }
            to { opacity: 1; }
        }

        .modal-content {
            background: linear-gradient(135deg, #12121a, #1a1a2e);
            padding: 50px;
            border-radius: 20px;
            text-align: center;
            border: 2px solid var(--neon-cyan);
            box-shadow: 0 0 100px rgba(0, 245, 255, 0.3);
            animation: modalSlide 0.4s ease;
            max-width: 90%;
        }

        @keyframes modalSlide {
            from { transform: scale(0.8) translateY(50px); opacity: 0; }
            to { transform: scale(1) translateY(0); opacity: 1; }
        }

        .modal h2 {
            font-family: 'Orbitron', sans-serif;
            font-size: 2.5rem;
            margin-bottom: 20px;
            background: linear-gradient(135deg, var(--neon-yellow), var(--neon-cyan));
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
            background-clip: text;
        }

        .modal-stats {
            display: flex;
            justify-content: center;
            gap: 30px;
            margin: 30px 0;
        }

        .leaderboard {
            margin-top: 40px;
            padding: 20px;
            background: rgba(0, 0, 0, 0.3);
            border-radius: 15px;
            border: 1px solid rgba(184, 41, 221, 0.3);
        }

        .leaderboard h3 {
            font-family: 'Orbitron', sans-serif;
            color: var(--neon-purple);
            margin-bottom: 15px;
            font-size: 1.2rem;
            letter-spacing: 3px;
        }

        .leaderboard-list {
            list-style: none;
        }

        .leaderboard-item {
            display: flex;
            justify-content: space-between;
            padding: 10px 15px;
            border-bottom: 1px solid rgba(255, 255, 255, 0.1);
            font-size: 0.95rem;
        }

        .leaderboard-item:last-child {
            border-bottom: none;
        }

        .rank {
            color: var(--neon-yellow);
            font-weight: 700;
            width: 30px;
        }

        .player-name {
            flex: 1;
            text-align: left;
            color: var(--neon-cyan);
        }

        .player-score {
            color: var(--neon-pink);
        }

        /* Start screen */
        .start-screen {
            text-align: center;
            padding: 60px 20px;
        }

        .start-screen input {
            font-family: 'Rajdhani', sans-serif;
            font-size: 1.2rem;
            padding: 15px 25px;
            border: 2px solid var(--neon-cyan);
            border-radius: 10px;
            background: rgba(0, 0, 0, 0.5);
            color: white;
            text-align: center;
            width: 100%;
            max-width: 300px;
            margin: 20px 0;
            transition: all 0.3s ease;
        }

        .start-screen input:focus {
            outline: none;
            box-shadow: 0 0 30px rgba(0, 245, 255, 0.4);
        }

        .start-screen input::placeholder {
            color: rgba(255, 255, 255, 0.4);
        }

        .difficulty-select {
            display: flex;
            justify-content: center;
            gap: 15px;
            margin: 25px 0;
            flex-wrap: wrap;
        }

        .difficulty-btn {
            padding: 12px 25px;
            border: 2px solid rgba(255, 255, 255, 0.2);
            border-radius: 8px;
            background: transparent;
            color: #888;
            cursor: pointer;
            transition: all 0.3s ease;
            font-family: 'Rajdhani', sans-serif;
            font-size: 1rem;
            font-weight: 600;
            letter-spacing: 1px;
        }

        .difficulty-btn:hover {
            border-color: var(--neon-cyan);
            color: var(--neon-cyan);
        }

        .difficulty-btn.active {
            border-color: var(--neon-pink);
            color: var(--neon-pink);
            box-shadow: 0 0 20px rgba(255, 45, 149, 0.3);
        }

        .game-container {
            display: none;
        }

        .game-container.active {
            display: block;
        }

        footer {
            text-align: center;
            padding: 40px;
            color: #555;
            font-size: 0.9rem;
        }

        footer a {
            color: var(--neon-cyan);
            text-decoration: none;
        }

        @media (max-width: 600px) {
            h1 { font-size: 2rem; letter-spacing: 4px; }
            .game-board { gap: 10px; padding: 20px; }
            .card-front, .card-back { font-size: 1.8rem; }
            .stats-bar { gap: 15px; }
            .stat { padding: 10px 20px; min-width: 100px; }
            .stat-value { font-size: 1.5rem; }
        }
    </style>
</head>
<body>
    <div class="bg-grid"></div>
    <div class="bg-glow glow-1"></div>
    <div class="bg-glow glow-2"></div>

    <div class="container">
        <header>
            <h1>Memory Match</h1>
            <p class="subtitle">Neon Edition</p>
        </header>

        <!-- Start Screen -->
        <div class="start-screen" id="startScreen">
            <input type="text" id="playerName" placeholder="Enter your name" maxlength="15">
            
            <p style="color: #888; margin-top: 20px; letter-spacing: 2px;">SELECT DIFFICULTY</p>
            <div class="difficulty-select">
                <button class="difficulty-btn active" data-pairs="6">EASY (6)</button>
                <button class="difficulty-btn" data-pairs="8">MEDIUM (8)</button>
                <button class="difficulty-btn" data-pairs="10">HARD (10)</button>
            </div>

            <button class="btn btn-primary" onclick="startGame()">START GAME</button>

            <div class="leaderboard" id="startLeaderboard">
                <h3>🏆 TOP PLAYERS</h3>
                <ul class="leaderboard-list" id="leaderboardList">
                    <li class="leaderboard-item" style="color: #555;">No scores yet. Be the first!</li>
                </ul>
            </div>
        </div>

        <!-- Game Container -->
        <div class="game-container" id="gameContainer">
            <div class="stats-bar">
                <div class="stat">
                    <div class="stat-value" id="movesCount">0</div>
                    <div class="stat-label">Moves</div>
                </div>
                <div class="stat">
                    <div class="stat-value" id="timerDisplay">0:00</div>
                    <div class="stat-label">Time</div>
                </div>
                <div class="stat">
                    <div class="stat-value" id="matchesCount">0</div>
                    <div class="stat-label">Matches</div>
                </div>
            </div>

            <div class="game-board" id="gameBoard"></div>

            <div class="controls">
                <button class="btn btn-secondary" onclick="restartGame()">RESTART</button>
                <button class="btn btn-secondary" onclick="goToMenu()">MENU</button>
            </div>
        </div>
    </div>

    <!-- Win Modal -->
    <div class="modal" id="winModal">
        <div class="modal-content">
            <h2>🎉 Victory!</h2>
            <p style="color: #aaa; font-size: 1.1rem;">You've matched all the cards!</p>
            
            <div class="modal-stats">
                <div class="stat">
                    <div class="stat-value" id="finalMoves">0</div>
                    <div class="stat-label">Moves</div>
                </div>
                <div class="stat">
                    <div class="stat-value" id="finalTime">0:00</div>
                    <div class="stat-label">Time</div>
                </div>
            </div>

            <div style="display: flex; gap: 15px; justify-content: center; flex-wrap: wrap;">
                <button class="btn btn-primary" onclick="restartGame()">PLAY AGAIN</button>
                <button class="btn btn-secondary" onclick="goToMenu()">MENU</button>
            </div>
        </div>
    </div>

    <footer>
        Built with 💜 using <a href="https://go.dev" target="_blank">Go</a>
    </footer>

    <script>
        const emojis = ['🚀', '⚡', '🔥', '💎', '🎯', '🎮', '👾', '🤖', '🛸', '🌟', '💫', '🎪'];
        
        let cards = [];
        let flippedCards = [];
        let matchedPairs = 0;
        let moves = 0;
        let timer = null;
        let seconds = 0;
        let totalPairs = 8;
        let playerName = 'Player';
        let gameStarted = false;

        // Difficulty selection
        document.querySelectorAll('.difficulty-btn').forEach(btn => {
            btn.addEventListener('click', () => {
                document.querySelectorAll('.difficulty-btn').forEach(b => b.classList.remove('active'));
                btn.classList.add('active');
                totalPairs = parseInt(btn.dataset.pairs);
            });
        });

        function shuffle(array) {
            for (let i = array.length - 1; i > 0; i--) {
                const j = Math.floor(Math.random() * (i + 1));
                [array[i], array[j]] = [array[j], array[i]];
            }
            return array;
        }

        function createBoard() {
            const board = document.getElementById('gameBoard');
            board.innerHTML = '';
            
            // Adjust grid based on pairs
            const cols = totalPairs <= 6 ? 3 : 4;
            board.style.gridTemplateColumns = ` + "`repeat(${cols}, 1fr)`" + `;
            
            const gameEmojis = shuffle([...emojis]).slice(0, totalPairs);
            cards = shuffle([...gameEmojis, ...gameEmojis]);

            cards.forEach((emoji, index) => {
                const card = document.createElement('div');
                card.className = 'card';
                card.innerHTML = ` + "`" + `
                    <div class="card-inner">
                        <div class="card-back"></div>
                        <div class="card-front">${emoji}</div>
                    </div>
                ` + "`" + `;
                card.addEventListener('click', () => flipCard(card, emoji, index));
                board.appendChild(card);
            });
        }

        function flipCard(card, emoji, index) {
            if (!gameStarted) {
                startTimer();
                gameStarted = true;
            }

            if (flippedCards.length >= 2 || card.classList.contains('flipped') || card.classList.contains('matched')) {
                return;
            }

            card.classList.add('flipped');
            flippedCards.push({ card, emoji, index });

            if (flippedCards.length === 2) {
                moves++;
                document.getElementById('movesCount').textContent = moves;

                if (flippedCards[0].emoji === flippedCards[1].emoji) {
                    // Match!
                    setTimeout(() => {
                        flippedCards.forEach(fc => fc.card.classList.add('matched'));
                        matchedPairs++;
                        document.getElementById('matchesCount').textContent = matchedPairs;
                        flippedCards = [];

                        if (matchedPairs === totalPairs) {
                            endGame();
                        }
                    }, 300);
                } else {
                    // No match
                    setTimeout(() => {
                        flippedCards.forEach(fc => fc.card.classList.remove('flipped'));
                        flippedCards = [];
                    }, 1000);
                }
            }
        }

        function startTimer() {
            timer = setInterval(() => {
                seconds++;
                const mins = Math.floor(seconds / 60);
                const secs = seconds % 60;
                document.getElementById('timerDisplay').textContent = ` + "`${mins}:${secs.toString().padStart(2, '0')}`" + `;
            }, 1000);
        }

        function stopTimer() {
            clearInterval(timer);
            timer = null;
        }

        function endGame() {
            stopTimer();
            
            const finalMins = Math.floor(seconds / 60);
            const finalSecs = seconds % 60;
            const timeStr = ` + "`${finalMins}:${finalSecs.toString().padStart(2, '0')}`" + `;
            
            document.getElementById('finalMoves').textContent = moves;
            document.getElementById('finalTime').textContent = timeStr;
            document.getElementById('winModal').classList.add('active');

            // Submit score
            submitScore();
        }

        async function submitScore() {
            try {
                await fetch('/api/score', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({
                        playerName: playerName,
                        moves: moves,
                        timeTaken: seconds
                    })
                });
                loadLeaderboard();
            } catch (e) {
                console.error('Failed to submit score:', e);
            }
        }

        async function loadLeaderboard() {
            try {
                const res = await fetch('/api/leaderboard');
                const data = await res.json();
                const list = document.getElementById('leaderboardList');
                
                if (data && data.length > 0) {
                    list.innerHTML = data.slice(0, 5).map((score, i) => ` + "`" + `
                        <li class="leaderboard-item">
                            <span class="rank">#${i + 1}</span>
                            <span class="player-name">${score.playerName}</span>
                            <span class="player-score">${score.moves} moves</span>
                        </li>
                    ` + "`" + `).join('');
                }
            } catch (e) {
                console.error('Failed to load leaderboard:', e);
            }
        }

        function startGame() {
            const nameInput = document.getElementById('playerName');
            playerName = nameInput.value.trim() || 'Player';
            
            document.getElementById('startScreen').style.display = 'none';
            document.getElementById('gameContainer').classList.add('active');
            
            resetGame();
            createBoard();
        }

        function resetGame() {
            stopTimer();
            flippedCards = [];
            matchedPairs = 0;
            moves = 0;
            seconds = 0;
            gameStarted = false;
            
            document.getElementById('movesCount').textContent = '0';
            document.getElementById('timerDisplay').textContent = '0:00';
            document.getElementById('matchesCount').textContent = '0';
            document.getElementById('winModal').classList.remove('active');
        }

        function restartGame() {
            resetGame();
            createBoard();
        }

        function goToMenu() {
            resetGame();
            document.getElementById('gameContainer').classList.remove('active');
            document.getElementById('startScreen').style.display = 'block';
            document.getElementById('winModal').classList.remove('active');
            loadLeaderboard();
        }

        // Load leaderboard on page load
        loadLeaderboard();
    </script>
</body>
</html>`
	w.Write([]byte(html))
}

func handleLeaderboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	mu.RLock()
	defer mu.RUnlock()

	json.NewEncoder(w).Encode(leaderboard)
}

func handleScore(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var score GameScore
	if err := json.NewDecoder(r.Body).Decode(&score); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	score.Timestamp = time.Now()

	mu.Lock()
	leaderboard = append(leaderboard, score)

	// Sort by moves (ascending), then by time (ascending)
	for i := 0; i < len(leaderboard)-1; i++ {
		for j := i + 1; j < len(leaderboard); j++ {
			if leaderboard[j].Moves < leaderboard[i].Moves ||
				(leaderboard[j].Moves == leaderboard[i].Moves && leaderboard[j].TimeTaken < leaderboard[i].TimeTaken) {
				leaderboard[i], leaderboard[j] = leaderboard[j], leaderboard[i]
			}
		}
	}

	// Keep only top 10
	if len(leaderboard) > 10 {
		leaderboard = leaderboard[:10]
	}
	mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}
