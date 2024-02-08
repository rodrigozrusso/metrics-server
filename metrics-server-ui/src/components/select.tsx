import * as React from 'react';
import { ChangeEvent } from 'react';

type SelectOption = {
  label: string;
  value: string;
};

type Props = {
  id?: string;
  value?: string;
  disabled?: boolean;
  className?: string;
  options: SelectOption[];
  onChange: (e: ChangeEvent<HTMLSelectElement>) => void;
};

const Select = ({
  id,
  value,
  disabled,
  className,
  options,
  onChange,
}: Props) => {
  const select = (
    <select id={id} className={className} disabled={disabled} onChange={onChange} value={value}>
      {options.map(({ value, label }) => (
        <option key={value} value={value}>
          {label}
        </option>
      ))}
    </select>
  );

  return <div>{select}</div>;
};

export { Select };
export type { SelectOption };