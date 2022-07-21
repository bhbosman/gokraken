package connection

//
//import (
//	"github.com/bhbosman/goCommonMarketData/fullMarketData"
//	marketDataStream "github.com/bhbosman/goMessages/marketData/stream"
//	"github.com/bhbosman/goerrors"
//	"github.com/emirpasic/gods/trees/avltree"
//	"github.com/emirpasic/gods/utils"
//	"hash/crc32"
//	"strings"
//)
//
//type OrderSide int8
//
//const BuySide OrderSide = 0
//const AskSide OrderSide = 1
//
////type PricePoint struct {
////	Price    float64
////	Volume   float64
////	Touched  bool
////	crcValue string
////}
//
//type FullMarketOrderBook struct {
//	InstrumentName     string
//	OrderSide          [2]*avltree.Tree
//	depth              int
//	SourceTimestamp    int64
//	SourceMessageCount int64
//	lastCheckSum       uint32
//}
//
//func (self *FullMarketOrderBook) VerifyChecksum() error {
//	if self.lastCheckSum != 0 {
//		localCheckSum := self.CalculateCheckSum()
//		if localCheckSum != self.lastCheckSum {
//			return goerrors.InvalidState
//		}
//	}
//	return nil
//}
//
//func (self *FullMarketOrderBook) CalculateCheckSum() uint32 {
//	crc := crc32.NewIEEE()
//	count := 0
//	for node := self.OrderSide[AskSide].Left(); node != nil && count < 10; node = node.Next() {
//		pp := node.Value.(*fullMarketData.PricePoint)
//		unk, _ := pp.List.Get(0)
//		ss := unk.(*fullMarketData.FullMarketOrder)
//		crc.Write([]byte(StateTrimmer(ss.ExtraData)))
//		count++
//	}
//	count = 0
//	for node := self.OrderSide[BuySide].Right(); node != nil && count < 10; node = node.Prev() {
//		pp := node.Value.(*fullMarketData.PricePoint)
//		unk, _ := pp.List.Get(0)
//		ss := unk.(*fullMarketData.FullMarketOrder)
//		crc.Write([]byte(StateTrimmer(ss.ExtraData)))
//		count++
//		count++
//	}
//	return crc.Sum32()
//}
//
//func (self *FullMarketOrderBook) Publish(forcePublish bool) *marketDataStream.PublishTop5 {
//	thereWasAChange := forcePublish
//	maxDepth := 10000
//	var bids []*marketDataStream.Point
//	if highBidNode := self.OrderSide[BuySide].Right(); highBidNode != nil {
//		count := 0
//		for node := highBidNode; node != nil && count < maxDepth; node = node.Prev() {
//			bidPrice := node.Key.(float64)
//			if pp, ok := node.Value.(*PricePoint); ok {
//				thereWasAChange = thereWasAChange || pp.Touched
//				pp.Touched = false
//				bids = append(bids, &marketDataStream.Point{
//					Price:  bidPrice,
//					Volume: pp.Volume,
//				})
//			}
//			count++
//		}
//	}
//	var asks []*marketDataStream.Point
//	if lowAskNode := self.OrderSide[AskSide].Left(); lowAskNode != nil {
//		count := 0
//		for node := lowAskNode; node != nil && count < maxDepth; node = node.Next() {
//			askPrice := node.Key.(float64)
//			if pp, ok := node.Value.(*PricePoint); ok {
//				thereWasAChange = thereWasAChange || pp.Touched
//				pp.Touched = false
//				asks = append(asks, &marketDataStream.Point{
//					Price:  askPrice,
//					Volume: pp.Volume,
//				})
//			}
//			count++
//		}
//	}
//	spread := 0.0
//	if len(asks) > 0 && len(bids) > 0 {
//		spread = asks[0].Price - bids[0].Price
//	}
//	if thereWasAChange {
//		if !forcePublish {
//			//self.UpdateCount++
//		}
//		top5 := &marketDataStream.PublishTop5{
//			Instrument: self.InstrumentName,
//			Spread:     spread,
//			//SourceTimeStamp:    self.FullMarketOrderBook.SourceTimestamp,
//			//SourceMessageCount: self.FullMarketOrderBook.SourceMessageCount,
//			//UpdateCount:        self.UpdateCount,
//			Bid: bids,
//			Ask: asks,
//		}
//		return top5
//	}
//	return nil
//}
//
//func StateTrimmer(s string) string {
//	addValue := false
//	ddd := func(r rune) bool {
//		switch r {
//		case '0':
//			return addValue
//		case '.':
//			return false
//		default:
//			addValue = true
//			return true
//		}
//	}
//	sb := strings.Builder{}
//	for _, r := range s {
//		if ddd(r) {
//			sb.WriteRune(r)
//		}
//	}
//	return sb.String()
//}
//
//func NewFullMarketOrderBook(instrumentName, bookName string) *FullMarketOrderBook {
//	var depth int
//	switch bookName {
//	case "book-10":
//		depth = 10
//	case "book-25":
//		depth = 25
//	case "book-100":
//		depth = 100
//	case "book-500":
//		depth = 500
//	case "book-1000":
//		depth = 1000
//	default:
//		depth = 10
//	}
//
//	return &FullMarketOrderBook{
//		InstrumentName: instrumentName,
//		OrderSide: [2]*avltree.Tree{
//			avltree.NewWith(utils.Float64Comparator),
//			avltree.NewWith(utils.Float64Comparator)},
//		depth: depth,
//	}
//}
