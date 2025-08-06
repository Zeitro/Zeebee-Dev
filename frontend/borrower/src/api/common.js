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
