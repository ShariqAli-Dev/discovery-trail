const unitChapters = document.getElementsByClassName("chapter-container");
validateChapters();

document.body.addEventListener("htmx:afterSwap", function (event) {
  validateChapters();
});

function validateChapters() {
  let allChapterQuestionsGenerated = true;
  for (let i = 0; i < unitChapters?.length; i++) {
    const chapter = unitChapters[i];
    for (let j = 0; j < chapter.children.length; j++) {
      const question = chapter.children[j];
      if (!question.classList.contains("completed")) {
        allChapterQuestionsGenerated = false;
      }
    }
  }

  if (allChapterQuestionsGenerated) {
    /** @type {?HTMLButtonElement}*/
    //@ts-ignore
    const generateCourseButton = document.getElementById(
      "button-generate-course"
    );
    if (generateCourseButton) {
      generateCourseButton.disabled = false;
    }
  }
}
