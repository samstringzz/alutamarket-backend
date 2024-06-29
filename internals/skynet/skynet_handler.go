package skynet

import "context"

type Handler struct {
	Service
}

func NewHandler(s Service) *Handler {
	return &Handler{
		Service: s,
	}
}

func (h *Handler) BuyAirtime(ctx context.Context, input *Airtime) (*string, error) {
	item, err := h.Service.BuyAirtime(ctx, input)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (h *Handler) BuyData(ctx context.Context, input *Data) (*string, error) {
	item, err := h.Service.BuyData(ctx, input)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (h *Handler) BuyTVSubscription(ctx context.Context, input *TVSubscription) (*string, error) {
	item, err := h.Service.BuyTVSubscription(ctx, input)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (h *Handler) BuyEducational(ctx context.Context, input *EducationPayment) (*string, error) {
	item, err := h.Service.BuyEducational(ctx, input)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (h *Handler) VerifySmartCard(ctx context.Context, serviceID, billersCode, cardType string) (*SmartcardVerificationResponse, error) {
	item, err := h.Service.VerifySmartCard(ctx, serviceID, billersCode, cardType)
	if err != nil {
		return nil, err
	}
	return item, nil
}
