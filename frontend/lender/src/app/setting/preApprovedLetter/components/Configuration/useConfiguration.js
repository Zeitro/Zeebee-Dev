import { formatAmount, formatPhone, isNumberFn } from '@repo/shared-methods/utils';
import _ from 'lodash';
const useConfiguration = ({ valueHtml, setValueHtml, data }) => {
	const formatValue = (format, value) => {
		switch (format) {
			case 'Currency': {
				// 当选中default项时，默认展示的'{amount}',所以需要单独处理
				const isNumber = isNumberFn(value);
				if (!isNumber) {
					return value;
				}
				return formatAmount(value);
			}
			case 'Phone':
				return formatPhone(value);
			default:
				break;
		}
		return value;
	};
	// 更新HTML中的变量
	const updateVariablesInHtml = (html, value, idKey, format, isHave = true) => {
		// 创建一个临时DOM元素来处理HTML
		const tempDiv = document.createElement('div');
		tempDiv.innerHTML = html;

		// 查找所有带有 data-variable 属性的元素
		const variableElements = tempDiv.querySelectorAll('[data-variable]');

		variableElements.forEach((element) => {
			const variableName = element.getAttribute('data-variable');

			// 检查是否在表单数据中有对应的值
			if (idKey === variableName) {
				// 更新元素的文本内容
				element.textContent = formatValue(format, value);
			}
		});

		// 更新预览内容
		if (isHave) {
			setValueHtml(tempDiv.innerHTML);
		}
		return tempDiv.innerHTML;
	};
	// 输入框修改更改html
	const handleFieldsChange = (changedValues, allValues) => {
		const formNameArr = _.cloneDeep(changedValues?.[0]?.name);
		let itemNameArr = formNameArr.slice(0, -1);
		if (['fixed_value', 'dynamicValue'].includes(formNameArr[formNameArr.length - 1])) {
			formNameArr[formNameArr.length - 1] = 'key';
		}
		const idKey = _.get(data, formNameArr);
		const item = _.get(data, itemNameArr);
		const changeItem = allValues.find((i) => _.isEqual(i.name, changedValues?.[0]?.name));
		updateVariablesInHtml(valueHtml, changeItem?.value, idKey, item.format);
	};

	return {
		handleFieldsChange,
		updateVariablesInHtml,
	};
};

export default useConfiguration;
