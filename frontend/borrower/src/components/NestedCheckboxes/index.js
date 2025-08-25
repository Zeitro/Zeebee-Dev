import React from 'react';
import { Checkbox, Col, Input, Row } from 'antd';
/**
 * 多层嵌套的多选框
 * @params paramsKey: 第一层checkbox的值对应的key（可与传给后端的key保持一致）
 */
const NestedCheckboxes = ({ value, onChange, options, paramsKey = 'list' }) => {
	/**
	 * 递归生成多层嵌套的多选框(并支持输入框和多选框的嵌套)
	 * @params list: 需要渲染的数据
	 * @params parentValue: 父级数据-用于判断父级是否勾选，子级才显示
	 * */
	const getCheckboxDom = (list, parentValue = []) => {
		return list.map((item, index) => {
			let checkDom = (nestedEvent) => (
				<Col span={24} key={index}>
					<Checkbox value={item.value} className={nestedEvent && '!pb-2'}>
						{item.label}
					</Checkbox>
					{nestedEvent}
				</Col>
			);
			if (item.extras?.eventType === 'input') {
				return checkDom(
					parentValue.includes(item.value) && (
						<Row gutter={[0, 16]}>
							<Col flex="24px"></Col>
							<Col flex="auto">
								<Input
									size="middle"
									className="w-full"
									value={value?.[item.extras?.paramsKey] || ''} //此处的paramsKey也可与后端的传值保持一致
									placeholder={item.extras?.placeholder}
								/>
							</Col>
						</Row>
					),
				);
			}
			if (item.extras?.eventType === 'checkbox') {
				return checkDom(
					parentValue.includes(item.value) && (
						<Checkbox.Group
							className="w-full !pl-6"
							value={value?.[item.extras?.paramsKey] || []}
							onChange={(e) =>
								onChange({
									...value,
									[item.extras?.paramsKey]: e,
								})
							}
						>
							<Row gutter={[0, 16]}>
								{getCheckboxDom(
									item.extras?.options,
									value?.[item.extras?.paramsKey],
								)}
							</Row>
						</Checkbox.Group>
					),
				);
			}
			return checkDom();
		});
	};

	return (
		<Checkbox.Group
			className="w-full"
			value={value?.[paramsKey] || []}
			onChange={(e) => {
				onChange({
					...value,
					[paramsKey]: e,
				});
			}}
		>
			<Row gutter={[0, 24]}>{getCheckboxDom(options, value?.[paramsKey])}</Row>
		</Checkbox.Group>
	);
};

export default NestedCheckboxes;
