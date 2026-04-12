import {create} from 'zustand';
import {getCityList, getStateList} from '@repo/shared-methods/api';

export const useAddressStore = create((set) => ({
	stateOptions: [],
	cityOptions: [],

	getStateList: async () => {
		const { list } = await getStateList();
		const stateOptions = list?.map((i) => ({
			label: i,
			value: i,
		}));
		set({ stateOptions });
	},
	getCityList: async (state) => {
		const { list } = await getCityList(state);
		const cityOptions = list?.map((i) => ({
			label: i,
			value: i,
		}));
		set({ cityOptions });
	},
}));
