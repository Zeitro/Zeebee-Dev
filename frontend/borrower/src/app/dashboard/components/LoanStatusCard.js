'use client';
const LoanStatusCard = ({ data, handleGotoLoans }) => {
	const list = [
		{
			title: 'All loans',
			icon: '',
			total: data?.total_loan_count,
			value: 'all',
		},
		{
			title: 'Active',
			icon: '',
			total: data?.active_loan_count,
			value: 'Active',
			dataType: 'updated_at',
		},
		{
			title: 'Pending',
			icon: '',
			total: data?.pending_loan_count,
			value: 'Pending',
			dataType: 'closing_date',
		},
		{
			title: 'Archived',
			icon: '',
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
						className="border rounded-lg border-base-border p-4 bg-background hover:shadow-lg cursor-pointer"
						key={item.value}
						onClick={() => handleGotoLoans(item?.value, item?.dataType)}
					>
						<div className="flex items-center justify-between pb-4">
							<div className="text-sm/[22px] text-base-label">{item.title}</div>
							{item.icon}
						</div>
						<div className="flex items-center justify-between gap-2">
							<span className="text-3xl/10 font-medium text-base-text">
								{item.total}
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
