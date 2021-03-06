package actions

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"

	"github.com/IEEESBITBA/Curso-de-Python-Sistemas/models"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop/v5"
	"github.com/pkg/errors"
)

// CategoriesIndex default implementation.
func CategoriesIndex(c buffalo.Context) error {
	catTitle := c.Param("cat_title")
	c.Logger().Debugf("accessed %s", catTitle)
	tx := c.Value("tx").(*pop.Connection)
	cat := &models.Category{}
	err := tx.Where("title = ?", catTitle).First(cat)
	if err != nil {
		c.Logger().Warnf("'where title = %s' FAILED!", catTitle)
		return c.Error(404, err)
	}
	c.Set("category", cat)
	forum := c.Value("forum").(*models.Forum)
	role := c.Value("role").(string)
	var ordering string
	var perPage, page int
	if forum.Defcon == "2" {
		page, perPage = setPagination(c.Params(), 100)
		ordering = "archived, created_at desc"
	} else if role == "admin" {
		page, perPage = setPagination(c.Params(), 12)
		ordering = "archived, created_at desc"
	} else {
		page, perPage = setPagination(c.Params(), 8)
		ordering = "created_at desc"
	}
	q := tx.BelongsTo(cat).Where("deleted IS false").Order(ordering).Paginate(page, perPage)

	topics := &models.Topics{}
	if err := q.All(topics); err != nil {
		c.Logger().Warnf("'tx.BelongsTo(cat).Order(%q).PaginateFromParams(c.Params())' FAILED!", ordering)
		return c.Error(404, err)
	}
	for i, t := range *topics {
		topic, err := loadTopic(c, t.ID.String())
		if err != nil {
			c.Logger().Errorf("'loadTopic(c, %s)' failed: %s", t.ID.String(), err)
			return errors.WithStack(err)
		}
		(*topics)[i] = *topic
	}

	var renderer render.Renderer
	if forum.Defcon == "2" {
		renderer = r.HTML("categories/index_voting.plush.html")
		sort.Sort(models.ByVotes(*topics))
	} else {
		renderer = r.HTML("categories/index.plush.html")
		if role == "admin" {
			sort.Sort(models.ByArchived(*topics))
		} else {
			sort.Sort(topics)
		}
	}

	c.Set("topics", topics)
	c.Set("pagination", q.Paginator)
	return c.Render(200, renderer)
	//return c.Render(http.StatusOK, r.HTML("categories/index.plush.html"))
}

func setPagination(params buffalo.ParamValues, defaultPerPage int) (page, perpage int) {
	if p, _ := strconv.Atoi(params.Get("page")); p < 1 {
		page = 1
	} else {
		page = p
	}
	if pp, _ := strconv.Atoi(params.Get("per_page")); pp < 1 {
		perpage = defaultPerPage
	} else {
		perpage = pp
	}
	return
}

// CategoriesCreateGet default implementation.
func CategoriesCreateGet(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.HTML("categories/create_get.plush.html"))
}

// CategoriesCreateOrEditPost default implementation.
func CategoriesCreateOrEditPost(c buffalo.Context) error {
	cat := &models.Category{}
	if err := c.Bind(cat); err != nil {
		c.Flash().Add("danger", "could not create category")
		return c.Error(500, err)
	}
	if !validURLDir(cat.Title) {
		c.Flash().Add("danger", T.Translate(c, "category-invalid-title"))
		return c.Redirect(302, "")
	}
	f := c.Value("forum").(*models.Forum)
	cat.ParentCategory = nulls.NewUUID(f.ID)
	tx := c.Value("tx").(*pop.Connection)
	q := tx.Where("title = ?", cat.Title)
	exist, err := q.Exists(&models.Forum{})
	if exist {
		c.Flash().Add("danger", "Category with same title already exists")
		return c.Redirect(302, "/")
	}
	if err != nil {
		return c.Error(500, err)
	}
	v, _ := cat.Validate(tx)
	if v.HasAny() {
		c.Flash().Add("danger", "Title should have something!")
		return c.Redirect(302, "/")
	}
	if len(cat.Description.String) > 255 {
		c.Flash().Add("danger", "Description too long, should be under 255 char (SQL varchar[255])!")
		return c.Redirect(302, "/")
	}
	nilIfNewCat := c.Value("category") // if we are editing a category this will be set
	if nilIfNewCat != nil {            // this branch edits an already existing category

		oldCat := nilIfNewCat.(*models.Category)
		oldCat.Title, oldCat.Description = cat.Title, cat.Description
		err = tx.Update(oldCat)

	} else { // this branch creates a new category
		err = tx.Save(cat)
	}

	if err != nil {
		c.Logger().Errorf("creating or editing category %s: %s", cat.Title, err)
		c.Flash().Add("danger", "Error creating category")
		return errors.WithStack(err)
	}
	u := c.Value("current_user").(*models.User)
	c.Logger().Infof("create category: %s, by %s", cat.Title, u.Email)
	c.Flash().Add("success", fmt.Sprintf("Category %s created", cat.Title))
	return c.Redirect(302, "forumPath()", render.Data{"forum_title": f.Title})
}

// SetCurrentCategory attempts to find a category and set context `category`
func SetCurrentCategory(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		tx := c.Value("tx").(*pop.Connection)
		cat := &models.Category{}
		title := c.Param("cat_title")
		if title == "" {
			return c.Redirect(302, "forumPath()")
		}
		q := tx.Where("title = ?", title)
		err := q.First(cat)
		if err != nil {
			c.Flash().Add("danger", "Error while seeking category")
			return c.Redirect(302, "forumPath()")
		}
		c.Set("inCat", true)
		c.Set("category", cat)
		return next(c)
	}
}
