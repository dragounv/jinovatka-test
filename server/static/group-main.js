// @ts-nocheck
// Script that hanldes button clicks on the group view.
// It mainly copies stuff to clipboard
document.getElementById("copy-urls").addEventListener("click", copyColumn);
document.getElementById("copy-ids").addEventListener("click", copyColumn);

const groupLink = document.getElementById("group-link");
document
  .getElementById("copy-group-link")
  .addEventListener("click", copyFrom(groupLink));

// Event handler.
// Copy to clipboard column from the table containig group info.
function copyColumn(e) {
  let targetColumn;
  if (this.id === "copy-urls") {
    targetColumn = 1;
  } else if (this.id === "copy-ids") {
    targetColumn = 2;
  } else {
    console.error("copyColumn was called on wrong element!");
    return;
  }
  const dataCells = Array.from(
    document.querySelectorAll(
      `#group-info-table > tbody > tr > td:nth-of-type(${targetColumn}) > a`
    )
  );
  const baseURL = window.location.href.origin;
  const columnData = dataCells.reduce(
    (accumulator, currentValue) =>
      accumulator + "\n" + new URL(currentValue.href, baseURL).toString()
  );
  const result = navigator.clipboard.writeText(columnData);
  result.catch((reason) => console.error(reason));
  result.then(() => showCopied(this));
}

// Closure returning event hanlder.
// Copy url of single anchor.
function copyFrom(target) {
  function handler(e) {
    const data = new URL(target.href, window.location.href.origin).toString();
    const result = navigator.clipboard.writeText(data);
    result.catch((reason) => console.error(reason));
    result.then(() => showCopied(this));
  }
  return handler;
}

function showCopied(target) {
  const tmp = target.textContent;
  target.textContent = "Zkopírováno";
  setTimeout(() => (target.textContent = tmp), 1000);
}
