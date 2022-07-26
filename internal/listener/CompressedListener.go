package listener

//func TextListener(
//	maxConnections int,
//	urlAsText string) fx.Option {
//	const TextListenerConnection = "TextListenerConnection"
//	return fx.Options(
//		fx.Provide(
//			fx.Annotated{
//				Group: "Apps",
//				Target: func(params struct {
//					fx.In
//					PubSub             *pubsub.PubSub `name:"Application"`
//					GoFunctionCounter  GoFunctionCounter.IService
//					NetAppFuncInParams common.NetAppFuncInParams
//				}) (messages.CreateAppCallback, error) {
//					textListenerUrl, err := url.Parse(urlAsText)
//					if err != nil {
//						return messages.CreateAppCallback{}, err
//					}
//					f := goCommsNetListener.NewNetListenApp(
//						TextListenerConnection,
//						TextListenerConnection,
//						false, nil,
//						textListenerUrl,
//						//goCommsDefinitions.TransportFactoryEmptyName,
//						common.MaxConnectionsSetting(maxConnections),
//						common.NewConnectionInstanceOptions(
//							goCommsDefinitions.ProvideTransportFactoryForEmptyName(
//								topStack.ProvideTopStack(),
//								bottom.Provide(),
//							),
//							PubSub.ProvidePubSubInstance("Application", params.PubSub),
//							fx.Provide(
//								fx.Annotated{
//									Target: func() GoFunctionCounter.IService {
//										return params.GoFunctionCounter
//									},
//								},
//							),
//							ProvideConnectionReactorFactory(),
//						),
//					)
//					return f(params.NetAppFuncInParams), nil
//				},
//			}),
//	)
//}
//
//func ProvideConnectionReactorFactory() fx.Option {
//	return fx.Options(
//		fx.Provide(
//			fx.Annotated{
//				Target: func(
//					params struct {
//						fx.In
//						CancelCtx              context.Context
//						CancelFunc             context.CancelFunc
//						ConnectionCancelFunc   model.ConnectionCancelFunc
//						Logger                 *zap.Logger
//						PubSub                 *pubsub.PubSub `name:"Application"`
//						UniqueReferenceService interfaces.IUniqueReferenceService
//						GoFunctionCounter      GoFunctionCounter.IService
//						FullMarketDataHelper   fullMarketDataHelper.IFullMarketDataHelper
//						FmdService             fullMarketDataManagerService.IFmdManagerService
//					},
//				) (intf.IConnectionReactor, error) {
//					return NewReactor(
//							params.Logger,
//							params.CancelCtx,
//							params.CancelFunc,
//							params.ConnectionCancelFunc,
//							func(m proto.Message) (goprotoextra.IReadWriterSize, error) {
//								bytes, err := json.Marshal(m)
//								if err != nil {
//									return nil, err
//								}
//								return gomessageblock.NewReaderWriterBlock(bytes), nil
//							},
//							params.PubSub,
//							params.UniqueReferenceService,
//							params.GoFunctionCounter,
//							params.FullMarketDataHelper,
//							params.FmdService,
//						),
//						nil
//				},
//			},
//		),
//	)
//}
