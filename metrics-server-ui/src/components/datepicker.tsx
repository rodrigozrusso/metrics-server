import { useState } from 'react'
import Datepicker from "react-tailwindcss-datepicker";

export default function DatePickerRange() {
  const startDate = new Date();
  startDate.setDate(startDate.getDate() - 7);

  const [value, setValue] = useState({
    startDate: startDate,
    endDate: new Date()
  });

  const handleValueChange = (newValue: any) => {
    setValue(newValue);
  }

  return (
    <div className="w-full">
      <label htmlFor="granularity" className="mb-2 text-sm font-medium text-gray-900 dark:text-white">Date range</label>
      <Datepicker
        inputId='date-range'
        value={value}
        onChange={handleValueChange}
      />
    </div>
  );
}
