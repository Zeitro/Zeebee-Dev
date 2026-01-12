package zeitrologic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/zeitro"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
)

type ZeitroListingSearchLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewZeitroListingSearchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ZeitroListingSearchLogic {
	return &ZeitroListingSearchLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ZeitroListingSearchLogic) ZeitroListingSearch(in *borrower.ZeitroListingSearchRequest) (resp *borrower.ZeitroListingSearchResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()

	// Call Zeitro client to search listings
	listings, err := l.svcCtx.ZeitroClient.SearchListings(l.ctx, zeitro.SearchRequest{
		Lat:  in.Lat,
		Lon:  in.Lon,
		Dist: in.Dist,
		Addr: in.Addr,
	})
	if err != nil {
		return nil, err
	}

	// slice.ForEach(listings, func(index int, item zeitro.PropertyListing) {
	// 	price := strings.TrimSpace(item.Price)
	// 	priceNum, _ := convertor.ToInt(price)
	// 	listings[index].Price = stringx.FormatUsDollarCurrency(float64(priceNum))
	// })

	// Convert domain model to protobuf using copierx.Trans
	result, err := copierx.Trans[[]zeitro.PropertyListing, []*borrower.ZeitroPropertyEntity](&listings)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	return &borrower.ZeitroListingSearchResponse{
		Data: *result,
	}, nil
}
