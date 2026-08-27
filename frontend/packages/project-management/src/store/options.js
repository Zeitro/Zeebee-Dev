import { create } from 'zustand';
import { getSelectOptions } from '@repo/shared-methods/api';

export const useOptionsStore = create((set) => ({
	priorityOptions: [], // project_management_priority
	stageOptions: [], // project_management_stage_phase
	getOptions: async (name, key) => {
		const { options } = await getSelectOptions(name);
		return (
			key &&
			set({
				[key]: options,
			})
		);
	},
}));
