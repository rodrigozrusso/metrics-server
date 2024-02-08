import { useEffect, useState, ChangeEvent } from 'react'
import { Select, SelectOption } from '@/components/select'

export enum GranularityValues {
  MINUTE = "minute",
  HOURLY = "hourly",
  DAILY = "daily",
}

export default function GranularitySelector() {

  const [loading, setLoading] = useState(true);
  const [returnMsg, setReturnMsg] = useState("");
  const [isResultSuccess, setIsResultSuccess] = useState(true);

  const [options, setOptions] = useState<SelectOption[]>([]);
  const [value, setValue] = useState('');

  const onChange = (event: ChangeEvent<HTMLSelectElement>) => {
    setValue(event.target.value);
  };

  const createGranularityOptions = () => {
    const options: SelectOption[] = [
      // { label: 'Select...', value: '' },
      ...Object.values(GranularityValues).map((granularity) => ({ label: granularity, value: granularity })),
    ];
    return options;

  }

  useEffect(() => {
    const fetchData = async () => {
      setOptions(createGranularityOptions());
      setLoading(false);
    };

    fetchData();
  }, []);


  return (
    <div className="w-full">
      <label htmlFor="granularity" className="mb-2 text-sm font-medium text-gray-900 dark:text-white">Granularity</label>
      {loading ? (
        <div>Loading...</div>
      ) : isResultSuccess ? (
        <div>
          <Select
            id="granularity"
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


