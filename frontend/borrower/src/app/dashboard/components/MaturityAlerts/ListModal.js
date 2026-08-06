import React, { useEffect, useState } from 'react';
import { NewModal, NewPagination } from '@repo/shared-methods';
import Item from './Item';
import { getPagination } from '@repo/shared-methods/utils';
import { Spin } from 'antd';
const ListModal = ({ setOpen, getLoanList, handleGotoLoanDetails }) => {
	const [pagination, setPagination] = useState(getPagination());
	const [loading, setLoading] = useState(false);
	const [maturityList, setMaturityList] = useState([]);
	const loadList = async ({ page, limit }) => {
		try {
			setLoading(true);
			const { total, list } = await getLoanList({ page, limit });
			setMaturityList(list);
			setPagination((v) => ({ ...v, total }));
			setLoading(false);
		} catch (error) {
			setLoading(false);
			throw error;
		}
	};

	const handleChangePagination = (page, pageSize) => {
		setPagination((v) => ({ ...v, current: page, pageSize: pageSize }));
		loadList({ page, limit: pageSize });
	};

	useEffect(() => {
		const loadFn = () => loadList({ page: pagination.current, limit: pagination.pageSize });
		loadFn(pagination);
	}, [pagination.current, pagination.pageSize]);
	return (
		<NewModal
			width={800}
			title="Maturity Alerts"
			onCancel={() => setOpen(false)}
			footer={null}
			className="[&_.ant-modal-content_.ant-modal-body]:!pt-0 [&_.ant-modal-header]:!mb-0"
		>
			<Spin spinning={loading}>
				<div className="flex flex-col gap-6 h-[60vh] min-h-[400px] max-h-[900px] overflow-hidden overflow-y-auto">
					{maturityList?.map((item, index) => (
						<Item
							key={index}
							item={item}
							handleGotoLoanDetails={handleGotoLoanDetails}
						/>
					))}
				</div>
				<NewPagination {...pagination} onChange={handleChangePagination} />
			</Spin>
		</NewModal>
	);
};

export default ListModal;
