package main

import (
	"Continuity/internal/data"
	"Continuity/internal/validate"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// createCollectionForm represents a parent collection that will create a sub collection
// within it.
type createCollectionForm struct {
	// This is the ID of the parent page within which the current collection
	// will be created. E.g.: for a parent marvel, Iron Man will be created
	// such that you won't enter into the Iron Man collection but will stay
	// in Marvel cause that's the ID this field refers to.
	CollectionID       string `form:"collection_id"`
	Name               string `form:"name"`
	IsRoot             bool   `form:"is_root"`
	validate.Validator `form:"_"`
}

// collectionPageData is all the data that should be shown on the page for the current
// collection.
type collectionPageData struct {
	CollectionID uuid.UUID
	IsRoot       bool
	Collections  []*data.Collection
}

func (b *backend) viewRootCollection(c echo.Context) error {
	root := b.contextGetRootCollection(c)

	colxnPage, err := b.collectionsPageData(c, root.ID)
	if err != nil {
		b.zlog.Err(err).
			Str("handler", "viewCollectionPage").
			Msg("get collections for parent id")
		return b.renderInternalServerErr(c)
	}
	colxnPage.IsRoot = true
	colxnForm := new(createCollectionForm)

	pd := b.NewPageData(c)
	pd.Form = colxnForm
	pd.Data = colxnPage
	return c.Render(http.StatusOK, "collections.gohtml", pd)
}

func (b *backend) viewCollectionPage(c echo.Context) error {
	colxnID, err := uuid.Parse(c.Param("collection-id"))
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/v1/collection")
	}
	colxnPage, err := b.collectionsPageData(c, colxnID)
	if err != nil {
		b.zlog.Err(err).
			Str("handler", "viewCollectionPage").
			Msg("get collections for parent id")
		return b.renderInternalServerErr(c)
	}
	root := b.contextGetRootCollection(c)

	if root.ID == colxnID {
		colxnPage.IsRoot = true
	} else {
		colxnPage.IsRoot = false
	}
	colxnForm := new(createCollectionForm)

	pd := b.NewPageData(c)
	pd.Form = colxnForm
	pd.Data = colxnPage
	return c.Render(http.StatusOK, "collections.gohtml", pd)
}

func (b *backend) postCreateCollectionForm(c echo.Context) error {
	var colxnForm createCollectionForm

	if err := b.decodePostForm(c, &colxnForm); err != nil {
		return b.renderBadRequestErr(c, "invalid form data")
	}
	pID, err := uuid.Parse(colxnForm.CollectionID)
	if err != nil {
		return b.renderBadRequestErr(c, "invalid parent id")
	}
	colxnPage, err := b.collectionsPageData(c, pID)
	if err != nil {
		b.zlog.Err(err).
			Str("handler", "postCreateCollectionForm").
			Msg("get collections page data")
		return b.renderInternalServerErr(c)
	}
	colxnPage.IsRoot = colxnForm.IsRoot
	userID := b.contextGetUser(c).ID

	colxn := &data.Collection{
		ID: uuid.New(), ParentID: uuid.NullUUID{
			UUID: pID, Valid: true,
		},
		UserID: userID,
		Name:   colxnForm.Name,
	}
	validate.ValidateCollection(&colxnForm.Validator, colxn)
	if !colxnForm.Valid() {
		return b.renderErrFormData(c, "collections.gohtml", &colxnForm, colxnPage)
	}
	if err = b.models.Collections.Insert(c.Request().Context(), colxn); err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateCollection):
			colxnForm.AddFieldError("name", "collection already exists")
			return b.renderErrFormData(c, "collections.gohtml", &colxnForm, colxnPage)
		default:
			b.zlog.Err(err).
				Str("handler", "postCreateCollectionForm").
				Msg("insert collection into db")
			return b.renderInternalServerErr(c)
		}
	}
	return c.Redirect(http.StatusSeeOther, "/v1/collection/"+colxnForm.CollectionID)
}

func (b *backend) deleteCollectionPage(c echo.Context) error {
	var colxnForm createCollectionForm
	if err := b.decodePostForm(c, &colxnForm); err != nil {
		b.zlog.Err(err).Msg("while decoding delete-form-data")
		return b.renderBadRequestErr(c, "invalid form data")
	}
	if colxnForm.CollectionID != c.Param("collection-id") {
		return b.renderBadRequestErr(c, "collection ids don't match")
	}
	cID, err := uuid.Parse(colxnForm.CollectionID)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/v1/collection")
	}
	colxn, err := b.models.Collections.GetById(c.Request().Context(), cID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			return b.renderNotFoundErr(c, "collection not found")
		default:
			return b.renderInternalServerErr(c)
		}
	}
	userID := b.contextGetUser(c).ID

	err = b.models.Collections.Delete(c.Request().Context(), cID, userID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrCollectionNotEmpty):
			colxnPage, err := b.collectionsPageData(c, cID)
			if err != nil {
				return b.renderInternalServerErr(c)
			}
			colxnPage.IsRoot = colxnForm.IsRoot

			colxnForm.AddError("collection", "collection is not empty")
			return b.renderErrFormData(c, "collections.gohtml", &colxnForm, colxnPage)

		case errors.Is(err, data.ErrRecordNotFound):
			return b.renderNotFoundErr(c, "collection not found")
		default:
			b.zlog.Err(err).
				Str("handler", "deleteCollectionPage").
				Msg("delete collection")
			return b.renderInternalServerErr(c)
		}
	}
	rootId := b.contextGetRootCollection(c).ID.String()
	parentID := colxn.ParentID.UUID.String()
	switch {
	case rootId == parentID:
		return c.Redirect(http.StatusSeeOther, "/v1/collection")
	default:
		return c.Redirect(http.StatusSeeOther, "/v1/collection/"+parentID)
	}
}

// collectionsPageData gets all the sub-collections for the provided collection
// id and returns an initialized collectionPageData struct or an error.
func (b *backend) collectionsPageData(c echo.Context, cID uuid.UUID) (*collectionPageData, error) {
	colxns, err := b.models.Collections.GetAllForParentId(c.Request().Context(), cID)
	if err != nil {
		return nil, fmt.Errorf("couldn't get collections. cID=%s: %w", cID.String(), err)
	}
	return &collectionPageData{
		Collections: colxns, CollectionID: cID,
	}, nil
}
