package main

import (
	"context"
	"log"

	"github.com/jxskiss/kitex-reflect/example/bookinfo/kitex_gen/cwg/bookinfo/details"
)

// DetailsServiceImpl implements the last service interface defined in the IDL.
type DetailsServiceImpl struct{}

// GetProduct implements the DetailsServiceImpl interface.
func (s *DetailsServiceImpl) GetProduct(ctx context.Context, req *details.GetProductReq) (resp *details.GetProductResp, err error) {
	log.Printf("get product details %s", req.ID)

	return &details.GetProductResp{
		Product: &details.Product{
			ID:          req.GetID(),
			Title:       "《Also sprach Zarathustra》",
			Author:      "Friedrich Nietzsche",
			Description: `Thus Spoke Zarathustra: A Book for All and None, also translated as Thus Spake Zarathustra, is a work of philosophical fiction written by German philosopher Friedrich Nietzsche between 1883 and 1885.`,
		},
	}, nil
}
