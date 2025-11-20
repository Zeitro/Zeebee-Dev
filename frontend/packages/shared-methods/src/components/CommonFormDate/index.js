import { DatePicker } from 'antd';
import { formatStr } from '@repo/shared-methods/utils';
import dayjs from 'dayjs';
export const CommonFormDate = ({
	value,
	onChange,
	format = formatStr,
	type = 'DatePicker',
	...rest
}) => {
	const handleChange = (date, dateString) => {
		onChange(dateString);
	};

	const props = {
		value: value ? dayjs(value) : null,
		onChange: handleChange,
		format,
		placeholder: 'Select a date',
		...rest,
		style: { width: '100%', ...(rest?.style || {}) },
	};

	const getContent = () => {
		let content = '';
		switch (type) {
			case 'DatePicker':
				content = <DatePicker {...props} />;
				break;

			default:
				break;
		}
		return content;
	};

	return getContent();
};
