// Go port of Coda Hale's Metrics library
//
// <https://yougam/libraries/rcrowley/go-metrics>
//
// Coda Hale's original work: <https://yougam/libraries/codahale/metrics>
package metrics

// UseNilMetrics is checked by the constructor functions for all of the
// standard metrics.  If it is true, the metric returned is a stub.
//
// This global kill-switch helps quantify the observer effect and makes
// for less cluttered pprof profiles.
var UseNilMetrics bool = false
