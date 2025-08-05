import { create } from 'zustand';
import { createRequest } from '@/lib/fetch';
const { request } = createRequest();
import { getSelectOptions,getLender } from '@/api/common';
export const useOptionsStore = create((set) => ({
	lenderOptions:[],
	loanPendingStatus: [], //pending中的状态，传值loan_pending_status
	loanStatusOptions: [], // loan状态,传值loan_stage
	loanPurposeOptions: [],//loan_purpose
	loanTypeOptions: [], 	//loan_type
	propertyTypeOptions: [], //property_type
	lienPositionOptions: [], //lien_position
	businessStrategyOptions:[], //business_strategy
	getOptions: async (name, key) => {
		const { options } = await getSelectOptions(name);
		return key && set({
			[key]: options,
		});
	},
	getLender: async()=>{
		const {lender_companies}= await getLender();
		const lenderOptions = lender_companies.map(i=>({
			label:i,
			value:i
		}))
		set({lenderOptions})
	}
}));
