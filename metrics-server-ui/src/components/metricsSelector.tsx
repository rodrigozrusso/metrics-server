import { useEffect, useState, ChangeEvent } from 'react'
import { Select, SelectOption } from '@/components/select'

export default function MetricsSelector() {

  const [loading, setLoading] = useState(true);
  const [returnMsg, setReturnMsg] = useState("");
  const [isResultSuccess, setIsResultSuccess] = useState(true);

  const [options, setOptions] = useState<SelectOption[]>([]);
  const [value, setValue] = useState('');

  const onChange = (event: ChangeEvent<HTMLSelectElement>) => {
    setValue(event.target.value);
  };

  const getMetricList = async () => {
    const metricServerURL = process.env.METRIC_SERVER_URL || 'http://localhost:8080';
    const metricsURL = `${metricServerURL}/v1/metrics/`
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

  const mapMetricList = async (metricList: string[]) => {
    const options: SelectOption[] = [
      // { label: 'Select...', value: '' },
      ...metricList.map((metricName) => ({ label: metricName, value: metricName })),
    ];
    return options;
  }

  useEffect(() => {
    const fetchData = async () => {
      const response = await getMetricList();
      setOptions(await mapMetricList(response as any));
      setLoading(false);
    };

    fetchData();
  }, []);


  return (
    <div className="w-full">
      <label htmlFor="metricName" className="mb-2 text-sm font-medium text-gray-900 dark:text-white">Metric</label>
      {loading ? (
        <div>Loading...</div>
      ) : isResultSuccess ? (
        <div>
          <Select
            id="metricName"
            value={value}
            options={options}
            onChange={onChange}
          />
        </div>
      ) : (
        <div>
          <h1>Failed</h1>
          <p>{returnMsg}</p>
        </div>
      )}

    </div>
  );
}


