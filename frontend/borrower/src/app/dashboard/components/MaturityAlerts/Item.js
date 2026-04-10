import {maturityAlertsStatusColor} from '@repo/shared-methods/enum';

const MaturityAlertsItem = ({item, handleGotoLoanDetails}) => {
    const statusText = {
        'Due Soon': `Due in ${item?.due_days} days`,
        Overdue: `Overdue ${item?.due_days} days`,
    };
    return (
        <div
            className="border-b last-of-type:border-b-0 border-base-border py-5 hover:bg-[#f5f5f5] rounded-xl px-1 cursor-pointer"
            onClick={() => handleGotoLoanDetails(item)}
        >
            <div className="flex items-center justify-between">
                <div className="flex text-base-primary font-bold gap-2">
                    {item?.id}
                    {item.source === 'Borrower' && (
                        <div className="px-2 bg-blue-light-1 text-blue-medium-default rounded-sm font-medium">
                            3rd party
                        </div>
                    )}
                </div>

                <span className="text-xl font-medium text-base-text">{item?.maturity_date}</span>
            </div>
            <div className="flex items-center justify-between gap-2">
                <span className="common-ellipsis">{item?.loan_property?.address?.address1}</span>
                {item?.due_status && (
                    <div
                        style={{
                            ...(maturityAlertsStatusColor[item?.due_status] || {}),
                            borderRadius: 2,
                            padding: '0 8px',
                            fontSize: 12,
                        }}
                    >
                        {statusText[item?.due_status]}
                    </div>
                )}
            </div>
        </div>
    );
};
export default MaturityAlertsItem;
