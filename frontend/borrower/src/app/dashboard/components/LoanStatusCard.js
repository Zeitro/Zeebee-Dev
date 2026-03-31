'use client';
import { Button } from 'antd';
import { ArrowRightOutlined } from '@ant-design/icons';
import { useRouter } from 'next/navigation';
import { sessionStorageUtil } from '@repo/shared-methods/utils';
const LoanStatusCard = ({ dayCount }) => {
	const router = useRouter();
	const list = [
		{
			title: 'All loans',
			icon: '',
			total: 5,
			value: 'all',
		},
		{
			title: 'Active',
			icon: '',
			total: 5,
			value: 'Active',
		},
		{
			title: 'Pending',
			icon: '',
			total: 5,
			value: 'Pending',
		},
		{
			title: 'Archived',
			icon: '',
			total: 5,
			value: 'Archived',
		},
	];

	const handleClick = (item) => {
		router.push('/loans');
		sessionStorageUtil.set('LOANS-FILTER', { stage: item.value, dayCount });
	};
	return (
		<div className="w-full lg:w-2/3 flex flex-wrap">
			{list.map((item) => {
				return (
					<div
						className="border rounded-lg border-base-border w-full md:w-1/2"
						key={item.value}
					>
						<div className="flex items-center justify-between">
							<div>{item.title}</div>

							{item.icon}
						</div>
						<div className="flex items-center justify-between">
							<span>{item.total}</span>
							<Button type="text" onClick={() => handleClick(item)}>
								View all
								<ArrowRightOutlined />
							</Button>
						</div>
					</div>
				);
			})}
		</div>
	);
};
export default LoanStatusCard;
