import { Col, Form, Input } from 'antd';
import { titleCol } from '@/enum/common';
import { getPathArr } from '@/lib/common';
import { validatedName, capitalizeFirstLetter } from '@/lib/format';
// 姓名组合（firstName + middleName + lastName）
export const InputNameCombination = ({
	layoutCol = {
		xs: 24,
		sm: 24,
		md: 8,
	},
	firstName = 'first_name',
	middleName = 'mid_name',
	lastName = 'last_name',
	firstLabel = 'First name',
	middleLabel = 'Middle name',
	lastLabel = 'Last name',
	title,
	style = { paddingBottom: 12, fontSize: 16 },
	config = {},
	// form,
	listName,
	disabled,
}) => {
	const form = Form.useFormInstance();
	const first = Form.useWatch(getPathArr(firstName, listName), form);
	const last = Form.useWatch(getPathArr(lastName, listName), form);
	const applicants = Form.useWatch(listName, form);
	const index = firstName?.[0] || 0;
	const getDependencies = (deleteName) => {
		const obj = {
			firstName: getPathArr(firstName, listName),
			middleName: getPathArr(middleName, listName),
			lastName: getPathArr(lastName, listName),
		};
		delete obj[deleteName];
		return Object.values(obj);
	};

	return (
		<>
			{title && (
				<Col {...titleCol} className="font-medium text-[#000000a6]" style={style}>
					{title}
				</Col>
			)}
			<Col {...layoutCol}>
				<Form.Item
					label={firstLabel}
					name={firstName}
					rules={[
						{ required: last || config?.firstNameRequired, message: 'Please Enter!' },
						{ validator: (_, value) => validatedName(value, index, applicants) },
					]}
					getValueFromEvent={(v) => {
						return capitalizeFirstLetter(v.target.value);
					}}
					dependencies={getDependencies('firstName')}
				>
					<Input placeholder={`Enter ${firstLabel.toLowerCase()}`} disabled={disabled} />
				</Form.Item>
			</Col>
			<Col {...layoutCol}>
				<Form.Item
					label={middleLabel}
					name={middleName}
					rules={[{ validator: (_, value) => validatedName(value, index, applicants) }]}
					dependencies={getDependencies('middleName')}
					getValueFromEvent={(v) => {
						return capitalizeFirstLetter(v.target.value);
					}}
				>
					<Input placeholder={`Enter ${middleLabel.toLowerCase()}`} disabled={disabled} />
				</Form.Item>
			</Col>
			<Col {...layoutCol}>
				<Form.Item
					label={lastLabel}
					name={lastName}
					rules={[
						{ required: first || config?.lastNameRequired, message: 'Please Enter!' },
						{ validator: (_, value) => validatedName(value, index, applicants) },
					]}
					dependencies={getDependencies('lastName')}
					getValueFromEvent={(v) => {
						return capitalizeFirstLetter(v.target.value);
					}}
				>
					<Input placeholder={`Enter ${lastLabel.toLowerCase()}`} disabled={disabled} />
				</Form.Item>
			</Col>
		</>
	);
};
