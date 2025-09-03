import React from 'react';
import { Col, Form, Input, Row } from 'antd';
import { titleCol } from '@/enum/common';
const InputName = ({
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
	gutter = 24,
	title,
	style = { paddingBottom: 12, fontSize: 16 },
}) => {
	return (
		<Row gutter={gutter}>
			{title && (
				<Col {...titleCol} className="font-medium text-[#000000a6]" style={style}>
					{title}
				</Col>
			)}
			<Col {...layoutCol}>
				<Form.Item
					label={firstLabel}
					name={firstName}
					rules={[{ required: true, message: 'Please Enter!' }]}
				>
					<Input placeholder={`Enter ${firstLabel.toLowerCase()}`} />
				</Form.Item>
			</Col>
			<Col {...layoutCol}>
				<Form.Item
					label={middleLabel}
					name={middleName}
					rules={[{ required: true, message: 'Please Enter!' }]}
				>
					<Input placeholder={`Enter ${middleLabel.toLowerCase()}`} />
				</Form.Item>
			</Col>
			<Col {...layoutCol}>
				<Form.Item
					label={lastLabel}
					name={lastName}
					rules={[{ required: true, message: 'Please Enter!' }]}
				>
					<Input placeholder={`Enter ${lastLabel.toLowerCase()}`} />
				</Form.Item>
			</Col>
		</Row>
	);
};
export default InputName;
