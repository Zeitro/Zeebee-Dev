/**
 * 格式化/验证相关方法或变量
 * */
// 禁止以空格、中点、连字符或下划线开头,禁止以这些特殊符号结尾,防止特殊符号连续出现,允许大小写字母、带重音符号的字母、空格、中点、连字符和撇号
export const nameRegex = /^(?![\s·\-_])(?!.*[\s·\-_]$)(?!.*[\s·\-_]{2})[A-Za-zÀ-ÿ\s·\-']+$/;
// Regular expression for validating phone numbers in various formats
export const phoneRegex = /^[+]?[(]?[0-9]{3}[)]?[-\s.]?[0-9]{3}[-\s.]?[0-9]{4,6}$/;
// 同时有大小写字母、数字至少8位以上
export const passwordRegex = /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d).{8,}$/;
// export const emailRegex = /a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;

export const emailRegex = /^\w+([.-]?\w+)*@\w+([.-]?\w+)*(\.\w{2,})+$/;
export const tinOrSSNRegex = /^\d{9}$/;

// 信用分（300-850分）
export const creditScoreRegex = /^([3-7]\d{2}|8[0-4]\d|850)(-\d{1,2})?$/;
export const minCreditScore = 300;
export const maxCreditScore = 850;

// 系统最大金额
export const maxAmount = 1000000000000;
// 最大单元数
export const maxNumberOfUnits = 100;
// 系统最大贷款期限（月）
export const maxTerm = 10000;
// ZeeBee,最小/最大贷款期限（月）
export const maxTerm_zeeBee = 24;
export const minTerm_zeeBee = 12;

// 通过property_type获取unit
export const changeUnitByPropertyType = (property_type) => {
	let unit = null,
		minUnit = 1,
		maxUnit = maxNumberOfUnits,
		disabled = false;
	switch (property_type) {
		case 'Single Family Residence':
			unit = minUnit = 1;
			disabled = true;
			break;
		case '2 unit':
			unit = minUnit = 2;
			disabled = true;
			break;
		case '3 unit':
			unit = minUnit = 3;
			disabled = true;
			break;
		case '4 unit':
			unit = minUnit = 4;
			disabled = true;
			break;
		case '4+ Multi':
			unit = minUnit = 5;
			break;
		default:
			break;
	}
	return { unit, minUnit, maxUnit, disabled };
};

// 格式化居住面积, 用于<InputNumber>中配置
export const livingAreaProps = {
	controls: false,
	min: 1,
	formatter: (value) => {
		// 保留两位小数
		const [int, dec] = value.split('.');
		const str = dec?.slice?.(0, 2);
		return `${int}` + (dec ? `.${str}` : '');
	},
	parser: (value) => parseFloat(value.replace(/[^\d.]/g, '')) || 0,
};

// 校验规则, 用于<Form.Item>rules中配置validator
export const validateRule = (value, regex, errorMsg = 'Invalid input') => {
	if (value && !regex.test(value)) {
		return Promise.reject(errorMsg);
	}
	return Promise.resolve();
};

//AssignmentFee 金额效验
export const validateAssignmentFeeAmount = (value, purchase_amount) => {
	// Assignment fee 最高只能达到purchase price的10%
	const max = purchase_amount * 0.1;
	if (!value) {
		return Promise.reject('Please enter assignment fee amount!');
	}
	if (value && value > max) {
		return Promise.reject(`Assignment fee cannot exceed purchase price's 10% ($${max})`);
	}
	return Promise.resolve();
};

export const validatedName = (value, index, list) => {
	// 有多个姓名时，比较是否相同
	if (list?.length) {
		const currentItem = list[index];
		const currentName =
			currentItem.first_name + (currentItem.mid_name || '') + currentItem.last_name;
		let indexArr = [index];
		list?.map((item, i) => {
			const name = item.first_name + (item.mid_name || '') + item.last_name;
			if (name === currentName && index !== i) {
				indexArr.push(i);
			}
			return item;
		}, {});
		if (indexArr?.length > 1) {
			const err = `${indexArr.map((item, i) => `Borrower ${i + 1} ${indexArr.length === i + 1 ? '' : 'and'}`).join(' ')} cannot be the same`;
			return Promise.reject(err);
		}
	}
	// 只有一个姓名时，校验格式
	// Todo
	return Promise.resolve();
};

