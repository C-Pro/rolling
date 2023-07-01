# Rolling

This go package provides rolling window aggregation container that calculates basic statistical metrics for a float64 time series data.

The idea is to calculate metrics on the fly without O(n) aggregation every time a new data point is added to the time series.

We still have to keep all values in memory, but we don't have to iterate over all of them every time we want to calculate a metric.
