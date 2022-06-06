package observer

import "fmt"

type ShopType int

const (
	ShopType_FlagShip ShopType = iota // 旗舰店
	ShopType_Normal                   // 普通店
)

// IMsgNotifier 消息触达者
type IMsgNotifier interface {
	GetName() string
	Notify(*Shop) error
}

type Shop struct {
	Id   int64
	Name string
	ShopType
}

type ShopNotifier struct {
	shop      *Shop
	notifiers []IMsgNotifier
}

func (s *ShopNotifier) NotifyAll() error {
	for _, notifier := range s.notifiers {
		if err := notifier.Notify(s.shop); err != nil {
			return err
		}
	}
	return nil
}

// 根据不同的店铺类型组装不同的消息触达者
func getShopNotifier(shop *Shop) *ShopNotifier {
	shopNotifier := &ShopNotifier{
		shop:      shop,
		notifiers: make([]IMsgNotifier, 0),
	}
	if shop == nil {
		return shopNotifier
	}
	switch shop.ShopType {
	case ShopType_FlagShip: // 旗舰店用短信+站内信通知
		shopNotifier.notifiers = append(shopNotifier.notifiers, &ShortMsgNotifier{}, &StationLetterNotifier{})
	case ShopType_Normal: // 普通店只用站内信通知
		shopNotifier.notifiers = append(shopNotifier.notifiers, &StationLetterNotifier{})
	}
	return shopNotifier
}

type ShortMsgNotifier struct{}

func (s *ShortMsgNotifier) GetName() string {
	return "短信"
}

func (s *ShortMsgNotifier) Notify(shop *Shop) error {
	fmt.Printf("【短信通知】店铺名称 %v\n", shop.Name)
	return nil
}

type StationLetterNotifier struct{}

func (s *StationLetterNotifier) GetName() string {
	return "站内信"
}

func (s *StationLetterNotifier) Notify(shop *Shop) error {
	fmt.Printf("【站内信通知】店铺名称 %v\n", shop.Name)
	return nil
}
