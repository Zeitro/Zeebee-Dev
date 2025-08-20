/**
 * 格式化/验证相关方法或变量
 * */

// Regular expression for validating phone numbers in various formats
export const phoneRegex = /^[+]?[(]?[0-9]{3}[)]?[-\s.]?[0-9]{3}[-\s.]?[0-9]{4,6}$/;

// export const emailRegex = /a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;

export const emailRegex = /^\w+([.-]?\w+)*@\w+([.-]?\w+)*(\.\w{2,})+$/;

// 信用分（300-850分）
export const creditScoreRegex = /^([3-7]\d{2}|8[0-4]\d|850)(-\d{1,2})?$/;
export const minCreditScore = 300;
export const maxCreditScore = 850;

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
	const arr = [address?.address1, address?.city, address?.state];
	if (isHaveZip && address?.zip_code) {
		arr.push(address?.zip_code);
	}
	const newArr = arr.filter((i) => i);
	return [...new Set(newArr)].join(',');
};
