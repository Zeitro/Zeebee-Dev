import {FieldTips} from '@repo/shared-methods';
import {assemblyAddress, formatAmount, isNumberFn} from '@repo/shared-methods/utils';

const useColumn = (data, stage) => {
    const {active_loan_summary, pending_loan_summary} = data;
    const titleInfo = {
        Pending: 'Pending loans',
        Active: 'Active loans',
    };
    const get3rdPartyMark = (v, {source}) => (
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
    const pendingColumn = [
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
            width: 260,
            render: (v, {loan_property}) => {
                return <FieldTips title={assemblyAddress(loan_property?.address)}>{v}</FieldTips>;
            },
        },
        {
            title: 'Target closing date',
            dataIndex: 'closing_date',
            key: 'closing_date',
            ellipsis: true,
            width: 170,
            render: (v) => v || '--',
        },
        {
            title: 'Total loan amount',
            dataIndex: 'total_amount',
            key: 'total_amount',
            ellipsis: true,
            width: 170,
            render: (v) => formatAmount(v),
        },
        {
            title: 'Rate',
            dataIndex: 'rate',
            key: 'rate',
            ellipsis: true,
            width: 90,
            render: (v) => (v ? `${v}%` : '--'),
        },
        {
            title: 'Hold back',
            dataIndex: 'holdback_amount',
            key: 'holdback_amount',
            ellipsis: true,
            width: 120,
            render: (v) => formatAmount(v),
        },
        {
            title: 'Loan type',
            dataIndex: 'loan_type',
            key: 'loan_type',
            ellipsis: true,
            width: 180,
        },
    ];
    const activeColumn = [
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
            width: 260,
        },
        {
            title: 'Total loan amount',
            dataIndex: 'total_amount',
            key: 'total_amount',
            ellipsis: true,
            width: 180,
            render: (v) => formatAmount(v),
        },
        {
            title: 'Rate',
            dataIndex: 'rate',
            key: 'rate',
            ellipsis: true,
            width: 90,
            render: (v) => (v ? `${v}%` : '--'),
        },
        {
            title: 'Outstanding Loan Amount',
            dataIndex: 'principal_amount',
            key: 'principal_amount',
            ellipsis: true,
            width: 220,
            render: (v) => formatAmount(v),
        },
        {
            title: 'Hold back (Release/Remaining)',
            dataIndex: 'holdback_amount',
            key: 'holdback_amount',
            ellipsis: true,
            render: (v, {holdback_released_amount, holdback_remaining_amount}) => (
                <>
                    <span>{formatAmount(v)}</span>
                    <div className="text-base-gray text-xs">
                        ({formatAmount(holdback_released_amount)}/
                        {formatAmount(holdback_remaining_amount)})
                    </div>
                </>
            ),
        },
        {
            title: 'Monthly Payment',
            dataIndex: 'monthly_payment',
            key: 'monthly_payment',
            ellipsis: true,
            width: 160,
            render: (v) => formatAmount(v),
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
    ];
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

    const getDataSourceByStage = () => {
        switch (stage) {
            case 'Pending':
                return pending_loan_summary?.list || [];
            case 'Active':
                return active_loan_summary?.list || [];
            default:
                return [];
        }
    };

    const totalList =
        {
            Pending: [
                {
                    label: 'Total loan amount',
                    value: formatAmount(pending_loan_summary?.total_loan_amount),
                },
                {
                    label: 'Average rate',
                    value: isNumberFn(pending_loan_summary?.average_rate)
                        ? `${pending_loan_summary?.average_rate} %`
                        : '--',
                },
                {
                    label: 'Estimated closing fund',
                    value: formatAmount(pending_loan_summary?.total_closing_fund_amount),
                },
            ],
            Active: [
                {
                    label: 'Total loan amount',
                    value: formatAmount(active_loan_summary?.total_loan_amount),
                },
                {
                    label: 'Outstanding loan amount',
                    value: formatAmount(active_loan_summary?.total_principal_amount),
                },
                {
                    label: 'Holdback(release/remaining)',
                    render: () => {
                        return (
                            <div className="flex flex-wrap">
								<span>
									{formatAmount(active_loan_summary?.total_holdback_amount)}
								</span>
                                <span className="text-base-gray text-lg">
									({formatAmount(active_loan_summary?.total_released_amount)}/
                                    {formatAmount(active_loan_summary?.total_remaining_amount)})
								</span>
                            </div>
                        );
                    },
                },
                {
                    label: 'Average rate',
                    value: isNumberFn(active_loan_summary?.average_rate)
                        ? `${active_loan_summary?.average_rate} %`
                        : '--',
                },
                {
                    label: 'Monthly payment',
                    value: formatAmount(active_loan_summary?.total_monthly_payment_amount),
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
