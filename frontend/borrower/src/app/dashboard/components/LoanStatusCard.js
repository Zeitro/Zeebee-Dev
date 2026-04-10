'use client';

import {BankOutlined, CheckCircleOutlined, ContainerOutlined, HistoryOutlined,} from '@ant-design/icons';

const LoanStatusCard = ({data, handleGotoLoans}) => {
    const list = [
        {
            title: 'All loans',
            icon: <BankOutlined style={{fontSize: 20, color: 'var(--base-primary)'}}/>,
            total: data?.total_loan_count,
            value: 'all',
            isNoClick: true,
        },
        {
            title: 'Active',
            icon: <CheckCircleOutlined style={{fontSize: 20, color: '#52C41A'}}/>,
            total: data?.active_loan_count,
            value: 'Active',
            dataType: 'updated_at',
        },
        {
            title: 'Pending',
            icon: <HistoryOutlined style={{fontSize: 20, color: '#FA8C16'}}/>,
            total: data?.pending_loan_count,
            value: 'Pending',
            dataType: 'closing_date',
        },
        {
            title: 'Archived',
            icon: <ContainerOutlined style={{fontSize: 20, color: '#44474C'}}/>,
            total: data?.archived_loan_count,
            value: 'Archived',
            dataType: 'updated_at',
        },
    ];

    return (
        <div className="lg:col-span-2 grid grid-cols-1 md:grid-cols-2 gap-2">
            {list.map((item) => {
                return (
                    <div
                        className="border rounded-lg border-base-border p-4 bg-background hover:shadow-lg"
                        style={{
                            cursor: item?.isNoClick ? 'no-drop' : 'pointer',
                        }}
                        key={item.value}
                        onClick={() => {
                            if (item?.isNoClick) return;
                            handleGotoLoans(item?.value, item?.dataType);
                        }}
                    >
                        <div className="flex items-center justify-between pb-4">
                            <div className="text-sm/[22px] text-base-label">{item.title}</div>
                            {item.icon}
                        </div>
                        <div className="flex items-center justify-between gap-2">
							<span className="text-3xl/10 font-medium text-base-text">
								{item.total || 0}
							</span>
                            {/* {!item?.isHide && (
								<Button
									type="text"
									onClick={() => handleGotoLoans(item?.value, item?.dataType)}
									size="small"
									className="!text-base-gray"
								>
									View all
									<ArrowRightOutlined />
								</Button>
							)} */}
                        </div>
                    </div>
                );
            })}
        </div>
    );
};
export default LoanStatusCard;
