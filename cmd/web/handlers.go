package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/markbates/goth/gothic"
	"github.com/shariqali-dev/discovery-trail/internal/gpt"
	"github.com/shariqali-dev/discovery-trail/internal/models"
	"github.com/shariqali-dev/discovery-trail/internal/types"
	"github.com/shariqali-dev/discovery-trail/internal/validator"
	"github.com/shariqali-dev/discovery-trail/ui/html/pages"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if app.isAuthenticated(r) {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	data, err := app.newTemplateData(r)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	homePage := pages.Home(data)
	app.render(w, r, http.StatusOK, homePage)
}

func (app *application) dashboard(w http.ResponseWriter, r *http.Request) {
	accountID, err := app.requestGetAccountID(r)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	courses, err := app.courses.All(accountID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	coursesWithUnitTitles := make([]types.CourseWithUnitTitles, len(courses))
	for cdx, course := range courses {
		unitTitles, err := app.units.GetCourseUnitTitles(course.ID)
		if err != nil {
			app.serverError(w, r, err)
		}
		coursesWithUnitTitles[cdx] = types.CourseWithUnitTitles{
			Course:     course,
			UnitTitles: unitTitles,
		}
	}

	data, err := app.newTemplateData(r)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	data.Courses = coursesWithUnitTitles
	dashboardPage := pages.Dashboard(data)
	app.render(w, r, http.StatusOK, dashboardPage)
}

func (app *application) create(w http.ResponseWriter, r *http.Request) {
	data, err := app.newTemplateData(r)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	createPage := pages.Create(data)
	app.render(w, r, http.StatusOK, createPage)
}

func (app *application) createCourse(w http.ResponseWriter, r *http.Request) {
	courseID := r.PathValue("courseID")
	if courseID == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	course, err := app.courses.Get(courseID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	if course.Processed {
		http.Redirect(w, r, fmt.Sprintf("/course/%s", courseID), http.StatusSeeOther)
		return
	}
	data, err := app.newTemplateData(r)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	courseUnits, err := app.units.GetCourseUnits(course.ID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	courseUnitsChapters := types.CourseWithUnitsWithChapters{
		Course: course,
		Units:  make([]types.UnitWithChapters, 0),
	}
	for _, unit := range courseUnits {
		unitChapters, err := app.chapters.GetUnitChapters(unit.ID)
		if err != nil {
			app.serverError(w, r, err)
			return
		}
		courseUnitsChapters.Units = append(courseUnitsChapters.Units, types.UnitWithChapters{
			Chapters: unitChapters,
			Unit:     unit,
		})
	}
	data.CourseUnitsChapters = courseUnitsChapters
	createCoursePage := pages.CreateCourse(data)
	app.render(w, r, http.StatusOK, createCoursePage)
}

func (app *application) createPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	var form types.CourseCreateForm
	err = app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	unitCountAsInt, err := strconv.Atoi(r.PostForm.Get("unit-count"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
	}
	form.UnitCount = unitCountAsInt
	form.UnitValues = make(map[string]string)
	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MinCount(form.UnitCount, 1), "unit-count", "Cant have less than 1 unit")
	form.CheckField(validator.MaxCount(form.UnitCount, 5), "unit-count", "Cant have more than 5 unit")
	for udx := range form.UnitCount {
		udx++
		unitString := fmt.Sprintf("unit-%d", udx)
		unitFormValue := r.PostForm.Get(unitString)
		form.UnitValues[unitString] = unitFormValue

		form.CheckField(validator.NotBlank(unitFormValue), unitString, "This field cannot be blank")
		form.CheckField(validator.MaxChars(unitFormValue, 40), unitString, "This field cannot be more than 40 characters long")
	}
	accountID, err := app.requestGetAccountID(r)
	if err != nil || accountID == "" {
		app.serverError(w, r, err)
		return
	}
	credits, err := app.accounts.GetCredits(accountID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	form.CheckField(validator.MinCount(credits, 1), "credits", "Invalid credits")
	// ** form validation completed **
	if !form.Valid() {
		courseCreateFormInputs := pages.CourseCreateFormInputs(form)
		app.render(w, r, http.StatusOK, courseCreateFormInputs)
		return
	}

	unsplashSearchTerm, err := gpt.GetImageSearchTermFromTitle(app.openAiClient, form.Title)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	unsplashResult, err := unsplashGetImage(strings.TrimSpace(unsplashSearchTerm.SearchTerm))
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	generatedCourseInformation, err := gpt.GenerateCourseTitleAndUnitChapters(app.openAiClient, form.Title, form.UnitValues)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	courseID, err := app.courses.Insert(generatedCourseInformation.Title, unsplashResult.Results[0].Images.Small, accountID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	for _, unit := range generatedCourseInformation.Units {
		unitID, err := app.units.Insert(unit.Title, courseID)
		if err != nil {
			app.serverError(w, r, err)
			return
		}
		for _, chapter := range unit.Chapters {
			youtubeVideoID, err := youtubeGetVideoIDFromSeachQuery(chapter.YouTubeSearchQuery)
			if err != nil {
				app.serverError(w, r, err)
				return
			}
			err = app.chapters.Insert(chapter.ChapterTitle, chapter.YouTubeSearchQuery, unitID, chapter.Summary, youtubeVideoID)
			if err != nil {
				app.serverError(w, r, err)
				return
			}
		}
	}
	_, err = app.accounts.DecrementCredits(accountID, credits)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.Header().Set("HX-Redirect", fmt.Sprintf("/create/%s", courseID))
	w.WriteHeader(http.StatusSeeOther)
}

func (app *application) callback(w http.ResponseWriter, r *http.Request) {
	account, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	accountIDBlob, err := json.Marshal(account.UserID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	cookieStore, err := app.store.Get(r, "discovery-trail")
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	currentSessionToken := cookieStore.Values["token"]
	if currentSessionToken != nil {
		err := app.sessions.Destroy(currentSessionToken.(string))
		if err != nil {
			app.serverError(w, r, err)
			return
		}
		cookieStore.Values["token"] = nil
	}

	exists, err := app.accounts.Exists(account.UserID)

	if err != nil {
		if !errors.Is(err, models.ErrorNoRecord) {
			app.serverError(w, r, err)
			return
		}
	}
	if !exists {
		err = app.accounts.Insert(account.UserID, fmt.Sprintf("%s %s", account.FirstName, account.LastName), account.Email)
		if err != nil {
			app.serverError(w, r, err)
			return
		}
	}

	sessionID, err := app.sessions.Create(accountIDBlob)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	cookieStore.Values["token"] = sessionID
	err = cookieStore.Save(r, w)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func (app *application) logout(w http.ResponseWriter, r *http.Request) {
	gothic.Logout(w, r)
	cookieStore, err := app.store.Get(r, "discovery-trail")
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	currentSessionToken := cookieStore.Values["token"]
	if currentSessionToken != nil {
		err := app.sessions.Destroy(currentSessionToken.(string))
		if err != nil {
			if !errors.Is(err, models.ErrorNoRecord) {
				app.serverError(w, r, err)
				return
			}
		}
	}

	cookieStore.Options.MaxAge = -1
	err = cookieStore.Save(r, w)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	account, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		gothic.BeginAuthHandler(w, r)
		return
	}
	accountIDBlob, err := json.Marshal(account.UserID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	cookieStore, err := app.store.Get(r, "discovery-trail")
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	currentSessionToken := cookieStore.Values["token"]
	if currentSessionToken != nil {
		err := app.sessions.Destroy(currentSessionToken.(string))
		if err != nil {
			app.serverError(w, r, err)
			return
		}
		cookieStore.Values["token"] = nil
	}

	exists, err := app.accounts.Exists(account.UserID)

	if err != nil {
		if !errors.Is(err, models.ErrorNoRecord) {
			app.serverError(w, r, err)
			return
		}
	}
	if !exists {
		err = app.accounts.Insert(account.UserID, fmt.Sprintf("%s %s", account.FirstName, account.LastName), account.Email)
		if err != nil {
			app.serverError(w, r, err)
			return
		}
	}

	sessionID, err := app.sessions.Create(accountIDBlob)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	cookieStore.Values["token"] = sessionID
	err = cookieStore.Save(r, w)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func (app *application) chapterStatusPost(w http.ResponseWriter, r *http.Request) {
	csrfToken := r.Header.Get("X-CSRF-TOKEN")
	chapterID := r.PathValue("chapterID")
	cdx := r.PathValue("cdx")
	chapterIDInt, err := strconv.Atoi(chapterID)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	cdxInt, err := strconv.Atoi(cdx)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	chapter, err := app.chapters.Get(int64(chapterIDInt))
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	generatedQuestion, err := gpt.GenerateQuestionFromChapter(app.openAiClient, chapter.Name)
	if err != nil {
		app.renderErrorChapter(w, r, chapter, cdxInt, csrfToken, err)
		return
	}
	questionID, err := app.questions.Insert(generatedQuestion, chapter.ID)
	if err != nil {
		app.renderErrorChapter(w, r, chapter, cdxInt, csrfToken, err)
		return
	}
	if err = app.chapters.UpdateChapterQuestionStatus(chapter.ID, models.QuestionStatuses[models.Completed]); err != nil {
		_ = app.questions.Delete(questionID)
		app.renderErrorChapter(w, r, chapter, cdxInt, csrfToken, err)
		return
	}

	chapter.QuestionsStatus = models.QuestionStatuses[models.Completed]
	chapterStatusPost := pages.ChapterQuestionsRender(cdxInt, chapter, csrfToken)
	app.render(w, r, http.StatusOK, chapterStatusPost)

}
func (app *application) renderErrorChapter(w http.ResponseWriter, r *http.Request, chapter models.Chapter, cdx int, csrfToken string, err error) {
	fmt.Println(err)
	_ = app.chapters.UpdateChapterQuestionStatus(chapter.ID, models.QuestionStatuses[models.Completed])
	chapter.QuestionsStatus = models.QuestionStatuses[models.Error]
	chapterStatusPost := pages.ChapterQuestionsRender(cdx, chapter, csrfToken)
	app.render(w, r, http.StatusOK, chapterStatusPost)
}

func (app *application) courseProcess(w http.ResponseWriter, r *http.Request) {
	courseID := r.PathValue("courseID")
	if err := app.courses.Process(courseID); err != nil {
		app.serverError(w, r, err)
		return
	}
	w.Header().Set("HX-Redirect", fmt.Sprintf("/course/%s", courseID))
	w.WriteHeader(http.StatusSeeOther)
}

func (app *application) courseUnitChapter(w http.ResponseWriter, r *http.Request) {
	courseID := r.PathValue("courseID")
	course, err := app.courses.Get(courseID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	units, err := app.units.GetCourseUnits(courseID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	courseUnitsChapters := types.CourseWithUnitsWithChapters{
		Course: course,
		Units:  make([]types.UnitWithChapters, 0),
	}

	for _, unit := range units {
		unitChapters, err := app.chapters.GetUnitChapters(unit.ID)
		if err != nil {
			app.serverError(w, r, err)
			return
		}
		courseUnitsChapters.Units = append(courseUnitsChapters.Units, types.UnitWithChapters{
			Chapters: unitChapters,
			Unit:     unit,
		})
	}
	data, err := app.newTemplateData(r)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	data.CourseUnitsChapters = courseUnitsChapters
	courseView := pages.Course(data)
	app.render(w, r, http.StatusOK, courseView)
}

type ChapterData struct {
	Name string `json:"name"`
}

func (app *application) ChapterInformation(w http.ResponseWriter, r *http.Request) {
	chapterIDStr := r.Header.Get("chapter-id")
	chapterIndexStr := r.Header.Get("chapter-index")
	unitIndexStr := r.Header.Get("unit-index")
	chapterIDInt, err := strconv.Atoi(chapterIDStr)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	chapter, err := app.chapters.Get(int64(chapterIDInt))
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	ChapterInformation := pages.ChapterInformation(chapter, unitIndexStr, chapterIndexStr)
	app.render(w, r, http.StatusOK, ChapterInformation)
}

func (app *application) deleteCourse(w http.ResponseWriter, r *http.Request) {
	courseID := r.PathValue("courseID")
	if err := app.courses.Delete(courseID); err != nil {
		app.logger.Error("error deleting course", "error", err)
	}
	app.logger.Info("deleted course completely fine")

}
