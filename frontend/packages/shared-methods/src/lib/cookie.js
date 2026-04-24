// 设置cookie
export const setCookie = (name, value, time) => {
	const Days = 30;
	const exp = new Date();
	exp.setTime(exp.getTime() + (time || Days * 24 * 60 * 60 * 1000));
	document.cookie = `${name}=${encodeURIComponent(value)};expires=${exp.toUTCString()};path=/`;
};

// 获取cookie
export const getCookie = (name) => {
	const reg = new RegExp('(^| )' + name + '=([^;]*)(;|$)');
	const arr = document.cookie.match(reg);
	return arr ? decodeURIComponent(arr[2]) : null;
};

// 删除指定cookie
export const delCookie = (name) => {
	document.cookie = `${name}=;expires=Thu, 01 Jan 1970 00:00:00 GMT;path=/`;
};

// 清除所有cookie（兼容所有路径）
export const clearAllCookie = () => {
	const cookies = document.cookie ? document.cookie.split('; ') : [];
	for (let i = 0; i < cookies.length; i++) {
		const eqPos = cookies[i].indexOf('=');
		const name = eqPos > -1 ? cookies[i].substring(0, eqPos) : cookies[i];
		// 删除根路径
		document.cookie = `${name}=;expires=Thu, 01 Jan 1970 00:00:00 GMT;path=/`;
		// 删除当前路径
		document.cookie = `${name}=;expires=Thu, 01 Jan 1970 00:00:00 GMT;`;
	}
};
