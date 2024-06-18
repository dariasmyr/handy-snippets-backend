package graph

import (
	"context"
	"pastebin/graph/model"
)

func (r *mutationResolver) CreateDocument(ctx context.Context, input model.CreateDocumentInput) (*int, error) {
	var id, err = r.DocumentService.CreateDocument(input.Value, *input.Title, input.AccessKey, *input.MaxViewCount, *input.TTLMs)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func (r *mutationResolver) DeleteDocument(ctx context.Context, id int, accessKey string) (*bool, error) {
	success, err := r.DocumentService.DeleteDocument(id, accessKey)
	if err != nil {
		return nil, err
	}
	return &success, nil
}

func (r *queryResolver) GetDocument(ctx context.Context, id int) (*model.Document, error) {
	return r.DocumentService.GetDocument(id)
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
