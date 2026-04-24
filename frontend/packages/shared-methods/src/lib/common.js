import _ from 'lodash';
import { matchRoute } from './commonRouter';

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

/**
 * 拼装key，可用于表单中，例如useWatch传入的name可以是字符串，也可以是数组
 * @param {*string | Array} name Form.Item中的name
 * @param {*string | Array} listName 当Form.Item用于Form.List时，Form.List中的name
 * @returns {*Array} 返回拼装的
 */
export const getPathArr = (name, listName) => {
	const listNameArr = Array.isArray(listName) ? listName : [listName];
	let pathArr = Array.isArray(name) ? name : [name];
	if (listName) {
		pathArr = [...listNameArr, ...pathArr];
	}
	return pathArr;
};

// 将多层嵌套数组扁平化(多层嵌套变为一维数组，不需要再递归处理)
export function flattenArray(list = [], childKey = 'children', isHaveChildParent) {
	return list?.reduce((acc, next) => {
		if (next?.[childKey]) {
			const childArr = flattenArray(next?.[childKey]);
			if (isHaveChildParent) {
				return [...acc, next, ...childArr];
			}
			return [...acc, ...childArr];
		}
		return [...acc, next];
	}, []);
}

/**
 * 可匹配动态路由
 * dynamicPath 动态路由/loans/123455
 * path 比较的固定路由/loans/:id
 * */
// export const matchRoute = (path, dynamicPath) => {
// 	if (path === dynamicPath) return true;
// 	const pattern = '^' + path.replace(/:[^/]+/g, '[^/]+') + '$';
// 	return new RegExp(pattern).test(dynamicPath);
// };

// 递归查找,返回查找到的某一项数据
export const findItem = ({
	list,
	key,
	keyName = 'key',
	childName = 'children',
	isPushIndex,
	isPushParent,
	parentKey,
	parentIndex,
}) => {
	let result = null;
	list?.map?.((i, index) => {
		if (isPushIndex) {
			i.index = index;
		}
		if (isPushParent) {
			i.parentKey = parentKey || null;
			i.parentIndex = parentIndex || null;
		}
		if (matchRoute(i?.[keyName], key)) {
			result = i;
		}
		if (i?.[childName]) {
			const childResult = findItem({
				list: i?.[childName],
				key,
				keyName,
				childName,
				isPushIndex,
				isPushParent,
				parentKey: i?.[keyName],
				parentIndex: index,
			});
			if (childResult) {
				result = childResult;
			}
		}
	});
	return result;
};

// 获取文件类型通过文件后缀
export const getFileTypeBySuffix = (suffix) => {
	const newSuffix = suffix?.toLowerCase();
	switch (newSuffix) {
		case 'jpg':
		case 'jpeg':
		case 'png':
			return 'image';
		case 'pdf':
			return 'pdf';
		case 'doc':
		case 'docx':
			return 'doc';
		default:
			return '';
	}
};

// 简单的复制到剪贴板函数
export const copyToClipboard = async (text) => {
	try {
		if (navigator.clipboard) {
			await navigator.clipboard.writeText(text);
			return true;
		} else {
			// 降级处理
			const textArea = document.createElement('textarea');
			textArea.value = text;
			document.body.appendChild(textArea);
			textArea.select();
			document.execCommand('copy');
			document.body.removeChild(textArea);
			return true;
		}
	} catch {
		return false;
	}
};

// 完整的html中，只将link/style标签插入head标签内，用于引入资源
export const injectHeadResources = (html) => {
	// 清理旧的
	document.querySelectorAll('[data-template-preview]').forEach((el) => el.remove());
	const doc = new window.DOMParser().parseFromString(html, 'text/html');
	const headNodes = doc.head ? Array.from(doc.head.children) : [];
	headNodes.forEach((node) => {
		if (['LINK', 'STYLE'].includes(node.tagName) || node.tagName.startsWith('FONT')) {
			const clone = node.cloneNode(true);
			clone.setAttribute('data-template-preview', '1');
			document.head.appendChild(clone);
		}
	});
};

// html中取出body的内容
export const getBodyInnerHTML = (html) => {
	try {
		const doc = new window.DOMParser().parseFromString(html, 'text/html');
		return doc.body.innerHTML;
	} catch {
		return html;
	}
};

// 完整的html，替换body内容
export function replaceBodyContent(html, newBodyHtml) {
	const parser = new DOMParser();
	const doc = parser.parseFromString(html, 'text/html');
	doc.body.innerHTML = newBodyHtml;
	return doc.documentElement.outerHTML;
}

/**
 * 自定义 storage 适配器
 * 根据字段名决定存储到 localStorage 还是 cookie
 */
// export const hybridStorageFn = (COOKIE_FIELDS) => {
// 	return {
// 		getItem: (name) => {
// 			try {
// 				// 先从 cookie 读取
// 				const cookieValue = Cookies.get(name);
// 				if (cookieValue) {
// 					return cookieValue;
// 				}
// 				// 再从 localStorage 读取（仅限客户端）
// 				if (typeof window !== 'undefined' && window.localStorage) {
// 					const localStorageValue = window.localStorage.getItem(name);
// 					return localStorageValue;
// 				}
// 				return null;
// 			} catch (error) {
// 				console.error('Error reading from storage:', error);
// 				return null;
// 			}
// 		},

// 		setItem: (name, value) => {
// 			try {
// 				const parsed = JSON.parse(value);

// 				// 检查是否包含需要存 cookie 的字段
// 				const hasCookieFields = Object.keys(parsed).some((key) =>
// 					COOKIE_FIELDS.includes(key),
// 				);

// 				if (hasCookieFields) {
// 					// 提取 cookie 字段
// 					const cookieData = {};
// 					const localStorageData = {};

// 					Object.entries(parsed).forEach(([key, val]) => {
// 						if (COOKIE_FIELDS.includes(key)) {
// 							cookieData[key] = val;
// 						} else {
// 							localStorageData[key] = val;
// 						}
// 					});
// 					// 动态设置 cookie 过期时间
// 					let expires = 7; // 默认7天
// 					if (cookieData.access && cookieData.access.access_token_expire) {
// 						expires = cookieData.access.access_token_expire / 86400;
// 					}
// 					Cookies.set(name, JSON.stringify(cookieData), {
// 						expires,
// 						secure: process.env.NODE_ENV === 'production',
// 						sameSite: 'strict',
// 					});

// 					// 存储到 localStorage
// 					if (Object.keys(localStorageData).length > 0) {
// 						if (typeof window !== 'undefined' && window.localStorage) {
// 							window.localStorage.setItem(name, JSON.stringify(localStorageData));
// 						}
// 					}
// 				} else {
// 					// 全部存 localStorage
// 					localStorage.setItem(name, value);
// 				}
// 			} catch (error) {
// 				console.error('Error writing to storage:', error);
// 			}
// 		},

// 		removeItem: (name) => {
// 			try {
// 				Cookies.remove(name);
// 				if (typeof window !== 'undefined' && window.localStorage) {
// 					window.localStorage.removeItem(name);
// 				}
// 			} catch (error) {
// 				console.error('Error removing from storage:', error);
// 			}
// 		},
// 	};
// };
