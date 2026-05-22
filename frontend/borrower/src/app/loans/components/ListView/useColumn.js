'use client';
import { formatAmount, formatTime } from '@repo/shared-methods/utils';
import {
	LoanStage,
	LoanPendingStatus,
} from '@repo/shared-methods/componentsUi/zeeBeeDetails/components/ChangeStageStatus';
import { FieldTips, TableOperation } from '@repo/shared-methods';
import { PERMISSIONS } from '@config/permissions';
import { useAuth } from '@repo/shared-methods/hooks';
import { SendOutlined } from '@ant-design/icons';
import QuoteBtn from '../Quote';

export default function useColumn({
	isNotAddress,
	successCallback,
	showDeleteModal,
	stage,
	searchQuery,
}) {
	const { hasPermission, getActiveByPermission } = useAuth();
	const get3rdPartyMark = (v, { source }) => (
		<div className="flex items-center">
			<p className="overflow-hidden text-ellipsis whitespace-nowrap" title={v}>
				{v}
			</p>
			{source === 'Borrower' && (
				<div className="ml-2 px-2 bg-blue-light-1 text-blue-medium-default rounded-xs">
					3rd party
				</div>
			)}
		</div>
	);
	const common = (isNotAddress) => {
		const addressArr = [
			{
				title: 'Loan ID',
				dataIndex: 'id',
				key: 'id',
				width: 160,
				fixed: 'left',
				render: get3rdPartyMark,
			},
			{
				title: 'Property address',
				dataIndex: ['loan_property', 'address', 'address1'],
				key: 'address',
				ellipsis: true,
				sorter: true,
				sortDirections: ['ascend', 'descend'],
				// defaultSortOrder: 'descend',
				fixed: 'left',
				width: 260,
			},
			{
				title: 'City',
				dataIndex: ['loan_property', 'address', 'city'],
				key: 'city',
				width: 120,
				ellipsis: true,
			},
			{
				title: 'State',
				dataIndex: ['loan_property', 'address', 'state'],
				key: 'state',
				width: 70,
				ellipsis: true,
			},
			{
				title: 'Zip code',
				dataIndex: ['loan_property', 'address', 'zip_code'],
				key: 'zip_code',
				ellipsis: true,
				width: 120,
			},
		];
		return [
			...(isNotAddress ? [] : addressArr),
			{
				title: 'Created time',
				dataIndex: 'created_at',
				key: 'created_at',
				ellipsis: true,
				sorter: true,
				width: 170,
				defaultSortOrder: searchQuery?.sorter?.column ? null : 'descend',
				render: (v) => formatTime(v) || '--',
			},
			{
				title: 'Borrower',
				dataIndex: 'applicants',
				key: 'applicants',
				ellipsis: true,
				render: (v, row) => {
					const value = v?.map((i) => i.full_name)?.join(',');
					if (isNotAddress) return get3rdPartyMark(value, row);
					return value;
				},
				width: isNotAddress ? 200 : 160,
			},
		];
	};

	const getButtonOptions = () => {
		// const commonBtn = [
		// 	{
		// 		permission: PERMISSIONS.LOAN.EDIT_3RD_PARTY,
		// 		width: 58,
		// 		type: '3rdParty',
		// 	},
		// 	{
		// 		permission: PERMISSIONS.LOAN.DELETE,
		// 		width: 74,
		// 		type: '3rdParty',
		// 	},
		// 	{
		// 		permission: PERMISSIONS.LOAN.DETAILS,
		// 		width: 77,
		// 		type: 'zeeBee',
		// 	},
		// ];
		const commonBtn = [
			PERMISSIONS.LOAN.EDIT_3RD_PARTY,
			PERMISSIONS.LOAN.DELETE,
			PERMISSIONS.LOAN.DETAILS,
		];
		switch (stage) {
			case 'Draft':
				return [
					PERMISSIONS.LOAN.EDIT_3RD_PARTY,
					PERMISSIONS.LOAN.EDIT_ZEEBEE_DRAFT,
					PERMISSIONS.LOAN.DELETE,
				];
			// return [
			// 	{
			// 		permission: PERMISSIONS.LOAN.EDIT_3RD_PARTY,
			// 		width: 58,
			// 		type: '3rdParty',
			// 	},
			// 	{
			// 		permission: PERMISSIONS.LOAN.EDIT_ZEEBEE_DRAFT,
			// 		width: 58,
			// 		type: 'zeeBee',
			// 	},
			// 	{
			// 		permission: PERMISSIONS.LOAN.DELETE,
			// 		width: 74,
			// 	},
			// ];
			case 'Pending':
				return [PERMISSIONS.LOAN.QUOTE, ...commonBtn];
			// return [
			// 	{
			// 		permission: PERMISSIONS.LOAN.QUOTE,
			// 		width: 125,
			// 		type: '3rdParty',
			// 	},
			// 	...commonBtn,
			// ];
			default:
				return commonBtn;
		}
	};

	const getStageOpt = (successCallback) => {
		return {
			title: 'Stage',
			dataIndex: 'stage',
			key: 'stage',
			ellipsis: true,
			width: 150,
			align: 'center',
			fixed: 'right',
			render: (v, record) => {
				let title = '';
				if (
					record?.source !== 'Borrower' &&
					!hasPermission(PERMISSIONS?.LOAN?.CHANGE_STAGE)
				) {
					title = "You don't have permission to change the stage.";
				} else if (record?.source !== 'Borrower') {
					title = 'Only the 3rd part can change the stage.';
				}
				return (
					<div onClick={(e) => e.stopPropagation()}>
						<LoanStage
							data={record}
							successCallback={successCallback}
							disabled={
								record?.source !== 'Borrower' ||
								!hasPermission(PERMISSIONS?.LOAN?.CHANGE_STAGE)
							}
							tooltipProps={{
								title,
							}}
						/>
					</div>
				);
			},
		};
	};
	const getOpt = (showDeleteModal) => {
		return getActiveByPermission({
			// width: 290,
			buttonOptions: getButtonOptions(),
			render: (_, record) => (
				<TableOperation
					record={record}
					getItems={(editItem, deleteItem, detailsItem) => {
						const arr = [];
						// Draft (zeebee和3rdParty都是编辑、删除按钮)
						if (stage === 'Draft') {
							// const url =
							// 	record?.source === 'Borrower'
							// 		? `/loans/${record?.id}`
							// 		: `/loans/createZeeBeeLoan?id=${record?.id}`;
							// arr.push({ ...editItem, url });
							// arr.push(deleteItem);
							if (
								hasPermission(PERMISSIONS?.LOAN?.EDIT_3RD_PARTY) &&
								record?.source === 'Borrower'
							) {
								arr.push({ ...editItem, url: `/loans/${record?.id}` });
							}
							if (
								hasPermission(PERMISSIONS?.LOAN?.EDIT_ZEEBEE_DRAFT) &&
								record?.source !== 'Borrower'
							) {
								arr.push({
									...editItem,
									url: `/loans/createZeeBeeLoan?id=${record?.id}`,
								});
							}
							if (hasPermission(PERMISSIONS?.LOAN?.DELETE)) {
								arr.push(deleteItem);
							}
							return arr;
						}
						// ZeeBee 除Draft状态下，只有详情按钮
						if (record?.source !== 'Borrower') {
							return hasPermission(PERMISSIONS?.LOAN?.DETAILS)
								? [{ ...detailsItem, url: `/loans/zeeBeeDetails?id=${record?.id}` }]
								: [];
						}

						// 3rdParty除Draft状态下，有询价、编辑、删除按钮
						const quoteDisabled = !(
							record?.stage === 'Pending' && record?.pending_status === 'Quoting'
						);
						if (
							hasPermission(PERMISSIONS?.LOAN?.QUOTE) &&
							record?.stage === 'Pending'
						) {
							arr.push({
								children:
									record?.source === 'Borrower' ? (
										<QuoteBtn
											key="quote"
											disabled={quoteDisabled}
											rowData={record}
											tooltipProps={{ placement: 'leftTop' }}
											type="link"
											ghost={false}
											size="small"
										>
											<SendOutlined />
											Request quote
										</QuoteBtn>
									) : (
										''
									),
								key: 'quote',
								disabled: quoteDisabled,
							});
						}
						if (hasPermission(PERMISSIONS?.LOAN?.EDIT_3RD_PARTY)) {
							arr.push({
								...editItem,
								url: `/loans/${record?.id}`,
							});
						}
						if (hasPermission(PERMISSIONS?.LOAN?.DELETE)) {
							arr.push(deleteItem);
						}

						return arr;
					}}
					onSelect={(key) => {
						if (key === 'delete') {
							showDeleteModal(record);
						}
					}}
				/>
			),
		});
	};

	const pendingColumns = (spliceIndex, spliceList, customColumnsFn) => {
		const status = {
			title: 'Status',
			dataIndex: 'pending_status',
			key: 'pending_status',
			ellipsis: true,
			width: 150,
			align: 'center',
			fixed: 'right',
			render: (v, record) =>
				hasPermission(PERMISSIONS?.LOAN?.CHANGE_PENDING_STATUS) ? (
					<LoanPendingStatus
						data={record}
						successCallback={successCallback}
						disabled={record?.source !== 'Borrower'}
						tooltipProps={{
							title: 'Zeebee-connected loans are managed automatically and cannot be edited manually by users.',
						}}
					/>
				) : (
					v
				),
		};
		let arr = [
			...common(isNotAddress),
			{
				title: 'Loan purpose',
				dataIndex: 'loan_purpose',
				key: 'loan_purpose',
				ellipsis: true,
				width: 120,
				render: (v) => v || '--',
			},
			{
				title: 'Target closing date',
				dataIndex: 'closing_date',
				key: 'closing_date',
				ellipsis: true,
				sorter: true,
				width: 170,
				defaultSortOrder: searchQuery?.sorter?.column === 'closing_date' ? 'descend' : null,
				render: (v) => v || '--',
			},
			{
				title: 'Business strategy',
				dataIndex: 'business_strategy',
				key: 'business_strategy',
				ellipsis: true,
				render: (v) => v?.join(',') || '--',
				width: 160,
			},
			{
				title: 'Target lender',
				dataIndex: 'lender_company',
				key: 'lender_company',
				ellipsis: true,
				width: 180,
			},
			{
				title: 'Total loan amount',
				dataIndex: 'total_amount',
				key: 'total_amount',
				ellipsis: true,
				sorter: true,
				width: 170,
				render: (v) => formatAmount(v),
			},
			{
				title: 'Initial loan',
				dataIndex: 'initial_amount',
				key: 'initial_amount',
				ellipsis: true,
				sorter: true,
				width: 120,
				render: (v) => formatAmount(v),
			},
			{
				title: 'Hold back',
				dataIndex: 'holdback_amount',
				key: 'holdback_amount',
				ellipsis: true,
				sorter: true,
				width: 120,
				render: (v) => formatAmount(v),
			},
			{
				title: 'Estimated Closing Fund',
				dataIndex: 'closing_fund',
				key: 'closing_fund',
				ellipsis: true,
				sorter: true,
				width: 200,
				render: (v) => formatAmount(v),
			},
			{
				title: 'Rate',
				dataIndex: 'rate',
				key: 'rate',
				ellipsis: true,
				sorter: true,
				width: 90,
				render: (v) => (v ? `${v}%` : '--'),
			},
			{
				title: 'Loan type',
				dataIndex: 'loan_type',
				key: 'loan_type',
				ellipsis: true,
				width: 180,
			},
			getStageOpt(successCallback),
			...(stage === 'Draft' ? [] : [status]),
			getOpt(showDeleteModal, stage),
		];
		if (spliceList) {
			arr.splice(spliceIndex, 0, ...spliceList);
		}
		if (customColumnsFn) {
			return customColumnsFn(arr);
		}
		return arr;
	};

	const activeColumns = (spliceIndex, spliceList, customColumnsFn) => {
		let arr = [
			...common(isNotAddress),
			{
				title: 'Lien Position',
				dataIndex: 'lien_position',
				key: 'lien_position',
				ellipsis: true,
				width: 120,
				render: (v) => v || '--',
			},
			{
				title: 'Lender',
				dataIndex: 'lender_company',
				key: 'lender_company',
				ellipsis: true,
				width: 180,
			},
			{
				title: 'Loan type',
				dataIndex: 'loan_type',
				key: 'loan_type',
				ellipsis: true,
				width: 180,
			},
			{
				title: 'Business strategy',
				dataIndex: 'business_strategy',
				key: 'business_strategy',
				ellipsis: true,
				render: (v) => v?.join(','),
				width: 160,
			},
			{
				title: 'Rate',
				dataIndex: 'rate',
				key: 'rate',
				ellipsis: true,
				sorter: true,
				width: 90,
				render: (v) => (v ? `${v}%` : '--'),
			},
			// {
			// 	title: 'Point',
			// 	dataIndex: 'point',
			// 	key: 'point',
			// 	ellipsis: true,
			// },
			{
				title: 'Maturity Date',
				dataIndex: 'maturity_date',
				key: 'maturity_date',
				ellipsis: true,
				sorter: true,
				width: 140,
				render: (v) => v || '--',
			},
			{
				title: 'Outstanding Loan Amount',
				dataIndex: 'principal_amount',
				key: 'principal_amount',
				ellipsis: true,
				sorter: true,
				width: 220,
				render: (v) => formatAmount(v),
			},
			{
				title: 'Total loan amount',
				dataIndex: 'total_amount',
				key: 'total_amount',
				ellipsis: true,
				sorter: true,
				width: 180,
				render: (v) => formatAmount(v),
			},
			{
				title: 'Hold back (Release/Remaining)',
				dataIndex: 'holdback_amount',
				key: 'holdback_amount',
				ellipsis: true,
				render: (v, { holdback_released_amount, holdback_remaining_amount }) => (
					<>
						<span>{formatAmount(v)}</span>
						<div className="text-base-gray text-xs">
							({formatAmount(holdback_released_amount)}/
							{formatAmount(holdback_remaining_amount)})
						</div>
					</>
				),
			},
			// {
			// 	title: 'Released Holdback',
			// 	dataIndex: 'releasedHoldback',
			// 	key: 'releasedHoldback',
			// 	ellipsis: true,
			// 	sorter: true,
			// 	sortDirections: ['ascend', 'descend'],
			// },
			{
				title: 'Monthly Payment',
				dataIndex: 'monthly_payment',
				key: 'monthly_payment',
				ellipsis: true,
				sorter: true,
				width: 160,
				render: (v) => formatAmount(v),
			},
			{
				title: 'Prepayment Penalty',
				dataIndex: 'prepayment_penalty',
				key: 'prepayment_penalty',
				ellipsis: true,
				width: 160,
				render: (v) => `${v} months interest`,
			},
			{
				title: 'Servicing Company',
				dataIndex: 'service_company',
				key: 'service_company',
				ellipsis: true,
				width: 160,
				render: (v) => v || '--',
			},
			getStageOpt(successCallback),
			getOpt(showDeleteModal),
		];
		if (spliceList) {
			arr.splice(spliceIndex, 0, ...spliceList);
		}
		if (customColumnsFn) {
			return customColumnsFn(arr);
		}
		return arr;
	};

	const archivedColumns = (spliceIndex, spliceList, customColumnsFn) => {
		let arr = [
			...common(isNotAddress),
			{
				title: 'Loan type',
				dataIndex: 'loan_type',
				key: 'loan_type',
				ellipsis: true,
				width: 180,
			},
			{
				title: 'Principal Balance',
				dataIndex: 'principal_amount',
				key: 'principal_amount',
				ellipsis: true,
				sorter: true,
				width: 180,
				render: (v) => formatAmount(v),
			},
			{
				title: 'Lender',
				dataIndex: 'lender_company',
				key: 'lender_company',
				ellipsis: true,
				width: 180,
			},
			{
				title: 'Rate',
				dataIndex: 'rate',
				key: 'rate',
				ellipsis: true,
				sorter: true,
				width: 90,
				render: (v) => (v ? `${v}%` : '--'),
			},
			{
				title: 'Paid Off Date',
				dataIndex: 'paid_off_date',
				key: 'paid_off_date',
				ellipsis: true,
				sorter: true,
				width: 160,
				render: (v) => v || '--',
			},
			getStageOpt(successCallback),
			getOpt(showDeleteModal),
		];
		if (spliceList) {
			arr.splice(spliceIndex, 0, ...spliceList);
		}
		if (customColumnsFn) {
			return customColumnsFn(arr);
		}
		return arr;
	};

	const getColumnsByStage = () => {
		switch (stage) {
			case 'Pending':
			case 'Draft':
				return pendingColumns();
			case 'Active':
				return activeColumns();
			case 'Archived':
				return archivedColumns();
			default:
				return [];
		}
	};

	return {
		pendingColumns,
		activeColumns,
		archivedColumns,
		get3rdPartyMark,
		getColumnsByStage,
	};
}
