import React from 'react';
import { Checkbox, Col, Input, Row } from 'antd';
/**
 * 多层嵌套的多选框
 * @params paramsKey: 第一层checkbox的值对应的key（可与传给后端的key保持一致）
 */
const NestedCheckboxes = ({ value, onChange, options, paramsKey = 'list', form, name, label }) => {
	// 递归删除,删除key对应的值，并删除所有子级的勾选和输入值
	const remove = (list, key, newValue) => {
		list.map((item) => {
			if (key && item.value === key) {
				delete newValue[item.extras?.paramsKey];
				if (item.extras?.options) {
					remove(item.extras?.options, null, newValue);
				}
			}
			if (!key) {
				delete newValue[item.extras?.paramsKey];
				if (item.extras?.options) {
					remove(item.extras?.options, null, newValue);
				}
			}
		});
		return newValue;
	};

	const validateStatus = (parentLabel, isCheck) => {
		if (!form || !name || !parentLabel) return;
		setTimeout(() => {
			form?.setFields([
				{
					name,
					errors: [
						isCheck
							? `Please check the corresponding options under '${parentLabel}'!`
							: `Please enter the corresponding content under '${parentLabel}'`,
					],
					validateStatus: 'error',
				},
			]);
		}, 500);
	};

	const handleChange = (v, paramsKey, options, isHave) => {
		let data = { ...value };
		// 当子级是数组时，父级取消勾选，子级清空
		if (options && v?.length < data?.[paramsKey]?.length) {
			const diff = data?.[paramsKey]?.find((i) => !v.includes(i));
			data = remove(options, diff, data);
		}
		// 有值赋值，无值或空数组清空
		if (isHave) {
			data[paramsKey] = v;
		} else {
			delete data[paramsKey];
		}
		onChange(data);
	};
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
									onChange={(e) =>
										handleChange(
											e.target.value,
											item.extras?.paramsKey,
											null,
											!!e.target.value,
										)
									}
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
							onChange={(e) => {
								handleChange(
									e,
									item.extras?.paramsKey,
									item.extras?.options,
									!!e.length,
								);
							}}
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
				handleChange(e, paramsKey, options, !!e.length);
			}}
		>
			<Row gutter={[0, 24]}>{getCheckboxDom(options, value?.[paramsKey])}</Row>
		</Checkbox.Group>
	);
};

export default NestedCheckboxes;
