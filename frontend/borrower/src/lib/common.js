// Regular expression for validating phone numbers in various formats
export const phoneRegex = /^[+]?[(]?[0-9]{3}[)]?[-\s.]?[0-9]{3}[-\s.]?[0-9]{4,6}$/;

// Regular expression for validating email addresses
export const emailRegex = /^\w+([.-]?\w+)*@\w+([.-]?\w+)*(\.\w{2,})+$/;

// 分页
export const getPagination = (pagi = {}) => {
	return {
		current: 1,
		pageSize: 10,
		total: 0,
		...pagi,
	};
};
/**
 * 将数组转为指定的键值对结构
 * sourceArray：源数组
 * labelField: 数组中原键显示的字段
 * valueField: 数组中原值的字段
 * outputText：输出的键
 * outputKey：输出的值
 * newFieldArr:新增的字段数组，与sourceArray顺序--对应
 * */
export const transformArrayStructure = (
	{
		sourceArray,
		entryKey = 'label',
		entryValue = 'value',
		outputKey = 'label',
		outputValue = 'value',
	},
	newFieldArr = [],
) => {
	return sourceArray.map((item, index) => {
		let obj = {};
		if (item[entryKey] !== undefined && item[entryValue] !== undefined) {
			obj[outputKey] = item[entryKey];
			obj[outputValue] = item[entryValue];
			const newField = newFieldArr?.[index] || {};
			obj = {
				...obj,
				...newField,
			};
		}
		return obj;
	});
};

// 格式化金额，千分位加逗号，可最多保留2位小数
export function formatAmount(amount) {
	if (!amount) return '$ 0';
	if (typeof amount === 'string') {
		amount = Number(amount);
	}

	let [integerPart, decimalPart = ''] = amount.toString().split('.');
	if (decimalPart.length > 2) {
		decimalPart = decimalPart.slice(0, 2);
	}
	// 格式化整数部分千分位
	integerPart = parseInt(integerPart, 10).toLocaleString('en-US');

	if (decimalPart.length === 0 || /^0*$/.test(decimalPart)) {
		// 没有有效小数部分，返回整数部分
		return `$ ${integerPart}`;
	}
	return `$ ${integerPart}.${decimalPart}`;
}

// 拼装地址
export const assemblyAddress = (address = {}, isHaveZip) => {
	const arr = [address?.address1,address?.city, address?.state];
	if (isHaveZip && address?.zip_code){
		arr.push(address?.zip_code);
	}
	const newArr = arr.filter(
		(i) => i,
	);
	return [...new Set(newArr)].join(',');
};
