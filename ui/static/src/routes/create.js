/** @type {?HTMLButtonElement}*/
//@ts-ignore
const buttonAddUnit = document.getElementById("button-add-unit");
/** @type {?HTMLButtonElement}*/
//@ts-ignore
const buttonRemoveUnit = document.getElementById("button-remove-unit");

// ADDING AND REMOVING UNITS
buttonAddUnit?.addEventListener("click", () => {
  const divUnits = document.getElementById("div-units");
  if (!divUnits) return;
  const newUnitCount = divUnits.childElementCount + 1;
  if (newUnitCount > 5) return;

  const newUnitDiv = document.createElement("div");
  newUnitDiv.className =
    "flex flex-col items-start w-full mt-4 sm:flex-row sm:items-center";

  const label = document.createElement("label");
  label.htmlFor = `unit-${newUnitCount}`;
  label.className = "flex-[1] text-xl label";
  label.textContent = `Unit ${newUnitCount}`;

  const input = document.createElement("input");
  input.id = `unit-${newUnitCount}`;
  input.name = `unit-${newUnitCount}`;
  input.className = "flex-[6] text-lg input input-bordered";
  input.placeholder = "Enter subtopic of the course (eg. Why learn htmx?)";

  newUnitDiv.appendChild(label);
  newUnitDiv.appendChild(input);

  divUnits?.appendChild(newUnitDiv);
  updateUnitCount();
});

buttonRemoveUnit?.addEventListener("click", () => {
  const divUnits = document.getElementById("div-units");
  if (divUnits?.lastChild && divUnits.childElementCount > 1) {
    divUnits.removeChild(divUnits.lastChild);
    updateUnitCount();
  }
});

function updateUnitCount() {
  const divUnits = document.getElementById("div-units");
  /** @type {?HTMLInputElement}*/
  //@ts-ignore
  const inputUnitCount = document.getElementById("input-unit-count");
  if (inputUnitCount) {
    const unitCount = divUnits?.childElementCount;
    inputUnitCount.value = `${unitCount}`;
  }
}
