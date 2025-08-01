import { create } from 'zustand';
export const useLoanStore = create((set) => ({
	loanPendingStatus: [
		{
			label: 'Quoting',
			value: 'quoting',
		},
		{
			label: 'Processing',
			value: 'processing',
		},
		{
			label: 'Closing',
			value: 'closing',
		},
	],
	loanStatusOptions: [
		{
			label: 'Pending',
			value: 'pending',
		},
		{
			label: 'Active',
			value: 'active',
		},
		{
			label: 'Archived',
			value: 'archived',
		},
	],
	loanPurposeOptions: [
		{
			label: 'Purchase',
			value: 'purchase',
		},
		{
			label: 'Rate-term refi',
			value: 'rateTermRefi',
		},
		{
			label: 'Cash-out refi',
			value: 'cashOutRefi',
		},
	],
	loanTypeOptions: [
		{
			label: 'Bridge',
			value: 'bridge',
		},
		{
			label: 'Rehab',
			value: 'rehab',
		},
		{
			label: 'Ground up construction',
			value: 'groundUpConstruction',
		},
	],
	getOptions: (userData) => set({ loanPendingStatus: [] }),
}));
