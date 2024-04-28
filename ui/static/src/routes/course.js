document.addEventListener("DOMContentLoaded", () => {
  // on page load, load the first element
  const courseSection = document.querySelector("#course-chapter-container");
  const firstMenuChapter = document.querySelector(".menu-chapter");
  const chapterID = firstMenuChapter?.id.split("menu-chapter-")[1];

  const hxHeader = JSON.parse(
    courseSection?.getAttribute("hx-headers") || "{}"
  );
  hxHeader["chapter-id"] = chapterID;
  hxHeader["chapter-index"] = 1;
  hxHeader["unit-index"] = 1;
  courseSection?.setAttribute("hx-headers", JSON.stringify(hxHeader));
  // @ts-expect-error
  htmx.trigger("#course-chapter-container", "custom-chapter");
});

const menuChapters = document.querySelectorAll(".menu-chapter");
menuChapters.forEach((chapter, cdx) => {
  if (chapter instanceof HTMLElement) {
    chapter.onclick = (event) => {
      // get the elements position
      const chapterParentElement = chapter.parentNode;
      if (!chapterParentElement) return;
      const chapterID = chapter.id.split("menu-chapter-")[1];
      const chapterIndex = getElementIndexOf(
        chapterParentElement.children,
        chapter
      );
      // disable if already active
      const activeChapterElement = document.querySelector(".active-chapter");
      if (chapterID === activeChapterElement?.id) return;

      const unitElement = chapter.parentElement?.parentElement;
      const unitElementParent = unitElement?.parentNode;
      if (!unitElementParent) return;
      const unitIndex = getElementIndexOf(
        unitElementParent.children,
        unitElement
      );

      // trigger the htmx call
      const courseSection = document.querySelector("#course-chapter-container");
      const hxHeader = JSON.parse(
        courseSection?.getAttribute("hx-headers") || "{}"
      );
      hxHeader["chapter-id"] = chapterID;
      hxHeader["chapter-index"] = chapterIndex + 1;
      hxHeader["unit-index"] = unitIndex + 1;
      courseSection?.setAttribute("hx-headers", JSON.stringify(hxHeader));
      // @ts-expect-error
      htmx.trigger("#course-chapter-container", "custom-chapter");

      // set all other elements to hidden
    };
  }
});

/**
 * Gets the index of the search element
 * @param {HTMLCollection} parentElementChildren - The parent elements children
 * @param {Element} searchElement - The search element
 */
function getElementIndexOf(parentElementChildren, searchElement) {
  return Array.prototype.indexOf.call(parentElementChildren, searchElement);
}
