import { create } from 'zustand';
import { createRequest } from '@/lib/fetch';
const { request } = createRequest();
import { getSelectOptions,getLender } from '@/api/common';
export const useLoanStore = create((set) => ({
	loanPendingStatus: [], //pending中的状态
	loanStatusOptions: [], // loan状态
	loanPurposeOptions: [],
	loanTypeOptions: [],
	lenderOptions:[],
	getOptions: async (name, key) => {
		const { options } = await getSelectOptions(name);
		return key && set({
			[key]: options,
		});
	},
	getLender: async()=>{
		const {lender_companies}= await getLender();
		set({lenderOptions: lender_companies})
	}
}));
