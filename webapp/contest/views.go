package contest

import (
	"github.com/flosch/pongo2"
	"github.com/paulCodes/pumpkin-voter/domain"
	"github.com/paulCodes/pumpkin-voter/pvhelpers"
	"github.com/paulCodes/pumpkin-voter/webtypes"
	"net/http"
	"strings"
	"github.com/gorilla/mux"
	"log"
	"fmt"
)

type ContestApp struct {
	Env webtypes.Env
}

func (app ContestApp) Contests(w http.ResponseWriter, r *http.Request) {
	contests, err := app.Env.Registry.Contest.All()
	if err != nil {
		panic("An error has occured accessing your contests." + err.Error())
	}

	models := []ContestLister{}
	for _, contest := range contests {
		models = append(models, ContestLister{
			Contest:  contest,
			Registry: app.Env.Registry,
		})
	}

	pvhelpers.RenderTemplate(w, r, "templates/contest/list.html",
		pongo2.Context{
			"point_to": "contest",
			"models":   models,
			"stub":     &ContestLister{},
		},
		"vote")
}

func (app ContestApp) Create(w http.ResponseWriter, r *http.Request) {
	session, _ := pvhelpers.Store.Get(r, "vote")
	contest := domain.Contest{}
	contest.Id = "-1"
	if r.Method == "POST" {
		contest.Title = strings.TrimSpace(r.FormValue("Title"))
		contest.CategoryIds = strings.TrimSpace(r.FormValue("CategoryIds"))
		contest.Active = strings.TrimSpace(r.FormValue("Active")) == "true"

		//TODO add validation
		contest.Id = pvhelpers.GenerateUUID()
		err := app.Env.Registry.Contest.Add(contest)

		if err != nil {
			pvhelpers.AddFlash(w,r,pvhelpers.FlashMessage{MsgType: "danger", Msg: "Error saving contest"}, "vote")
		} else {
		pvhelpers.AddFlash(w, r, pvhelpers.FlashMessage{MsgType: "success", Msg: "Contest created successfully"}, "vote")
		session.Save(r, w)
		http.Redirect(w, r, "/voter", http.StatusFound)
		return
		}
	}

	session.Save(r, w)
	ctx := pongo2.Context{
		"point_to":         "contest",
		"model":            ContestLister{Contest: contest, Registry: app.Env.Registry},
		"id":               contest.Id,
	}
	pvhelpers.RenderTemplate(w, r, "templates/contest/create.html", ctx, "voter")
}

func (app ContestApp) Edit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	contestId := vars["contestId"]
	session, _ := pvhelpers.Store.Get(r, "vote")
	contest, err := app.Env.Registry.Contest.GetID(contestId)
	if err != nil {
		pvhelpers.LogErrorObject(err)
		pvhelpers.AddFlash(w, r, pvhelpers.FlashMessage{MsgType: "danger", Msg: "Could not find contest"}, "vote")
		session.Save(r, w)
		http.Redirect(w, r, "/voter/contest", http.StatusFound)
		return
	}
	if r.Method == "POST" {
		contest.Title = strings.TrimSpace(r.FormValue("Title"))
		contest.CategoryIds = strings.TrimSpace(r.FormValue("CategoryIds"))
		contest.Active = strings.TrimSpace(r.FormValue("Active")) == "true"

		//TODO add validation
		//contest.Id = pvhelpers.GenerateUUID()
		err := app.Env.Registry.Contest.Replace(contest)

		if err != nil {
			pvhelpers.AddFlash(w,r,pvhelpers.FlashMessage{MsgType: "danger", Msg: "Error saving contest"}, "vote")
		} else {
			pvhelpers.AddFlash(w, r, pvhelpers.FlashMessage{MsgType: "success", Msg: "Contest created successfully"}, "vote")
			session.Save(r, w)
			http.Redirect(w, r, "/voter/contest", http.StatusFound)
			return
		}
	}

	session.Save(r, w)
	ctx := pongo2.Context{
		"point_to":         "contest",
		"model":            ContestLister{Contest: contest, Registry: app.Env.Registry},
		"id":               contest.Id,
	}
	pvhelpers.RenderTemplate(w, r, "templates/contest/create.html", ctx, "voter")
}

func (app ContestApp) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	contestId := vars["contestId"]
	session, _ := pvhelpers.Store.Get(r, "vote")
	contest, err := app.Env.Registry.Contest.GetID(contestId)
	if err != nil {
		pvhelpers.LogErrorObject(err)
		pvhelpers.AddFlash(w, r, pvhelpers.FlashMessage{MsgType: "danger", Msg: "Could not find contest"}, "vote")
		session.Save(r, w)
		http.Redirect(w, r, "/voter/contest", http.StatusFound)
		return
	}

	err = app.Env.Registry.Contest.Delete(contest)
	if err != nil {
		pvhelpers.LogErrorObject(err)
		pvhelpers.AddFlash(w, r, pvhelpers.FlashMessage{MsgType: "danger", Msg: "Error deleting contest"}, "vote")
		session.Save(r, w)
		http.Redirect(w, r, "/voter/contest", http.StatusFound)
		return
	}

	session.Save(r, w)
	http.Redirect(w, r, "/voter/contest", http.StatusFound)
}

