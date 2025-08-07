import { createRequest } from '@/lib/fetch';
const { request } = createRequest();
// 公共的下拉字段
export const getSelectOptions = (name) => {
	const params = new URLSearchParams({ name });
	return request(`/select_options?${params}`, {
		method: 'GET',
	});
};
// lender下拉
export const getLenders = () => {
	return request('/loan/lender_companies', {
		method: 'GET',
	});
};
// borrower下拉
export const getBorrowers = () => {
	return request('/applicant/borrower', {
		method: 'POST',
	});
};

// 获取地址state列表
export const getStateList = () => {
	return request('/us_states', {
		method: 'GET',
	});
};

// 获取地址city列表通过state
export const getCityList = (state) => {
	const params = new URLSearchParams({ state });
	return request(`/us_cities?${params}`, {
		method: 'GET',
	});
};