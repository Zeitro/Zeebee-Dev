/**
 * 格式化/验证相关方法或变量
 * */

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
// 最大贷款期限（月）
export const maxTerm = 10000;

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
	// 移除所有非数字字符
	const cleaned = phoneNumber.replace(/\D/g, '');

	// 截取前10位数字
	const truncated = cleaned.substring(0, 10);

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
	const { first_name, last_name } = userInfo;
	const name1 = first_name?.[0]?.toUpperCase();
	const name2 = last_name?.[0]?.toUpperCase();
	return `${name1}${name2}`;
};

// 首字母大写，有空格时，空格后的字母大写
export function capitalizeFirstLetter(str) {
	return str
		.split(/(?:-|\s)+/)
		.map((word) => word.charAt(0).toUpperCase() + word.slice(1))
		.join(' ');
}