func (app ContestApp) Results(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	contestId := vars["contestId"]
	session, _ := pvhelpers.Store.Get(r, "vote")
	contest, err := app.Env.Registry.Contest.GetID(contestId)
	if err != nil {
		pvhelpers.LogErrorObject(err)
		pvhelpers.AddFlash(w, r, pvhelpers.FlashMessage{MsgType: "danger", Msg: "Could not find contest"}, "vote")
		session.Save(r, w)
		http.Redirect(w, r, "/voter/contest", http.StatusFound)
		return
	}

	tempcategoryIds, err := app.Env.Registry.Entry.FindAllCategoryIdFromContest(contest.Id)
	log.Printf("tenpCategoryIds %v %v", tempcategoryIds,err)

	var categoryIds []string
	for _, catIds := range tempcategoryIds {
		cids := strings.Split(catIds, ",")
		for _, cid := range cids {
			if !stringInSlice(cid, categoryIds) {
				categoryIds = append(categoryIds, cid)
			}
		}
	}
	var contestResults domain.ContestResults
	var results []domain.CategoryVoteCalc
	for _, catId := range categoryIds {
		var categoryVoteCalc domain.CategoryVoteCalc
		voteCalcs := app.Env.Registry.Contest.ResultsByContestIdAndCategoryId(contest.Id, catId)
		category, _ := app.Env.Registry.Category.GetID(catId)

log.Printf("voteCalcs %v", voteCalcs)
		categoryVoteCalc.Category = category
		categoryVoteCalc.VoteCalcs = voteCalcs
		results = append(results, categoryVoteCalc)
	}
	log.Printf("results %v", results)
	contestResults.ContestId = contest.Id
	contestResults.ContestTitle = contest.Title
	contestResults.Results = results

	session.Save(r, w)
	ctx := pongo2.Context{
		"point_to":         "contest",
		"model":            contestResults,
		"id":               contest.Id,
	}
	pvhelpers.RenderTemplate(w, r, "templates/vote/voteResults.html", ctx, "voter")
}

func (app ContestApp) Vote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	contestId := vars["contestId"]
	session, _ := pvhelpers.Store.Get(r, "vote")
	contest, err := app.Env.Registry.Contest.GetID(contestId)
	if err != nil {
		pvhelpers.LogErrorObject(err)
		pvhelpers.AddFlash(w, r, pvhelpers.FlashMessage{MsgType: "danger", Msg: "Could not find contest"}, "vote")
		session.Save(r, w)
		http.Redirect(w, r, "/voter/contest", http.StatusFound)
		return
	}

	// build out vote form
	if r.Method == "GET" {
		votingForm := app.buildContestVotingForm(contest)
		ctx := pongo2.Context{
			"point_to":         "contest",
			"model":            votingForm,
			"id":               contest.Id,
		}
		pvhelpers.RenderTemplate(w, r, "templates/vote/voteForm.html", ctx, "voter")
		return
	}

	// save votes
	if r.Method == "POST" {
		log.Printf("HELLS YEA POST TIME")
		r.ParseForm()
		for k := range r.Form {
			var vote domain.Vote
			vote.Id = pvhelpers.GenerateUUID()
			vote.ContestId = contestId
			vote.CategoryId = k
			vote.EntryId = r.Form.Get(k)
			log.Printf("vote thing || %v %v %v", vote.EntryId, vote.CategoryId, vote.ContestId)
			err = app.Env.Registry.Vote.Add(vote)
			if err != nil {
				log.Printf("vote save failed %v", err)
			}
		}

		session.Save(r, w)
		http.Redirect(w, r, fmt.Sprintf("/voter"), http.StatusFound)
		return
	}

	session.Save(r, w)
	http.Redirect(w, r, fmt.Sprintf("/voter/contest/vote/%v", contestId), http.StatusFound)
}

func (app ContestApp) buildContestVotingForm(contest domain.Contest) (votingForm domain.VoteForm) {
	var categoryEntries []domain.CategoryEntries
	tempcategoryIds, err := app.Env.Registry.Entry.FindAllCategoryIdFromContest(contest.Id)
	log.Printf("tenpCategoryIds %v %v", tempcategoryIds,err)

	var categoryIds []string
	for _, catIds := range tempcategoryIds {
		cids := strings.Split(catIds, ",")
		for _, cid := range cids {
			if !stringInSlice(cid, categoryIds) {
				categoryIds = append(categoryIds, cid)
			}
		}
	}

	log.Printf("CategoryIds %v", tempcategoryIds)

	for _, categoryId := range categoryIds {
		var categoryEntry domain.CategoryEntries
		category, err := app.Env.Registry.Category.GetID(categoryId)
		if err != nil {
			log.Printf("%v --build voter form failed : %v", categoryId, err)
		} else {
			entries, err := app.Env.Registry.Entry.FindAllForCategoryId(categoryId)
			if err != nil {
				log.Printf("fuck !! %v", err)
			}
			categoryEntry.Category = category
			categoryEntry.Entries = entries

			categoryEntries = append(categoryEntries, categoryEntry)
		}

	}

	votingForm.ContestTitle = contest.Title
	votingForm.ContestId = contest.Id
	votingForm.EntryByCategory = categoryEntries
	return

}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}