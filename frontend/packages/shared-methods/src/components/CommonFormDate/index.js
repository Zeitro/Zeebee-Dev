import { DatePicker } from 'antd';
import { dayjs, formatStr } from '../../lib/date.js';

const { RangePicker } = DatePicker;

const toDayjs = (value, format) => {
	if (!value) return null;
	const parsed = dayjs(value, format, true);
	return parsed.isValid() ? parsed : null;
};

const toPickerValue = (value, format, isRange) => {
	if (!value) return null;
	if (isRange) {
		if (!Array.isArray(value) || value.length < 2) return null;
		const start = toDayjs(value[0], format);
		const end = toDayjs(value[1], format);
		if (!start || !end) return null;
		return [start, end];
	}
	return toDayjs(value, format);
};

export const CommonFormDate = ({
	value,
	onChange,
	format = formatStr,
	type = 'DatePicker',
	placeholder,
	...rest
}) => {
	const isRange = type === 'RangePicker';

	const handleChange = (dates, dateStrings) => {
		if (!onChange) return;
		if (isRange) {
			if (!dates || !dateStrings?.[0] || !dateStrings?.[1]) {
				onChange(null);
				return;
			}
			onChange(dateStrings);
			return;
		}
		onChange(dateStrings || null);
	};

	const pickerValue = toPickerValue(value, format, isRange);

	const sharedProps = {
		onChange: handleChange,
		format,
		allowClear: true,
		...rest,
		value: pickerValue,
		style: { width: '100%', ...(rest?.style || {}) },
	};

	if (isRange) {
		return (
			<RangePicker {...sharedProps} placeholder={placeholder ?? ['Start date', 'End date']} />
		);
	}

	return <DatePicker {...sharedProps} placeholder={placeholder ?? 'Select date'} />;
};
