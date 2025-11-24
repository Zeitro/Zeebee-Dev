import { create } from 'zustand';
import { getSelectOptions, getLenders, getBorrowers } from '@repo/shared-methods/api';
export const useOptionsStore = create((set) => ({
	lenderOptions: [],
	borrowerOptions: [],
	loanPendingStatus: [], //pending中的状态，传值loan_pending_status
	loanStatusOptions: [], // loan状态,传值loan_stage
	loanPurposeOptions: [], //loan_purpose
	loanTypeOptions: [], //loan_type
	propertyTypeOptions: [], //property_type
	lienPositionOptions: [], //lien_position
	businessStrategyOptions: [], //business_strategy
	prepaymentPenaltyOptions: [], //loan_prepayment_penalty
	borrowerTypeOptions: [], //applicant_type
	maritalStatusOptions: [], //marital_status
	serviceTypeOptions: [], //loan_third_service_type
	sexOptions: [], //applicant_sex
	bankruptcyChapterOptions: [], //bankruptcy_chapter
	getOptions: async (name, key) => {
		const { options } = await getSelectOptions(name);
		return (
			key &&
			set({
				[key]: options,
			})
		);
	},
	getLender: async () => {
		const { lender_companies } = await getLenders();
		const lenderOptions = lender_companies.map((i) => ({
			label: i,
			value: i,
		}));
		set({ lenderOptions });
	},
	getBorrowerOptions: async () => {
		const { borrower_names } = await getBorrowers();
		const borrowerOptions = borrower_names.map((i) => ({
			label: i,
			value: i,
		}));
		set({ borrowerOptions });
	},
}));
