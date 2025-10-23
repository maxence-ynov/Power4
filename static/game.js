document.addEventListener("DOMContentLoaded", () => {
  const boardForm = document.querySelector("form[action='/play']");
  const boardGrid = document.querySelector(".board-grid");
  if (!boardForm || !boardGrid) {
    return;
  }
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
      const lastCol = parseInt(nextGrid.getAttribute("data-last-col"), 10);
      const lastRow = parseInt(nextGrid.getAttribute("data-last-row"), 10);
      syncBoardAttributes(boardGrid, nextGrid);
      boardGrid.innerHTML = nextGrid.innerHTML;
      highlightExactCell(lastCol, lastRow);
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

  function highlightExactCell(colIndex, rowIndex) {
    if (Number.isNaN(colIndex) || Number.isNaN(rowIndex) || colIndex < 0 || rowIndex < 0) {
      return;
    }
    const columns = boardGrid.querySelectorAll(".board-column");
    const targetColumn = columns[colIndex];
    if (!targetColumn) {
      return;
    }
    const cells = targetColumn.querySelectorAll(".cell");
    const targetCell = cells[rowIndex];
    if (!targetCell) {
      return;
    }
    boardGrid.querySelectorAll(".cell.last").forEach((cell) => cell.classList.remove("last"));
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
    const lastCol = source.getAttribute("data-last-col");
    if (lastCol !== null) {
      target.setAttribute("data-last-col", lastCol);
    } else {
      target.removeAttribute("data-last-col");
    }
    const lastRow = source.getAttribute("data-last-row");
    if (lastRow !== null) {
      target.setAttribute("data-last-row", lastRow);
    } else {
      target.removeAttribute("data-last-row");
    }
  }
});

function updateSection(selector, doc) {
  const next = doc.querySelector(selector);
  const current = document.querySelector(selector);
  if (next && current) {
    current.innerHTML = next.innerHTML;
  }
}
