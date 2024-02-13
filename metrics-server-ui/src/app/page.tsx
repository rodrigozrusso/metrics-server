'use client'
import { useEffect, useState } from 'react'
import { LineChart } from '@carbon/charts-react'
import type { ChartTabularData, LineChartOptions } from '@carbon/charts'
import DatePickerRange from '@/components/datepicker'
import MetricsSelector from '@/components/metricsSelector'
import GranularitySelector, { GranularityValues } from '@/components/granularitySelector'


const granularityBackendMap = {
  [GranularityValues.MINUTE]: "minute",
  [GranularityValues.HOURLY]: "hour",
  [GranularityValues.DAILY]: "day",
}

export default function Home() {
  const [loading, setLoading] = useState(true);
  const [isResultSuccess, setIsResultSuccess] = useState(true);
  const [returnMsg, setReturnMsg] = useState("");
  const [data, setData] = useState()

  const defaultGranularity = GranularityValues.MINUTE;

  const handleSubmit = async (values: any) => {
    const granularity = values.granularity || (document.getElementById('granularity') as HTMLInputElement)?.value;
    const dateRange = values.dateRange || (document.getElementById('date-range') as HTMLInputElement)?.value;
    const data = {
      granularity: granularityBackendMap[granularity as GranularityValues],
      metricName: values.metricName || (document.getElementById('metricName') as HTMLInputElement)?.value,
      startDate: new Date(dateRange.split(" ~ ")[0]).toISOString().split("T")[0],
      endDate: new Date(dateRange.split(" ~ ")[1]).toISOString().split("T")[0]
    };

    const response = await getData(data.metricName, data.granularity, data.startDate, data.endDate);
    setData(mapData(data.metricName, response.data));
    options.timeScale.timeInterval = granularity;
  };

  const options = {
    title: 'Metric (time series)',
    axes: {
      bottom: {
        title: 'Time',
        mapsTo: 'date',
        scaleType: 'time',
      },
      left: {
        title: 'Avg',
        mapsTo: 'value',
        scaleType: 'linear'
      }
    },
    curve: 'curveMonotoneX',
    height: '400px',
    timeScale: {
      timeInterval: defaultGranularity,
    }

  }

  const mapData = (metricName: string, data: any) => {
    return data.map((d: any) => {
      return {
        group: metricName,
        date: new Date(d.timeFrame),
        value: d.avg,
      }
    })
  }

  const getData = async (metricName: string, granularity: string, startDate: string, endDate: string) => {
    const metricServerURL = process.env.METRICS_SERVER_URL || 'http://localhost:8080';
    const metricsURL = `${metricServerURL}/v1/metrics/${metricName}/${granularity}/${startDate}/${endDate}`
    const response = await fetch(metricsURL, {
      method: "GET",
    });

    if (response.status === 200) {
      return response.json();
    } else {
      const result = await response.json();
      setLoading(false);
      setIsResultSuccess(false);
      setReturnMsg(result.data);
    }
  };

  useEffect(() => {
    const fetchData = async () => {
      const metricName = "metric_a";
      const response = await getData(metricName, "hour", "2024-01-01", "2024-02-05");
      setData(mapData(metricName, response.data));
      setLoading(false);
    };

    fetchData();
  }, []);

  return (
    <div className="App">
      {/* filter */}
      <section className="bg-white dark:bg-gray-900">
        <div className="py-8 px-4 mx-auto max-w-2xl lg:py-8">
          <form>
            <div className="grid gap-4 sm:grid-cols-2 sm:gap-6">
              <DatePickerRange />
              <GranularitySelector />
              <MetricsSelector />
              <div className="-mx-3 mt-6 flex flex-wrap">
                <div className="w-full px-3">
                  <button
                    type="button"
                    onClick={handleSubmit}
                    className="btn w-full bg-blue-600 text-white hover:bg-blue-700"
                  >
                    Filter
                  </button>
                </div>
              </div>
            </div>
          </form>
        </div>
      </section>

      {/* chart */}
      {loading ? (
        <div>Loading chart...</div>
      ) : isResultSuccess ? (
        <div>
          <LineChart data={data as unknown as ChartTabularData} options={options as unknown as LineChartOptions} />
        </div>
      ) : (
        <div>
          <h1>Failed</h1>
          <p>{returnMsg}</p>
        </div>
      )}

    </div>);
}
