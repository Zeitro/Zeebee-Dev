import { Button, Space, Spin } from 'antd';
import MaturityAlertsItem from './Item';
import ListModal from './ListModal';
import SettingModal from './SettingModal';
import { useCallback, useEffect, useState } from 'react';
import { SettingOutlined } from '@ant-design/icons';
import { createRequest } from '@repo/shared-methods/utils';
import { NoData } from '@repo/shared-methods';

const MaturityAlerts = ({ list, source, handleGotoLoanDetails }) => {
	const [open, setOpen] = useState(false);
	const [openSetting, setOpenSetting] = useState(false);
	const [loading, setLoading] = useState(false);
	const [maturityList, setMaturityList] = useState([]);
	const [total, setTotal] = useState(0);
	const [dueDay, setDueDay] = useState(60); // 默认60天来计算到期时间
	const { request } = createRequest();
	// const [pagination, setPagination] = useState(getPagination());
	const getLoanList = useCallback(
		async ({ page, limit }) => {
			const res = await request('/loan/page', {
				method: 'POST',
				body: {
					page,
					limit,
					// sorter: {
					// 	column: 'maturity_date',
					// 	sort: 'desc',
					// },
					filter: {
						source,
						due_in_days: dueDay,
						date_type: 'maturity_date',
					},
				},
			});
			return res;
		},
		[dueDay, source],
	);

	const loadInitList = async () => {
		try {
			setLoading(true);
			const { total, list } = await getLoanList({ page: 1, limit: 2 });
			setMaturityList(list);
			setTotal(total);
			setLoading(false);
		} catch (error) {
			setLoading(false);
			throw error;
		}
	};

	useEffect(() => {
		const loadFn = () => loadInitList();
		loadFn();
		// eslint-disable-next-line react-hooks/exhaustive-deps
	}, [source]);

	return (
		<div className="border rounded-lg border-base-border lg:col-span-1 p-4 bg-background">
			<div className="flex items-center justify-between">
				<div className="flex text-xl font-medium text-base-text gap-2">
					Maturity alerts({total})
				</div>
				<Space size={8}>
					{total > 0 && (
						<Button size="small" onClick={() => setOpen(true)}>
							View all
						</Button>
					)}
					<Button
						size="small"
						icon={<SettingOutlined />}
						onClick={() => setOpenSetting(true)}
					></Button>
				</Space>
			</div>
			<Spin spinning={loading}>
				{total > 0 ? (
					maturityList.map((item) => (
						<MaturityAlertsItem
							key={item.id}
							item={item}
							handleGotoLoanDetails={handleGotoLoanDetails}
						/>
					))
				) : (
					<NoData imageProps={{ style: { width: 80, height: 80 } }} />
				)}
			</Spin>
			{open && (
				<ListModal
					setOpen={setOpen}
					list={list}
					source={source}
					getLoanList={getLoanList}
					handleGotoLoanDetails={handleGotoLoanDetails}
				/>
			)}
			{openSetting && (
				<SettingModal
					setOpen={setOpenSetting}
					dueDay={dueDay}
					setDueDay={setDueDay}
					handleSearch={loadInitList}
				/>
			)}
		</div>
	);
};
export default MaturityAlerts;
