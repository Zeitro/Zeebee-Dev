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
