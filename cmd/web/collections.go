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

type createCollectionForm struct {
	Name               string `form:"name"`
	ParentID           string `form:"parent_id"`
	validate.Validator `form:"_"`
}

type collectionPageData struct {
	Collections []*data.Collection
	ParentID    uuid.UUID
}

func (b *backend) viewRootCollection(c echo.Context) error {
	root := b.contextGetRootCollection(c)

	colxns, err := b.collectionsPageData(c, root.ID)
	if err != nil {
		return b.renderInternalServerErr(c)
	}
	colxnForm := new(createCollectionForm)
	colxnData := collectionPageData{
		Collections: colxns.Collections, ParentID: colxns.ParentID,
	}
	pd := b.NewPageData(c)
	pd.Form = colxnForm
	pd.Data = colxnData
	return c.Render(http.StatusOK, "collections.gohtml", pd)
}

func (b *backend) viewCreateCollectionForm(c echo.Context) error {
	pd := b.NewPageData(c)
	pd.Form = new(createCollectionForm)
	return c.Render(http.StatusOK, "collections.gohtml", pd)
}

func (b *backend) postCreateCollectionForm(c echo.Context) error {
	var colxnForm createCollectionForm

	if err := b.decodePostForm(c, &colxnForm); err != nil {
		b.zlog.Err(err).Msg("decode post form data")
		return b.renderBadRequestErr(c, "invalid form data")
	}
	pID, err := uuid.Parse(colxnForm.ParentID)
	if err != nil {
		return b.renderBadRequestErr(c, "invalid parent id")
	}
	colxnData, err := b.collectionsPageData(c, pID)
	if err != nil {
		b.zlog.Err(err).
			Str("handler", "postCreateCollectionForm").
			Msg("get collections page data")
		return b.renderInternalServerErr(c)
	}
	userID := b.contextGetUser(c).ID

	colxn := &data.Collection{
		ID: uuid.New(), ParentID: uuid.NullUUID{
			UUID: pID, Valid: true,
		},
		UserID: userID,
		Name:   colxnForm.Name,
	}
	validate.ValidateCollection(&colxnForm.Validator, colxn)
	if !colxnForm.Validator.Valid() {
		return b.renderErrFormData(c, "collections.gohtml", &colxnForm, colxnData)
	}
	if err = b.models.Collections.Insert(c.Request().Context(), colxn); err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateCollection):
			colxnForm.AddFieldError("name", "collection already exists")
			return b.renderErrFormData(c, "collections.gohtml", &colxnForm, colxnData)
		default:
			b.zlog.Err(err).
				Str("handler", "postCreateCollectionForm").
				Msg("insert collection into db")
			return b.renderInternalServerErr(c)
		}
	}
	return c.Redirect(http.StatusSeeOther, "/v1/collection")
}

func (b *backend) collectionsPageData(c echo.Context, pID uuid.UUID) (*collectionPageData, error) {
	colxns, err := b.models.Collections.GetAllForParentId(c.Request().Context(), pID)
	if err != nil {
		return nil, fmt.Errorf("couldn't get collections. pID=%s: %w", pID.String(), err)
	}
	return &collectionPageData{
		Collections: colxns, ParentID: pID,
	}, nil
}
