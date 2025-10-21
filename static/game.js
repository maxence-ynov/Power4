// 🎮 Power4 - Animation réelle des jetons
document.addEventListener("DOMContentLoaded", () => {
  const boardForm = document.querySelector("form[action='/play']");
  if (!boardForm) return;

  const board = document.querySelector(".board");

  const columnCount = board.dataset.columns || "7";
  board.style.setProperty("--column-count", columnCount);
  board.style.gridTemplateColumns = `repeat(${columnCount}, 1fr)`;

  let currentBoardState = captureBoardState();

  attachColumnListeners();

  function attachColumnListeners() {
    const columns = board.querySelectorAll(".col");
    columns.forEach((col) => {
      col.addEventListener("click", handleColumnClick, { once: true });
    });
  }

  async function handleColumnClick(e) {
    e.preventDefault();
    const colButton = e.currentTarget;
    const columnIndex = colButton.value;

    // Envoi du coup au serveur
    await fetch("/play", {
      method: "POST",
      headers: { "Content-Type": "application/x-www-form-urlencoded" },
      body: `column=${columnIndex}`,
    });

    // Récupère le nouveau plateau HTML
    const response = await fetch("/");
    const html = await response.text();

    // Parse le nouvel HTML pour extraire le plateau mis à jour
    const parser = new DOMParser();
    const doc = parser.parseFromString(html, "text/html");
    const newBoard = doc.querySelector(".board");

    if (!newBoard) {
      attachColumnListeners();
      return;
    }

    const newState = captureBoardState(newBoard);

    // Trouve la différence (le nouveau jeton)
    const diff = findBoardDifference(currentBoardState, newState);
    if (diff) {
      const { colIndex, rowIndex, player } = diff;
      animateTokenFall(colIndex, rowIndex, player);
    }

    // Met à jour le plateau après animation
    setTimeout(() => {
      const newColumnCount = newBoard.dataset.columns || board.dataset.columns || "7";
      board.dataset.columns = newColumnCount;
      board.style.gridTemplateColumns = `repeat(${newColumnCount}, 1fr)`;
      board.innerHTML = newBoard.innerHTML;
      currentBoardState = newState;
      attachColumnListeners();
    }, 500);
  }
});

// Capture le plateau sous forme de tableau de valeurs 0/1/2
function captureBoardState(targetBoard) {
  const board = targetBoard || document.querySelector(".board");
  const state = [];
  const columns = board.querySelectorAll(".col");

  columns.forEach((col) => {
    const cells = Array.from(col.querySelectorAll(".cell"));
    state.push(
      cells.map((cell) => {
        if (cell.classList.contains("player1")) return 1;
        if (cell.classList.contains("player2")) return 2;
        return 0;
      }),
    );
  });
  return state;
}

// Trouve la différence entre deux plateaux
function findBoardDifference(oldBoard, newBoard) {
  for (let c = 0; c < newBoard.length; c++) {
    for (let r = 0; r < newBoard[c].length; r++) {
      if (oldBoard[c][r] !== newBoard[c][r]) {
        const player = newBoard[c][r];
        if (player === 1 || player === 2) {
          return { colIndex: c, rowIndex: r, player };
        }
      }
    }
  }
  return null;
}

// Anime la chute du jeton
function animateTokenFall(colIndex, rowIndex, player) {
  const board = document.querySelector(".board");
  const col = board.querySelectorAll(".col")[colIndex];
  const cells = col.querySelectorAll(".cell");
  const targetCell = cells[cells.length - 1 - rowIndex];

  if (!targetCell) return;

  // Crée un pion animé temporaire
  const token = document.createElement("div");
  token.classList.add("cell", player === 1 ? "player1" : "player2");
  token.style.position = "absolute";
  token.style.left = `${col.offsetLeft + targetCell.offsetLeft}px`;
  token.style.top = `${col.offsetTop - 400}px`;
  token.style.transition = "transform 0.5s cubic-bezier(0.34,1.56,0.64,1)";
  token.style.zIndex = "1000";
  token.style.transform = "translateY(0px)";

  document.body.appendChild(token);

  // Calcul de la chute
  const targetY =
    col.offsetTop +
    targetCell.getBoundingClientRect().top -
    col.getBoundingClientRect().top;

  requestAnimationFrame(() => {
    token.style.transform = `translateY(${targetY}px)`;
  });

  // Nettoyage après animation
  setTimeout(() => {
    token.remove();
  }, 600);
}
