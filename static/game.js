document.addEventListener("DOMContentLoaded", () => {
  const boardForm = document.querySelector("form[action='/play']");
  const boardGrid = document.querySelector(".board-grid");

  if (!boardForm || !boardGrid) {
    return;
  }

  let currentBoardState = captureBoardState(boardGrid);
  let isAnimating = false;

  attachColumnListeners();

  function attachColumnListeners() {
    const columns = boardGrid.querySelectorAll(".board-column");
    columns.forEach((column) => {
      column.addEventListener("click", handleColumnClick);
    });
  }

  async function handleColumnClick(event) {
    const columnButton = event.currentTarget;

    if (columnButton.disabled) {
      return;
    }

    if (isAnimating) {
      event.preventDefault();
      return;
    }

    event.preventDefault();

    const columnIndex = Number(columnButton.value);
    if (Number.isNaN(columnIndex)) {
      boardForm.submit();
      return;
    }

    isAnimating = true;

    try {
      const playResponse = await fetch("/play", {
        method: "POST",
        headers: {
          "Content-Type": "application/x-www-form-urlencoded",
        },
        body: `column=${columnIndex}`,
      });

      if (!playResponse.ok) {
        throw new Error("move_failed");
      }

      const refreshResponse = await fetch("/", {
        headers: {
          "X-Requested-With": "fetch",
        },
      });

      if (!refreshResponse.ok) {
        throw new Error("refresh_failed");
      }

      const html = await refreshResponse.text();
      const parser = new DOMParser();
      const doc = parser.parseFromString(html, "text/html");
      const nextGrid = doc.querySelector(".board-grid");

      if (!nextGrid) {
        throw new Error("grid_missing");
      }

      const nextState = captureBoardState(nextGrid);
      const diff = findBoardDifference(currentBoardState, nextState);

      syncBoardAttributes(boardGrid, nextGrid);
      boardGrid.innerHTML = nextGrid.innerHTML;
      currentBoardState = nextState;

      highlightLastMove(diff);
      updateSection(".status-card", doc);
      updateSection(".legend", doc);
      updateSection(".game-actions", doc);

      attachColumnListeners();
    } catch (error) {
      boardForm.submit();
    } finally {
      isAnimating = false;
    }
  }

  function highlightLastMove(diff) {
    if (!diff) {
      return;
    }

    const columns = boardGrid.querySelectorAll(".board-column");
    const targetColumn = columns[diff.colIndex];
    if (!targetColumn) {
      return;
    }

    const cells = targetColumn.querySelectorAll(".cell");
    const targetCell = cells[diff.rowIndex];
    if (!targetCell) {
      return;
    }

    boardGrid.querySelectorAll(".cell.last").forEach((cell) => {
      cell.classList.remove("last");
    });

    requestAnimationFrame(() => {
      targetCell.classList.add("drop-in", "last");
      setTimeout(() => {
        targetCell.classList.remove("drop-in");
      }, 650);
    });
  }

  function syncBoardAttributes(target, source) {
    const dataColumns = source.getAttribute("data-columns");
    if (dataColumns !== null) {
      target.setAttribute("data-columns", dataColumns);
    } else {
      target.removeAttribute("data-columns");
    }

    const columnCount = source.style.getPropertyValue("--column-count");
    if (columnCount) {
      target.style.setProperty("--column-count", columnCount);
    }
  }
});

function captureBoardState(gridElement) {
  const state = [];
  const columns = gridElement.querySelectorAll(".board-column");

  columns.forEach((column) => {
    const columnCells = [];
    column.querySelectorAll(".cell").forEach((cell) => {
      if (cell.classList.contains("player1")) {
        columnCells.push(1);
      } else if (cell.classList.contains("player2")) {
        columnCells.push(2);
      } else {
        columnCells.push(0);
      }
    });
    state.push(columnCells);
  });

  return state;
}

function findBoardDifference(previous, next) {
  for (let column = 0; column < next.length; column += 1) {
    for (let row = 0; row < next[column].length; row += 1) {
      if (previous[column][row] !== next[column][row]) {
        const value = next[column][row];
        if (value === 1 || value === 2) {
          return { colIndex: column, rowIndex: row, player: value };
        }
      }
    }
  }
  return null;
}

function updateSection(selector, doc) {
  const next = doc.querySelector(selector);
  const current = document.querySelector(selector);

  if (next && current) {
    current.innerHTML = next.innerHTML;
  }
}
