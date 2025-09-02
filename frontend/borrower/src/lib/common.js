import _ from 'lodash';
/**
 * 公共方法封装（格式化或验证相关写在format文件）
 * */

// 分页
export const getPagination = (pagi = {}) => {
	return {
		current: 1,
		pageSize: 10,
		total: 0,
		showTotal: (total) => `▶  Total ${total}`,
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

/**
 * 将多层嵌套对象转换为二维路径数组(拿到字段的键，可用于表单中)
 * @param obj 目标对象
 * @param parentPath 当前路径（递归内部使用）
 * @returns 二维路径数组，如 [['a', 'b'], ['a', 0, 'c']]
 */
export function flattenObjectPaths(obj, parentPath = []) {
	let paths = [];

	_.forOwn(obj, (value, key) => {
		const currentPath = [...parentPath, key];

		if (_.isPlainObject(value) || _.isArray(value)) {
			// 递归处理对象或数组
			paths.push(...flattenObjectPaths(value, currentPath));
		} else {
			// 叶子节点，直接添加路径
			paths.push(currentPath);
		}
	});

	return paths;
}