// 格式化手机号 (xxx) xxx-xxxx
export const formatPhone = (phoneNumber) => {
	if (!phoneNumber) return '';
	// 移除所有非数字字符
	const cleaned = phoneNumber?.replace(/\D/g, '');

	// 截取前10位数字
	const truncated = cleaned?.substring(0, 10);

	// 根据输入长度格式化显示
	if (truncated.length <= 3) {
		return truncated;
	} else if (truncated.length <= 6) {
		return `(${truncated.substring(0, 3)}) ${truncated.substring(3)}`;
	} else {
		return `(${truncated.substring(0, 3)}) ${truncated.substring(3, 6)}-${truncated.substring(6)}`;
	}
};

// 格式化金额，千分位加逗号，可最多保留2位小数
export function formatAmount(amount) {
	if (!amount) return '$ --';
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

// 格式化TIN,包括SSN/ITIN/EIN, 不同类型税务ID
export const formatTINOrSSN = (taxId, type = 'SSN') => {
	const cleaned = taxId.replace(/\D/g, '');
	switch (type) {
		case 'SSN':
		case 'ITIN':
			// ITIN与SSN格式一致(公民): XXX-XX-XXXX (9位数字，只是ITIN第一位数是9开头的，表示非美国公民)
			if (cleaned.length <= 3) return cleaned;
			if (cleaned.length <= 5) return `${cleaned.slice(0, 3)}-${cleaned.slice(3)}`;
			return `${cleaned.slice(0, 3)}-${cleaned.slice(3, 5)}-${cleaned.slice(5, 9)}`;
		case 'EIN':
		case 'TIN':
			// EIN格式(企业): XX-XXXXXXX (9位数字)
			if (cleaned.length <= 2) return cleaned;
			return `${cleaned.slice(0, 2)}-${cleaned.slice(2, 9)}`;
		default:
			return cleaned;
	}
};

// 拼装地址
export const assemblyAddress = (address = {}, isHaveZip) => {
	const arr = [address?.address1, address?.city, address?.state];
	if (isHaveZip && address?.zip_code) {
		arr.push(address?.zip_code);
	}
	const newArr = arr.filter((i) => i);
	return [...new Set(newArr)].join(',');
};

// 简写姓名：只显示姓名首字母，并大写
export const formatFirstLetterName = (userInfo) => {
	if (!userInfo) return '';
	const { first_name, last_name } = userInfo;
	const name1 = first_name?.[0]?.toUpperCase();
	const name2 = last_name?.[0]?.toUpperCase();
	return `${name1}${name2}`;
};

/**
 * 格式化名字（分隔符只有 空格, -, ·, ', ’），首字母大写，分隔符后字母大写
 * @param {*} value 名字
 * @returns 符合禁止以空格、中点、连字符或下划线开头,防止特殊符号连续出现,允许大小写字母、带重音符号的字母、空格、中点、连字符和撇号
 * Input输入值，做不到禁止以这些特殊符号结尾
 */
export function capitalizeFirstLetter(value) {
	if (!value && value !== 0) return '';
	let s = String(value);

	// 1. 移除下划线
	s = s.replace(/_/g, '');

	// 2. 删除非法字符（保留拉丁字母、空格、中点、连字符和撇号）
	// 包含常见的弯撇号 ’
	s = s.replace(/[^A-Za-zÀ-ÖØ-öø-ÿ\s·\-'’]/g, '');

	// 3. 折叠同类连续符号
	s = s
		.replace(/·{2,}/g, '·')
		.replace(/-{2,}/g, '-')
		.replace(/'{2,}/g, "'")
		.replace(/’{2,}/g, '’')
		.replace(/\s{2,}/g, ' ');

	// 4. 折叠混合连续的分隔符为第一个字符（避免 "- ' " 这类混合）
	s = s.replace(/([ \-·'’]){2,}/g, (m, p1) => p1);

	// 5. 去除开头的分隔符
	if (/([ \-·'’])+/g.test(s[0])) {
		s = s.slice(1);
	}

	// 6. 首字母大写逻辑：遇到分隔符后下一个字母大写，其余小写
	const separators = new Set([' ', '-', '·', "'", '’']);
	let needUpper = true;
	let out = '';
	for (const ch of s) {
		if (separators.has(ch)) {
			out += ch;
			needUpper = true;
		} else {
			out += needUpper ? ch.toUpperCase() : ch.toLowerCase();
			needUpper = false;
		}
	}
	return out;
}

/**
 * 最终严格格式化名字（用于 Input.onBlur 或提交前）
 */
export function finalizeName(name) {
	let str = capitalizeFirstLetter(name);
	// 去除开头和结尾的分隔符
	str = str.replace(/^[ \-·'’]+|[ \-·'’]+$/g, '');
	return str.trim();
}
