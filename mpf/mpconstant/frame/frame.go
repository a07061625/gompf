/**
 * 框架常量
 * User: 姜伟
 * Date: 2019/12/24 0024
 * Time: 15:39
 */
package frame

const (
    MiddlewareEventRequestBefore  = 1
    MiddlewareEventRequestAfter   = 2
    MiddlewareEventRouteBefore    = 3
    MiddlewareEventRouteAfter     = 4
    MiddlewareEventActionBefore   = 5
    MiddlewareEventActionAfter    = 6
    MiddlewareEventResponseBefore = 7
    MiddlewareEventResponseAfter  = 8
)

var (
    TotalMiddlewareEvent map[int]string
)

func init() {
    TotalMiddlewareEvent = make(map[int]string)
    TotalMiddlewareEvent[MiddlewareEventRequestBefore] = "request_before"
    TotalMiddlewareEvent[MiddlewareEventRequestAfter] = "request_after"
    TotalMiddlewareEvent[MiddlewareEventRouteBefore] = "route_before"
    TotalMiddlewareEvent[MiddlewareEventRouteAfter] = "route_after"
    TotalMiddlewareEvent[MiddlewareEventActionBefore] = "action_before"
    TotalMiddlewareEvent[MiddlewareEventActionAfter] = "action_after"
    TotalMiddlewareEvent[MiddlewareEventResponseBefore] = "response_before"
    TotalMiddlewareEvent[MiddlewareEventResponseAfter] = "response_after"
}
