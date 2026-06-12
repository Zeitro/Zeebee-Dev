import { useState } from 'react';
import { createRequest } from '@repo/shared-methods/utils';
const useData = () => {
	const { request } = createRequest();
	const [loading, setLoading] = useState(false);
	const [apiTreeData, setApiTreeData] = useState([]);

	const changeTreeLabel = (list) => {
		return list.map((item) => {
			item.label = item?.path ? (
				<>
					{item.name} (<span className="text-base-label"> {item.path}</span>)
				</>
			) : (
				item.name
			);
			item.value = item.id;
			if (item.children) {
				item.children = changeTreeLabel(item.children);
			}
			return item;
		});
	};

	const getApiTree = async (v) => {
		try {
			setLoading(true);
			const payload = new URLSearchParams({ type: v });
			const api = `/api/tree?${payload}`;
			const res = await request(api);
			// const apiList = changeTreeLabel(res?.list);
			setApiTreeData(res?.list);
			setLoading(false);
		} catch (error) {
			setLoading(false);
			throw error;
		}
	};
	return { getApiTree, apiTreeData, loading };
};

export default useData;
