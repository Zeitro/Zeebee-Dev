import { useOptionsStore } from '@/store/options';
import { useCallback, useEffect, useRef } from 'react';
/**
 * list = [{
 *      key: loan_purpose,
 *      valueName: loanPurposeOptions
 * }]
 */
export const useGetOptions = (list = []) => {
	const optionsInfo = useOptionsStore();
	const isInitGet = useRef(true);
	const { getOptions, getLender, ...state } = optionsInfo;
	const getInitOptions = useCallback(async () => {
		await Promise.all(
			list.map(async (item) => {
				try {
					if (!optionsInfo?.[item.valueName]?.length && item.key) {
						return await optionsInfo?.getOptions(item.key, item.valueName);
					}
					return Promise.resolve();
				} catch (error) {
					return Promise.resolve();
				}
			}),
		);
	}, [getOptions, list, state]);

	useEffect(() => {
        console.log('options-get-init');
        
		if (list.length && isInitGet.current) {
			getInitOptions();
			isInitGet.current = false;
		}
	}, [JSON.stringify(list)]);

	// useEffect(() => {
	// 	if (!loanPurposeOptions?.length && loanPurposeName) {
	// 		getOptions('loan_purpose', loanPurposeName);
	// 	}
	// }, [loanPurposeOptions, loanPurposeName]);

	// useEffect(() => {
	// 	if (!loanTypeOptions?.length && loanTypeName) {
	// 		getOptions('loan_type', loanTypeName);
	// 	}
	// }, [loanTypeOptions, loanTypeName]);

	// useEffect(() => {
	// 	if (!loanPendingStatus?.length && pendingStatusName) {
	// 		getOptions('loan_pending_status', 'loanPendingStatus');
	// 	}
	// }, [loanPendingStatus, pendingStatusName]);

	// useEffect(() => {
	// 	!loanStatusOptions?.length && getOptions('loan_stage', 'loanStatusOptions');
	// }, [loanStatusOptions]);
	return optionsInfo;
};
