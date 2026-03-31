const useColumn = (data, stage) => {
	const titleInfo = {
		Pending: 'Pending loans',
		Active: 'Active loans',
	};
	const pendingColumn = [];
	const activeColumn = [];
	const getColumnsByStage = (stage) => {
		switch (stage) {
			case 'Pending':
				return pendingColumn;
			case 'Active':
				return activeColumn;
			default:
				return [];
		}
	};

	const getDataSourceByStage = (stage) => {
		switch (stage) {
			case 'Pending':
				return data?.pendingDataSource || [];
			case 'Active':
				return data?.activeDataSource || [];
			default:
				return [];
		}
	};

	const totalList =
		{
			Pending: [
				{
					label: 'Total loan amount',
					value: 111,
				},
				{
					label: 'Average rate',
					value: 8,
				},
				{
					label: 'Estimated closing fund',
					value: 1111,
				},
			],
			Active: [
				{
					label: 'Total loan amount',
					value: 111,
				},
				{
					label: 'Outstanding loan amount',
					value: 111,
				},
				{
					label: 'Holdback(release/remaining)',
					value: 111,
				},
				{
					label: 'Average rate',
					value: 8,
				},
				{
					label: 'Monthly payment',
					value: 1111,
				},
			],
		}?.[stage] || [];
	return {
		pendingColumn,
		activeColumn,
		getColumnsByStage,
		getDataSourceByStage,
		titleInfo,
		totalList,
	};
};
export default useColumn;
