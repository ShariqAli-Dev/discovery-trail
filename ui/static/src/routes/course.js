document.addEventListener("DOMContentLoaded", () => {
  const chapterID = localStorage.getItem("course-chapter-id");
  const unitID = localStorage.getItem("course-unit-id");
  const courseSection = document.querySelector("#course-chapter-container");

  const hxHeader = JSON.parse(
    courseSection?.getAttribute("hx-headers") || "{}"
  );
  hxHeader["chapter-id"] = chapterID;
  hxHeader["unit-id"] = unitID;
  courseSection?.setAttribute("hx-headers", JSON.stringify(hxHeader));

  // @ts-expect-error
  htmx.trigger("#course-chapter-container", "custom-chapter");
  console.log("made it to the trigger");
});

const menuChapters = document.querySelectorAll(".menu-chapter");
menuChapters.forEach((chapter) => {
  if (chapter instanceof HTMLElement) {
    chapter.onclick = (event) => {
      const chapterId = chapter.id.split("menu-chapter-")[1];
      setAllOtherChaptersToHidden(chapterId);
    };
  }
  return;
});

/**
 * Sets all chapters except the specified one to hidden.
 * @param {string} chapterID - The ID of the chapter to remain visible.
 */
function setAllOtherChaptersToHidden(chapterID) {
  const courseSection = document.querySelector("#course-chapter-container");
  const unitID = localStorage.getItem("course-unit-id");

  const hxHeader = JSON.parse(
    courseSection?.getAttribute("hx-headers") || "{}"
  );
  hxHeader["chapter-id"] = chapterID;
  hxHeader["unit-id"] = unitID;
  courseSection?.setAttribute("hx-headers", JSON.stringify(hxHeader));

  // @ts-expect-error
  htmx.trigger("#course-chapter-container", "custom-chapter");
}
