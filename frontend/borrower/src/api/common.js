import { createRequest } from '@/lib/fetch';
const { request } = createRequest();
// 公共的下拉字段
export const getSelectOptions = (name) =>{
	const params = new URLSearchParams({ name })
	return request(`/select_options?${params}`, {
		method: 'GET',
	});
}
	
export const getLender = () =>{
	return request('/loan/lender_companies', {
		method: 'GET',
	});
}
	
